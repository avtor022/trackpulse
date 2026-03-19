package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RCModelTypeRepository handles data access for RC model types
type RCModelTypeRepository struct {
	db *sql.DB
}

// NewRCModelTypeRepository creates a new RC model type repository
func NewRCModelTypeRepository(db *sql.DB) *RCModelTypeRepository {
	return &RCModelTypeRepository{db: db}
}

// GetAll returns all RC model types
func (r *RCModelTypeRepository) GetAll() ([]models.RCModelType, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_types
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rc_model_types: %w", err)
	}
	defer rows.Close()

	var types []models.RCModelType
	for rows.Next() {
		var t models.RCModelType
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model type: %w", err)
		}

		createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err == nil {
			t.CreatedAt = createdAtTime
		}
		updatedAtTime, err := time.Parse(time.RFC3339, updatedAtStr)
		if err == nil {
			t.UpdatedAt = updatedAtTime
		}

		types = append(types, t)
	}

	return types, rows.Err()
}

// GetByName returns a model type by name
func (r *RCModelTypeRepository) GetByName(name string) (*models.RCModelType, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_types
		WHERE name = ?
	`, name)

	var t models.RCModelType
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&t.ID,
		&t.Name,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model type: %w", err)
	}

	createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
	if err == nil {
		t.CreatedAt = createdAtTime
	}
	updatedAtTime, err := time.Parse(time.RFC3339, updatedAtStr)
	if err == nil {
		t.UpdatedAt = updatedAtTime
	}

	return &t, nil
}

// Create inserts a new model type
func (r *RCModelTypeRepository) Create(name string) (*models.RCModelType, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO rc_model_types (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create model type: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	t := &models.RCModelType{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return t, nil
}

// GetOrCreate gets a model type by name or creates it if it doesn't exist
func (r *RCModelTypeRepository) GetOrCreate(name string) (*models.RCModelType, error) {
	// Try to get existing model type
	t, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// Create new model type
	return r.Create(name)
}

// Delete deletes a model type by name
func (r *RCModelTypeRepository) Delete(name string) error {
	_, err := r.db.Exec(`DELETE FROM rc_model_types WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete model type: %w", err)
	}
	return nil
}

