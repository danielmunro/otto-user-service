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
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DBNAME"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"))
}

func CreateConnection(host string, port string, dbname string, user string, password string) *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
			host,
			port,
			dbname,
			user,
			password))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
