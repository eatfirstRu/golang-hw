-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id            VARCHAR(255) PRIMARY KEY,
    title         VARCHAR(255)             NOT NULL,
    date_time     TIMESTAMP WITH TIME ZONE NOT NULL,
    duration      BIGINT                   NOT NULL DEFAULT 0,
    description   TEXT                     NOT NULL DEFAULT '',
    user_id       VARCHAR(255)             NOT NULL,
    notify_before BIGINT                   NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_events_date_time ON events (date_time);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events (user_id);

-- +goose Down
DROP TABLE IF EXISTS events;
