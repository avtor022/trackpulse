package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// LapHistory represents a record in the lap_history table
type LapHistory struct {
	ID                  string    `db:"id" json:"id"`
	RaceParticipantID   string    `db:"race_participant_id" json:"race_participant_id"`
	LapNumber           int       `db:"lap_number" json:"lap_number"`
	LapTimeMs           int       `db:"lap_time_ms" json:"lap_time_ms"`
	StartTime           string    `db:"start_time" json:"start_time"`
	EndTime             string    `db:"end_time" json:"end_time"`
	IsValid             bool      `db:"is_valid" json:"is_valid"`
	InvalidationReason  *string   `db:"invalidation_reason" json:"invalidation_reason,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name
func (LapHistory) TableName() string {
	return "lap_history"
}

// Create creates a new record
func (lh *LapHistory) Create(db *sql.DB) error {
	lh.ID = uuid.New().String()
	lh.IsValid = true
	lh.CreatedAt = time.Now()

	query := `INSERT INTO lap_history (id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, lh.ID, lh.RaceParticipantID, lh.LapNumber, lh.LapTimeMs, lh.StartTime, lh.EndTime, lh.IsValid, lh.InvalidationReason, lh.CreatedAt)
	return err
}

// GetByID gets a record by ID
func (lh *LapHistory) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at 
	          FROM lap_history WHERE id = ?`
	
	return db.QueryRow(query, id).Scan(&lh.ID, &lh.RaceParticipantID, &lh.LapNumber, &lh.LapTimeMs, &lh.StartTime, &lh.EndTime, &lh.IsValid, &lh.InvalidationReason, &lh.CreatedAt)
}

// Update updates a record
func (lh *LapHistory) Update(db *sql.DB) error {
	query := `UPDATE lap_history SET race_participant_id = ?, lap_number = ?, lap_time_ms = ?, start_time = ?, end_time = ?, is_valid = ?, invalidation_reason = ? 
	          WHERE id = ?`
	
	_, err := db.Exec(query, lh.RaceParticipantID, lh.LapNumber, lh.LapTimeMs, lh.StartTime, lh.EndTime, lh.IsValid, lh.InvalidationReason, lh.ID)
	return err
}

// Delete deletes a record
func (lh *LapHistory) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM lap_history WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (lh *LapHistory) GetAll(db *sql.DB) ([]LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, end_time, is_valid, invalidation_reason, created_at 
	          FROM lap_history ORDER BY lap_number, start_time`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []LapHistory
	for rows.Next() {
		var history LapHistory
		err := rows.Scan(&history.ID, &history.RaceParticipantID, &history.LapNumber, &history.LapTimeMs, &history.StartTime, &history.EndTime, &history.IsValid, &history.InvalidationReason, &history.CreatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}