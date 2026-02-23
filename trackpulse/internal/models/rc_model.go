package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// RCModel represents an RC model
type RCModel struct {
	ID         string    `db:"id" json:"id"`
	Brand      string    `db:"brand" json:"brand"`
	ModelName  string    `db:"model_name" json:"model_name"`
	Scale      string    `db:"scale" json:"scale"`
	ModelType  string    `db:"model_type" json:"model_type"`
	MotorType  *string   `db:"motor_type" json:"motor_type,omitempty"`
	DriveType  *string   `db:"drive_type" json:"drive_type,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for rc_models
func (RCModel) TableName() string {
	return "rc_models"
}

// Create inserts a new RC model into the database
func (rm *RCModel) Create(db *sql.DB) error {
	if rm.ID == "" {
		rm.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	rm.CreatedAt = now
	rm.UpdatedAt = now

	query := `INSERT INTO rc_models (id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rm.ID, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, rm.CreatedAt, rm.UpdatedAt)
	return err
}

// GetByID retrieves an RC model by ID
func (rm *RCModel) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at 
	          FROM rc_models WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rm.ID, &rm.Brand, &rm.ModelName, &rm.Scale, &rm.ModelType, &rm.MotorType, &rm.DriveType, &rm.CreatedAt, &rm.UpdatedAt,
	)
	return err
}

// Update updates an existing RC model in the database
func (rm *RCModel) Update(db *sql.DB) error {
	rm.UpdatedAt = time.Now().UTC()

	query := `UPDATE rc_models SET brand = ?, model_name = ?, scale = ?, model_type = ?, 
	                 motor_type = ?, drive_type = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, rm.UpdatedAt, rm.ID)
	return err
}

// Delete removes an RC model from the database
func (rm *RCModel) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM rc_models WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all RC models from the database
func (rm *RCModel) GetAll(db *sql.DB) ([]RCModel, error) {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at 
	          FROM rc_models ORDER BY brand, model_name`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []RCModel
	for rows.Next() {
		var model RCModel
		err := rows.Scan(
			&model.ID, &model.Brand, &model.ModelName, &model.Scale, &model.ModelType, &model.MotorType, &model.DriveType, 
			&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

// Validate checks if the RC model data is valid
func (rm *RCModel) Validate() error {
	if rm.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if rm.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	if rm.Scale == "" {
		return fmt.Errorf("scale is required")
	}
	if rm.ModelType == "" {
		return fmt.Errorf("model type is required")
	}
	return nil
}