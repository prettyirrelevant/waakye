-- +goose Up
CREATE TABLE IF NOT EXISTS playlists_conversion_history (
    id TEXT PRIMARY KEY,
    playlist_url TEXT NOT NULL,
    source TEXT NOT NULL,
    destination TEXT NOT NULL,
    conversion_count INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
CREATE UNIQUE INDEX idx_unique_conversion_constraint ON playlists_conversion_history (playlist_url, source, destination);
CREATE INDEX IF NOT EXISTS idx_source ON playlists_conversion_history (source);
CREATE INDEX IF NOT EXISTS idx_destination ON playlists_conversion_history (destination);
CREATE INDEX IF NOT EXISTS idx_playlist_url ON playlists_conversion_history (playlist_url);
-- +goose Down
DROP TABLE IF EXISTS playlists_conversion_history;
