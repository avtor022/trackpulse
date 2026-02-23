package models

import (
	"database/sql"
	"time"
)

// RaceParticipant represents a record in the race_participants table
type RaceParticipant struct {
	ID              string    `db:"id" json:"id"`
	RaceID          string    `db:"race_id" json:"race_id"`
	RacerModelID    string    `db:"racer_model_id" json:"racer_model_id"`
	GridPosition    *int      `db:"grid_position" json:"grid_position,omitempty"`
	IsFinished      bool      `db:"is_finished" json:"is_finished"`
	Disqualified    bool      `db:"disqualified" json:"disqualified"`
	DNFReason       *string   `db:"dnf_reason" json:"dnf_reason,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (RaceParticipant) TableName() string {
	return "race_participants"
}

// Create creates a new record
func (rp *RaceParticipant) Create(db *sql.DB) error {
	query := `INSERT INTO race_participants (id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, rp.ID, rp.RaceID, rp.RacerModelID, rp.GridPosition, rp.IsFinished, 
	                  rp.Disqualified, rp.DNFReason, now, now)
	if err != nil {
		return err
	}
	
	rp.CreatedAt = now
	rp.UpdatedAt = now
	
	return nil
}

// GetByID gets a record by ID
func (rp *RaceParticipant) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at 
	          FROM race_participants WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rp.ID, &rp.RaceID, &rp.RacerModelID, &rp.GridPosition, &rp.IsFinished, &rp.Disqualified, 
		&rp.DNFReason, &rp.CreatedAt, &rp.UpdatedAt,
	)
	return err
}

// Update updates a record
func (rp *RaceParticipant) Update(db *sql.DB) error {
	query := `UPDATE race_participants SET race_id = ?, racer_model_id = ?, grid_position = ?, is_finished = ?, 
	                                   disqualified = ?, dnf_reason = ?, updated_at = ? 
	          WHERE id = ?`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, rp.RaceID, rp.RacerModelID, rp.GridPosition, rp.IsFinished, 
	                  rp.Disqualified, rp.DNFReason, now, rp.ID)
	if err != nil {
		return err
	}
	
	rp.UpdatedAt = now
	return nil
}

// Delete deletes a record
func (rp *RaceParticipant) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM race_participants WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (rp *RaceParticipant) GetAll(db *sql.DB) ([]RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at 
	          FROM race_participants ORDER BY grid_position`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var participants []RaceParticipant
	for rows.Next() {
		var participant RaceParticipant
		err := rows.Scan(
			&participant.ID, &participant.RaceID, &participant.RacerModelID, &participant.GridPosition, 
			&participant.IsFinished, &participant.Disqualified, &participant.DNFReason, 
			&participant.CreatedAt, &participant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}
	
	return participants, nil
}