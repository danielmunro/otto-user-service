package service

import "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

type AuthResponseType int

const (
	Unknown              		AuthResponseType = iota + 1
	ChallengeNewPassword 		AuthResponseType = iota
	SessionAuthenticated 	    AuthResponseType = iota
	SessionFailedAuthentication AuthResponseType = iota
)

type AuthResponse struct {
	AuthResponse AuthResponseType
	Token        *string
}

func createSuccessfulRefreshResponse(response *cognitoidentityprovider.AdminInitiateAuthOutput) *AuthResponse {
	return &AuthResponse{
		AuthResponse: SessionAuthenticated,
		Token:        response.AuthenticationResult.AccessToken,
	}
}
