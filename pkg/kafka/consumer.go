package kafka

import (
    "context"
    "go-server/internal/service"
    "github.com/segmentio/kafka-go"
    "go.uber.org/zap"
)

type Consumer struct {
    reader  *kafka.Reader
    service *service.OrderService
    logger  *zap.Logger
}

func NewConsumer(brokers []string, topic string, groupID string, service *service.OrderService, logger *zap.Logger) *Consumer {
    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:  brokers,
            Topic:    topic,
            GroupID:  groupID,
            MinBytes: 10e3,
            MaxBytes: 10e6,
        }),
        service: service,
        logger:  logger,
    }
}

func (c *Consumer) Consume(ctx context.Context) error {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            c.logger.Error("Failed to fetch message", zap.Error(err))
            continue
        }
        if err := c.service.ProcessOrder(ctx, msg.Value); err != nil {
            c.logger.Error("Failed to process message", zap.Error(err))
            continue
        }
        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            c.logger.Error("Failed to commit message", zap.Error(err))
        }
    }
}
