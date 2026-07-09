package app

import (
	"context"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

func New(logger Logger, store storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: store,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) ListEventsDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsDay(ctx, date)
}

func (a *App) ListEventsWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsWeek(ctx, date)
}

func (a *App) ListEventsMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsMonth(ctx, date)
}
