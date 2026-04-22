-- Competitor Models (transponders) table
CREATE TABLE IF NOT EXISTS competitor_models (
    id TEXT PRIMARY KEY NOT NULL,
    competitor_id TEXT NOT NULL REFERENCES competitors(id) ON DELETE CASCADE,
    rc_model_id TEXT NOT NULL REFERENCES rc_models(id) ON DELETE CASCADE,
    transponder_number TEXT NOT NULL UNIQUE,
    transponder_type TEXT DEFAULT 'RFID',
    is_active BOOLEAN DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_competitor_models_transponder ON competitor_models(transponder_number);
CREATE INDEX IF NOT EXISTS idx_competitor_models_racer ON competitor_models(competitor_id);
CREATE INDEX IF NOT EXISTS idx_competitor_models_model ON competitor_models(rc_model_id);
