package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// AthleteRepository handles data access for athletes
type AthleteRepository struct {
	db *sql.DB
}

// NewAthleteRepository creates a new athlete repository
func NewAthleteRepository(db *sql.DB) *AthleteRepository {
	return &AthleteRepository{db: db}
}

// GetAll returns all athletes
func (r *AthleteRepository) GetAll() ([]models.Athlete, error) {
	rows, err := r.db.Query(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM athletes
		ORDER BY racer_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query athletes: %w", err)
	}
	defer rows.Close()

	var athletes []models.Athlete
	for rows.Next() {
		var athlete models.Athlete
		var birthday sql.NullString
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&athlete.ID,
			&athlete.RacerNumber,
			&athlete.FullName,
			&birthday,
			&athlete.Country,
			&athlete.City,
			&athlete.Rating,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan athlete: %w", err)
		}

		if birthday.Valid {
			t, err := time.Parse(time.RFC3339, birthday.String)
			if err == nil {
				athlete.Birthday = &t
			}
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			athlete.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			athlete.UpdatedAt = t
		}

		athletes = append(athletes, athlete)
	}

	return athletes, rows.Err()
}

// GetByID returns an athlete by ID
func (r *AthleteRepository) GetByID(id string) (*models.Athlete, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM athletes
		WHERE id = ?
	`, id)

	var athlete models.Athlete
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&athlete.ID,
		&athlete.RacerNumber,
		&athlete.FullName,
		&birthday,
		&athlete.Country,
		&athlete.City,
		&athlete.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			athlete.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		athlete.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		athlete.UpdatedAt = t
	}

	return &athlete, nil
}

// GetByNumber returns an athlete by racer number
func (r *AthleteRepository) GetByNumber(number int) (*models.Athlete, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM athletes
		WHERE racer_number = ?
	`, number)

	var athlete models.Athlete
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&athlete.ID,
		&athlete.RacerNumber,
		&athlete.FullName,
		&birthday,
		&athlete.Country,
		&athlete.City,
		&athlete.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get athlete by number: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			athlete.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		athlete.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		athlete.UpdatedAt = t
	}

	return &athlete, nil
}

// Create inserts a new athlete
func (r *AthleteRepository) Create(athlete *models.Athlete) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if athlete.Birthday != nil {
		birthdayStr = sql.NullString{String: athlete.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		INSERT INTO athletes (id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		athlete.ID,
		athlete.RacerNumber,
		athlete.FullName,
		birthdayStr,
		athlete.Country,
		athlete.City,
		athlete.Rating,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create athlete: %w", err)
	}

	// Verify the row was inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating athlete")
	}

	return nil
}

// Update updates an existing athlete
func (r *AthleteRepository) Update(athlete *models.Athlete) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if athlete.Birthday != nil {
		birthdayStr = sql.NullString{String: athlete.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		UPDATE athletes
		SET racer_number = ?, full_name = ?, birthday = ?, country = ?, city = ?, rating = ?, updated_at = ?
		WHERE id = ?
	`,
		athlete.RacerNumber,
		athlete.FullName,
		birthdayStr,
		athlete.Country,
		athlete.City,
		athlete.Rating,
		now,
		athlete.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update athlete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete not found")
	}

	return nil
}

// Delete removes an athlete by ID
func (r *AthleteRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM athletes WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete athlete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete not found")
	}

	return nil
}

// Count returns total number of athletes
func (r *AthleteRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM athletes`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count athletes: %w", err)
	}
	return count, nil
}
