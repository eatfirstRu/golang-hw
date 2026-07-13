package queue

import "context"

type Producer interface {
	SendNotification(ctx context.Context, data []byte) error
	Close() error
}

type Consumer interface {
	ReadMessages(ctx context.Context, handler func(ctx context.Context, data []byte) error) error
	Close() error
}
