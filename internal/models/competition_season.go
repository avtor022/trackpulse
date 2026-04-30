package models

import "time"

// CompetitionSeason represents a competition season in the catalog
type CompetitionSeason struct {
	ID        string    `json:"id" db:"id"`
	Season    string    `json:"season" db:"season"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
