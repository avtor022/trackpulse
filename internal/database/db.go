package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// InitializeDB creates the database tables if they don't exist
func InitializeDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set PRAGMA settings
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-64000;", // 64MB cache
		"PRAGMA temp_store=MEMORY;",
	}

	for _, pragma := range pragmas {
		_, err := db.Exec(pragma)
		if err != nil {
			return nil, fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	// Create tables
	tables := []string{
		// Racers table
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
		
		// RC Models table
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
		
		// Racer Models (bindings) table
		`CREATE TABLE IF NOT EXISTS racer_models (
			id TEXT PRIMARY KEY NOT NULL,
			racer_id TEXT NOT NULL,
			rc_model_id TEXT NOT NULL,
			transponder_number TEXT UNIQUE NOT NULL,
			transponder_type TEXT DEFAULT 'RFID',
			is_active BOOLEAN DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY(racer_id) REFERENCES racers(id) ON DELETE CASCADE,
			FOREIGN KEY(rc_model_id) REFERENCES rc_models(id) ON DELETE CASCADE
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_racer_models_transponder ON racer_models(transponder_number);`,
		`CREATE INDEX IF NOT EXISTS idx_racer_models_racer ON racer_models(racer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_racer_models_model ON racer_models(rc_model_id);`,
		
		// Races table
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
		
		// Race Participants table
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
			FOREIGN KEY(race_id) REFERENCES races(id) ON DELETE CASCADE,
			FOREIGN KEY(racer_model_id) REFERENCES racer_models(id) ON DELETE CASCADE
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_participants_race ON race_participants(race_id);`,
		`CREATE INDEX IF NOT EXISTS idx_participants_racer_model ON race_participants(racer_model_id);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_participants_unique ON race_participants(race_id, racer_model_id);`,
		
		// Race Laps table
		`CREATE TABLE IF NOT EXISTS race_laps (
			id TEXT PRIMARY KEY NOT NULL,
			race_participant_id TEXT UNIQUE NOT NULL,
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
			FOREIGN KEY(race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_race_laps_participant ON race_laps(race_participant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_race_laps_laps ON race_laps(number_of_laps DESC);`,
		
		// Lap History table
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
			FOREIGN KEY(race_participant_id) REFERENCES race_participants(id) ON DELETE CASCADE
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_lap_history_participant ON lap_history(race_participant_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lap_history_race_lap ON lap_history(race_participant_id, lap_number);`,
		`CREATE INDEX IF NOT EXISTS idx_lap_history_time ON lap_history(end_time);`,
		
		// Raw Scans table
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
			FOREIGN KEY(linked_racer_model_id) REFERENCES racer_models(id) ON DELETE SET NULL
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_timestamp ON raw_scans(timestamp);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_tag ON raw_scans(tag_value);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_processed ON raw_scans(is_processed);`,
		`CREATE INDEX IF NOT EXISTS idx_raw_scans_tag_timestamp ON raw_scans(tag_value, timestamp);`,
		
		// System Settings table
		`CREATE TABLE IF NOT EXISTS system_settings (
			key TEXT PRIMARY KEY NOT NULL,
			value TEXT NOT NULL,
			value_type TEXT DEFAULT 'string',
			description TEXT,
			updated_at TEXT NOT NULL
		);`,
		
		`CREATE INDEX IF NOT EXISTS idx_settings_key ON system_settings(key);`,
		
		// Audit Log table
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

	for _, tableSQL := range tables {
		_, err := db.Exec(tableSQL)
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Check if system_settings table is empty (first run)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM system_settings").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check system_settings: %w", err)
	}

	if count == 0 {
		// Insert default settings including auth credentials
		password := "admin"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash default password: %w", err)
		}

		defaultSettings := map[string]string{
			"ui_language":             "ru",
			"db_path":                 "./track_pulse.db",
			"hardware_com_port":       "",
			"hardware_reader_type":    "EM4095",
			"hardware_debounce_ms":    "2000",
			"log_retention_years":     "1",
			"auth_user":               "admin",
			"auth_password_hash":      string(hashedPassword),
		}

		for key, value := range defaultSettings {
			_, err := db.Exec(`
				INSERT INTO system_settings (key, value, value_type, updated_at) 
				VALUES (?, ?, ?, datetime('now'))
			`, key, value, getValueType(value))
			if err != nil {
				return nil, fmt.Errorf("failed to insert default setting %s: %w", key, err)
			}
		}
		
		log.Println("Database initialized with default settings and admin credentials (admin/admin)")
	}

	return db, nil
}

// getValueType determines the type of a value for system_settings
func getValueType(value string) string {
	// Simple heuristic to determine type based on content
	switch value {
	case "true", "false":
		return "bool"
	}
	
	// Try to parse as integer
	for _, r := range value {
		if r < '0' || r > '9' {
			return "string"
		}
	}
	return "int"
}