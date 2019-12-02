package mapper

import (
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/model"
)

func MapUserEntityToModel(user *entity.User) *model.User {
	return &model.User{
		Uuid: user.Uuid.String(),
		CurrentEmail: user.CurrentEmail,
	}
}
