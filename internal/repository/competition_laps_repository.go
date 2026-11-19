package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// CompetitionLapsRepository handles database operations for competition laps
type CompetitionLapsRepository struct {
	db *sql.DB
}

// NewCompetitionLapsRepository creates a new competition laps repository
func NewCompetitionLapsRepository(db *sql.DB) *CompetitionLapsRepository {
	return &CompetitionLapsRepository{db: db}
}

// GetByParticipantID gets competition laps by participant ID
func (r *CompetitionLapsRepository) GetByParticipantID(participantID string) (*models.CompetitionLaps, error) {
	query := `
		SELECT id, competition_participant_id, time_start, time_finish, 
		       number_of_laps, best_lap_time_ms, best_lap_number, 
		       last_lap_time_ms, last_pass_time, total_competition_time_ms,
		       created_at, updated_at
		FROM competition_laps
		WHERE competition_participant_id = ?
	`

	var laps models.CompetitionLaps
	var timeStart, createdAt, updatedAt string
	var timeFinish, lastPassTime sql.NullString

	err := r.db.QueryRow(query, participantID).Scan(
		&laps.ID,
		&laps.CompetitionParticipantID,
		&timeStart,
		&timeFinish,
		&laps.NumberOfLaps,
		&laps.BestLapTimeMs,
		&laps.BestLapNumber,
		&laps.LastLapTimeMs,
		&lastPassTime,
		&laps.TotalCompetitionTimeMs,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competition laps: %w", err)
	}
	if laps.TimeStart, err = time.Parse(time.RFC3339Nano, timeStart); err != nil {
		return nil, fmt.Errorf("invalid lap start time: %w", err)
	}
	if laps.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt); err != nil {
		return nil, fmt.Errorf("invalid lap creation time: %w", err)
	}
	if laps.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt); err != nil {
		return nil, fmt.Errorf("invalid lap update time: %w", err)
	}

	if timeFinish.Valid {
		t, err := time.Parse(time.RFC3339Nano, timeFinish.String)
		if err != nil {
			return nil, fmt.Errorf("invalid lap finish time: %w", err)
		}
		laps.TimeFinish = &t
	}

	if lastPassTime.Valid {
		t, err := time.Parse(time.RFC3339Nano, lastPassTime.String)
		if err != nil {
			return nil, fmt.Errorf("invalid last passage time: %w", err)
		}
		laps.LastPassTime = &t
	}

	return &laps, nil
}

