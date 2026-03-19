package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// CompetitorModelRepository handles data access for competitor models (transponders)
type CompetitorModelRepository struct {
	db *sql.DB
}

// NewCompetitorModelRepository creates a new competitor model repository
func NewCompetitorModelRepository(db *sql.DB) *CompetitorModelRepository {
	return &CompetitorModelRepository{db: db}
}

// GetAll returns all competitor models
func (r *CompetitorModelRepository) GetAll() ([]models.CompetitorModel, error) {
	rows, err := r.db.Query(`
		SELECT id, competitor_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM competitor_models
		ORDER BY transponder_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competitor_models: %w", err)
	}
	defer rows.Close()

	var competitorModels []models.CompetitorModel
	for rows.Next() {
		var cm models.CompetitorModel
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&cm.ID,
			&cm.CompetitorID,
			&cm.RCModelID,
			&cm.TransponderNumber,
			&cm.TransponderType,
			&cm.IsActive,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competitor model: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			cm.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			cm.UpdatedAt = t
		}

		competitorModels = append(competitorModels, cm)
	}

	return competitorModels, rows.Err()
}

// GetByID returns a competitor model by ID
func (r *CompetitorModelRepository) GetByID(id string) (*models.CompetitorModel, error) {
	row := r.db.QueryRow(`
		SELECT id, competitor_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM competitor_models
		WHERE id = ?
	`, id)

	var cm models.CompetitorModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&cm.ID,
		&cm.CompetitorID,
		&cm.RCModelID,
		&cm.TransponderNumber,
		&cm.TransponderType,
		&cm.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor model: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		cm.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		cm.UpdatedAt = t
	}

	return &cm, nil
}

// GetByTransponderNumber returns a competitor model by transponder number
func (r *CompetitorModelRepository) GetByTransponderNumber(number string) (*models.CompetitorModel, error) {
	row := r.db.QueryRow(`
		SELECT id, competitor_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM competitor_models
		WHERE transponder_number = ?
	`, number)

	var cm models.CompetitorModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&cm.ID,
		&cm.CompetitorID,
		&cm.RCModelID,
		&cm.TransponderNumber,
		&cm.TransponderType,
		&cm.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor model by transponder number: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		cm.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		cm.UpdatedAt = t
	}

	return &cm, nil
}

// Create inserts a new competitor model
func (r *CompetitorModelRepository) Create(cm *models.CompetitorModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		INSERT INTO competitor_models (id, competitor_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		cm.ID,
		cm.CompetitorID,
		cm.RCModelID,
		cm.TransponderNumber,
		cm.TransponderType,
		cm.IsActive,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create competitor model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating competitor model")
	}

	return nil
}

// Update updates an existing competitor model
func (r *CompetitorModelRepository) Update(cm *models.CompetitorModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		UPDATE competitor_models
		SET competitor_id = ?, rc_model_id = ?, transponder_number = ?, transponder_type = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`,
		cm.CompetitorID,
		cm.RCModelID,
		cm.TransponderNumber,
		cm.TransponderType,
		cm.IsActive,
		now,
		cm.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update competitor model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competitor model not found")
	}

	return nil
}

// Delete removes a competitor model by ID
func (r *CompetitorModelRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM competitor_models WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete competitor model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competitor model not found")
	}

	return nil
}

// Count returns total number of competitor models
func (r *CompetitorModelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM competitor_models`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count competitor models: %w", err)
	}
	return count, nil
}
