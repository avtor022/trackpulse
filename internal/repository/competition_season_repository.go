package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitionSeasonRepository handles data access for competition seasons
type CompetitionSeasonRepository struct {
	db *sql.DB
}

// NewCompetitionSeasonRepository creates a new competition season repository
func NewCompetitionSeasonRepository(db *sql.DB) *CompetitionSeasonRepository {
	return &CompetitionSeasonRepository{db: db}
}

// GetAll returns all competition seasons
func (r *CompetitionSeasonRepository) GetAll() ([]models.CompetitionSeason, error) {
	rows, err := r.db.Query(`
		SELECT id, season, created_at, updated_at
		FROM competition_seasons
		ORDER BY season ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition_seasons: %w", err)
	}
	defer rows.Close()

	var seasons []models.CompetitionSeason
	for rows.Next() {
		var season models.CompetitionSeason
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&season.ID,
			&season.Season,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan season: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			season.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			season.UpdatedAt = t
		}

		seasons = append(seasons, season)
	}

	return seasons, rows.Err()
}

// GetBySeason returns a season by its value
func (r *CompetitionSeasonRepository) GetBySeason(season string) (*models.CompetitionSeason, error) {
	row := r.db.QueryRow(`
		SELECT id, season, created_at, updated_at
		FROM competition_seasons
		WHERE season = ?
	`, season)

	var s models.CompetitionSeason
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&s.ID,
		&s.Season,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get season: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		s.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		s.UpdatedAt = t
	}

	return &s, nil
}

// Create inserts a new season
func (r *CompetitionSeasonRepository) Create(season string) (*models.CompetitionSeason, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO competition_seasons (id, season, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, season, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create season: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	s := &models.CompetitionSeason{
		ID:        id,
		Season:    season,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s, nil
}

// GetOrCreate gets a season by value or creates it if it doesn't exist
func (r *CompetitionSeasonRepository) GetOrCreate(season string) (*models.CompetitionSeason, error) {
	// Try to get existing season
	s, err := r.GetBySeason(season)
	if err != nil {
		return nil, err
	}
	if s != nil {
		return s, nil
	}

	// Create new season
	return r.Create(season)
}

// Delete deletes a season by value
func (r *CompetitionSeasonRepository) Delete(season string) error {
	_, err := r.db.Exec(`DELETE FROM competition_seasons WHERE season = ?`, season)
	if err != nil {
		return fmt.Errorf("failed to delete season: %w", err)
	}
	return nil
}
