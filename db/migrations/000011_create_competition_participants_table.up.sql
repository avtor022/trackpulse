-- Competition Participants table
CREATE TABLE IF NOT EXISTS competition_participants (
    id TEXT PRIMARY KEY NOT NULL,
    competition_id TEXT NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
    competitor_model_id TEXT NOT NULL REFERENCES competitor_models(id) ON DELETE CASCADE,
    grid_position INTEGER,
    is_finished BOOLEAN DEFAULT 0,
    disqualified BOOLEAN DEFAULT 0,
    dnf_reason TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(competition_id, competitor_model_id)
);
CREATE INDEX IF NOT EXISTS idx_participants_competition ON competition_participants(competition_id);
CREATE INDEX IF NOT EXISTS idx_participants_competitor_model ON competition_participants(competitor_model_id);
