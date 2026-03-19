package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// AthleteModelRepository handles data access for athlete models (transponders)
type AthleteModelRepository struct {
	db *sql.DB
}

// NewAthleteModelRepository creates a new athlete model repository
func NewAthleteModelRepository(db *sql.DB) *AthleteModelRepository {
	return &AthleteModelRepository{db: db}
}

// GetAll returns all athlete models
func (r *AthleteModelRepository) GetAll() ([]models.AthleteModel, error) {
	rows, err := r.db.Query(`
		SELECT id, athlete_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM athlete_models
		ORDER BY transponder_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query athlete_models: %w", err)
	}
	defer rows.Close()

	var athleteModels []models.AthleteModel
	for rows.Next() {
		var am models.AthleteModel
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&am.ID,
			&am.AthleteID,
			&am.RCModelID,
			&am.TransponderNumber,
			&am.TransponderType,
			&am.IsActive,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan athlete model: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			am.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			am.UpdatedAt = t
		}

		athleteModels = append(athleteModels, am)
	}

	return athleteModels, rows.Err()
}

// GetByID returns an athlete model by ID
func (r *AthleteModelRepository) GetByID(id string) (*models.AthleteModel, error) {
	row := r.db.QueryRow(`
		SELECT id, athlete_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM athlete_models
		WHERE id = ?
	`, id)

	var am models.AthleteModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&am.ID,
		&am.AthleteID,
		&am.RCModelID,
		&am.TransponderNumber,
		&am.TransponderType,
		&am.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete model: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		am.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		am.UpdatedAt = t
	}

	return &am, nil
}

// GetByTransponderNumber returns an athlete model by transponder number
func (r *AthleteModelRepository) GetByTransponderNumber(number string) (*models.AthleteModel, error) {
	row := r.db.QueryRow(`
		SELECT id, athlete_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at
		FROM athlete_models
		WHERE transponder_number = ?
	`, number)

	var am models.AthleteModel
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&am.ID,
		&am.AthleteID,
		&am.RCModelID,
		&am.TransponderNumber,
		&am.TransponderType,
		&am.IsActive,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete model by transponder number: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		am.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		am.UpdatedAt = t
	}

	return &am, nil
}

// Create inserts a new athlete model
func (r *AthleteModelRepository) Create(am *models.AthleteModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		INSERT INTO athlete_models (id, athlete_id, rc_model_id, transponder_number, transponder_type, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		am.ID,
		am.AthleteID,
		am.RCModelID,
		am.TransponderNumber,
		am.TransponderType,
		am.IsActive,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create athlete model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating athlete model")
	}

	return nil
}

// Update updates an existing athlete model
func (r *AthleteModelRepository) Update(am *models.AthleteModel) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		UPDATE athlete_models
		SET athlete_id = ?, rc_model_id = ?, transponder_number = ?, transponder_type = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`,
		am.AthleteID,
		am.RCModelID,
		am.TransponderNumber,
		am.TransponderType,
		am.IsActive,
		now,
		am.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update athlete model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete model not found")
	}

	return nil
}

// Delete removes an athlete model by ID
func (r *AthleteModelRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM athlete_models WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete athlete model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete model not found")
	}

	return nil
}

// Count returns total number of athlete models
func (r *AthleteModelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM athlete_models`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count athlete models: %w", err)
	}
	return count, nil
}
