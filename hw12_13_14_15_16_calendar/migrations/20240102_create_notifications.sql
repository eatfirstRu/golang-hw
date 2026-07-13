-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
    id         SERIAL PRIMARY KEY,
    event_id   VARCHAR(255)             NOT NULL,
    title      VARCHAR(255)             NOT NULL,
    date_time  TIMESTAMP WITH TIME ZONE NOT NULL,
    user_id    VARCHAR(255)             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_event_id ON notifications (event_id);

-- +goose Down
DROP TABLE IF EXISTS notifications;
