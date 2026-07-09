package sqlstorage

import (
	"context"
	"fmt"
	"time"

	// PostgreSQL driver for database/sql via pgx.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

type dbEvent struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	DateTime     time.Time `db:"date_time"`
	Duration     int64     `db:"duration"`
	Description  string    `db:"description"`
	UserID       string    `db:"user_id"`
	NotifyBefore int64     `db:"notify_before"`
}

func New(dsn string) (*Storage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	var count int
	err := s.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM events WHERE user_id=$1 AND date_time=$2",
		event.UserID, event.DateTime)
	if err != nil {
		return fmt.Errorf("check date busy: %w", err)
	}
	if count > 0 {
		return storage.ErrDateBusy
	}

	d := toDBEvent(event)
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO events (id, title, date_time, duration, description, user_id, notify_before)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		d.ID, d.Title, d.DateTime, d.Duration, d.Description, d.UserID, d.NotifyBefore)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	var exists bool
	err := s.db.GetContext(ctx, &exists,
		"SELECT EXISTS(SELECT 1 FROM events WHERE id=$1)", id)
	if err != nil {
		return fmt.Errorf("check event exists: %w", err)
	}
	if !exists {
		return storage.ErrEventNotFound
	}

	var count int
	err = s.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM events WHERE user_id=$1 AND date_time=$2 AND id<>$3",
		event.UserID, event.DateTime, id)
	if err != nil {
		return fmt.Errorf("check date busy: %w", err)
	}
	if count > 0 {
		return storage.ErrDateBusy
	}

	d := toDBEvent(event)
	_, err = s.db.ExecContext(ctx,
		`UPDATE events
		 SET title=$1, date_time=$2, duration=$3, description=$4, user_id=$5, notify_before=$6
		 WHERE id=$7`,
		d.Title, d.DateTime, d.Duration, d.Description, d.UserID, d.NotifyBefore, id)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return storage.ErrEventNotFound
	}
	return nil
}

func (s *Storage) ListEventsDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := truncateToDay(date)
	end := start.AddDate(0, 0, 1)
	return s.listEvents(ctx, start, end)
}

func (s *Storage) ListEventsWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := truncateToDay(date)
	end := start.AddDate(0, 0, 7)
	return s.listEvents(ctx, start, end)
}

func (s *Storage) ListEventsMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := truncateToDay(date)
	end := start.AddDate(0, 1, 0)
	return s.listEvents(ctx, start, end)
}

func (s *Storage) listEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var rows []dbEvent
	err := s.db.SelectContext(ctx, &rows,
		`SELECT id, title, date_time, duration, description, user_id, notify_before
		 FROM events
		 WHERE date_time >= $1 AND date_time < $2
		 ORDER BY date_time`,
		start, end)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	events := make([]storage.Event, 0, len(rows))
	for _, r := range rows {
		events = append(events, r.toEvent())
	}
	return events, nil
}

func toDBEvent(e storage.Event) dbEvent {
	return dbEvent{
		ID:           e.ID,
		Title:        e.Title,
		DateTime:     e.DateTime,
		Duration:     int64(e.Duration),
		Description:  e.Description,
		UserID:       e.UserID,
		NotifyBefore: int64(e.NotifyBefore),
	}
}

func (d dbEvent) toEvent() storage.Event {
	return storage.Event{
		ID:           d.ID,
		Title:        d.Title,
		DateTime:     d.DateTime,
		Duration:     time.Duration(d.Duration),
		Description:  d.Description,
		UserID:       d.UserID,
		NotifyBefore: time.Duration(d.NotifyBefore),
	}
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
