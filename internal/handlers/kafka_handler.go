package handlers

import (
	"fmt"

	"github.com/tvbondar/go-server/internal/usecases"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartKafkaConsumer(usecase *usecases.ProcessOrderUseCase) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{"orders"}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			usecase.Execute(msg.Value)
			c.CommitMessage(msg) // Подтверждение для Kafka, чтобы не терять сообщения
		} else {
			fmt.Printf("Error reading message: %v\n", err)
		}
	}
}
