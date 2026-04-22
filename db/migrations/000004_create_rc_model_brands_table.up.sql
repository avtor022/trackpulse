-- RC Model Brands dictionary table
CREATE TABLE IF NOT EXISTS rc_model_brands (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_brands_name ON rc_model_brands(name);
