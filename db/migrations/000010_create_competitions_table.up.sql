-- Competitions table
CREATE TABLE IF NOT EXISTS competitions (
    id TEXT PRIMARY KEY NOT NULL,
    competition_title TEXT NOT NULL,
    competition_type TEXT DEFAULT 'qualifying',
    model_type TEXT,
    model_scale TEXT,
    track_name TEXT,
    lap_count_target INTEGER,
    time_limit_minutes INTEGER,
    time_start TEXT,
    time_finish TEXT,
    status TEXT DEFAULT 'scheduled',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_competitions_status ON competitions(status);
CREATE INDEX IF NOT EXISTS idx_competitions_time_start ON competitions(time_start);
