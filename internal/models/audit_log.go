package models

import (
	"time"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID          string    `db:"id" json:"id"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
	ActionType  string    `db:"action_type" json:"action_type"`
	EntityType  string    `db:"entity_type" json:"entity_type"`
	EntityID    *string   `db:"entity_id" json:"entity_id,omitempty"`
	UserName    *string   `db:"user_name" json:"user_name,omitempty"`
	IPAddress   *string   `db:"ip_address" json:"ip_address,omitempty"`
	Details     *string   `db:"details" json:"details,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}