package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// CreateNewUserV1 - Create a new user
func CreateNewUserV1(w http.ResponseWriter, r *http.Request) {
	newUserModel := model.DecodeRequestToNewUser(r)
	user, err := service.CreateDefaultUserService().CreateUser(newUserModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)
}

// GetUserByUsernameV1 - Get a user by username
func GetUserByUsernameV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	user, err := service.CreateDefaultUserService().GetUserFromUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	_, _ = w.Write(data)

}

// UpdateUserV1 - Update a user
func UpdateUserV1(w http.ResponseWriter, r *http.Request) {
	userModel := model.DecodeRequestToUser(r)
	err := service.CreateDefaultUserService().UpdateUser(userModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(userModel)
	_, _ = w.Write(data)
}
