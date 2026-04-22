-- Drop tables in reverse order of creation (respecting foreign key dependencies)
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS raw_scans;
DROP TABLE IF EXISTS lap_history;
DROP TABLE IF EXISTS competition_laps;
DROP TABLE IF EXISTS competition_participants;
DROP TABLE IF EXISTS competitions;
DROP TABLE IF EXISTS competitor_models;
DROP TABLE IF EXISTS competition_tracks;
DROP TABLE IF EXISTS rc_models;
DROP TABLE IF EXISTS rc_model_types;
DROP TABLE IF EXISTS rc_model_scales;
DROP TABLE IF EXISTS rc_model_brands;
DROP TABLE IF EXISTS competitors;
