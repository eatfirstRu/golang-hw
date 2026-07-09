package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkDateBusy(event, ""); err != nil {
		return err
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotFound
	}
	if err := s.checkDateBusy(event, id); err != nil {
		return err
	}
	event.ID = id
	s.events[id] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) ListEventsDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := truncateToDay(date)
	end := start.AddDate(0, 0, 1)
	return s.filterByDateRange(start, end), nil
}

func (s *Storage) ListEventsWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := truncateToDay(date)
	end := start.AddDate(0, 0, 7)
	return s.filterByDateRange(start, end), nil
}

func (s *Storage) ListEventsMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	start := truncateToDay(date)
	end := start.AddDate(0, 1, 0)
	return s.filterByDateRange(start, end), nil
}

func (s *Storage) filterByDateRange(start, end time.Time) []storage.Event {
	var result []storage.Event
	for _, e := range s.events {
		if !e.DateTime.Before(start) && e.DateTime.Before(end) {
			result = append(result, e)
		}
	}
	return result
}

func (s *Storage) checkDateBusy(event storage.Event, excludeID string) error {
	for _, e := range s.events {
		if e.ID == excludeID {
			continue
		}
		if e.UserID == event.UserID && e.DateTime.Equal(event.DateTime) {
			return storage.ErrDateBusy
		}
	}
	return nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
