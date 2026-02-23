package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=on", path))
	if err != nil {
		return nil, err
	}

	// Set connection pool settings for better concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create tables
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS racers (
			id TEXT PRIMARY KEY NOT NULL,
			racer_number INTEGER UNIQUE NOT NULL,
			full_name TEXT NOT NULL,
			birthday TEXT,
			country TEXT,
			city TEXT,
			rating INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_racers_number ON racers(racer_number);`,
		`CREATE INDEX IF NOT EXISTS idx_racers_name ON racers(full_name);`,

		`CREATE TABLE IF NOT EXISTS rc_models (
			id TEXT PRIMARY KEY NOT NULL,
			brand TEXT NOT NULL,
			model_name TEXT NOT NULL,
			scale TEXT NOT NULL,
			model_type TEXT NOT NULL,
			motor_type TEXT,
			drive_type TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,

		`CREATE INDEX IF NOT EXISTS idx_models_brand ON rc_models(brand);`,
		`CREATE INDEX IF NOT EXISTS idx_models_type ON rc_models(model_type);`,
		`CREATE INDEX IF NOT EXISTS idx_models_scale ON rc_models(scale);`,

		`CREATE TABLE IF NOT EXISTS racer_models (
			id TEXT PRIMARY KEY NOT NULL,
			racer_id TEXT NOT NULL,
			rc_model_id TEXT NOT NULL,
			transponder_number TEXT UNIQUE NOT NULL,
			transponder_type TEXT DEFAULT 'RFID',
			is_active BOOLEAN DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (racer_id) REFERENCES racers(id) ON DELETE CASCADE,
			UNIQUE(racer_id, rc_model_id)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_racer_models_transponder ON racer_models(transponder_number);`,
		`CREATE INDEX IF NOT EXISTS idx_racer_models_racer ON racer_models(racer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_racer_models_model ON racer_models(rc_model_id);`,

		`CREATE TABLE IF NOT EXISTS races (
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
		);`,

		`CREATE INDEX IF NOT EXISTS idx_races_status ON races(status);`,
		`CREATE INDEX IF NOT EXISTS idx_races_time_start ON races(time_start);`,

		`CREATE TABLE IF NOT EXISTS race_participants (
			id TEXT PRIMARY KEY NOT NULL,
			race_id TEXT NOT NULL,
			racer_model_id TEXT NOT NULL,
			grid_position INTEGER,
			is_finished BOOLEAN DEFAULT 0,
			disqualified BOOLEAN DEFAULT 0,
			dnf_reason TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE,
			FOREIGN KEY (racer_model_id) REFERENCES racer_models(id) ON DELETE CASCADE,
			UNIQUE(race_id, racer_model_id)
		);`,

		`CREATE INDEX IF NOT EXISTS idx_participants_race ON race_participants(race_id);`,
		`CREATE INDEX IF NOT EXISTS idx_participants_racer_model ON race_participants(racer_model_id);`,

		`CREATE TABLE IF NOT EXISTS race_laps (
			id TEXT PRIMARY KEY NOT NULL,
			race_participant_id TEXT NOT NULL UNIQUE,
			time_start TEXT NOT NULL,
			time_finish TEXT,
			number_of_laps INTEGER DEFAULT 0,
			best_lap_time_ms INTEGER DEFAULT 0,
			best_lap_number INTEGER DEFAULT 0,
			last_lap_time_ms INTEGER DEFAULT 0,
			last_pass_time TEXT,
			total_race_time_ms INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_race_laps_participant ON race_laps(race_participant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_race_laps_laps ON race_laps(number_of_laps);`,

		`CREATE TABLE IF NOT EXISTS lap_history (
			id TEXT PRIMARY KEY NOT NULL,
			race_participant_id TEXT NOT NULL,
			lap_number INTEGER NOT NULL,
			lap_time_ms INTEGER NOT NULL,
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			is_valid BOOLEAN DEFAULT 1,
			invalidation_reason TEXT,
			created_at TEXT NOT NULL,
			FOREIGN KEY (race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
		);`,

		`CREATE INDEX IF NOT EXISTS idx_lap_history_participant ON lap_history(race_participant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lap_history_race_lap ON lap_history(race_participant_id, lap_number);`,
		`CREATE INDEX IF NOT EXISTS idx_lap_history_time ON lap_history(start_time);`,

		`CREATE TABLE IF NOT EXISTS raw_scans (
			id TEXT PRIMARY KEY NOT NULL,
			timestamp TEXT NOT NULL,
			tag_value TEXT NOT NULL,
			reader_type TEXT NOT NULL,
			com_port TEXT,
			signal_strength INTEGER,
			is_processed BOOLEAN DEFAULT 0,
			linked_racer_model_id TEXT,
			created_at TEXT NOT NULL,
			FOREIGN KEY (linked_racer_model_id) REFERENCES racer_models(id) ON DELETE SET NULL
		);`,

		`CREATE INDEX IF NOT EXISTS idx_raw_scans_timestamp ON raw_scans(timestamp);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_tag ON raw_scans(tag_value);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_processed ON raw_scans(is_processed);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_tag_timestamp ON raw_scans(tag_value, timestamp);`,

		`CREATE TABLE IF NOT EXISTS system_settings (
			key TEXT PRIMARY KEY NOT NULL,
			value TEXT NOT NULL,
			value_type TEXT DEFAULT 'string',
			description TEXT,
			updated_at TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS audit_log (
			id TEXT PRIMARY KEY NOT NULL,
			timestamp TEXT NOT NULL,
			action_type TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id TEXT,
			user_name TEXT,
			ip_address TEXT,
			details TEXT,
			created_at TEXT NOT NULL
		);`,

		`CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_log(timestamp);`,
		`CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_log(action_type);`,
		`CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_log(entity_type, entity_id);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v, error: %v", query, err)
		}
	}

	return nil
}