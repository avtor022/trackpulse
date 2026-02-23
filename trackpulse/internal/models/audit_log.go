package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// AuditLog represents a record in the audit_log table
type AuditLog struct {
	ID          string    `db:"id" json:"id"`
	Timestamp   string    `db:"timestamp" json:"timestamp"`
	ActionType  string    `db:"action_type" json:"action_type"`
	EntityType  string    `db:"entity_type" json:"entity_type"`
	EntityID    *string   `db:"entity_id" json:"entity_id,omitempty"`
	UserName    *string   `db:"user_name" json:"user_name,omitempty"`
	IPAddress   *string   `db:"ip_address" json:"ip_address,omitempty"`
	Details     *string   `db:"details" json:"details,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name
func (AuditLog) TableName() string {
	return "audit_log"
}

// Create creates a new record
func (al *AuditLog) Create(db *sql.DB) error {
	al.ID = uuid.New().String()
	al.CreatedAt = time.Now()

	query := `INSERT INTO audit_log (id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, al.ID, al.Timestamp, al.ActionType, al.EntityType, al.EntityID, al.UserName, al.IPAddress, al.Details, al.CreatedAt)
	return err
}

// GetByID gets a record by ID
func (al *AuditLog) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at 
	          FROM audit_log WHERE id = ?`
	
	return db.QueryRow(query, id).Scan(&al.ID, &al.Timestamp, &al.ActionType, &al.EntityType, &al.EntityID, &al.UserName, &al.IPAddress, &al.Details, &al.CreatedAt)
}

// GetAll gets all records
func (al *AuditLog) GetAll(db *sql.DB) ([]AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at 
	          FROM audit_log ORDER BY timestamp DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(&log.ID, &log.Timestamp, &log.ActionType, &log.EntityType, &log.EntityID, &log.UserName, &log.IPAddress, &log.Details, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}