package repository

import (
	"database/sql"
	"time"
)

// Setting represents a system setting
type Setting struct {
	Key         string
	Value       string
	ValueType   string
	Description string
	UpdatedAt   time.Time
}

// SettingsRepository handles database operations for settings
type SettingsRepository struct {
	db *sql.DB
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// Get retrieves a setting by key
func (r *SettingsRepository) Get(key string) (*Setting, error) {
	query := `SELECT key, value, value_type, description, updated_at 
			  FROM system_settings WHERE key = ?`
	
	row := r.db.QueryRow(query, key)
	
	var s Setting
	var updatedAt string
	err := row.Scan(&s.Key, &s.Value, &s.ValueType, &s.Description, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	s.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &s, nil
}

// GetValue retrieves the value of a setting by key
func (r *SettingsRepository) GetValue(key string) (string, error) {
	query := `SELECT value FROM system_settings WHERE key = ?`
	
	var value string
	err := r.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	
	return value, nil
}

// Set updates or inserts a setting
func (r *SettingsRepository) Set(key, value, valueType, description string) error {
	now := time.Now().Format(time.RFC3339)
	
	query := `INSERT OR REPLACE INTO system_settings (key, value, value_type, description, updated_at)
			  VALUES (?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, key, value, valueType, description, now)
	return err
}

// UpdateValue updates only the value of an existing setting
func (r *SettingsRepository) UpdateValue(key, value string) error {
	now := time.Now().Format(time.RFC3339)
	
	query := `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ?`
	
	_, err := r.db.Exec(query, value, now, key)
	return err
}
