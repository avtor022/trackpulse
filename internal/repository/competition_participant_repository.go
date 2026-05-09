package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// CompetitionParticipantRepository handles data access for competition participants
type CompetitionParticipantRepository struct {
	db *sql.DB
}

// NewCompetitionParticipantRepository creates a new competition participant repository
func NewCompetitionParticipantRepository(db *sql.DB) *CompetitionParticipantRepository {
	return &CompetitionParticipantRepository{db: db}
}

// GetAll returns all competition participants
func (r *CompetitionParticipantRepository) GetAll() ([]models.CompetitionParticipant, error) {
	rows, err := r.db.Query(`
		SELECT id, competition_id, competitor_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at
		FROM competition_participants
		ORDER BY competition_id, grid_position ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition_participants: %w", err)
	}
	defer rows.Close()

	var participants []models.CompetitionParticipant
	for rows.Next() {
		var p models.CompetitionParticipant
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&p.ID,
			&p.CompetitionID,
			&p.CompetitorModelID,
			&p.GridPosition,
			&p.IsFinished,
			&p.Disqualified,
			&p.DNFReason,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competition participant: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			p.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			p.UpdatedAt = t
		}

		participants = append(participants, p)
	}

	return participants, rows.Err()
}

// GetByID returns a competition participant by ID
func (r *CompetitionParticipantRepository) GetByID(id string) (*models.CompetitionParticipant, error) {
	row := r.db.QueryRow(`
		SELECT id, competition_id, competitor_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at
		FROM competition_participants
		WHERE id = ?
	`, id)

	var p models.CompetitionParticipant
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&p.ID,
		&p.CompetitionID,
		&p.CompetitorModelID,
		&p.GridPosition,
		&p.IsFinished,
		&p.Disqualified,
		&p.DNFReason,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competition participant: %w", err)
	}

	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		p.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		p.UpdatedAt = t
	}

	return &p, nil
}

// GetByCompetitionID returns all participants for a specific competition
func (r *CompetitionParticipantRepository) GetByCompetitionID(competitionID string) ([]models.CompetitionParticipant, error) {
	rows, err := r.db.Query(`
		SELECT id, competition_id, competitor_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at
		FROM competition_participants
		WHERE competition_id = ?
		ORDER BY grid_position ASC
	`, competitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition participants by competition_id: %w", err)
	}
	defer rows.Close()

	var participants []models.CompetitionParticipant
	for rows.Next() {
		var p models.CompetitionParticipant
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&p.ID,
			&p.CompetitionID,
			&p.CompetitorModelID,
			&p.GridPosition,
			&p.IsFinished,
			&p.Disqualified,
			&p.DNFReason,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competition participant: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			p.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			p.UpdatedAt = t
		}

		participants = append(participants, p)
	}

	return participants, rows.Err()
}

// Create inserts a new competition participant
func (r *CompetitionParticipantRepository) Create(participant *models.CompetitionParticipant) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		INSERT INTO competition_participants (id, competition_id, competitor_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		participant.ID,
		participant.CompetitionID,
		participant.CompetitorModelID,
		participant.GridPosition,
		participant.IsFinished,
		participant.Disqualified,
		participant.DNFReason,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create competition participant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating competition participant")
	}

	return nil
}

// Update updates an existing competition participant
func (r *CompetitionParticipantRepository) Update(participant *models.CompetitionParticipant) error {
	now := time.Now().Format(time.RFC3339)

	result, err := r.db.Exec(`
		UPDATE competition_participants
		SET competition_id = ?, competitor_model_id = ?, grid_position = ?, is_finished = ?, disqualified = ?, dnf_reason = ?, updated_at = ?
		WHERE id = ?
	`,
		participant.CompetitionID,
		participant.CompetitorModelID,
		participant.GridPosition,
		participant.IsFinished,
		participant.Disqualified,
		participant.DNFReason,
		now,
		participant.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update competition participant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competition participant not found")
	}

	return nil
}

// Delete removes a competition participant by ID
func (r *CompetitionParticipantRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM competition_participants WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete competition participant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competition participant not found")
	}

	return nil
}

// DeleteByCompetitionAndCompetitorModel removes a competition participant by competition_id and competitor_model_id
func (r *CompetitionParticipantRepository) DeleteByCompetitionAndCompetitorModel(competitionID, competitorModelID string) error {
	result, err := r.db.Exec(`DELETE FROM competition_participants WHERE competition_id = ? AND competitor_model_id = ?`, competitionID, competitorModelID)
	if err != nil {
		return fmt.Errorf("failed to delete competition participant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competition participant not found")
	}

	return nil
}

// Exists checks if a participant already exists for the given competition and competitor model
func (r *CompetitionParticipantRepository) Exists(competitionID, competitorModelID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM competition_participants 
		WHERE competition_id = ? AND competitor_model_id = ?
	`, competitionID, competitorModelID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of competition participant: %w", err)
	}
	return count > 0, nil
}

// Count returns total number of competition participants
func (r *CompetitionParticipantRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM competition_participants`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count competition participants: %w", err)
	}
	return count, nil
}
