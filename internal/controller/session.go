package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/service"
	"net/http"
)

// CreateSessionV1 - Create a new session
func CreateSessionV1(w http.ResponseWriter, r *http.Request) {
	newSessionModel := model.DecodeRequestToNewSession(r)
	result := service.CreateDefaultUserService().CreateSession(newSessionModel)
	if result.Token == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

// RespondToChallengeV1 - Respond to an authentication challenge with a password reset
func RespondToChallengeV1(w http.ResponseWriter, r *http.Request) {
	passwordResetModel := model.DecodeRequestToPasswordReset(r)
	result := service.CreateDefaultUserService().ProvideChallengeResponse(passwordResetModel)
	_, _ = w.Write(result.ToJson())
}

// GetSessionV1 - validate a session token
func GetSessionV1(w http.ResponseWriter, r *http.Request) {
	sessionToken := model.DecodeRequestToSessionToken(r)
	session, err := service.CreateDefaultUserService().GetSession(sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data, _ := json.Marshal(session)
	_, _ = w.Write(data)
}

// RefreshSessionV1 - refresh a session token
func RefreshSessionV1(w http.ResponseWriter, r *http.Request) {
	sessionToken := model.DecodeRequestToSessionRefresh(r)
	response := service.CreateDefaultUserService().RefreshSession(sessionToken)
	_, _ = w.Write(response.ToJson())
}