// Create creates a new competition laps record
func (r *CompetitionLapsRepository) Create(laps *models.CompetitionLaps) error {
	query := `
		INSERT INTO competition_laps (
			id, competition_participant_id, time_start, time_finish,
			number_of_laps, best_lap_time_ms, best_lap_number,
			last_lap_time_ms, last_pass_time, total_competition_time_ms,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var timeFinish, lastPassTime sql.NullString
	if laps.TimeFinish != nil {
		timeFinish = sql.NullString{String: laps.TimeFinish.Format(time.RFC3339Nano), Valid: true}
	}
	if laps.LastPassTime != nil {
		lastPassTime = sql.NullString{String: laps.LastPassTime.Format(time.RFC3339Nano), Valid: true}
	}

	_, err := r.db.Exec(query,
		laps.ID,
		laps.CompetitionParticipantID,
		laps.TimeStart.Format(time.RFC3339Nano),
		timeFinish,
		laps.NumberOfLaps,
		laps.BestLapTimeMs,
		laps.BestLapNumber,
		laps.LastLapTimeMs,
		lastPassTime,
		laps.TotalCompetitionTimeMs,
		laps.CreatedAt.Format(time.RFC3339Nano),
		laps.UpdatedAt.Format(time.RFC3339Nano),
	)

	if err != nil {
		return fmt.Errorf("failed to create competition laps: %w", err)
	}

	return nil
}

// Update updates an existing competition laps record
func (r *CompetitionLapsRepository) Update(laps *models.CompetitionLaps) error {
	query := `
		UPDATE competition_laps
		SET time_start = ?, time_finish = ?,
		    number_of_laps = ?, best_lap_time_ms = ?, best_lap_number = ?,
		    last_lap_time_ms = ?, last_pass_time = ?, total_competition_time_ms = ?,
		    updated_at = ?
		WHERE competition_participant_id = ?
	`

	var timeFinish, lastPassTime sql.NullString
	if laps.TimeFinish != nil {
		timeFinish = sql.NullString{String: laps.TimeFinish.Format(time.RFC3339Nano), Valid: true}
	}
	if laps.LastPassTime != nil {
		lastPassTime = sql.NullString{String: laps.LastPassTime.Format(time.RFC3339Nano), Valid: true}
	}

	_, err := r.db.Exec(query,
		laps.TimeStart.Format(time.RFC3339Nano),
		timeFinish,
		laps.NumberOfLaps,
		laps.BestLapTimeMs,
		laps.BestLapNumber,
		laps.LastLapTimeMs,
		lastPassTime,
		laps.TotalCompetitionTimeMs,
		laps.UpdatedAt.Format(time.RFC3339Nano),
		laps.CompetitionParticipantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update competition laps: %w", err)
	}

	return nil
}

// Upsert creates or updates competition laps record
func (r *CompetitionLapsRepository) Upsert(laps *models.CompetitionLaps) error {
	existing, err := r.GetByParticipantID(laps.CompetitionParticipantID)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(laps)
	}

	return r.Update(laps)
}

// GetAllByCompetitionID gets all competition laps for a competition
func (r *CompetitionLapsRepository) GetAllByCompetitionID(competitionID string) ([]*models.CompetitionLaps, error) {
	query := `
		SELECT cl.id, cl.competition_participant_id, cl.time_start, cl.time_finish,
		       cl.number_of_laps, cl.best_lap_time_ms, cl.best_lap_number,
		       cl.last_lap_time_ms, cl.last_pass_time, cl.total_competition_time_ms,
		       cl.created_at, cl.updated_at
		FROM competition_laps cl
		JOIN competition_participants cp ON cl.competition_participant_id = cp.id
		WHERE cp.competition_id = ?
	`

	rows, err := r.db.Query(query, competitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query competition laps: %w", err)
	}
	defer rows.Close()

	var lapsList []*models.CompetitionLaps
	for rows.Next() {
		var laps models.CompetitionLaps
		var timeStart, createdAt, updatedAt string
		var timeFinish, lastPassTime sql.NullString

		err := rows.Scan(
			&laps.ID,
			&laps.CompetitionParticipantID,
			&timeStart,
			&timeFinish,
			&laps.NumberOfLaps,
			&laps.BestLapTimeMs,
			&laps.BestLapNumber,
			&laps.LastLapTimeMs,
			&lastPassTime,
			&laps.TotalCompetitionTimeMs,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competition laps: %w", err)
		}
		if laps.TimeStart, err = time.Parse(time.RFC3339Nano, timeStart); err != nil {
			return nil, fmt.Errorf("invalid lap start time: %w", err)
		}
		if laps.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt); err != nil {
			return nil, fmt.Errorf("invalid lap creation time: %w", err)
		}
		if laps.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt); err != nil {
			return nil, fmt.Errorf("invalid lap update time: %w", err)
		}

		if timeFinish.Valid {
			t, err := time.Parse(time.RFC3339Nano, timeFinish.String)
			if err != nil {
				return nil, fmt.Errorf("invalid lap finish time: %w", err)
			}
			laps.TimeFinish = &t
		}

		if lastPassTime.Valid {
			t, err := time.Parse(time.RFC3339Nano, lastPassTime.String)
			if err != nil {
				return nil, fmt.Errorf("invalid last passage time: %w", err)
			}
			laps.LastPassTime = &t
		}

		lapsList = append(lapsList, &laps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating competition laps: %w", err)
	}

	return lapsList, nil
}

// GetStandings returns registered competitors ordered by laps and fastest lap.
func (r *CompetitionLapsRepository) GetStandings(competitionID string) ([]models.CompetitionStanding, error) {
	rows, err := r.db.Query(`
		SELECT cp.id, c.competitor_number, c.full_name,
		       COALESCE(cl.number_of_laps, 0), COALESCE(cl.best_lap_time_ms, 0)
		FROM competition_participants cp
		JOIN competitor_models cm ON cm.id = cp.competitor_model_id
		JOIN competitors c ON c.id = cm.competitor_id
		LEFT JOIN competition_laps cl ON cl.competition_participant_id = cp.id
		WHERE cp.competition_id = ?
		ORDER BY COALESCE(cl.number_of_laps, 0) DESC,
		         CASE WHEN COALESCE(cl.best_lap_time_ms, 0) > 0 THEN 0 ELSE 1 END,
		         COALESCE(cl.best_lap_time_ms, 0) ASC, c.competitor_number ASC, cp.id ASC`, competitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query standings: %w", err)
	}
	defer rows.Close()
	var standings []models.CompetitionStanding
	for rows.Next() {
		var standing models.CompetitionStanding
		if err := rows.Scan(&standing.ParticipantID, &standing.CompetitorNumber, &standing.FullName, &standing.NumberOfLaps, &standing.BestLapTimeMs); err != nil {
			return nil, fmt.Errorf("failed to scan standings: %w", err)
		}
		standings = append(standings, standing)
	}
	return standings, rows.Err()
}
