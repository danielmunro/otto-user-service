package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
)

func GetReader() *kafka.Consumer {
	log.Print(
		fmt.Printf(
			"debug kafka reader :: %s %s %s %s %s",
			os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
			os.Getenv("KAFKA_SECURITY_PROTOCOL"),
			os.Getenv("KAFKA_SASL_MECHANISMS"),
			os.Getenv("KAFKA_SASL_USERNAME"),
			os.Getenv("KAFKA_SASL_PASSWORD"),
		),
	)
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":     os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_SASL_PASSWORD"),
		"group.id":          "otto",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"images"}, nil)
	return c
}
