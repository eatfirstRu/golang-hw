package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date is already busy by another event")
)

type Event struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	DateTime     time.Time     `db:"date_time"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, id string, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEventsDay(ctx context.Context, date time.Time) ([]Event, error)
	ListEventsWeek(ctx context.Context, date time.Time) ([]Event, error)
	ListEventsMonth(ctx context.Context, date time.Time) ([]Event, error)
}
