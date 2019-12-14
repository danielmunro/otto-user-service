package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/mapper"
	"log"
)

const challengeNewPasswordString = "ChallengeNewPassword"

func getAuthResponseFromChallenge(response string) AuthResponseType {
	if response == AuthResponseChallenge {
		return ChallengeNewPassword
	}
	return Unknown
}

func createSessionResponse(user *entity.User, response *cognitoidentityprovider.AdminInitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		Token: response.AuthenticationResult.AccessToken,
		User: mapper.MapUserEntityToPublicUser(user),
	}
}

func createChallengeSessionResponse(user *entity.User, response *cognitoidentityprovider.AdminInitiateAuthOutput) *AuthResponse {
	log.Print(user, mapper.MapUserEntityToPublicUser(user))
	return &AuthResponse{
		getAuthResponseFromChallenge(*response.ChallengeName),
		response.Session,
		mapper.MapUserEntityToPublicUser(user),
	}
}

func createAuthFailedSessionResponse() *AuthResponse {
	return &AuthResponse{
		AuthResponse: SessionFailedAuthentication,
	}
}

func getChallengeString(authResponse AuthResponseType) string {
	if authResponse == ChallengeNewPassword {
		return challengeNewPasswordString
	}
	return ""
}

func (c *AuthResponse) ToJson() []byte {
	data, _ := json.Marshal(map[string]string{
		"authResponse": getChallengeString(c.AuthResponse),
		"token": *c.Token,
	})
	return data
}
