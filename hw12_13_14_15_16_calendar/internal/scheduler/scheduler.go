package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/queue"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Scheduler struct {
	storage  storage.Storage
	producer queue.Producer
	logger   Logger
	interval time.Duration
}

func New(logger Logger, store storage.Storage, producer queue.Producer, interval time.Duration) *Scheduler {
	return &Scheduler{
		storage:  store,
		producer: producer,
		logger:   logger,
		interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("scheduler started", "interval", s.interval.String())

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return nil
		case <-ticker.C:
			if err := s.tick(ctx); err != nil {
				s.logger.Error("scheduler tick failed", "error", err)
			}
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) error {
	now := time.Now()

	events, err := s.storage.GetEventsToNotify(ctx, now)
	if err != nil {
		return fmt.Errorf("get events to notify: %w", err)
	}

	for _, e := range events {
		n := storage.Notification{
			EventID:  e.ID,
			Title:    e.Title,
			DateTime: e.DateTime,
			UserID:   e.UserID,
		}

		data, err := json.Marshal(n)
		if err != nil {
			s.logger.Error("marshal notification", "error", err)
			continue
		}

		if err := s.producer.SendNotification(ctx, data); err != nil {
			s.logger.Error("send notification", "error", err, "event_id", e.ID)
			continue
		}
		s.logger.Info("notification sent", "event_id", e.ID, "title", e.Title)
	}

	oneYearAgo := now.AddDate(-1, 0, 0)
	deleted, err := s.storage.DeleteOldEvents(ctx, oneYearAgo)
	if err != nil {
		return fmt.Errorf("delete old events: %w", err)
	}
	if deleted > 0 {
		s.logger.Info("old events deleted", "count", deleted)
	}

	return nil
}
