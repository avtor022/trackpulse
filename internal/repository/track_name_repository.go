package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// TrackNameRepository handles data access for track names
type TrackNameRepository struct {
	db *sql.DB
}

// NewTrackNameRepository creates a new track name repository
func NewTrackNameRepository(db *sql.DB) *TrackNameRepository {
	return &TrackNameRepository{db: db}
}

// GetAll returns all track names
func (r *TrackNameRepository) GetAll() ([]models.TrackName, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM track_names
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query track_names: %w", err)
	}
	defer rows.Close()

	var tracks []models.TrackName
	for rows.Next() {
		var t models.TrackName
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track name: %w", err)
		}

		createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err == nil {
			t.CreatedAt = createdAtTime
		}
		updatedAtTime, err := time.Parse(time.RFC3339, updatedAtStr)
		if err == nil {
			t.UpdatedAt = updatedAtTime
		}

		tracks = append(tracks, t)
	}

	return tracks, rows.Err()
}

// GetByName returns a track name by name
func (r *TrackNameRepository) GetByName(name string) (*models.TrackName, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM track_names
		WHERE name = ?
	`, name)

	var t models.TrackName
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
		return nil, fmt.Errorf("failed to get track name: %w", err)
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

// Create inserts a new track name
func (r *TrackNameRepository) Create(name string) (*models.TrackName, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO track_names (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create track name: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	t := &models.TrackName{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return t, nil
}

// GetOrCreate gets a track name by name or creates it if it doesn't exist
func (r *TrackNameRepository) GetOrCreate(name string) (*models.TrackName, error) {
	// Try to get existing track name
	t, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// Create new track name
	return r.Create(name)
}

// Delete deletes a track name by name
func (r *TrackNameRepository) Delete(name string) error {
	_, err := r.db.Exec(`DELETE FROM track_names WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete track name: %w", err)
	}
	return nil
}

