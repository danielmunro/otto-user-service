package controller

import (
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/service"
	"net/http"
)

// CreateNewUserV1 - Create a new user
func CreateNewUserV1(w http.ResponseWriter, r *http.Request) {
	newUserModel := model.DecodeRequestToNewUser(r)
	_, err := service.CreateDefaultUserService().CreateUser(newUserModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetUserV1 - Get a user
func GetUserV1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

