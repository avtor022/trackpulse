package models

import (
	"database/sql"
	"time"
)

// SystemSetting represents a record in the system_settings table
type SystemSetting struct {
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	ValueType   string    `db:"value_type" json:"value_type"`
	Description *string   `db:"description" json:"description,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (SystemSetting) TableName() string {
	return "system_settings"
}

// Create creates a new record
func (ss *SystemSetting) Create(db *sql.DB) error {
	query := `INSERT INTO system_settings (key, value, value_type, description, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, ss.Key, ss.Value, ss.ValueType, ss.Description, now)
	if err != nil {
		return err
	}
	
	ss.UpdatedAt = now
	
	return nil
}

// GetByKey gets a record by key
func (ss *SystemSetting) GetByKey(db *sql.DB, key string) error {
	query := `SELECT key, value, value_type, description, updated_at 
	          FROM system_settings WHERE key = ?`
	
	err := db.QueryRow(query, key).Scan(
		&ss.Key, &ss.Value, &ss.ValueType, &ss.Description, &ss.UpdatedAt,
	)
	return err
}

// Update updates a record
func (ss *SystemSetting) Update(db *sql.DB) error {
	query := `UPDATE system_settings SET value = ?, value_type = ?, description = ?, updated_at = ? 
	          WHERE key = ?`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, ss.Value, ss.ValueType, ss.Description, now, ss.Key)
	if err != nil {
		return err
	}
	
	ss.UpdatedAt = now
	return nil
}

// Delete deletes a record
func (ss *SystemSetting) Delete(db *sql.DB, key string) error {
	query := `DELETE FROM system_settings WHERE key = ?`
	_, err := db.Exec(query, key)
	return err
}

// GetAll gets all records
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