-- RC Models table
CREATE TABLE IF NOT EXISTS rc_models (
    id TEXT PRIMARY KEY NOT NULL,
    brand TEXT NOT NULL REFERENCES rc_model_brands(name) ON DELETE RESTRICT,
    model_name TEXT NOT NULL,
    scale TEXT NOT NULL,
    model_type TEXT NOT NULL,
    motor_type TEXT,
    drive_type TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_models_brand ON rc_models(brand);
CREATE INDEX IF NOT EXISTS idx_models_type ON rc_models(model_type);
CREATE INDEX IF NOT EXISTS idx_models_scale ON rc_models(scale);
