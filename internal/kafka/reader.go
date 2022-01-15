package kafka

import (
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"os"
	"time"
)

func GetReader() *kafka.Reader {
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_SASL_USERNAME"),
		Password: os.Getenv("KAFKA_SASL_PASSWORD"),
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")},
		Topic:     string(constants.Images),
		GroupID: "user_service",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		Dialer: dialer,
	})
	return r
}
