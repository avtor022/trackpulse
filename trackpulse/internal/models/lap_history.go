package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// LapHistory represents individual lap records
type LapHistory struct {
	ID                  string    `db:"id" json:"id"`
	RaceParticipantID   string    `db:"race_participant_id" json:"race_participant_id"`
	LapNumber           int       `db:"lap_number" json:"lap_number"`
	LapTimeMS           int       `db:"lap_time_ms" json:"lap_time_ms"`
	StartTime           time.Time `db:"start_time" json:"start_time"`
	EndTime             time.Time `db:"end_time" json:"end_time"`
	IsValid             bool      `db:"is_valid" json:"is_valid"`
	InvalidationReason  *string   `db:"invalidation_reason" json:"invalidation_reason,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name for lap_history
func (LapHistory) TableName() string {
	return "lap_history"
}

// Create inserts a new lap history record into the database
func (lh *LapHistory) Create(db *sql.DB) error {
	if lh.ID == "" {
		lh.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	lh.CreatedAt = now

	query := `INSERT INTO lap_history (id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, lh.ID, lh.RaceParticipantID, lh.LapNumber, lh.LapTimeMS, lh.StartTime, lh.EndTime, lh.IsValid, lh.InvalidationReason, lh.CreatedAt)
	return err
}

// GetByID retrieves a lap history record by ID
func (lh *LapHistory) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at 
	          FROM lap_history WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&lh.ID, &lh.RaceParticipantID, &lh.LapNumber, &lh.LapTimeMS, &lh.StartTime, &lh.EndTime, &lh.IsValid, &lh.InvalidationReason, &lh.CreatedAt,
	)
	return err
}

// Update updates an existing lap history record in the database
func (lh *LapHistory) Update(db *sql.DB) error {
	query := `UPDATE lap_history SET lap_number = ?, lap_time_ms = ?, start_time = ?, end_time = ?, 
	                 is_valid = ?, invalidation_reason = ? WHERE id = ?`
	
	_, err := db.Exec(query, lh.LapNumber, lh.LapTimeMS, lh.StartTime, lh.EndTime, lh.IsValid, lh.InvalidationReason, lh.ID)
	return err
}

// Delete removes a lap history record from the database
func (lh *LapHistory) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM lap_history WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all lap history records from the database
func (lh *LapHistory) GetAll(db *sql.DB) ([]LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at 
	          FROM lap_history ORDER BY lap_number ASC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lapHistories []LapHistory
	for rows.Next() {
		var lapHistory LapHistory
		err := rows.Scan(
			&lapHistory.ID, &lapHistory.RaceParticipantID, &lapHistory.LapNumber, &lapHistory.LapTimeMS, 
			&lapHistory.StartTime, &lapHistory.EndTime, &lapHistory.IsValid, &lapHistory.InvalidationReason, 
			&lapHistory.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		lapHistories = append(lapHistories, lapHistory)
	}

	return lapHistories, nil
}

// Validate checks if the lap history data is valid
func (lh *LapHistory) Validate() error {
	if lh.RaceParticipantID == "" {
		return fmt.Errorf("race participant ID is required")
	}
	if lh.LapNumber <= 0 {
		return fmt.Errorf("lap number must be positive")
	}
	if lh.LapTimeMS < 0 {
		return fmt.Errorf("lap time cannot be negative")
	}
	if lh.StartTime.After(lh.EndTime) {
		return fmt.Errorf("start time cannot be after end time")
	}
	return nil
}