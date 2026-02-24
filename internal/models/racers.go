package models

import (
	"time"
)

// Racer represents a racer in the database
type Racer struct {
	ID          string    `db:"id" json:"id"`
	RacerNumber int       `db:"racer_number" json:"racer_number"`
	FullName    string    `db:"full_name" json:"full_name"`
	Birthday    *string   `db:"birthday" json:"birthday,omitempty"`
	Country     *string   `db:"country" json:"country,omitempty"`
	City        *string   `db:"city" json:"city,omitempty"`
	Rating      int       `db:"rating" json:"rating"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}