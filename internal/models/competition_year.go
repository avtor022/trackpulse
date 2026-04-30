package models

import "time"

// CompetitionYear represents a competition year in the catalog
type CompetitionYear struct {
	ID        string    `json:"id" db:"id"`
	Year      int       `json:"year" db:"year"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
