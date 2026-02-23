package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// RaceLap represents a record in the race_laps table
type RaceLap struct {
	ID                   string    `db:"id" json:"id"`
	RaceParticipantID    string    `db:"race_participant_id" json:"race_participant_id"`
	TimeStart            string    `db:"time_start" json:"time_start"`
	TimeFinish           *string   `db:"time_finish" json:"time_finish,omitempty"`
	NumberOfLaps         int       `db:"number_of_laps" json:"number_of_laps"`
	BestLapTimeMs        int       `db:"best_lap_time_ms" json:"best_lap_time_ms"`
	BestLapNumber        int       `db:"best_lap_number" json:"best_lap_number"`
	LastLapTimeMs        int       `db:"last_lap_time_ms" json:"last_lap_time_ms"`
	LastPassTime         *string   `db:"last_pass_time" json:"last_pass_time,omitempty"`
	TotalRaceTimeMs      int       `db:"total_race_time_ms" json:"total_race_time_ms"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (RaceLap) TableName() string {
	return "race_laps"
}

// Create creates a new record
func (rl *RaceLap) Create(db *sql.DB) error {
	rl.ID = uuid.New().String()
	rl.NumberOfLaps = 0
	rl.BestLapTimeMs = 0
	rl.BestLapNumber = 0
	rl.LastLapTimeMs = 0
	rl.TotalRaceTimeMs = 0
	rl.CreatedAt = time.Now()
	rl.UpdatedAt = time.Now()

	query := `INSERT INTO race_laps (id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rl.ID, rl.RaceParticipantID, rl.TimeStart, rl.TimeFinish, rl.NumberOfLaps, rl.BestLapTimeMs, rl.BestLapNumber, rl.LastLapTimeMs, rl.LastPassTime, rl.TotalRaceTimeMs, rl.CreatedAt, rl.UpdatedAt)
	return err
}

// GetByID gets a record by ID
func (rl *RaceLap) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at 
	          FROM race_laps WHERE id = ?`
	
	return db.QueryRow(query, id).Scan(&rl.ID, &rl.RaceParticipantID, &rl.TimeStart, &rl.TimeFinish, &rl.NumberOfLaps, &rl.BestLapTimeMs, &rl.BestLapNumber, &rl.LastLapTimeMs, &rl.LastPassTime, &rl.TotalRaceTimeMs, &rl.CreatedAt, &rl.UpdatedAt)
}

// Update updates a record
func (rl *RaceLap) Update(db *sql.DB) error {
	rl.UpdatedAt = time.Now()

	query := `UPDATE race_laps SET race_participant_id = ?, time_start = ?, time_finish = ?, number_of_laps = ?, best_lap_time_ms = ?, best_lap_number = ?, last_lap_time_ms = ?, last_pass_time = ?, total_race_time_ms = ?, updated_at = ? 
	          WHERE id = ?`
	
	_, err := db.Exec(query, rl.RaceParticipantID, rl.TimeStart, rl.TimeFinish, rl.NumberOfLaps, rl.BestLapTimeMs, rl.BestLapNumber, rl.LastLapTimeMs, rl.LastPassTime, rl.TotalRaceTimeMs, rl.UpdatedAt, rl.ID)
	return err
}

// Delete deletes a record
func (rl *RaceLap) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM race_laps WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (rl *RaceLap) GetAll(db *sql.DB) ([]RaceLap, error) {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, total_race_time_ms, created_at, updated_at 
	          FROM race_laps ORDER BY number_of_laps DESC, best_lap_time_ms ASC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var laps []RaceLap
	for rows.Next() {
		var lap RaceLap
		err := rows.Scan(&lap.ID, &lap.RaceParticipantID, &lap.TimeStart, &lap.TimeFinish, &lap.NumberOfLaps, &lap.BestLapTimeMs, &lap.BestLapNumber, &lap.LastLapTimeMs, &lap.LastPassTime, &lap.TotalRaceTimeMs, &lap.CreatedAt, &lap.UpdatedAt)
		if err != nil {
			return nil, err
		}
		laps = append(laps, lap)
	}

	return laps, nil
}