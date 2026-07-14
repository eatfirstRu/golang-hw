package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func getDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=calendar password=calendar dbname=calendar sslmode=disable"
	}
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		t.Fatalf("connect to db: %v", err)
	}
	return db
}

func TestStorerSavesNotifications(t *testing.T) {
	now := time.Now().Add(5 * time.Minute).Truncate(time.Second)

	evt := createEvent(t, eventRequest{
		Title:        "Storer Test Event",
		DateTime:     now.Format(time.RFC3339),
		Duration:     "30m",
		UserID:       "integ-user-storer",
		NotifyBefore: "10m",
	})

	_ = evt

	db := getDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	var count int
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting for notification in DB, last count: %d", count)
		default:
		}

		err := db.GetContext(ctx, &count,
			"SELECT COUNT(*) FROM notifications WHERE title = $1 AND user_id = $2",
			"Storer Test Event", "integ-user-storer")
		if err != nil {
			t.Logf("query error (retrying): %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if count > 0 {
			t.Logf("notification found in DB after polling")
			return
		}

		fmt.Println("notification not yet in DB, waiting...")
		time.Sleep(3 * time.Second)
	}
}
