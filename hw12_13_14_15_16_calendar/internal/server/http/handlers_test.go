package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type mockLogger struct{}

func (m *mockLogger) Info(_ string, _ ...any)  {}
func (m *mockLogger) Error(_ string, _ ...any) {}
func (m *mockLogger) Warn(_ string, _ ...any)  {}
func (m *mockLogger) Debug(_ string, _ ...any) {}

type mockApp struct {
	events map[string]storage.Event
}

func newMockApp() *mockApp {
	return &mockApp{events: make(map[string]storage.Event)}
}

func (m *mockApp) CreateEvent(_ context.Context, event storage.Event) error {
	for _, e := range m.events {
		if e.UserID == event.UserID && e.DateTime.Equal(event.DateTime) {
			return storage.ErrDateBusy
		}
	}
	if event.ID == "" {
		event.ID = "generated-id"
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockApp) UpdateEvent(_ context.Context, id string, event storage.Event) error {
	if _, ok := m.events[id]; !ok {
		return storage.ErrEventNotFound
	}
	event.ID = id
	m.events[id] = event
	return nil
}

func (m *mockApp) DeleteEvent(_ context.Context, id string) error {
	if _, ok := m.events[id]; !ok {
		return storage.ErrEventNotFound
	}
	delete(m.events, id)
	return nil
}

func (m *mockApp) ListEventsDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	var result []storage.Event
	for _, e := range m.events {
		if e.DateTime.Year() == date.Year() && e.DateTime.YearDay() == date.YearDay() {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockApp) ListEventsWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	end := date.AddDate(0, 0, 7)
	var result []storage.Event
	for _, e := range m.events {
		if !e.DateTime.Before(date) && e.DateTime.Before(end) {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockApp) ListEventsMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	end := date.AddDate(0, 1, 0)
	var result []storage.Event
	for _, e := range m.events {
		if !e.DateTime.Before(date) && e.DateTime.Before(end) {
			result = append(result, e)
		}
	}
	return result, nil
}

func newTestServer() (*Server, *mockApp) {
	app := newMockApp()
	srv := NewServer(&mockLogger{}, app, "localhost", 0)
	return srv, app
}

func TestCreateEvent(t *testing.T) {
	srv, _ := newTestServer()

	body := `{"title":"Meeting","date_time":"2024-06-15T10:00:00Z","duration":"1h","user_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	srv.eventsHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp eventResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Title != "Meeting" {
		t.Fatalf("expected title 'Meeting', got %q", resp.Title)
	}
}

func TestCreateEventBadRequest(t *testing.T) {
	srv, _ := newTestServer()

	body := `{"title":"","date_time":"2024-06-15T10:00:00Z","user_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	srv.eventsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateEventConflict(t *testing.T) {
	srv, app := newTestServer()

	app.events["existing"] = storage.Event{
		ID:       "existing",
		Title:    "Existing",
		DateTime: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}

	body := `{"title":"New","date_time":"2024-06-15T10:00:00Z","duration":"1h","user_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	srv.eventsHandler(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateEvent(t *testing.T) {
	srv, app := newTestServer()

	app.events["evt-1"] = storage.Event{
		ID:       "evt-1",
		Title:    "Old Title",
		DateTime: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}

	body := `{"title":"New Title","date_time":"2024-06-16T10:00:00Z","duration":"2h","user_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPut, "/events/evt-1", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	srv.eventByIDHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp eventResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Title != "New Title" {
		t.Fatalf("expected title 'New Title', got %q", resp.Title)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	srv, _ := newTestServer()

	body := `{"title":"X","date_time":"2024-06-15T10:00:00Z","duration":"1h","user_id":"user-1"}`
	req := httptest.NewRequest(http.MethodPut, "/events/nonexistent", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	srv.eventByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteEvent(t *testing.T) {
	srv, app := newTestServer()

	app.events["evt-1"] = storage.Event{ID: "evt-1", Title: "To Delete"}

	req := httptest.NewRequest(http.MethodDelete, "/events/evt-1", nil)
	w := httptest.NewRecorder()

	srv.eventByIDHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	if _, exists := app.events["evt-1"]; exists {
		t.Fatal("event should have been deleted")
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	srv, _ := newTestServer()

	req := httptest.NewRequest(http.MethodDelete, "/events/nonexistent", nil)
	w := httptest.NewRecorder()

	srv.eventByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestListEventsDay(t *testing.T) {
	srv, app := newTestServer()

	app.events["evt-1"] = storage.Event{
		ID:       "evt-1",
		Title:    "Day Event",
		DateTime: time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}
	app.events["evt-2"] = storage.Event{
		ID:       "evt-2",
		Title:    "Other Day",
		DateTime: time.Date(2024, 6, 16, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}

	req := httptest.NewRequest(http.MethodGet, "/events/day?date=2024-06-15", nil)
	w := httptest.NewRecorder()

	srv.listEventsDayHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp eventsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.Events))
	}
	if resp.Events[0].Title != "Day Event" {
		t.Fatalf("expected title 'Day Event', got %q", resp.Events[0].Title)
	}
}

func TestListEventsWeek(t *testing.T) {
	srv, app := newTestServer()

	app.events["evt-1"] = storage.Event{
		ID:       "evt-1",
		Title:    "This Week",
		DateTime: time.Date(2024, 6, 17, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}
	app.events["evt-2"] = storage.Event{
		ID:       "evt-2",
		Title:    "Next Week",
		DateTime: time.Date(2024, 6, 25, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}

	req := httptest.NewRequest(http.MethodGet, "/events/week?date=2024-06-15", nil)
	w := httptest.NewRecorder()

	srv.listEventsWeekHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp eventsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.Events))
	}
}

func TestListEventsMonth(t *testing.T) {
	srv, app := newTestServer()

	app.events["evt-1"] = storage.Event{
		ID:       "evt-1",
		Title:    "June Event",
		DateTime: time.Date(2024, 6, 20, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}
	app.events["evt-2"] = storage.Event{
		ID:       "evt-2",
		Title:    "July Event",
		DateTime: time.Date(2024, 7, 5, 10, 0, 0, 0, time.UTC),
		UserID:   "user-1",
	}

	req := httptest.NewRequest(http.MethodGet, "/events/month?date=2024-06-01", nil)
	w := httptest.NewRecorder()

	srv.listEventsMonthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp eventsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.Events))
	}
}

func TestListEventsMissingDate(t *testing.T) {
	srv, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/events/day", nil)
	w := httptest.NewRecorder()

	srv.listEventsDayHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestListEventsInvalidDate(t *testing.T) {
	srv, _ := newTestServer()

	req := httptest.NewRequest(http.MethodGet, "/events/day?date=invalid", nil)
	w := httptest.NewRecorder()

	srv.listEventsDayHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
