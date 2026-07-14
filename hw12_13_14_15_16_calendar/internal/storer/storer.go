package storer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/metrics"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/queue"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Storer struct {
	storage  storage.Storage
	consumer queue.Consumer
	logger   Logger
}

func New(logger Logger, store storage.Storage, consumer queue.Consumer) *Storer {
	return &Storer{
		storage:  store,
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Storer) Run(ctx context.Context) error {
	s.logger.Info("storer started")

	return s.consumer.ReadMessages(ctx, func(ctx context.Context, data []byte) error {
		var n storage.Notification
		if err := json.Unmarshal(data, &n); err != nil {
			s.logger.Error("unmarshal notification", "error", err)
			return nil
		}

		if err := s.storage.SaveNotification(ctx, n); err != nil {
			return fmt.Errorf("save notification: %w", err)
		}

		metrics.NotificationsSaved.Inc()
		s.logger.Info("notification saved", "event_id", n.EventID, "title", n.Title)
		return nil
	})
}
