package service

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"os"
)

type UserService struct {
	cognitoUserPool     string
	cognitoClientID     string
	cognitoClientSecret string
	cognito             *cognitoidentityprovider.CognitoIdentityProvider
	awsRegion           string
	userRepository      *repository.UserRepository
}

const AuthFlowAdminNoSRP = "ADMIN_NO_SRP_AUTH"
const AuthFlowRefreshToken = "REFRESH_TOKEN_AUTH"
const AuthResponseChallenge = "NEW_PASSWORD_REQUIRED"
const JwkTokenUrl = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"

func CreateDefaultUserService() *UserService {
	return CreateUserService(repository.CreateUserRepository(db.GetConnection(
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DBNAME"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"))))
}

func CreateUserService(userRepository *repository.UserRepository) *UserService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &UserService{
		cognito:             cognitoidentityprovider.New(sess),
		cognitoUserPool:     os.Getenv("USER_POOL_ID"),
		cognitoClientID:     os.Getenv("COGNITO_CLIENT_ID"),
		cognitoClientSecret: os.Getenv("COGNITO_CLIENT_SECRET"),
		awsRegion:           os.Getenv("AWS_REGION"),
		userRepository:      userRepository,
	}
}

func (s *UserService) CreateUser(newUser *model.NewUser) (*entity.User, error) {
	response, err := s.cognito.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		Username:  aws.String(newUser.Email),
		TemporaryPassword: aws.String(newUser.Password),
		UserPoolId: aws.String(s.cognitoUserPool),
	})

	if err != nil {
		return nil, errors.New(err.Error())
	}

	user := entity.CreateUser(newUser, *response.User.Username)
	s.userRepository.Create(user)
	return user, nil
}

func (s *UserService) CreateSession(newSession *model.NewSession) *AuthResponse {
	response, err := s.cognito.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:          aws.String(AuthFlowAdminNoSRP),
		AuthParameters:    map[string]*string{
			"USERNAME": aws.String(newSession.Email),
			"PASSWORD": aws.String(newSession.Password),
		},
		ClientId:          aws.String(s.cognitoClientID),
		UserPoolId:        aws.String(s.cognitoUserPool),
	})

	if err != nil {
		return createAuthFailedSessionResponse()
	}

	user := s.userRepository.GetUserFromEmail(newSession.Email)

	if response.AuthenticationResult != nil {
		s.updateUserTokens(user, response.AuthenticationResult)
		return createSessionResponse(response)
	}

	s.updateUserWithCreateSessionResult(user, response)
	return createChallengeSessionResponse(response)
}

func (s *UserService) ProvideChallengeResponse(passwordReset *model.PasswordReset) *AuthResponse {
	user := s.userRepository.GetUserFromEmail(passwordReset.Email)

	if user == nil {
		return createAuthFailedSessionResponse()
	}

	data := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName:      aws.String(AuthResponseChallenge),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(passwordReset.Email),
			"NEW_PASSWORD": aws.String(passwordReset.Password),
		},
		ClientId:           aws.String(s.cognitoClientID),
		Session:            aws.String(user.LastSessionToken),
		UserPoolId:         aws.String(s.cognitoUserPool),
	}

	response, err := s.cognito.AdminRespondToAuthChallenge(data)

	if err != nil {
		return createAuthFailedSessionResponse()
	}

	if response.AuthenticationResult != nil {
		s.updateUserTokens(user, response.AuthenticationResult)
	}

	return createChallengeResponse(response)
}

func (s *UserService) ValidateSessionToken(sessionToken *model.SessionToken) bool {
	keySet, jwkErr := jwk.Fetch(fmt.Sprintf(JwkTokenUrl, s.awsRegion, s.cognitoUserPool))
	if jwkErr != nil {
		return false
	}

	token, parseErr := jwt.Parse(sessionToken.Token, func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		keys := keySet.LookupKeyID(kid)
		return keys[0].Materialize()
	})
	if parseErr != nil {
		log.Print(parseErr)
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := claims.Valid(); err != nil || claims.VerifyAudience(s.cognitoClientID, false) == false {
		log.Print(err)
		return false
	}

	response, err := s.cognito.GetUser(&cognitoidentityprovider.GetUserInput{ AccessToken: aws.String(sessionToken.Token) })
	if err != nil {
		log.Print(err)
		return false
	}

	user := s.userRepository.GetUserFromSessionToken(sessionToken.Token)
	if user == nil {
		log.Print("user does not match jwt")
		return false
	}

	return user.CognitoId.String() == *response.Username
}

func (s *UserService) RefreshSession(sessionRefresh *model.SessionRefresh) *AuthResponse {
	user := s.userRepository.GetUserFromSessionToken(sessionRefresh.Token)

	if user == nil {
		return createAuthFailedSessionResponse()
	}

	result, err := s.cognito.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:          aws.String(AuthFlowRefreshToken),
		AuthParameters:    map[string]*string{
			"REFRESH_TOKEN": aws.String(user.LastRefreshToken),
			"DEVICE_KEY": 	 aws.String(user.DeviceKey),
		},
		ClientId:          aws.String(s.cognitoClientID),
		UserPoolId:        aws.String(s.cognitoUserPool),
	})

	if err != nil {
		return createAuthFailedSessionResponse()
	}

	s.updateUserTokens(user, result.AuthenticationResult)
	return createSuccessfulRefreshResponse(result)
}

func (s *UserService) updateUserWithCreateSessionResult(user *entity.User, result *cognitoidentityprovider.AdminInitiateAuthOutput) {
	user.SRP = *result.ChallengeParameters["USER_ID_FOR_SRP"]
	user.LastSessionToken = *result.Session
	s.userRepository.Save(user)
}

func (s *UserService) updateUserTokens(user *entity.User, result *cognitoidentityprovider.AuthenticationResultType) {
	if result.NewDeviceMetadata != nil {
		user.DeviceGroupKey = *result.NewDeviceMetadata.DeviceGroupKey
		user.DeviceKey = *result.NewDeviceMetadata.DeviceKey
	}
	user.LastAccessToken = *result.AccessToken
	user.LastIdToken = *result.IdToken
	if result.RefreshToken != nil {
		user.LastRefreshToken = *result.RefreshToken
	}
	s.userRepository.Save(user)
}
