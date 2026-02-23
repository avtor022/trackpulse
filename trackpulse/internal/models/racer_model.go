package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// RacerModel represents the link between a racer and their model
type RacerModel struct {
	ID                string    `db:"id" json:"id"`
	RacerID           string    `db:"racer_id" json:"racer_id"`
	RCModelID         string    `db:"rc_model_id" json:"rc_model_id"`
	TransponderNumber string    `db:"transponder_number" json:"transponder_number"`
	TransponderType   string    `db:"transponder_type" json:"transponder_type"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for racer_models
func (RacerModel) TableName() string {
	return "racer_models"
}

// Create inserts a new racer-model link into the database
func (rm *RacerModel) Create(db *sql.DB) error {
	if rm.ID == "" {
		rm.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	rm.CreatedAt = now
	rm.UpdatedAt = now

	query := `INSERT INTO racer_models (id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rm.ID, rm.RacerID, rm.RCModelID, rm.TransponderNumber, rm.TransponderType, rm.IsActive, rm.CreatedAt, rm.UpdatedAt)
	return err
}

// GetByID retrieves a racer-model link by ID
func (rm *RacerModel) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at 
	          FROM racer_models WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rm.ID, &rm.RacerID, &rm.RCModelID, &rm.TransponderNumber, &rm.TransponderType, &rm.IsActive, &rm.CreatedAt, &rm.UpdatedAt,
	)
	return err
}

// Update updates an existing racer-model link in the database
func (rm *RacerModel) Update(db *sql.DB) error {
	rm.UpdatedAt = time.Now().UTC()

	query := `UPDATE racer_models SET racer_id = ?, rc_model_id = ?, transponder_number = ?, 
	                 transponder_type = ?, is_active = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, rm.RacerID, rm.RCModelID, rm.TransponderNumber, rm.TransponderType, rm.IsActive, rm.UpdatedAt, rm.ID)
	return err
}

// Delete removes a racer-model link from the database
func (rm *RacerModel) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM racer_models WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all racer-model links from the database
func (rm *RacerModel) GetAll(db *sql.DB) ([]RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at 
	          FROM racer_models ORDER BY transponder_number`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racerModels []RacerModel
	for rows.Next() {
		var racerModel RacerModel
		err := rows.Scan(
			&racerModel.ID, &racerModel.RacerID, &racerModel.RCModelID, &racerModel.TransponderNumber, 
			&racerModel.TransponderType, &racerModel.IsActive, &racerModel.CreatedAt, &racerModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		racerModels = append(racerModels, racerModel)
	}

	return racerModels, nil
}

// Validate checks if the racer-model link data is valid
func (rm *RacerModel) Validate() error {
	if rm.RacerID == "" {
		return fmt.Errorf("racer ID is required")
	}
	if rm.RCModelID == "" {
		return fmt.Errorf("RC model ID is required")
	}
	if rm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}
	return nil
}