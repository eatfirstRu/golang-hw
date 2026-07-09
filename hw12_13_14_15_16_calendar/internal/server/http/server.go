package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type Application interface{}

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
