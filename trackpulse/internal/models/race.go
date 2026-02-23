package models

import (
	"database/sql"
	"time"
)

// Race represents a record in the races table
type Race struct {
	ID               string     `db:"id" json:"id"`
	RaceTitle        string     `db:"race_title" json:"race_title"`
	RaceType         string     `db:"race_type" json:"race_type"`
	ModelType        *string    `db:"model_type" json:"model_type,omitempty"`
	ModelScale       *string    `db:"model_scale" json:"model_scale,omitempty"`
	TrackName        *string    `db:"track_name" json:"track_name,omitempty"`
	LapCountTarget   *int       `db:"lap_count_target" json:"lap_count_target,omitempty"`
	TimeLimitMinutes *int       `db:"time_limit_minutes" json:"time_limit_minutes,omitempty"`
	TimeStart        *string    `db:"time_start" json:"time_start,omitempty"`
	TimeFinish       *string    `db:"time_finish" json:"time_finish,omitempty"`
	Status           string     `db:"status" json:"status"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (Race) TableName() string {
	return "races"
}

// Create creates a new record
func (r *Race) Create(db *sql.DB) error {
	query := `INSERT INTO races (id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, time_start, time_finish, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, r.ID, r.RaceTitle, r.RaceType, r.ModelType, r.ModelScale, r.TrackName, 
	                  r.LapCountTarget, r.TimeLimitMinutes, r.TimeStart, r.TimeFinish, r.Status, now, now)
	if err != nil {
		return err
	}
	
	r.CreatedAt = now
	r.UpdatedAt = now
	
	return nil
}

// GetByID gets a record by ID
func (r *Race) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, 
	                 time_start, time_finish, status, created_at, updated_at 
	          FROM races WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&r.ID, &r.RaceTitle, &r.RaceType, &r.ModelType, &r.ModelScale, &r.TrackName, &r.LapCountTarget, 
		&r.TimeLimitMinutes, &r.TimeStart, &r.TimeFinish, &r.Status, &r.CreatedAt, &r.UpdatedAt,
	)
	return err
}

// Update updates a record
func (r *Race) Update(db *sql.DB) error {
	query := `UPDATE races SET race_title = ?, race_type = ?, model_type = ?, model_scale = ?, track_name = ?, 
	                         lap_count_target = ?, time_limit_minutes = ?, time_start = ?, time_finish = ?, 
	                         status = ?, updated_at = ? 
	          WHERE id = ?`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, r.RaceTitle, r.RaceType, r.ModelType, r.ModelScale, r.TrackName, 
	                  r.LapCountTarget, r.TimeLimitMinutes, r.TimeStart, r.TimeFinish, r.Status, now, r.ID)
	if err != nil {
		return err
	}
	
	r.UpdatedAt = now
	return nil
}

// Delete deletes a record
func (r *Race) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM races WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (r *Race) GetAll(db *sql.DB) ([]Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, time_limit_minutes, 
	                 time_start, time_finish, status, created_at, updated_at 
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