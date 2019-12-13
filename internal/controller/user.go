package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/service"
	"github.com/danielmunro/otto-user-service/internal/util"
	"net/http"
)

// CreateNewUserV1 - Create a new user
func CreateNewUserV1(w http.ResponseWriter, r *http.Request) {
	newUserModel := model.DecodeRequestToNewUser(r)
	user, err := service.CreateDefaultUserService().CreateUser(newUserModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)
}

// GetUserV1 - Get a user
func GetUserV1(w http.ResponseWriter, r *http.Request) {
	userUuid := util.GetUuidFromPathSecondPosition(r.URL.Path)
	user, err := service.CreateDefaultUserService().GetUser(userUuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)
}

