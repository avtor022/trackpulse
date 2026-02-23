package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// RcModel represents a record in the rc_models table
type RcModel struct {
	ID          string    `db:"id" json:"id"`
	Brand       string    `db:"brand" json:"brand"`
	ModelName   string    `db:"model_name" json:"model_name"`
	Scale       string    `db:"scale" json:"scale"`
	ModelType   string    `db:"model_type" json:"model_type"`
	MotorType   *string   `db:"motor_type" json:"motor_type,omitempty"`
	DriveType   *string   `db:"drive_type" json:"drive_type,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name
func (RcModel) TableName() string {
	return "rc_models"
}

// Create creates a new record
func (rm *RcModel) Create(db *sql.DB) error {
	rm.ID = uuid.New().String()
	rm.CreatedAt = time.Now()
	rm.UpdatedAt = time.Now()

	query := `INSERT INTO rc_models (id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rm.ID, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, rm.CreatedAt, rm.UpdatedAt)
	return err
}

// GetByID gets a record by ID
func (rm *RcModel) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at 
	          FROM rc_models WHERE id = ?`
	
	return db.QueryRow(query, id).Scan(&rm.ID, &rm.Brand, &rm.ModelName, &rm.Scale, &rm.ModelType, &rm.MotorType, &rm.DriveType, &rm.CreatedAt, &rm.UpdatedAt)
}

// Update updates a record
func (rm *RcModel) Update(db *sql.DB) error {
	rm.UpdatedAt = time.Now()

	query := `UPDATE rc_models SET brand = ?, model_name = ?, scale = ?, model_type = ?, motor_type = ?, drive_type = ?, updated_at = ? 
	          WHERE id = ?`
	
	_, err := db.Exec(query, rm.Brand, rm.ModelName, rm.Scale, rm.ModelType, rm.MotorType, rm.DriveType, rm.UpdatedAt, rm.ID)
	return err
}

// Delete deletes a record
func (rm *RcModel) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM rc_models WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (rm *RcModel) GetAll(db *sql.DB) ([]RcModel, error) {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at 
	          FROM rc_models ORDER BY brand, model_name`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []RcModel
	for rows.Next() {
		var model RcModel
		err := rows.Scan(&model.ID, &model.Brand, &model.ModelName, &model.Scale, &model.ModelType, &model.MotorType, &model.DriveType, &model.CreatedAt, &model.UpdatedAt)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}