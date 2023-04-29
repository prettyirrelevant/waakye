-- +goose Up
CREATE TABLE IF NOT EXISTS oauth_credentials (
    id TEXT PRIMARY KEY,
    platform TEXT NOT NULL UNIQUE,
    credentials TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS oauth_credentials;
