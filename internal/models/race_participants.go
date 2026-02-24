package models

import (
	"time"
)

// RaceParticipant represents a participant in a race
type RaceParticipant struct {
	ID              string    `db:"id" json:"id"`
	RaceID          string    `db:"race_id" json:"race_id"`
	RacerModelID    string    `db:"racer_model_id" json:"racer_model_id"`
	GridPosition    *int      `db:"grid_position" json:"grid_position,omitempty"`
	IsFinished      bool      `db:"is_finished" json:"is_finished"`
	Disqualified    bool      `db:"disqualified" json:"disqualified"`
	DNFReason       *string   `db:"dnf_reason" json:"dnf_reason,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}