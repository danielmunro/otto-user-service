package kafka

import (
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"os"
	"time"
)

func CreateWriter() *kafka.Writer {
	log.Print("server", os.Getenv("KAFKA_BOOTSTRAP_SERVER"))
	log.Print("user", os.Getenv("KAFKA_SASL_USERNAME"))
	log.Print("pw", os.Getenv("KAFKA_SASL_PASSWORD"))
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_SASL_USERNAME"),
		Password: os.Getenv("KAFKA_SASL_PASSWORD"),
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")},
		Topic: string(constants.Users),
		Dialer: dialer,
	})
}
