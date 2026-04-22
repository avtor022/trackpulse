-- Competition Laps (aggregated) table
CREATE TABLE IF NOT EXISTS competition_laps (
    id TEXT PRIMARY KEY NOT NULL,
    competition_participant_id TEXT NOT NULL UNIQUE REFERENCES competition_participants(id) ON DELETE CASCADE,
    time_start TEXT NOT NULL,
    time_finish TEXT,
    number_of_laps INTEGER DEFAULT 0,
    best_lap_time_ms INTEGER DEFAULT 0,
    best_lap_number INTEGER DEFAULT 0,
    last_lap_time_ms INTEGER DEFAULT 0,
    last_pass_time TEXT,
    total_competition_time_ms INTEGER DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_competition_laps_participant ON competition_laps(competition_participant_id);
CREATE INDEX IF NOT EXISTS idx_competition_laps_laps ON competition_laps(number_of_laps DESC);
