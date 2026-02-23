package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// RaceParticipant represents a participant in a race
type RaceParticipant struct {
	ID             string    `db:"id" json:"id"`
	RaceID         string    `db:"race_id" json:"race_id"`
	RacerModelID   string    `db:"racer_model_id" json:"racer_model_id"`
	GridPosition   *int      `db:"grid_position" json:"grid_position,omitempty"`
	IsFinished     bool      `db:"is_finished" json:"is_finished"`
	Disqualified   bool      `db:"disqualified" json:"disqualified"`
	DNFReason      *string   `db:"dnf_reason" json:"dnf_reason,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for race_participants
func (RaceParticipant) TableName() string {
	return "race_participants"
}

// Create inserts a new race participant into the database
func (rp *RaceParticipant) Create(db *sql.DB) error {
	if rp.ID == "" {
		rp.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	rp.CreatedAt = now
	rp.UpdatedAt = now

	query := `INSERT INTO race_participants (id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rp.ID, rp.RaceID, rp.RacerModelID, rp.GridPosition, rp.IsFinished, rp.Disqualified, rp.DNFReason, rp.CreatedAt, rp.UpdatedAt)
	return err
}

// GetByID retrieves a race participant by ID
func (rp *RaceParticipant) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at 
	          FROM race_participants WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rp.ID, &rp.RaceID, &rp.RacerModelID, &rp.GridPosition, &rp.IsFinished, &rp.Disqualified, &rp.DNFReason, &rp.CreatedAt, &rp.UpdatedAt,
	)
	return err
}

// Update updates an existing race participant in the database
func (rp *RaceParticipant) Update(db *sql.DB) error {
	rp.UpdatedAt = time.Now().UTC()

	query := `UPDATE race_participants SET race_id = ?, racer_model_id = ?, grid_position = ?, is_finished = ?, 
	                 disqualified = ?, dnf_reason = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, rp.RaceID, rp.RacerModelID, rp.GridPosition, rp.IsFinished, rp.Disqualified, rp.DNFReason, rp.UpdatedAt, rp.ID)
	return err
}

// Delete removes a race participant from the database
func (rp *RaceParticipant) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM race_participants WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all race participants from the database
func (rp *RaceParticipant) GetAll(db *sql.DB) ([]RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at 
	          FROM race_participants ORDER BY race_id, grid_position`
	
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

// Validate checks if the race participant data is valid
func (rp *RaceParticipant) Validate() error {
	if rp.RaceID == "" {
		return fmt.Errorf("race ID is required")
	}
	if rp.RacerModelID == "" {
		return fmt.Errorf("racer model ID is required")
	}
	return nil
}