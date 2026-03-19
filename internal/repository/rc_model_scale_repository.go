package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RCModelScaleRepository handles data access for RC model scales
type RCModelScaleRepository struct {
	db *sql.DB
}

// NewRCModelScaleRepository creates a new RC model scale repository
func NewRCModelScaleRepository(db *sql.DB) *RCModelScaleRepository {
	return &RCModelScaleRepository{db: db}
}

// GetAll returns all RC model scales
func (r *RCModelScaleRepository) GetAll() ([]models.RCModelScale, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_scales
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rc_model_scales: %w", err)
	}
	defer rows.Close()

	var scales []models.RCModelScale
	for rows.Next() {
		var scale models.RCModelScale
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&scale.ID,
			&scale.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scale: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			scale.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			scale.UpdatedAt = t
		}

		scales = append(scales, scale)
	}

	return scales, rows.Err()
}

// GetByName returns a scale by name
func (r *RCModelScaleRepository) GetByName(name string) (*models.RCModelScale, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_scales
		WHERE name = ?
	`, name)

	var scale models.RCModelScale
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&scale.ID,
		&scale.Name,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scale: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		scale.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		scale.UpdatedAt = t
	}

	return &scale, nil
}

// Create inserts a new scale
func (r *RCModelScaleRepository) Create(name string) (*models.RCModelScale, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO rc_model_scales (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create scale: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	scale := &models.RCModelScale{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return scale, nil
}

// GetOrCreate gets a scale by name or creates it if it doesn't exist
func (r *RCModelScaleRepository) GetOrCreate(name string) (*models.RCModelScale, error) {
	// Try to get existing scale
	scale, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if scale != nil {
		return scale, nil
	}

	// Create new scale
	return r.Create(name)
}

// Delete deletes a scale by name
func (r *RCModelScaleRepository) Delete(name string) error {
	_, err := r.db.Exec(`DELETE FROM rc_model_scales WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete scale: %w", err)
	}
	return nil
}

