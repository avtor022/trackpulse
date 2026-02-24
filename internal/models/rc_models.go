package models

import (
	"time"
)

// RCModel represents an RC model in the database
type RCModel struct {
	ID         string    `db:"id" json:"id"`
	Brand      string    `db:"brand" json:"brand"`
	ModelName  string    `db:"model_name" json:"model_name"`
	Scale      string    `db:"scale" json:"scale"`
	ModelType  string    `db:"model_type" json:"model_type"`
	MotorType  *string   `db:"motor_type" json:"motor_type,omitempty"`
	DriveType  *string   `db:"drive_type" json:"drive_type,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}