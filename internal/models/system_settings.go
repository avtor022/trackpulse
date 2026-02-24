package models

import (
	"time"
)

// SystemSetting represents a system configuration setting
type SystemSetting struct {
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	ValueType   string    `db:"value_type" json:"value_type"`
	Description *string   `db:"description" json:"description,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}