package kafka

import (
	"context"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) ReadMessages(ctx context.Context, handler func(ctx context.Context, data []byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("kafka read: %w", err)
		}

		if err := handler(ctx, msg.Value); err != nil {
			return fmt.Errorf("handle message: %w", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
