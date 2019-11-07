package main

import (
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/entity"
)

func main() {
	conn := db.GetConnection()
	conn.DropTableIfExists(&entity.User{})
	conn.DropTableIfExists(&entity.Password{})
	conn.DropTableIfExists(&entity.Email{})
	conn.CreateTable(&entity.User{})
	conn.CreateTable(&entity.Password{})
	conn.CreateTable(&entity.Email{})
}