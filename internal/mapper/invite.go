package mapper

import (
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/danielmunro/otto-user-service/internal/model"
)

func MapInviteEntityToModel(invite *entity.Invite) *model.Invite {
	return &model.Invite{
		Code: invite.Code,
	}
}
