package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// RacerModelRepository handles data access for racer models (transponders)
type RacerModelRepository struct {
	db *sql.DB
}

// NewRacerModelRepository creates a new racer model repository
func NewRacerModelRepository(db *sql.DB) *RacerModelRepository {
	return &RacerModelRepository{db: db}
}

// GetAll returns all racer models
func (r *RacerModelRepository) GetAll() ([]models.RacerModel, error) {
	rows, err := r.db.Query(`
		SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM racer_models
		ORDER BY transponder_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query racer_models: %w", err)
	}
	defer rows.Close()

	var racerModels []models.RacerModel
	for rows.Next() {
		var rm models.RacerModel
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&rm.ID,
			&rm.RacerID,
			&rm.RCModelID,
			&rm.TransponderNumber,
			&rm.TransponderType,
			&rm.IsActive,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan racer model: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			rm.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			rm.UpdatedAt = t
		}

		racerModels = append(racerModels, rm)
	}

	return racerModels, rows.Err()
}

// GetByID returns a racer model by ID
func (r *RacerModelRepository) GetByID(id string) (*models.RacerModel, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM racer_models
		WHERE id = ?
	`, id)

	var rm models.RacerModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&rm.ID,
		&rm.RacerID,
		&rm.RCModelID,
		&rm.TransponderNumber,
		&rm.TransponderType,
		&rm.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get racer model: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		rm.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		rm.UpdatedAt = t
	}

	return &rm, nil
}

// GetByTransponderNumber returns a racer model by transponder number
func (r *RacerModelRepository) GetByTransponderNumber(number string) (*models.RacerModel, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM racer_models
		WHERE transponder_number = ?
	`, number)

	var rm models.RacerModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&rm.ID,
		&rm.RacerID,
		&rm.RCModelID,
		&rm.TransponderNumber,
		&rm.TransponderType,
		&rm.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get racer model by transponder number: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		rm.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		rm.UpdatedAt = t
	}

	return &rm, nil
}

// Create inserts a new racer model
func (r *RacerModelRepository) Create(rm *models.RacerModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		INSERT INTO racer_models (id, racer_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		rm.ID,
		rm.RacerID,
		rm.RCModelID,
		rm.TransponderNumber,
		rm.TransponderType,
		rm.IsActive,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create racer model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating racer model")
	}

	return nil
}

// Update updates an existing racer model
func (r *RacerModelRepository) Update(rm *models.RacerModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		UPDATE racer_models
		SET racer_id = ?, rc_model_id = ?, transponder_number = ?, transponder_type = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`,
		rm.RacerID,
		rm.RCModelID,
		rm.TransponderNumber,
		rm.TransponderType,
		rm.IsActive,
		now,
		rm.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update racer model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("racer model not found")
	}

	return nil
}

// Delete removes a racer model by ID
func (r *RacerModelRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM racer_models WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete racer model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("racer model not found")
	}

	return nil
}

// Count returns total number of racer models
func (r *RacerModelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM racer_models`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count racer models: %w", err)
	}
	return count, nil
}
