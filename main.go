/*
 * Otto user service
 */

package main

import (
	"github.com/danielmunro/otto-user-service/internal"
	"github.com/danielmunro/otto-user-service/internal/middleware"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
)

func main() {
	router := internal.NewRouter()
	log.Fatal(http.ListenAndServe("localhost:8080", middleware.ContentTypeMiddleware(router)))
}
