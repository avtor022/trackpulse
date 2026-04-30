package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitionYearRepository handles data access for competition years
type CompetitionYearRepository struct {
	db *sql.DB
}

// NewCompetitionYearRepository creates a new competition year repository
func NewCompetitionYearRepository(db *sql.DB) *CompetitionYearRepository {
	return &CompetitionYearRepository{db: db}
}

// GetAll returns all competition years
func (r *CompetitionYearRepository) GetAll() ([]models.CompetitionYear, error) {
	rows, err := r.db.Query(`
		SELECT id, year, created_at, updated_at
		FROM competition_years
		ORDER BY year DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition_years: %w", err)
	}
	defer rows.Close()

	var years []models.CompetitionYear
	for rows.Next() {
		var year models.CompetitionYear
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&year.ID,
			&year.Year,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan year: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			year.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			year.UpdatedAt = t
		}

		years = append(years, year)
	}

	return years, rows.Err()
}

// GetByYear returns a year by its value
func (r *CompetitionYearRepository) GetByYear(year int) (*models.CompetitionYear, error) {
	row := r.db.QueryRow(`
		SELECT id, year, created_at, updated_at
		FROM competition_years
		WHERE year = ?
	`, year)

	var y models.CompetitionYear
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&y.ID,
		&y.Year,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get year: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		y.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		y.UpdatedAt = t
	}

	return &y, nil
}

// Create inserts a new year
func (r *CompetitionYearRepository) Create(year int) (*models.CompetitionYear, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO competition_years (id, year, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, year, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create year: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	y := &models.CompetitionYear{
		ID:        id,
		Year:      year,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return y, nil
}

// GetOrCreate gets a year by value or creates it if it doesn't exist
func (r *CompetitionYearRepository) GetOrCreate(year int) (*models.CompetitionYear, error) {
	// Try to get existing year
	y, err := r.GetByYear(year)
	if err != nil {
		return nil, err
	}
	if y != nil {
		return y, nil
	}

	// Create new year
	return r.Create(year)
}

// Delete deletes a year by value
func (r *CompetitionYearRepository) Delete(year int) error {
	_, err := r.db.Exec(`DELETE FROM competition_years WHERE year = ?`, year)
	if err != nil {
		return fmt.Errorf("failed to delete year: %w", err)
	}
	return nil
}
