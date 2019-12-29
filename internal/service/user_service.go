package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/mapper"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/segmentio/kafka-go"
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
	kafkaWriter         *kafka.Writer
}

const AuthFlowAdminNoSRP = "ADMIN_NO_SRP_AUTH"
const AuthFlowRefreshToken = "REFRESH_TOKEN_AUTH"
const AuthResponseChallenge = "NEW_PASSWORD_REQUIRED"
const JwkTokenUrl = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"

func CreateDefaultUserService() *UserService {
	return CreateUserService(
		repository.CreateUserRepository(db.CreateDefaultConnection()),
		kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{os.Getenv("KAFKA_HOST")},
			Topic: string(constants.Users),
			Balancer: &kafka.LeastBytes{},
		}))
}

func CreateUserService(userRepository *repository.UserRepository, kafkaWriter *kafka.Writer) *UserService {
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
		kafkaWriter:         kafkaWriter,
	}
}

func (s *UserService) GetUserFromUsername(username string) (*model.PublicUser, error) {
	userEntity, err := s.userRepository.GetUserFromUsername(username)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserEntityToPublicUser(userEntity), nil
}

func (s *UserService) GetUserFromUuid(userUuid uuid.UUID) (*model.PublicUser, error) {
	userEntity, err := s.userRepository.GetUserFromUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserEntityToPublicUser(userEntity), nil
}

func (s *UserService) CreateUser(newUser *model.NewUser) (*model.User, error) {
	response, err := s.cognito.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		Username:  aws.String(newUser.Email),
		TemporaryPassword: aws.String(newUser.Password),
		UserPoolId: aws.String(s.cognitoUserPool),
	})

	if err != nil {
		return nil, err
	}

	user := mapper.MapNewUserModelToEntity(newUser, uuid.MustParse(*response.User.Attributes[0].Value))
	s.userRepository.Create(user)
	userModel := mapper.MapUserEntityToModel(user)
	userData, _ := json.Marshal(userModel)
	_ = s.kafkaWriter.WriteMessages(
		context.Background(),
		kafka.Message{Value: userData})
	return userModel, nil
}

func (s *UserService) UpdateUser(userModel *model.User) error {
	userEntity, err := s.userRepository.GetUserFromUsername(userModel.Username)
	if err != nil {
		return err
	}
	userEntity.UpdateUserProfileFromModel(userModel)
	data, _ := json.Marshal(userModel)
	s.userRepository.Save(userEntity)
	_ = s.kafkaWriter.WriteMessages(
		context.Background(),
		kafka.Message{Value: data})
	return nil
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
		return createAuthFailedSessionResponse("auth failed")
	}

	user, err := s.userRepository.GetUserFromEmail(newSession.Email)
	if err != nil {
		return createAuthFailedSessionResponse("user not found")
	}

	if response.AuthenticationResult != nil {
		s.updateUserTokens(user, response.AuthenticationResult)
		return createSessionResponse(user, response)
	}

	s.updateUserWithCreateSessionResult(user, response)
	log.Print("created session from AWS: ", response.String())
	return createChallengeSessionResponse(user, response)
}

func (s *UserService) ProvideChallengeResponse(passwordReset *model.PasswordReset) *AuthResponse {
	user, err := s.userRepository.GetUserFromEmail(passwordReset.Email)

	if err != nil {
		return createAuthFailedSessionResponse("user not found")
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
		return createAuthFailedSessionResponse("auth failed")
	}

	if response.AuthenticationResult != nil {
		s.updateUserTokens(user, response.AuthenticationResult)
	}

	return createChallengeResponse(response)
}

func (s *UserService) GetSession(sessionToken *model.SessionToken) (*model.Session, error) {
	keySet, jwkErr := jwk.Fetch(fmt.Sprintf(JwkTokenUrl, s.awsRegion, s.cognitoUserPool))
	if jwkErr != nil {
		log.Print("error fetching jwk: ", jwkErr)
		return nil, errors.New("jwk fetch error")
	}

	token, parseErr := jwt.Parse(sessionToken.Token, func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		keys := keySet.LookupKeyID(kid)
		if len(keys) > 0 {
			return keys[0].Materialize()
		}
		log.Print("error finding user session")
		return nil, errors.New("no session found")
	})
	if parseErr != nil {
		log.Print("jwt parse error")
		return nil, parseErr
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := claims.Valid(); err != nil || claims.VerifyAudience(s.cognitoClientID, false) == false {
		log.Print("token verification failed with: ", err)
		return nil, errors.New("verification failed")
	}

	response, err := s.cognito.GetUser(&cognitoidentityprovider.GetUserInput{ AccessToken: aws.String(sessionToken.Token) })
	if err != nil {
		log.Print("error retrieving user: ", err)
		return nil, err
	}
	user := s.userRepository.GetUserFromSessionToken(sessionToken.Token)
	if user == nil || user.CognitoId.String() != *response.Username {
		log.Print("user does not match jwt: ", response.String(), " and user: ", user)
		return nil, errors.New("user does not match jwt")
	}
	return model.CreateSession(mapper.MapUserEntityToPublicUser(user), sessionToken.Token), nil
}

func (s *UserService) RefreshSession(sessionRefresh *model.SessionRefresh) *AuthResponse {
	user := s.userRepository.GetUserFromSessionToken(sessionRefresh.Token)

	if user == nil {
		return createAuthFailedSessionResponse("auth failed")
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
		return createAuthFailedSessionResponse("auth failed")
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
