package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// RacerRepository handles data access for racers
type RacerRepository struct {
	db *sql.DB
}

// NewRacerRepository creates a new racer repository
func NewRacerRepository(db *sql.DB) *RacerRepository {
	return &RacerRepository{db: db}
}

// GetAll returns all racers
func (r *RacerRepository) GetAll() ([]models.Racer, error) {
	rows, err := r.db.Query(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM racers
		ORDER BY racer_number ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query racers: %w", err)
	}
	defer rows.Close()

	var racers []models.Racer
	for rows.Next() {
		var racer models.Racer
		var birthday sql.NullString
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&racer.ID,
			&racer.RacerNumber,
			&racer.FullName,
			&birthday,
			&racer.Country,
			&racer.City,
			&racer.Rating,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan racer: %w", err)
		}

		if birthday.Valid {
			t, err := time.Parse(time.RFC3339, birthday.String)
			if err == nil {
				racer.Birthday = &t
			}
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			racer.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			racer.UpdatedAt = t
		}

		racers = append(racers, racer)
	}

	return racers, rows.Err()
}

// GetByID returns a racer by ID
func (r *RacerRepository) GetByID(id string) (*models.Racer, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM racers
		WHERE id = ?
	`, id)

	var racer models.Racer
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&racer.ID,
		&racer.RacerNumber,
		&racer.FullName,
		&birthday,
		&racer.Country,
		&racer.City,
		&racer.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get racer: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			racer.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		racer.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		racer.UpdatedAt = t
	}

	return &racer, nil
}

// GetByNumber returns a racer by racer number
func (r *RacerRepository) GetByNumber(number int) (*models.Racer, error) {
	row := r.db.QueryRow(`
		SELECT id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at
		FROM racers
		WHERE racer_number = ?
	`, number)

	var racer models.Racer
	var birthday sql.NullString
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&racer.ID,
		&racer.RacerNumber,
		&racer.FullName,
		&birthday,
		&racer.Country,
		&racer.City,
		&racer.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get racer by number: %w", err)
	}

	if birthday.Valid {
		t, err := time.Parse(time.RFC3339, birthday.String)
		if err == nil {
			racer.Birthday = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		racer.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		racer.UpdatedAt = t
	}

	return &racer, nil
}

// Create inserts a new racer
func (r *RacerRepository) Create(racer *models.Racer) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if racer.Birthday != nil {
		birthdayStr = sql.NullString{String: racer.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		INSERT INTO racers (id, racer_number, full_name, birthday, country, city, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		racer.ID,
		racer.RacerNumber,
		racer.FullName,
		birthdayStr,
		racer.Country,
		racer.City,
		racer.Rating,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create racer: %w", err)
	}

	// Verify the row was inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating racer")
	}

	return nil
}

// Update updates an existing racer
func (r *RacerRepository) Update(racer *models.Racer) error {
	now := time.Now().Format(time.RFC3339)
	var birthdayStr sql.NullString
	if racer.Birthday != nil {
		birthdayStr = sql.NullString{String: racer.Birthday.Format(time.RFC3339), Valid: true}
	}

	result, err := r.db.Exec(`
		UPDATE racers
		SET racer_number = ?, full_name = ?, birthday = ?, country = ?, city = ?, rating = ?, updated_at = ?
		WHERE id = ?
	`,
		racer.RacerNumber,
		racer.FullName,
		birthdayStr,
		racer.Country,
		racer.City,
		racer.Rating,
		now,
		racer.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update racer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("racer not found")
	}

	return nil
}

// Delete removes a racer by ID
func (r *RacerRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM racers WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete racer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("racer not found")
	}

	return nil
}

// Count returns total number of racers
func (r *RacerRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM racers`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count racers: %w", err)
	}
	return count, nil
}
