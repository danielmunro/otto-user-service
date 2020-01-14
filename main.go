/*
 * Otto user service
 */

package main

import (
	"github.com/danielmunro/otto-user-service/internal"
	"github.com/danielmunro/otto-user-service/internal/kafka"
	"github.com/danielmunro/otto-user-service/internal/middleware"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

func main() {
	go serveHttp()
	readKafka()
}

func readKafka() {
	kafkaHost := os.Getenv("KAFKA_HOST")
	log.Print("connecting to kafka", kafkaHost)
	kafka.InitializeAndRunLoop(kafkaHost)
	log.Print("exit kafka loop")
}

func serveHttp() {
	log.Print("Listening on 8080")
	router := internal.NewRouter()
	log.Fatal(http.ListenAndServe(":8080",
		middleware.CorsMiddleware(middleware.ContentTypeMiddleware(router))))
}
