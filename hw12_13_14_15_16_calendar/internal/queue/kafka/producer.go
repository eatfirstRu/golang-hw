package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafkago.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Producer{writer: w}
}

func (p *Producer) SendNotification(ctx context.Context, data []byte) error {
	err := p.writer.WriteMessages(ctx, kafkago.Message{
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
