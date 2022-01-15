package kafka

import (
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"os"
	"time"
)

func CreateWriter() *kafka.Writer {
	log.Print("server", os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
	log.Print("user", os.Getenv("KAFKA_SASL_USERNAME"))
	log.Print("pw", os.Getenv("KAFKA_SASL_PASSWORD"))
	mechanism, err := scram.Mechanism(
		scram.SHA512,
		os.Getenv("KAFKA_SASL_USERNAME"),
		os.Getenv("KAFKA_SASL_PASSWORD"))
	if err != nil {
		log.Panic(err)
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{os.Getenv("KAFKA_BOOTSTRAP_SERVERS")},
		Topic: string(constants.Users),
		Dialer: dialer,
	})
}
