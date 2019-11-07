package entity

import (
	"github.com/danielmunro/otto-user-service/internal/enum"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Uuid uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CognitoId uuid.UUID
	SRP string
	LastSessionToken string
	LastAccessToken string
	LastIdToken string
	LastRefreshToken string
	DeviceGroupKey string
	DeviceKey string
	Name string
	CurrentEmail string `gorm:"unique;not null"`
	CurrentPassword string `gorm:"not null"`
	CurrentStatus string
	Birthday string
	Verified bool `gorm:"not null"`
	Emails []*Email
	Passwords []*Password
}

func CreateUser(newUser *model.NewUser, cognitoId string) *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	return &User{
		CurrentEmail: newUser.Email,
		CurrentPassword: string(hashedPassword),
		CognitoId: uuid.MustParse(cognitoId),
		CurrentStatus: string(enum.EmailStatusUnverified),
		Emails: []*Email{
			CreateEmail(newUser.Email),
		},
	}
}
