package handlers

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/tvbondar/go-server/internal/usecases"
)

func StartKafkaConsumer(usecase *usecases.ProcessOrderUseCase) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		GroupID:  "my-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	ctx := context.Background()
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			fmt.Printf("Error fetching message: %v\n", err)
			continue
		}
		if err := usecase.Execute(msg.Value); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
			continue
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			fmt.Printf("Error committing message: %v\n", err)
		}
	}
}
