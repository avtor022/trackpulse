package models

import (
	"time"
)

// LapHistory represents individual lap records for a race participant
type LapHistory struct {
	ID                   string    `db:"id" json:"id"`
	RaceParticipantID    string    `db:"race_participant_id" json:"race_participant_id"`
	LapNumber            int       `db:"lap_number" json:"lap_number"`
	LapTimeMs            int       `db:"lap_time_ms" json:"lap_time_ms"`
	StartTime            time.Time `db:"start_time" json:"start_time"`
	EndTime              time.Time `db:"end_time" json:"end_time"`
	IsValid              bool      `db:"is_valid" json:"is_valid"`
	InvalidationReason   *string   `db:"invalidation_reason" json:"invalidation_reason,omitempty"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
}