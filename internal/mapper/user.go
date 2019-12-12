package mapper

import (
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/google/uuid"
)

func MapUserEntityToModel(user *entity.User) *model.User {
	return &model.User{
		Uuid:            user.Uuid.String(),
		Name:            user.Name,
		Username:        user.Username,
		CurrentEmail:    user.CurrentEmail,
		CurrentPassword: user.CurrentPassword,
		ProfilePic:      user.ProfilePic,
		AddressCity:     user.AddressCity,
		AddressStreet:   user.AddressStreet,
		AddressZip:      user.AddressZip,
		BioMessage:      user.BioMessage,
		Birthday:        user.Birthday,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func MapNewUserModelToEntity(user *model.NewUser, cognitoId uuid.UUID) *entity.User {
	return &entity.User{
		Name:      user.Name,
		Username:  user.Username,
		CurrentEmail: user.Email,
		CurrentPassword: user.Password,
		CognitoId: cognitoId,
	}
}
