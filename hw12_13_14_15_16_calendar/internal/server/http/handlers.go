package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

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

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) eventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.createEventHandler(w, r)
}

func (s *Server) eventByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/events/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "event id is required")
		return
	}

	switch r.Method {
	case http.MethodPut:
		s.updateEventHandler(w, r, id)
	case http.MethodDelete:
		s.deleteEventHandler(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	event, err := parseEventRequest(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.app.CreateEvent(r.Context(), event); err != nil {
		if errors.Is(err, storage.ErrDateBusy) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toEventResponse(event))
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request, id string) {
	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	event, err := parseEventRequest(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.app.UpdateEvent(r.Context(), id, event); err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, storage.ErrDateBusy) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	event.ID = id
	writeJSON(w, http.StatusOK, toEventResponse(event))
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.app.DeleteEvent(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listEventsDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.listEventsHandler(w, r, func(ctx context.Context, date time.Time) ([]storage.Event, error) {
		return s.app.ListEventsDay(ctx, date)
	})
}

func (s *Server) listEventsWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.listEventsHandler(w, r, func(ctx context.Context, date time.Time) ([]storage.Event, error) {
		return s.app.ListEventsWeek(ctx, date)
	})
}

func (s *Server) listEventsMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.listEventsHandler(w, r, func(ctx context.Context, date time.Time) ([]storage.Event, error) {
		return s.app.ListEventsMonth(ctx, date)
	})
}

func (s *Server) listEventsHandler(
	w http.ResponseWriter,
	r *http.Request,
	fn func(context.Context, time.Time) ([]storage.Event, error),
) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "date query parameter is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	events, err := fn(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := eventsResponse{Events: make([]eventResponse, 0, len(events))}
	for _, e := range events {
		resp.Events = append(resp.Events, toEventResponse(e))
	}

	writeJSON(w, http.StatusOK, resp)
}

func parseEventRequest(req eventRequest) (storage.Event, error) {
	if req.Title == "" {
		return storage.Event{}, errors.New("title is required")
	}
	if req.DateTime == "" {
		return storage.Event{}, errors.New("date_time is required")
	}
	if req.UserID == "" {
		return storage.Event{}, errors.New("user_id is required")
	}

	dateTime, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		return storage.Event{}, errors.New("invalid date_time format, use RFC3339")
	}

	var duration time.Duration
	if req.Duration != "" {
		duration, err = time.ParseDuration(req.Duration)
		if err != nil {
			return storage.Event{}, errors.New("invalid duration format")
		}
	}

	var notifyBefore time.Duration
	if req.NotifyBefore != "" {
		notifyBefore, err = time.ParseDuration(req.NotifyBefore)
		if err != nil {
			return storage.Event{}, errors.New("invalid notify_before format")
		}
	}

	return storage.Event{
		Title:        req.Title,
		DateTime:     dateTime,
		Duration:     duration,
		Description:  req.Description,
		UserID:       req.UserID,
		NotifyBefore: notifyBefore,
	}, nil
}

func toEventResponse(e storage.Event) eventResponse {
	resp := eventResponse{
		ID:       e.ID,
		Title:    e.Title,
		DateTime: e.DateTime.Format(time.RFC3339),
		UserID:   e.UserID,
	}
	if e.Duration != 0 {
		resp.Duration = e.Duration.String()
	}
	if e.Description != "" {
		resp.Description = e.Description
	}
	if e.NotifyBefore != 0 {
		resp.NotifyBefore = e.NotifyBefore.String()
	}
	return resp
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, errorResponse{Error: msg})
}
