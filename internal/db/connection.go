package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

func CreateDefaultConnection() *gorm.DB {
	return CreateConnection(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
}

func CreateConnection(host string, port string, dbname string, user string, password string) *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s dbname=user_service user=%s password=%s sslmode=disable",
			host,
			port,
			user,
			password))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
