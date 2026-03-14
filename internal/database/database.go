package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection with WAL mode enabled
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &DB{db}, nil
}

// Initialize creates all required tables if they don't exist
func (db *DB) Initialize() error {
	schema := `
	-- Racers table
	CREATE TABLE IF NOT EXISTS racers (
		id TEXT PRIMARY KEY NOT NULL,
		racer_number INTEGER UNIQUE NOT NULL,
		full_name TEXT NOT NULL,
		birthday TEXT,
		country TEXT,
		city TEXT,
		rating INTEGER DEFAULT 0,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_racers_number ON racers(racer_number);
	CREATE INDEX IF NOT EXISTS idx_racers_name ON racers(full_name);

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

	-- RC Model Brands dictionary table
	CREATE TABLE IF NOT EXISTS rc_model_brands (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT UNIQUE NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_brands_name ON rc_model_brands(name);

	-- Racer Models (transponders) table
	CREATE TABLE IF NOT EXISTS racer_models (
		id TEXT PRIMARY KEY NOT NULL,
		racer_id TEXT NOT NULL REFERENCES racers(id) ON DELETE CASCADE,
		rc_model_id TEXT NOT NULL REFERENCES rc_models(id) ON DELETE CASCADE,
		transponder_number TEXT NOT NULL UNIQUE,
		transponder_type TEXT DEFAULT 'RFID',
		is_active BOOLEAN DEFAULT 1,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_racer_models_transponder ON racer_models(transponder_number);
	CREATE INDEX IF NOT EXISTS idx_racer_models_racer ON racer_models(racer_id);
	CREATE INDEX IF NOT EXISTS idx_racer_models_model ON racer_models(rc_model_id);

	-- Races table
	CREATE TABLE IF NOT EXISTS races (
		id TEXT PRIMARY KEY NOT NULL,
		race_title TEXT NOT NULL,
		race_type TEXT DEFAULT 'qualifying',
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
	CREATE INDEX IF NOT EXISTS idx_races_status ON races(status);
	CREATE INDEX IF NOT EXISTS idx_races_time_start ON races(time_start);

	-- Race Participants table
	CREATE TABLE IF NOT EXISTS race_participants (
		id TEXT PRIMARY KEY NOT NULL,
		race_id TEXT NOT NULL REFERENCES races(id) ON DELETE CASCADE,
		racer_model_id TEXT NOT NULL REFERENCES racer_models(id) ON DELETE CASCADE,
		grid_position INTEGER,
		is_finished BOOLEAN DEFAULT 0,
		disqualified BOOLEAN DEFAULT 0,
		dnf_reason TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		UNIQUE(race_id, racer_model_id)
	);
	CREATE INDEX IF NOT EXISTS idx_participants_race ON race_participants(race_id);
	CREATE INDEX IF NOT EXISTS idx_participants_racer_model ON race_participants(racer_model_id);

	-- Race Laps (aggregated) table
	CREATE TABLE IF NOT EXISTS race_laps (
		id TEXT PRIMARY KEY NOT NULL,
		race_participant_id TEXT NOT NULL UNIQUE REFERENCES race_participants(id) ON DELETE CASCADE,
		time_start TEXT NOT NULL,
		time_finish TEXT,
		number_of_laps INTEGER DEFAULT 0,
		best_lap_time_ms INTEGER DEFAULT 0,
		best_lap_number INTEGER DEFAULT 0,
		last_lap_time_ms INTEGER DEFAULT 0,
		last_pass_time TEXT,
		total_race_time_ms INTEGER DEFAULT 0,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_race_laps_participant ON race_laps(race_participant_id);
	CREATE INDEX IF NOT EXISTS idx_race_laps_laps ON race_laps(number_of_laps DESC);

	-- Lap History table
	CREATE TABLE IF NOT EXISTS lap_history (
		id TEXT PRIMARY KEY NOT NULL,
		race_participant_id TEXT NOT NULL REFERENCES race_participants(id) ON DELETE CASCADE,
		lap_number INTEGER NOT NULL,
		lap_time_ms INTEGER NOT NULL,
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		is_valid BOOLEAN DEFAULT 1,
		invalidation_reason TEXT,
		created_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_lap_history_participant ON lap_history(race_participant_id);
	CREATE INDEX IF NOT EXISTS idx_lap_history_race_lap ON lap_history(race_participant_id, lap_number);
	CREATE INDEX IF NOT EXISTS idx_lap_history_time ON lap_history(end_time);

	-- Raw Scans table
	CREATE TABLE IF NOT EXISTS raw_scans (
		id TEXT PRIMARY KEY NOT NULL,
		timestamp TEXT NOT NULL,
		tag_value TEXT NOT NULL,
		reader_type TEXT NOT NULL,
		com_port TEXT,
		signal_strength INTEGER,
		is_processed BOOLEAN DEFAULT 0,
		linked_racer_model_id TEXT REFERENCES racer_models(id) ON DELETE SET NULL,
		created_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_raw_scans_timestamp ON raw_scans(timestamp);
	CREATE INDEX IF NOT EXISTS idx_raw_scans_tag ON raw_scans(tag_value);
	CREATE INDEX IF NOT EXISTS idx_raw_scans_processed ON raw_scans(is_processed);
	CREATE INDEX IF NOT EXISTS idx_raw_scans_tag_timestamp ON raw_scans(tag_value, timestamp);

	-- System Settings table
	CREATE TABLE IF NOT EXISTS system_settings (
		key TEXT PRIMARY KEY NOT NULL,
		value TEXT NOT NULL,
		value_type TEXT DEFAULT 'string',
		description TEXT,
		updated_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_settings_key ON system_settings(key);

	-- Audit Log table
	CREATE TABLE IF NOT EXISTS audit_log (
		id TEXT PRIMARY KEY NOT NULL,
		timestamp TEXT NOT NULL,
		action_type TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id TEXT,
		user_name TEXT,
		ip_address TEXT,
		details TEXT,
		created_at TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_log(timestamp);
	CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action_type);
	CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_log(entity_type, entity_id);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Insert default settings if not exists
	defaultSettings := []struct {
		key         string
		value       string
		valueType   string
		description string
	}{
		{"web_interface_enabled", "false", "bool", "Enable web interface for spectators"},
		{"ui_language", "en", "string", "UI language (ru, en)"},
		{"db_path", "", "string", "Path to database file"},
		{"db_api_user", "admin", "string", "API user"},
		{"db_api_password_hash", "", "string", "API password hash (bcrypt)"},
		{"hardware_com_port", "", "string", "COM port for Arduino"},
		{"hardware_reader_type", "EM4095", "string", "RFID reader type"},
		{"hardware_debounce_ms", "2000", "int", "Debounce delay in ms"},
		{"log_retention_years", "2", "int", "Log retention period for manual cleanup"},
		{"last_log_cleanup_date", "", "string", "Last manual log cleanup date"},
		{"app_version", "1.0.0", "string", "Application version"},
		{"schema_version", "1", "int", "Database schema version"},
	}

	now := time.Now().Format(time.RFC3339)
	for _, setting := range defaultSettings {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO system_settings (key, value, value_type, description, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
			setting.key, setting.value, setting.valueType, setting.description, now)
		if err != nil {
			return fmt.Errorf("failed to insert default setting %s: %w", setting.key, err)
		}
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
