package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// SystemSetting represents a system setting
type SystemSetting struct {
	Key        string    `db:"key" json:"key"`
	Value      string    `db:"value" json:"value"`
	ValueType  string    `db:"value_type" json:"value_type"`
	Description *string  `db:"description" json:"description,omitempty"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for system_settings
func (SystemSetting) TableName() string {
	return "system_settings"
}

// Create inserts a new system setting into the database
func (ss *SystemSetting) Create(db *sql.DB) error {
	now := time.Now().UTC()
	ss.UpdatedAt = now
	if ss.ValueType == "" {
		ss.ValueType = "string"
	}

	query := `INSERT INTO system_settings (key, value, value_type, description, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, ss.Key, ss.Value, ss.ValueType, ss.Description, ss.UpdatedAt)
	return err
}

// GetByKey retrieves a system setting by key
func (ss *SystemSetting) GetByKey(db *sql.DB, key string) error {
	query := `SELECT key, value, value_type, description, updated_at 
	          FROM system_settings WHERE key = ?`
	
	err := db.QueryRow(query, key).Scan(
		&ss.Key, &ss.Value, &ss.ValueType, &ss.Description, &ss.UpdatedAt,
	)
	return err
}

// Update updates an existing system setting in the database
func (ss *SystemSetting) Update(db *sql.DB) error {
	ss.UpdatedAt = time.Now().UTC()

	query := `UPDATE system_settings SET value = ?, value_type = ?, description = ?, updated_at = ? WHERE key = ?`
	
	_, err := db.Exec(query, ss.Value, ss.ValueType, ss.Description, ss.UpdatedAt, ss.Key)
	return err
}

// Delete removes a system setting from the database
func (ss *SystemSetting) Delete(db *sql.DB, key string) error {
	query := `DELETE FROM system_settings WHERE key = ?`
	_, err := db.Exec(query, key)
	return err
}

// GetAll retrieves all system settings from the database
func (ss *SystemSetting) GetAll(db *sql.DB) ([]SystemSetting, error) {
	query := `SELECT key, value, value_type, description, updated_at 
	          FROM system_settings ORDER BY key`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []SystemSetting
	for rows.Next() {
		var setting SystemSetting
		err := rows.Scan(
			&setting.Key, &setting.Value, &setting.ValueType, &setting.Description, &setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

// Validate checks if the system setting data is valid
func (ss *SystemSetting) Validate() error {
	if ss.Key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}