package logger

import "testing"

func TestNew(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "unknown"}
	for _, level := range levels {
		l := New(level)
		if l == nil {
			t.Fatalf("expected non-nil logger for level %q", level)
		}
		if l.logger == nil {
			t.Fatalf("expected non-nil log.Logger for level %q", level)
		}
	}
}

func TestLoggerDoesNotPanic(t *testing.T) {
	l := New("debug")
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")
	l.Info("structured", "key", "value", "num", 42)
}
