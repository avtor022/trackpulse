-- Competitors table
CREATE TABLE IF NOT EXISTS competitors (
    id TEXT PRIMARY KEY NOT NULL,
    competitor_number INTEGER UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    birthday TEXT,
    country TEXT,
    city TEXT,
    rating INTEGER DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_competitors_number ON competitors(competitor_number);
CREATE INDEX IF NOT EXISTS idx_competitors_name ON competitors(full_name);
