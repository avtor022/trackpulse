package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// CompetitorRepository handles data access for competitors
type CompetitorRepository struct {
	db *sql.DB
}

// NewCompetitorRepository creates a new competitor repository
func NewCompetitorRepository(db *sql.DB) *CompetitorRepository {
	return &CompetitorRepository{db: db}
}

// GetAll returns all competitors
func (r *CompetitorRepository) GetAll() ([]models.Competitor, error) {
	rows, err := r.db.Query(`
		SELECT id, competitor_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM competitors
		ORDER BY competitor_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competitors: %w", err)
	}
	defer rows.Close()

	var competitors []models.Competitor
	for rows.Next() {
		var competitor models.Competitor
		var birthday sql.NullString
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&competitor.ID,
			&competitor.CompetitorNumber,
			&competitor.FullName,
			&birthday,
			&competitor.Country,
			&competitor.City,
			&competitor.Rating,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competitor: %w", err)
		}

		if birthday.Valid {
			t, err := time.Parse(time.RFC3339, birthday.String)
			if err == nil {
				competitor.Birthday = &t
			}
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			competitor.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			competitor.UpdatedAt = t
		}

		competitors = append(competitors, competitor)
	}

	return competitors, rows.Err()
}

// GetByID returns a competitor by ID
func (r *CompetitorRepository) GetByID(id string) (*models.Competitor, error) {
	row := r.db.QueryRow(`
		SELECT id, competitor_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM competitors
		WHERE id = ?
	`, id)

	var competitor models.Competitor
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&competitor.ID,
		&competitor.CompetitorNumber,
		&competitor.FullName,
		&birthday,
		&competitor.Country,
		&competitor.City,
		&competitor.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			competitor.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		competitor.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		competitor.UpdatedAt = t
	}

	return &competitor, nil
}

// GetByNumber returns a competitor by competitor number
func (r *CompetitorRepository) GetByNumber(number int) (*models.Competitor, error) {
	row := r.db.QueryRow(`
		SELECT id, competitor_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM competitors
		WHERE competitor_number = ?
	`, number)

	var competitor models.Competitor
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&competitor.ID,
		&competitor.CompetitorNumber,
		&competitor.FullName,
		&birthday,
		&competitor.Country,
		&competitor.City,
		&competitor.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor by number: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			competitor.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		competitor.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		competitor.UpdatedAt = t
	}

	return &competitor, nil
}

// Create inserts a new competitor
func (r *CompetitorRepository) Create(competitor *models.Competitor) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if competitor.Birthday != nil {
		birthdayStr = sql.NullString{String: competitor.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		INSERT INTO competitors (id, competitor_number, full_name, birthday, country, city, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		competitor.ID,
		competitor.CompetitorNumber,
		competitor.FullName,
		birthdayStr,
		competitor.Country,
		competitor.City,
		competitor.Rating,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create competitor: %w", err)
	}

	// Verify the row was inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating competitor")
	}

	return nil
}

// Update updates an existing competitor
func (r *CompetitorRepository) Update(competitor *models.Competitor) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if competitor.Birthday != nil {
		birthdayStr = sql.NullString{String: competitor.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		UPDATE competitors
		SET competitor_number = ?, full_name = ?, birthday = ?, country = ?, city = ?, rating = ?, updated_at = ?
		WHERE id = ?
	`,
		competitor.CompetitorNumber,
		competitor.FullName,
		birthdayStr,
		competitor.Country,
		competitor.City,
		competitor.Rating,
		now,
		competitor.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update competitor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competitor not found")
	}

	return nil
}

// Delete removes a competitor by ID
func (r *CompetitorRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM competitors WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete competitor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competitor not found")
	}

	return nil
}

// Count returns total number of competitors
func (r *CompetitorRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM competitors`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count competitors: %w", err)
	}
	return count, nil
}
