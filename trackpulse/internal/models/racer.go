package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// Racer represents a single racer
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

// TableName returns the table name for racers
func (Racer) TableName() string {
	return "racers"
}

// Create inserts a new racer into the database
func (r *Racer) Create(db *sql.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	r.CreatedAt = now
	r.UpdatedAt = now

	query := `INSERT INTO racers (id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, r.ID, r.RacerNumber, r.FullName, r.Birthday, r.Country, r.City, r.Rating, r.CreatedAt, r.UpdatedAt)
	return err
}

// GetByID retrieves a racer by ID
func (r *Racer) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at 
	          FROM racers WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&r.ID, &r.RacerNumber, &r.FullName, &r.Birthday, &r.Country, &r.City, &r.Rating, &r.CreatedAt, &r.UpdatedAt,
	)
	return err
}

// Update updates an existing racer in the database
func (r *Racer) Update(db *sql.DB) error {
	r.UpdatedAt = time.Now().UTC()

	query := `UPDATE racers SET racer_number = ?, full_name = ?, birthday = ?, country = ?, city = ?, 
	                 rating = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, r.RacerNumber, r.FullName, r.Birthday, r.Country, r.City, r.Rating, r.UpdatedAt, r.ID)
	return err
}

// Delete removes a racer from the database
func (r *Racer) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM racers WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all racers from the database
func (r *Racer) GetAll(db *sql.DB) ([]Racer, error) {
	query := `SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at 
	          FROM racers ORDER BY racer_number`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racers []Racer
	for rows.Next() {
		var racer Racer
		err := rows.Scan(
			&racer.ID, &racer.RacerNumber, &racer.FullName, &racer.Birthday, &racer.Country, &racer.City, 
			&racer.Rating, &racer.CreatedAt, &racer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		racers = append(racers, racer)
	}

	return racers, nil
}

// Validate checks if the racer data is valid
func (r *Racer) Validate() error {
	if r.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if r.RacerNumber <= 0 {
		return fmt.Errorf("racer number must be positive")
	}
	return nil
}