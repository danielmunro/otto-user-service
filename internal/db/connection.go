package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func GetConnection() *gorm.DB {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=54321 user=otto dbname=user_service password=otto sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
