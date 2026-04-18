package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitionTrackRepository handles data access for competition tracks
type CompetitionTrackRepository struct {
	db *sql.DB
}

// NewCompetitionTrackRepository creates a new competition track repository
func NewCompetitionTrackRepository(db *sql.DB) *CompetitionTrackRepository {
	return &CompetitionTrackRepository{db: db}
}

// GetAll returns all competition tracks
func (r *CompetitionTrackRepository) GetAll() ([]models.CompetitionTrack, error) {
	rows, err := r.db.Query(`
		SELECT id, name, created_at, updated_at
		FROM competition_tracks
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition_tracks: %w", err)
	}
	defer rows.Close()

	var tracks []models.CompetitionTrack
	for rows.Next() {
		var track models.CompetitionTrack
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&track.ID,
			&track.Name,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			track.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			track.UpdatedAt = t
		}

		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

// GetByName returns a track by name
func (r *CompetitionTrackRepository) GetByName(name string) (*models.CompetitionTrack, error) {
	row := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM competition_tracks
		WHERE name = ?
	`, name)

	var track models.CompetitionTrack
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&track.ID,
		&track.Name,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		track.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		track.UpdatedAt = t
	}

	return &track, nil
}

// Create inserts a new track
func (r *CompetitionTrackRepository) Create(name string) (*models.CompetitionTrack, error) {
	now := time.Now().Format(time.RFC3339)
	id := uuid.New().String()

	result, err := r.db.Exec(`
		INSERT INTO competition_tracks (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, id, name, now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create track: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	track := &models.CompetitionTrack{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return track, nil
}

// GetOrCreate gets a track by name or creates it if it doesn't exist
func (r *CompetitionTrackRepository) GetOrCreate(name string) (*models.CompetitionTrack, error) {
	// Try to get existing track
	track, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}
	if track != nil {
		return track, nil
	}

	// Create new track
	return r.Create(name)
}

// Delete deletes a track by name
func (r *CompetitionTrackRepository) Delete(name string) error {
	_, err := r.db.Exec(`DELETE FROM competition_tracks WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}
	return nil
}

