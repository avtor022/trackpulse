package models

import (
	"time"
)

// RacerModel represents the binding between a racer, model, and transponder
type RacerModel struct {
	ID                string    `db:"id" json:"id"`
	RacerID           string    `db:"racer_id" json:"racer_id"`
	RCModelID         string    `db:"rc_model_id" json:"rc_model_id"`
	TransponderNumber string    `db:"transponder_number" json:"transponder_number"`
	TransponderType   string    `db:"transponder_type" json:"transponder_type"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}