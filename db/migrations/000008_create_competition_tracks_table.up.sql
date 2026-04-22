-- Competition Tracks dictionary table
CREATE TABLE IF NOT EXISTS competition_tracks (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_tracks_name ON competition_tracks(name);
