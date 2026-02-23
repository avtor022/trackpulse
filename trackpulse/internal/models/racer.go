package models

import (
	"database/sql"
	"time"
)

// Racer represents a record in the racers table
type Racer struct {
	ID          string     `db:"id" json:"id"`
	RacerNumber int        `db:"racer_number" json:"racer_number"`
	FullName    string     `db:"full_name" json:"full_name"`
	Birthday    *string    `db:"birthday" json:"birthday,omitempty"`
	Country     *string    `db:"country" json:"country,omitempty"`
	City        *string    `db:"city" json:"city,omitempty"`
	Rating      int        `db:"rating" json:"rating"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (Racer) TableName() string {
	return "racers"
}

// Create creates a new record
func (r *Racer) Create(db *sql.DB) error {
	query := `INSERT INTO racers (id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, r.ID, r.RacerNumber, r.FullName, r.Birthday, r.Country, r.City, r.Rating, now, now)
	if err != nil {
		return err
	}
	
	r.CreatedAt = now
	r.UpdatedAt = now
	
	return nil
}

// GetByID gets a record by ID
func (r *Racer) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at 
	          FROM racers WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&r.ID, &r.RacerNumber, &r.FullName, &r.Birthday, &r.Country, &r.City, &r.Rating, &r.CreatedAt, &r.UpdatedAt,
	)
	return err
}

// Update updates a record
func (r *Racer) Update(db *sql.DB) error {
	query := `UPDATE racers SET racer_number = ?, full_name = ?, birthday = ?, country = ?, city = ?, rating = ?, updated_at = ? 
	          WHERE id = ?`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, r.RacerNumber, r.FullName, r.Birthday, r.Country, r.City, r.Rating, now, r.ID)
	if err != nil {
		return err
	}
	
	r.UpdatedAt = now
	return nil
}

// Delete deletes a record
func (r *Racer) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM racers WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
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
			&racer.ID, &racer.RacerNumber, &racer.FullName, &racer.Birthday, &racer.Country, &racer.City, &racer.Rating, 
			&racer.CreatedAt, &racer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		racers = append(racers, racer)
	}
	
	return racers, nil
}