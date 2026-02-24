package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RCModelsRepository interface {
	GetAll() ([]models.RCModel, error)
	GetByID(id string) (*models.RCModel, error)
	Create(model *models.RCModel) error
	Update(model *models.RCModel) error
	Delete(id string) error
}

type rcModelsRepo struct {
	db *sql.DB
}

func NewRCModelsRepository(db *sql.DB) RCModelsRepository {
	return &rcModelsRepo{db: db}
}

func (r *rcModelsRepo) GetAll() ([]models.RCModel, error) {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, 
						created_at, updated_at 
					FROM rc_models ORDER BY brand, model_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modelsList []models.RCModel
	for rows.Next() {
		var model models.RCModel
		var motorType, driveType *string

		err := rows.Scan(
			&model.ID,
			&model.Brand,
			&model.ModelName,
			&model.Scale,
			&model.ModelType,
			&motorType,
			&driveType,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		model.MotorType = motorType
		model.DriveType = driveType

		modelsList = append(modelsList, model)
	}

	return modelsList, nil
}

func (r *rcModelsRepo) GetByID(id string) (*models.RCModel, error) {
	query := `SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, 
						created_at, updated_at 
					FROM rc_models WHERE id = ?`

	var model models.RCModel
	var motorType, driveType *string

	err := r.db.QueryRow(query, id).Scan(
		&model.ID,
		&model.Brand,
		&model.ModelName,
		&model.Scale,
		&model.ModelType,
		&motorType,
		&driveType,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("model not found")
		}
		return nil, err
	}

	// Handle nullable fields
	model.MotorType = motorType
	model.DriveType = driveType

	return &model, nil
}

func (r *rcModelsRepo) Create(model *models.RCModel) error {
	query := `INSERT INTO rc_models (id, brand, model_name, scale, model_type, motor_type, drive_type, 
																		created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		model.ID,
		model.Brand,
		model.ModelName,
		model.Scale,
		model.ModelType,
		model.MotorType,
		model.DriveType,
		model.CreatedAt,
		model.UpdatedAt,
	)
	return err
}

func (r *rcModelsRepo) Update(model *models.RCModel) error {
	query := `UPDATE rc_models SET brand = ?, model_name = ?, scale = ?, model_type = ?, 
						motor_type = ?, drive_type = ?, updated_at = ?
					WHERE id = ?`

	_, err := r.db.Exec(query,
		model.Brand,
		model.ModelName,
		model.Scale,
		model.ModelType,
		model.MotorType,
		model.DriveType,
		model.UpdatedAt,
		model.ID,
	)
	return err
}

func (r *rcModelsRepo) Delete(id string) error {
	query := `DELETE FROM rc_models WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("model not found")
	}

	return nil
}