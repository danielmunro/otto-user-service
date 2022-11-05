package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/service"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var numbers = []rune("0123456789")
var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// CreateInviteV1 -- create new invites for new users
func CreateInviteV1(w http.ResponseWriter, r *http.Request) {
	userService := service.CreateDefaultUserService()
	sessionToken := getSessionToken(r)
	sessionModel := &model.SessionToken{
		Token: sessionToken,
	}
	session, err := userService.GetSession(sessionModel)
	if err != nil || session.User.Role == model.USER {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	code := generateCode()
	attempt := 0
	for {
		_, err = userService.GetInvite(code)
		if err.Error() == "no invite found" {
			break
		}
		code = generateCode()
		attempt += 1
		if attempt > 5 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	invite, err := userService.CreateInviteFromCode(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(invite)
	_, _ = w.Write(data)
}

func generateCode() string {
	b := make([]rune, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	for i := range b {
		b[i] = numbers[rand.Intn(len(letters))]
	}
	return string(b)
}
