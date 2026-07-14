package memorystorage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

var ctx = context.Background()

func newEvent(id, userID string, dt time.Time) storage.Event {
	return storage.Event{
		ID:       id,
		Title:    "Event " + id,
		DateTime: dt,
		Duration: time.Hour,
		UserID:   userID,
	}
}

func TestCreateEvent(t *testing.T) {
	s := New()
	event := newEvent("1", "user1", time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(s.events))
	}
}

func TestCreateEventDateBusy(t *testing.T) {
	s := New()
	dt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "user1", dt))
	err := s.CreateEvent(ctx, newEvent("2", "user1", dt))

	if !errors.Is(err, storage.ErrDateBusy) {
		t.Fatalf("expected ErrDateBusy, got: %v", err)
	}
}

func TestCreateEventSameTimeDifferentUser(t *testing.T) {
	s := New()
	dt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "user1", dt))
	err := s.CreateEvent(ctx, newEvent("2", "user2", dt))

	if err != nil {
		t.Fatalf("different users should be able to have events at the same time: %v", err)
	}
}

func TestUpdateEvent(t *testing.T) {
	s := New()
	dt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "user1", dt))

	updated := storage.Event{
		Title:    "Updated",
		DateTime: dt.Add(time.Hour),
		UserID:   "user1",
	}

	err := s.UpdateEvent(ctx, "1", updated)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.events["1"].Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", s.events["1"].Title)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	s := New()

	err := s.UpdateEvent(ctx, "nonexistent", storage.Event{})
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got: %v", err)
	}
}

func TestUpdateEventDateBusy(t *testing.T) {
	s := New()
	dt1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	dt2 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "user1", dt1))
	_ = s.CreateEvent(ctx, newEvent("2", "user1", dt2))

	updated := storage.Event{DateTime: dt1, UserID: "user1"}
	err := s.UpdateEvent(ctx, "2", updated)
	if !errors.Is(err, storage.ErrDateBusy) {
		t.Fatalf("expected ErrDateBusy, got: %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := New()
	dt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "user1", dt))

	err := s.DeleteEvent(ctx, "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(s.events))
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	s := New()

	err := s.DeleteEvent(ctx, "nonexistent")
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got: %v", err)
	}
}

func TestListEventsDay(t *testing.T) {
	s := New()
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	_ = s.CreateEvent(ctx, newEvent("1", "u1", base.Add(2*time.Hour)))
	_ = s.CreateEvent(ctx, newEvent("2", "u2", base.Add(5*time.Hour)))
	_ = s.CreateEvent(ctx, newEvent("3", "u3", base.AddDate(0, 0, 1)))

	events, err := s.ListEventsDay(ctx, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events for day, got %d", len(events))
	}
}

func TestListEventsWeek(t *testing.T) {
	s := New()
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 10; i++ {
		_ = s.CreateEvent(ctx, newEvent(
			fmt.Sprintf("e%d", i), "u1",
			base.AddDate(0, 0, i),
		))
	}

	events, err := s.ListEventsWeek(ctx, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 7 {
		t.Fatalf("expected 7 events for week, got %d", len(events))
	}
}

func TestListEventsMonth(t *testing.T) {
	s := New()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 35; i++ {
		_ = s.CreateEvent(ctx, newEvent(
			fmt.Sprintf("e%d", i), "u1",
			base.AddDate(0, 0, i),
		))
	}

	events, err := s.ListEventsMonth(ctx, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 31 {
		t.Fatalf("expected 31 events for month, got %d", len(events))
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			event := newEvent(
				fmt.Sprintf("event-%d", i),
				fmt.Sprintf("user-%d", i),
				time.Date(2024, 1, 1, i%24, i%60, 0, 0, time.UTC),
			)
			_ = s.CreateEvent(ctx, event)
		}(i)
	}
	wg.Wait()

	var wg2 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg2.Add(2)
		go func(i int) {
			defer wg2.Done()
			_, _ = s.ListEventsDay(ctx, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
		}(i)
		go func(i int) {
			defer wg2.Done()
			_ = s.DeleteEvent(ctx, fmt.Sprintf("event-%d", i))
		}(i)
	}
	wg2.Wait()
}
