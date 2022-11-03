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
		Role:            model.Role(user.Role),
		IsBanned:        user.IsBanned,
		AddressCity:     user.AddressCity,
		AddressStreet:   user.AddressStreet,
		AddressZip:      user.AddressZip,
		BioMessage:      user.BioMessage,
		Birthday:        user.Birthday,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func MapUserEntityToPublicUser(user *entity.User) *model.PublicUser {
	return &model.PublicUser{
		Uuid:          user.Uuid.String(),
		Name:          user.Name,
		Username:      user.Username,
		ProfilePic:    user.ProfilePic,
		IsBanned:      user.IsBanned,
		Role:          model.Role(user.Role),
		AddressCity:   user.AddressCity,
		AddressStreet: user.AddressStreet,
		AddressZip:    user.AddressZip,
		BioMessage:    user.BioMessage,
		Birthday:      user.Birthday,
		CreatedAt:     user.CreatedAt,
	}
}

func MapNewUserModelToEntity(user *model.NewUser, cognitoId uuid.UUID) *entity.User {
	return &entity.User{
		Name:         user.Name,
		Username:     user.Username,
		CurrentEmail: user.Email,
		CognitoId:    cognitoId,
	}
}
