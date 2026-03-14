package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// RCModelRepository handles data access for RC models
type RCModelRepository struct {
	db *sql.DB
}

// NewRCModelRepository creates a new RC model repository
func NewRCModelRepository(db *sql.DB) *RCModelRepository {
	return &RCModelRepository{db: db}
}

// GetAll returns all RC models
func (r *RCModelRepository) GetAll() ([]models.RCModel, error) {
	rows, err := r.db.Query(`
		SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at
		FROM rc_models
		ORDER BY brand ASC, model_name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rc_models: %w", err)
	}
	defer rows.Close()

	var rcModels []models.RCModel
	for rows.Next() {
		var model models.RCModel
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&model.ID,
			&model.Brand,
			&model.ModelName,
			&model.Scale,
			&model.ModelType,
			&model.MotorType,
			&model.DriveType,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			model.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			model.UpdatedAt = t
		}

		rcModels = append(rcModels, model)
	}

	return rcModels, rows.Err()
}

// GetByID returns an RC model by ID
func (r *RCModelRepository) GetByID(id string) (*models.RCModel, error) {
	row := r.db.QueryRow(`
		SELECT id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at
		FROM rc_models
		WHERE id = ?
	`, id)

	var model models.RCModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&model.ID,
		&model.Brand,
		&model.ModelName,
		&model.Scale,
		&model.ModelType,
		&model.MotorType,
		&model.DriveType,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		model.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		model.UpdatedAt = t
	}

	return &model, nil
}

// Create inserts a new RC model
func (r *RCModelRepository) Create(model *models.RCModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		INSERT INTO rc_models (id, brand, model_name, scale, model_type, motor_type, drive_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		model.ID,
		model.Brand,
		model.ModelName,
		model.Scale,
		model.ModelType,
		model.MotorType,
		model.DriveType,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating model")
	}

	return nil
}

// Update updates an existing RC model
func (r *RCModelRepository) Update(model *models.RCModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		UPDATE rc_models
		SET brand = ?, model_name = ?, scale = ?, model_type = ?, motor_type = ?, drive_type = ?, updated_at = ?
		WHERE id = ?
	`,
		model.Brand,
		model.ModelName,
		model.Scale,
		model.ModelType,
		model.MotorType,
		model.DriveType,
		now,
		model.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("model not found")
	}

	return nil
}

// Delete removes an RC model by ID
func (r *RCModelRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM rc_models WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("model not found")
	}

	return nil
}

// Count returns total number of RC models
func (r *RCModelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM rc_models`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count models: %w", err)
	}
	return count, nil
}

// GetUniqueBrands returns a list of unique brand names from all models
func (r *RCModelRepository) GetUniqueBrands() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT brand
		FROM rc_models
		ORDER BY brand ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query unique brands: %w", err)
	}
	defer rows.Close()

	var brands []string
	for rows.Next() {
		var brand string
		err := rows.Scan(&brand)
		if err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, rows.Err()
}

// GetUniqueModelNames returns a list of unique model names from all models
func (r *RCModelRepository) GetUniqueModelNames() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT model_name
		FROM rc_models
		ORDER BY model_name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query unique model names: %w", err)
	}
	defer rows.Close()

	var modelNames []string
	for rows.Next() {
		var modelName string
		err := rows.Scan(&modelName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model name: %w", err)
		}
		modelNames = append(modelNames, modelName)
	}

	return modelNames, rows.Err()
}

// GetAllModelNames returns all unique model names
func (r *RCModelRepository) GetAllModelNames() ([]string, error) {
	return r.GetUniqueModelNames()
}
