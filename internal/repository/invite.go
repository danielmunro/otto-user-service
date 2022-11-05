package repository

import (
	"errors"
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/jinzhu/gorm"
)

type InviteRepository struct {
	conn *gorm.DB
}

func CreateInviteRepository(conn *gorm.DB) *InviteRepository {
	return &InviteRepository{conn}
}

func (r *InviteRepository) FindOneByCode(code string) (*entity.Invite, error) {
	invite := &entity.Invite{}
	r.conn.Where("code = ?", code).Find(invite)
	if invite.ID == 0 {
		return nil, errors.New("no invite found")
	}
	return invite, nil
}

func (r *InviteRepository) Create(invite *entity.Invite) *gorm.DB {
	return r.conn.Create(invite)
}
