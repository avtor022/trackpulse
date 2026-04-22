-- Lap History table
CREATE TABLE IF NOT EXISTS lap_history (
    id TEXT PRIMARY KEY NOT NULL,
    competition_participant_id TEXT NOT NULL REFERENCES competition_participants(id) ON DELETE CASCADE,
    lap_number INTEGER NOT NULL,
    lap_time_ms INTEGER NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    is_valid BOOLEAN DEFAULT 1,
    invalidation_reason TEXT,
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_lap_history_participant ON lap_history(competition_participant_id);
CREATE INDEX IF NOT EXISTS idx_lap_history_competition_lap ON lap_history(competition_participant_id, lap_number);
CREATE INDEX IF NOT EXISTS idx_lap_history_time ON lap_history(end_time);
