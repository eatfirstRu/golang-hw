package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

var calendarURL string

func TestMain(m *testing.M) {
	calendarURL = os.Getenv("CALENDAR_URL")
	if calendarURL == "" {
		calendarURL = "http://localhost:8888"
	}

	for i := 0; i < 30; i++ {
		resp, err := http.Get(calendarURL + "/hello")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		time.Sleep(time.Second)
	}

	os.Exit(m.Run())
}

type eventRequest struct {
	Title        string `json:"title"`
	DateTime     string `json:"date_time"`
	Duration     string `json:"duration"`
	Description  string `json:"description"`
	UserID       string `json:"user_id"`
	NotifyBefore string `json:"notify_before"`
}

type eventResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	DateTime     string `json:"date_time"`
	Duration     string `json:"duration"`
	Description  string `json:"description"`
	UserID       string `json:"user_id"`
	NotifyBefore string `json:"notify_before"`
}

type eventsResponse struct {
	Events []eventResponse `json:"events"`
}

func createEvent(t *testing.T, req eventRequest) eventResponse {
	t.Helper()
	body, _ := json.Marshal(req)
	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create event request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, string(b))
	}

	var result eventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return result
}

func TestCreateEvent(t *testing.T) {
	resp := createEvent(t, eventRequest{
		Title:    "Integration Test Event",
		DateTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		Duration: "1h",
		UserID:   "integ-user-1",
	})

	if resp.Title != "Integration Test Event" {
		t.Fatalf("expected title 'Integration Test Event', got %q", resp.Title)
	}
}

func TestCreateEventValidation(t *testing.T) {
	body, _ := json.Marshal(eventRequest{
		Title:    "",
		DateTime: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		UserID:   "integ-user-1",
	})

	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestCreateEventDateBusy(t *testing.T) {
	dt := time.Now().Add(48 * time.Hour).Truncate(time.Second).Format(time.RFC3339)

	createEvent(t, eventRequest{
		Title:    "First Event",
		DateTime: dt,
		Duration: "1h",
		UserID:   "integ-user-busy",
	})

	body, _ := json.Marshal(eventRequest{
		Title:    "Second Event",
		DateTime: dt,
		Duration: "1h",
		UserID:   "integ-user-busy",
	})

	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}
}

func TestListEventsDay(t *testing.T) {
	tomorrow := time.Now().Add(72 * time.Hour).Truncate(24 * time.Hour)
	dt := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 14, 0, 0, 0, time.UTC)

	createEvent(t, eventRequest{
		Title:    "Day List Event",
		DateTime: dt.Format(time.RFC3339),
		Duration: "30m",
		UserID:   "integ-user-day",
	})

	url := fmt.Sprintf("%s/events/day?date=%s", calendarURL, dt.Format("2006-01-02"))
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result eventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	found := false
	for _, e := range result.Events {
		if e.Title == "Day List Event" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("event not found in day listing")
	}
}

func TestListEventsWeek(t *testing.T) {
	base := time.Now().Add(96 * time.Hour).Truncate(24 * time.Hour)
	dt := time.Date(base.Year(), base.Month(), base.Day(), 10, 0, 0, 0, time.UTC)

	createEvent(t, eventRequest{
		Title:    "Week List Event",
		DateTime: dt.Format(time.RFC3339),
		Duration: "1h",
		UserID:   "integ-user-week",
	})

	url := fmt.Sprintf("%s/events/week?date=%s", calendarURL, dt.Format("2006-01-02"))
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result eventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	found := false
	for _, e := range result.Events {
		if e.Title == "Week List Event" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("event not found in week listing")
	}
}

func TestListEventsMonth(t *testing.T) {
	base := time.Now().Add(120 * time.Hour).Truncate(24 * time.Hour)
	dt := time.Date(base.Year(), base.Month(), base.Day(), 16, 0, 0, 0, time.UTC)

	createEvent(t, eventRequest{
		Title:    "Month List Event",
		DateTime: dt.Format(time.RFC3339),
		Duration: "2h",
		UserID:   "integ-user-month",
	})

	firstOfMonth := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
	url := fmt.Sprintf("%s/events/month?date=%s", calendarURL, firstOfMonth.Format("2006-01-02"))
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result eventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	found := false
	for _, e := range result.Events {
		if e.Title == "Month List Event" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("event not found in month listing")
	}
}
