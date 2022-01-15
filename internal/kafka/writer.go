package kafka

import (
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"os"
	"time"
)

func CreateWriter() *kafka.Writer {
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
		Balancer: &kafka.LeastBytes{},
		Dialer: dialer,
	})
}
