package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RCModelBrandRepository handles data access for RC model brands
type RCModelBrandRepository struct {
	db *sql.DB
}

// NewRCModelBrandRepository creates a new RC model brand repository
func NewRCModelBrandRepository(db *sql.DB) *RCModelBrandRepository {
	return &RCModelBrandRepository{db: db}
}

// GetAll returns all RC model brands
func (r *RCModelBrandRepository) GetAll() ([]models.RCModelBrand, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_brands
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rc_model_brands: %w", err)
	}
	defer rows.Close()

	var brands []models.RCModelBrand
	for rows.Next() {
		var brand models.RCModelBrand
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&brand.ID,
			&brand.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			brand.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			brand.UpdatedAt = t
		}

		brands = append(brands, brand)
	}

	return brands, rows.Err()
}

// GetByName returns a brand by name
func (r *RCModelBrandRepository) GetByName(name string) (*models.RCModelBrand, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_brands
		WHERE name = ?
	`, name)

	var brand models.RCModelBrand
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&brand.ID,
		&brand.Name,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		brand.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		brand.UpdatedAt = t
	}

	return &brand, nil
}

// Create inserts a new brand
func (r *RCModelBrandRepository) Create(name string) (*models.RCModelBrand, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO rc_model_brands (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create brand: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	brand := &models.RCModelBrand{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return brand, nil
}

// GetOrCreate gets a brand by name or creates it if it doesn't exist
func (r *RCModelBrandRepository) GetOrCreate(name string) (*models.RCModelBrand, error) {
	// Try to get existing brand
	brand, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if brand != nil {
		return brand, nil
	}

	// Create new brand
	return r.Create(name)
}
