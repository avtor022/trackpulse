package models

import (
	"time"
)

// RaceLap represents aggregated lap data for a race participant
type RaceLap struct {
	ID                  string    `db:"id" json:"id"`
	RaceParticipantID   string    `db:"race_participant_id" json:"race_participant_id"`
	TimeStart           time.Time `db:"time_start" json:"time_start"`
	TimeFinish          *time.Time `db:"time_finish" json:"time_finish,omitempty"`
	NumberOfLaps        int       `db:"number_of_laps" json:"number_of_laps"`
	BestLapTimeMs       int       `db:"best_lap_time_ms" json:"best_lap_time_ms"`
	BestLapNumber       int       `db:"best_lap_number" json:"best_lap_number"`
	LastLapTimeMs       int       `db:"last_lap_time_ms" json:"last_lap_time_ms"`
	LastPassTime        *time.Time `db:"last_pass_time" json:"last_pass_time,omitempty"`
	TotalRaceTimeMs     int       `db:"total_race_time_ms" json:"total_race_time_ms"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}