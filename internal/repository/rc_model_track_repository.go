package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RCModelTrackRepository handles data access for RC model tracks
type RCModelTrackRepository struct {
	db *sql.DB
}

// NewRCModelTrackRepository creates a new RC model track repository
func NewRCModelTrackRepository(db *sql.DB) *RCModelTrackRepository {
	return &RCModelTrackRepository{db: db}
}

// GetAll returns all RC model tracks
func (r *RCModelTrackRepository) GetAll() ([]models.RCModelTrack, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_tracks
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query rc_model_tracks: %w", err)
	}
	defer rows.Close()

	var tracks []models.RCModelTrack
	for rows.Next() {
		var t models.RCModelTrack
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model track: %w", err)
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

// GetByName returns a model track by name
func (r *RCModelTrackRepository) GetByName(name string) (*models.RCModelTrack, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM rc_model_tracks
		WHERE name = ?
	`, name)

	var t models.RCModelTrack
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
		return nil, fmt.Errorf("failed to get model track: %w", err)
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

// Create inserts a new model track
func (r *RCModelTrackRepository) Create(name string) (*models.RCModelTrack, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO rc_model_tracks (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create model track: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	t := &models.RCModelTrack{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return t, nil
}

// GetOrCreate gets a model track by name or creates it if it doesn't exist
func (r *RCModelTrackRepository) GetOrCreate(name string) (*models.RCModelTrack, error) {
	// Try to get existing model track
	t, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// Create new model track
	return r.Create(name)
}

// Delete deletes a model track by name
func (r *RCModelTrackRepository) Delete(name string) error {
	_, err := r.db.Exec(`DELETE FROM rc_model_tracks WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete model track: %w", err)
	}
	return nil
}

