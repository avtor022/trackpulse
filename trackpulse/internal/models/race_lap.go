package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// RaceLap represents aggregated lap data for a participant
type RaceLap struct {
	ID                 string     `db:"id" json:"id"`
	RaceParticipantID  string     `db:"race_participant_id" json:"race_participant_id"`
	TimeStart          time.Time  `db:"time_start" json:"time_start"`
	TimeFinish         *time.Time `db:"time_finish" json:"time_finish,omitempty"`
	NumberOfLaps       int        `db:"number_of_laps" json:"number_of_laps"`
	BestLapTimeMS      int        `db:"best_lap_time_ms" json:"best_lap_time_ms"`
	BestLapNumber      int        `db:"best_lap_number" json:"best_lap_number"`
	LastLapTimeMS      int        `db:"last_lap_time_ms" json:"last_lap_time_ms"`
	LastPassTime       *time.Time `db:"last_pass_time" json:"last_pass_time,omitempty"`
	TotalRaceTimeMS    int        `db:"total_race_time_ms" json:"total_race_time_ms"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for race_laps
func (RaceLap) TableName() string {
	return "race_laps"
}

// Create inserts new race lap data into the database
func (rl *RaceLap) Create(db *sql.DB) error {
	if rl.ID == "" {
		rl.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	rl.CreatedAt = now
	rl.UpdatedAt = now

	query := `INSERT INTO race_laps (id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rl.ID, rl.RaceParticipantID, rl.TimeStart, rl.TimeFinish, rl.NumberOfLaps, rl.BestLapTimeMS, rl.BestLapNumber, rl.LastLapTimeMS, rl.LastPassTime, rl.TotalRaceTimeMS, rl.CreatedAt, rl.UpdatedAt)
	return err
}

// GetByID retrieves race lap data by ID
func (rl *RaceLap) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at 
	          FROM race_laps WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rl.ID, &rl.RaceParticipantID, &rl.TimeStart, &rl.TimeFinish, &rl.NumberOfLaps, &rl.BestLapTimeMS, &rl.BestLapNumber, &rl.LastLapTimeMS, &rl.LastPassTime, &rl.TotalRaceTimeMS, &rl.CreatedAt, &rl.UpdatedAt,
	)
	return err
}

// Update updates existing race lap data in the database
func (rl *RaceLap) Update(db *sql.DB) error {
	rl.UpdatedAt = time.Now().UTC()

	query := `UPDATE race_laps SET time_start = ?, time_finish = ?, number_of_laps = ?, best_lap_time_ms = ?, 
	                 best_lap_number = ?, last_lap_time_ms = ?, last_pass_time = ?, total_race_time_ms = ?, 
	                 updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, rl.TimeStart, rl.TimeFinish, rl.NumberOfLaps, rl.BestLapTimeMS, rl.BestLapNumber, 
	                 rl.LastLapTimeMS, rl.LastPassTime, rl.TotalRaceTimeMS, rl.UpdatedAt, rl.ID)
	return err
}

// Delete removes race lap data from the database
func (rl *RaceLap) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM race_laps WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all race lap data from the database
func (rl *RaceLap) GetAll(db *sql.DB) ([]RaceLap, error) {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at 
	          FROM race_laps ORDER BY number_of_laps DESC, best_lap_time_ms ASC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var raceLaps []RaceLap
	for rows.Next() {
		var raceLap RaceLap
		err := rows.Scan(
			&raceLap.ID, &raceLap.RaceParticipantID, &raceLap.TimeStart, &raceLap.TimeFinish, &raceLap.NumberOfLaps, 
			&raceLap.BestLapTimeMS, &raceLap.BestLapNumber, &raceLap.LastLapTimeMS, &raceLap.LastPassTime, 
			&raceLap.TotalRaceTimeMS, &raceLap.CreatedAt, &raceLap.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		raceLaps = append(raceLaps, raceLap)
	}

	return raceLaps, nil
}

// Validate checks if the race lap data is valid
func (rl *RaceLap) Validate() error {
	if rl.RaceParticipantID == "" {
		return fmt.Errorf("race participant ID is required")
	}
	if rl.NumberOfLaps < 0 {
		return fmt.Errorf("number of laps cannot be negative")
	}
	if rl.BestLapTimeMS < 0 {
		return fmt.Errorf("best lap time cannot be negative")
	}
	if rl.LastLapTimeMS < 0 {
		return fmt.Errorf("last lap time cannot be negative")
	}
	if rl.TotalRaceTimeMS < 0 {
		return fmt.Errorf("total race time cannot be negative")
	}
	return nil
}