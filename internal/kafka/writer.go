package kafka

import (
	"github.com/danielmunro/otto-user-service/internal/constants"
	"github.com/segmentio/kafka-go"
)

func CreateWriter(kafkaHost string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaHost},
		Topic: string(constants.Users),
		Balancer: &kafka.LeastBytes{},
	})
}
