package kafka

import (
	"context"
	"github.com/danielmunro/otto-user-service/internal/db"
	"github.com/danielmunro/otto-user-service/internal/model"
	"github.com/danielmunro/otto-user-service/internal/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
)

func InitializeAndRunLoop(kafkaHost string) {
	reader := GetReader(kafkaHost)
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	err := loopKafkaReader(userRepository, reader)
	if err != nil {
		log.Fatal(err)
	}
}

func loopKafkaReader(userRepository *repository.UserRepository, reader *kafka.Reader) error {
	for {
		log.Print("kafka ready to consume image messages")
		data, err := reader.ReadMessage(context.Background())
		if err != nil  {
			log.Print(err)
			return nil
		}
		log.Print("consuming image message ", string(data.Value))
		image, err := model.DecodeMessageToImage(data.Value)
		if err != nil {
			log.Print("error decoding message to user, skipping", string(data.Value))
			continue
		}
		userEntity, err := userRepository.GetUserFromUuid(uuid.MustParse(image.User.Uuid))
		if err != nil {
			log.Print("user not found when updating profile pic")
			continue
		}
		log.Print("update user with s3 key", userEntity.Uuid.String(), image.S3Key)
		userEntity.ProfilePic = image.S3Key
		userRepository.Update(userEntity)
	}
}
