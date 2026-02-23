package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         string    `db:"id" json:"id"`
	Timestamp  time.Time `db:"timestamp" json:"timestamp"`
	ActionType string    `db:"action_type" json:"action_type"`
	EntityType string    `db:"entity_type" json:"entity_type"`
	EntityID   *string   `db:"entity_id" json:"entity_id,omitempty"`
	UserName   *string   `db:"user_name" json:"user_name,omitempty"`
	IPAddress  *string   `db:"ip_address" json:"ip_address,omitempty"`
	Details    *string   `db:"details" json:"details,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name for audit_log
func (AuditLog) TableName() string {
	return "audit_log"
}

// Create inserts a new audit log entry into the database
func (al *AuditLog) Create(db *sql.DB) error {
	if al.ID == "" {
		al.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	al.Timestamp = now
	al.CreatedAt = now

	query := `INSERT INTO audit_log (id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, al.ID, al.Timestamp, al.ActionType, al.EntityType, al.EntityID, al.UserName, al.IPAddress, al.Details, al.CreatedAt)
	return err
}

// GetByID retrieves an audit log entry by ID
func (al *AuditLog) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at 
	          FROM audit_log WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&al.ID, &al.Timestamp, &al.ActionType, &al.EntityType, &al.EntityID, &al.UserName, &al.IPAddress, &al.Details, &al.CreatedAt,
	)
	return err
}

// GetAll retrieves all audit log entries from the database
func (al *AuditLog) GetAll(db *sql.DB) ([]AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at 
	          FROM audit_log ORDER BY timestamp DESC LIMIT 1000`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditLogs []AuditLog
	for rows.Next() {
		var auditLog AuditLog
		err := rows.Scan(
			&auditLog.ID, &auditLog.Timestamp, &auditLog.ActionType, &auditLog.EntityType, &auditLog.EntityID, 
			&auditLog.UserName, &auditLog.IPAddress, &auditLog.Details, &auditLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		auditLogs = append(auditLogs, auditLog)
	}

	return auditLogs, nil
}