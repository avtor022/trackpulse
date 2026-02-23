package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// Race represents a single race event
type Race struct {
	ID               string     `db:"id" json:"id"`
	RaceTitle        string     `db:"race_title" json:"race_title"`
	RaceType         string     `db:"race_type" json:"race_type"`
	ModelType        *string    `db:"model_type" json:"model_type,omitempty"`
	ModelScale       *string    `db:"model_scale" json:"model_scale,omitempty"`
	TrackName        *string    `db:"track_name" json:"track_name,omitempty"`
	LapCountTarget   *int       `db:"lap_count_target" json:"lap_count_target,omitempty"`
	TimeLimitMinutes *int       `db:"time_limit_minutes" json:"time_limit_minutes,omitempty"`
	TimeStart        *time.Time `db:"time_start" json:"time_start,omitempty"`
	TimeFinish       *time.Time `db:"time_finish" json:"time_finish,omitempty"`
	Status           string     `db:"status" json:"status"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for races
func (Race) TableName() string {
	return "races"
}

// Create inserts a new race into the database
func (r *Race) Create(db *sql.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	r.CreatedAt = now
	r.UpdatedAt = now
	if r.Status == "" {
		r.Status = "scheduled"
	}

	query := `INSERT INTO races (id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, time_start, time_finish, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, r.ID, r.RaceTitle, r.RaceType, r.ModelType, r.ModelScale, r.TrackName, r.LapCountTarget, r.TimeLimitMinutes, r.TimeStart, r.TimeFinish, r.Status, r.CreatedAt, r.UpdatedAt)
	return err
}

// GetByID retrieves a race by ID
func (r *Race) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, time_start, time_finish, status, created_at, updated_at 
	          FROM races WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&r.ID, &r.RaceTitle, &r.RaceType, &r.ModelType, &r.ModelScale, &r.TrackName, &r.LapCountTarget, &r.TimeLimitMinutes, 
		&r.TimeStart, &r.TimeFinish, &r.Status, &r.CreatedAt, &r.UpdatedAt,
	)
	return err
}

// Update updates an existing race in the database
func (r *Race) Update(db *sql.DB) error {
	r.UpdatedAt = time.Now().UTC()

	query := `UPDATE races SET race_title = ?, race_type = ?, model_type = ?, model_scale = ?, track_name = ?, 
	                 lap_count_target = ?, time_limit_minutes = ?, time_start = ?, time_finish = ?, status = ?, 
	                 updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, r.RaceTitle, r.RaceType, r.ModelType, r.ModelScale, r.TrackName, r.LapCountTarget, 
	                 r.TimeLimitMinutes, r.TimeStart, r.TimeFinish, r.Status, r.UpdatedAt, r.ID)
	return err
}

// Delete removes a race from the database
func (r *Race) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM races WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all races from the database
func (r *Race) GetAll(db *sql.DB) ([]Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, time_start, time_finish, status, created_at, updated_at 
	          FROM races ORDER BY created_at DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []Race
	for rows.Next() {
		var race Race
		err := rows.Scan(
			&race.ID, &race.RaceTitle, &race.RaceType, &race.ModelType, &race.ModelScale, &race.TrackName, 
			&race.LapCountTarget, &race.TimeLimitMinutes, &race.TimeStart, &race.TimeFinish, &race.Status, 
			&race.CreatedAt, &race.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		races = append(races, race)
	}

	return races, nil
}

// Validate checks if the race data is valid
func (r *Race) Validate() error {
	if r.RaceTitle == "" {
		return fmt.Errorf("race title is required")
	}
	
	allowedTypes := map[string]bool{
		"qualifying": true,
		"main":       true,
		"final":      true,
	}
	
	if r.RaceType != "" && !allowedTypes[r.RaceType] {
		return fmt.Errorf("race type must be qualifying, main, or final")
	}
	
	allowedStatuses := map[string]bool{
		"scheduled":  true,
		"active":     true,
		"finished":   true,
		"cancelled":  true,
	}
	
	if r.Status != "" && !allowedStatuses[r.Status] {
		return fmt.Errorf("status must be scheduled, active, finished, or cancelled")
	}
	
	return nil
}