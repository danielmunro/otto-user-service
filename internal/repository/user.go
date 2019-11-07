package repository

import (
	"github.com/danielmunro/otto-user-service/internal/entity"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	conn *gorm.DB
}

func CreateUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{ conn }
}

func (r *UserRepository) GetUserFromEmail(email string) *entity.User {
	user := &entity.User{}
	r.conn.Where("current_email = ?", email).Find(&user)
	return user
}

func (r *UserRepository) GetUserFromSessionToken(token string) *entity.User {
	user := &entity.User{}
	r.conn.Where("last_access_token = ?", token).Find(&user)
	return user
}

func (r *UserRepository) Create(user *entity.User) {
	r.conn.Create(user)
}

func (r *UserRepository) Save(user *entity.User) {
	r.conn.Save(user)
}
