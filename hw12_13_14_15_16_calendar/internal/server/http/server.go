package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEventsDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type Server struct {
	httpServer *http.Server
	logger     Logger
	app        Application
}

func NewServer(logger Logger, app Application, host string, port int) *Server {
	s := &Server{
		logger: logger,
		app:    app,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/events", s.eventsHandler)
	mux.HandleFunc("/events/", s.eventByIDHandler)
	mux.HandleFunc("/events/day", s.listEventsDayHandler)
	mux.HandleFunc("/events/week", s.listEventsWeekHandler)
	mux.HandleFunc("/events/month", s.listEventsMonthHandler)

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(host, fmt.Sprintf("%d", port)),
		Handler: loggingMiddleware(logger, mux),
	}

	return s
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("starting HTTP server", "addr", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}
