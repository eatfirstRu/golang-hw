package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/queue/kafka"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/scheduler"
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(cfg.Logger.Level)

	store, err := sqlstorage.New(cfg.Database.DSN())
	if err != nil {
		logg.Error("failed to connect to database", "error", err)
		os.Exit(1) //nolint:gocritic
	}
	defer store.Close(context.Background()) //nolint:errcheck

	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer producer.Close() //nolint:errcheck

	interval, err := time.ParseDuration(cfg.Scheduler.Interval)
	if err != nil {
		logg.Error("invalid scheduler interval", "error", err)
		os.Exit(1) //nolint:gocritic
	}

	sched := scheduler.New(logg, store, producer, interval)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_scheduler is running...")

	if err := sched.Run(ctx); err != nil {
		logg.Error("scheduler error", "error", err)
		os.Exit(1) //nolint:gocritic
	}
}
