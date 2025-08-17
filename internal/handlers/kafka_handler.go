package handlers

import (
	"internal/usecases"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func KafkaConsumer(usecase *usecases.ProcessOrderUseCase) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"orders"}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			usecase.Execute(msg.Value)
		}
	}
}
