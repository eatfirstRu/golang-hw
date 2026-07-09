package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
)

type Logger struct {
	level  int
	logger *log.Logger
}

func New(level string) *Logger {
	var lvl int
	switch strings.ToLower(level) {
	case "debug":
		lvl = levelDebug
	case "info":
		lvl = levelInfo
	case "warn", "warning":
		lvl = levelWarn
	case "error":
		lvl = levelError
	default:
		lvl = levelInfo
	}

	return &Logger{
		level:  lvl,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= levelDebug {
		l.log("DEBUG", msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= levelInfo {
		l.log("INFO", msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= levelWarn {
		l.log("WARN", msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= levelError {
		l.log("ERROR", msg, args...)
	}
}

func (l *Logger) log(level, msg string, args ...interface{}) {
	if len(args) > 0 {
		kvs := make([]string, 0, len(args)/2)
		for i := 0; i+1 < len(args); i += 2 {
			kvs = append(kvs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		}
		l.logger.Printf("[%s] %s %s", level, msg, strings.Join(kvs, " "))
	} else {
		l.logger.Printf("[%s] %s", level, msg)
	}
}
