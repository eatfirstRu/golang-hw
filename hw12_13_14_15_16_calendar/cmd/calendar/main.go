package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/server/http"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(cfg.Logger.Level)

	var store storage.Storage
	switch cfg.Storage.Type {
	case "sql":
		sqlStore, err := sqlstorage.New(cfg.Database.DSN())
		if err != nil {
			logg.Error("failed to connect to database", "error", err)
			os.Exit(1) //nolint:gocritic
		}
		defer sqlStore.Close(context.Background()) //nolint:errcheck
		store = sqlStore
	default:
		store = memorystorage.New()
	}

	calendar := app.New(logg, store)

	server := internalhttp.NewServer(logg, calendar, cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server", "error", err)
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server", "error", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
