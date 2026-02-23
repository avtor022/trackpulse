package models

import (
	"database/sql"
	"time"
)

// RCModel represents a record in the rc_models table
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

// TableName returns the table name
func (RCModel) TableName() string {
	return "rc_models"
}

// Create creates a new record
func (rm *RCModel) Create(db *sql.DB) error {
	query := `INSERT INTO rc_models (id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, rm.ID, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, now, now)
	if err != nil {
		return err
	}
	
	rm.CreatedAt = now
	rm.UpdatedAt = now
	
	return nil
}

// GetByID gets a record by ID
func (rm *RCModel) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at 
	          FROM rc_models WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rm.ID, &rm.Brand, &rm.ModelName, &rm.Scale, &rm.ModelType, &rm.MotorType, &rm.DriveType, 
		&rm.CreatedAt, &rm.UpdatedAt,
	)
	return err
}

// Update updates a record
func (rm *RCModel) Update(db *sql.DB) error {
	query := `UPDATE rc_models SET brand = ?, model_name = ?, scale = ?, model_type = ?, motor_type = ?, drive_type = ?, updated_at = ? 
	          WHERE id = ?`
	
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(query, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, now, rm.ID)
	if err != nil {
		return err
	}
	
	rm.UpdatedAt = now
	return nil
}

// Delete deletes a record
func (rm *RCModel) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM rc_models WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
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