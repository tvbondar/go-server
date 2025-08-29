package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/tvbondar/go-server/internal/usecases"
)

func StartKafkaConsumer(usecase *usecases.ProcessOrderUseCase) {
	config := kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "orders",
		GroupID:  "my-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}

	var reader *kafka.Reader
	var err error

	for i := 0; i < 10; i++ {
		reader = kafka.NewReader(config)
		var conn *kafka.Conn
		conn, err = kafka.Dial("tcp", "kafka:9092")
		if err == nil {
			conn.Close()
			break
		}
		fmt.Printf("Failed to connect to Kafka: %v\n", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Printf("Failed to connect to Kafka after retries: %v\n", err)
		return
	}
	defer reader.Close()
	fmt.Println("Successfully connected to Kafka topic 'orders'")

	ctx := context.Background()
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			fmt.Printf("Error fetching message: %v\n", err)
			continue
		}
		fmt.Printf("Received message with key: %s\n", string(msg.Key))
		if err := usecase.Execute(msg.Value); err != nil {
			fmt.Printf("Error processing message: %v\n", err)
			continue
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			fmt.Printf("Error committing message: %v\n", err)
		} else {
			fmt.Printf("Message processed and committed\n")
		}
	}
}
