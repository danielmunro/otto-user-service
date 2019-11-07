package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const challengeNewPasswordString = "ChallengeNewPassword"

func getAuthResponseFromChallenge(response string) AuthResponseType {
	if response == AuthResponseChallenge {
		return ChallengeNewPassword
	}
	return Unknown
}

func createSessionResponse(response *cognitoidentityprovider.AdminInitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		Token: response.AuthenticationResult.AccessToken,
	}
}

func createChallengeSessionResponse(response *cognitoidentityprovider.AdminInitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		getAuthResponseFromChallenge(*response.ChallengeName),
		response.Session,
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
