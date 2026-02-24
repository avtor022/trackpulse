package models

import (
	"time"
)

// Race represents a race event in the database
type Race struct {
	ID                 string     `db:"id" json:"id"`
	RaceTitle          string     `db:"race_title" json:"race_title"`
	RaceType           string     `db:"race_type" json:"race_type"`
	ModelType          *string    `db:"model_type" json:"model_type,omitempty"`
	ModelScale         *string    `db:"model_scale" json:"model_scale,omitempty"`
	TrackName          *string    `db:"track_name" json:"track_name,omitempty"`
	LapCountTarget     *int       `db:"lap_count_target" json:"lap_count_target,omitempty"`
	TimeLimitMinutes   *int       `db:"time_limit_minutes" json:"time_limit_minutes,omitempty"`
	TimeStart          *time.Time `db:"time_start" json:"time_start,omitempty"`
	TimeFinish         *time.Time `db:"time_finish" json:"time_finish,omitempty"`
	Status             string     `db:"status" json:"status"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}