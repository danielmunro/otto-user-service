package entity

import (
	"encoding/json"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Uuid             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CognitoId        uuid.UUID
	SRP              string
	LastSessionToken string
	LastAccessToken  string
	LastIdToken      string
	LastRefreshToken string
	DeviceGroupKey   string
	DeviceKey        string
	Name             string
	Username         string `gorm:"unique;not null"`
	ProfilePic       string
	BioMessage       string
	Role             string `gorm:"default:'user'"`
	IsBanned         bool   `gorm:"default:false"`
	AddressStreet    string
	AddressCity      string
	AddressZip       string
	CurrentEmail     string `gorm:"unique;not null"`
	CurrentPassword  string `gorm:"not null"`
	Birthday         string
	Verified         bool `gorm:"not null"`
	Emails           []*Email
	Passwords        []*Password
}

func CreateUser(newUser *model.NewUser, cognitoId string) *User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	return &User{
		CurrentEmail:    newUser.Email,
		CurrentPassword: string(hashedPassword),
		CognitoId:       uuid.MustParse(cognitoId),
		Emails: []*Email{
			CreateEmail(newUser.Email),
		},
	}
}

func (u *User) UpdateUserProfileFromModel(user *model.User) {
	u.Name = user.Name
	u.ProfilePic = user.ProfilePic
	u.BioMessage = user.BioMessage
	u.Birthday = user.Birthday
}

func (u *User) ToJson() []byte {
	data, _ := json.Marshal(u)
	return data
}
