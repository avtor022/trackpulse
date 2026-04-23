package repository

import (
	"database/sql"
	"fmt"
	"time"

	"trackpulse/internal/models"
)

// CompetitionRepository handles data access for competitions
type CompetitionRepository struct {
	db *sql.DB
}

// NewCompetitionRepository creates a new competition repository
func NewCompetitionRepository(db *sql.DB) *CompetitionRepository {
	return &CompetitionRepository{db: db}
}

// GetAll returns all competitions
func (r *CompetitionRepository) GetAll() ([]models.Competition, error) {
	rows, err := r.db.Query(`
		SELECT id, competition_title, competition_type, model_type, model_scale, track_name, 
		       lap_count_target, time_limit_minutes, time_start, time_finish, status, 
		       season_name, competition_year, created_at, updated_at
		FROM competitions
		ORDER BY time_start DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query competitions: %w", err)
	}
	defer rows.Close()

	var competitions []models.Competition
	for rows.Next() {
		var c models.Competition
		var lapCountTarget, timeLimitMinutes sql.NullInt64
		var timeStartStr, timeFinishStr sql.NullString
		var seasonName sql.NullString
		var competitionYear sql.NullInt64
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&c.ID,
			&c.CompetitionTitle,
			&c.CompetitionType,
			&c.ModelType,
			&c.ModelScale,
			&c.TrackName,
			&lapCountTarget,
			&timeLimitMinutes,
			&timeStartStr,
			&timeFinishStr,
			&c.Status,
			&seasonName,
			&competitionYear,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competition: %w", err)
		}

		if lapCountTarget.Valid {
			lct := int(lapCountTarget.Int64)
			c.LapCountTarget = &lct
		}
		if timeLimitMinutes.Valid {
			tlm := int(timeLimitMinutes.Int64)
			c.TimeLimitMinutes = &tlm
		}
		if timeStartStr.Valid {
			if t, err := time.Parse(time.RFC3339, timeStartStr.String); err == nil {
				c.TimeStart = &t
			}
		}
		if timeFinishStr.Valid {
			if t, err := time.Parse(time.RFC3339, timeFinishStr.String); err == nil {
				c.TimeFinish = &t
			}
		}
		if seasonName.Valid {
			sn := seasonName.String
			c.SeasonName = &sn
		}
		if competitionYear.Valid {
			cy := int(competitionYear.Int64)
			c.CompetitionYear = &cy
		}
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			c.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			c.UpdatedAt = t
		}

		competitions = append(competitions, c)
	}

	return competitions, rows.Err()
}

// GetByID returns a competition by ID
func (r *CompetitionRepository) GetByID(id string) (*models.Competition, error) {
	row := r.db.QueryRow(`
		SELECT id, competition_title, competition_type, model_type, model_scale, track_name, 
		       lap_count_target, time_limit_minutes, time_start, time_finish, status, 
		       season_name, competition_year, created_at, updated_at
		FROM competitions
		WHERE id = ?
	`, id)

	var c models.Competition
	var lapCountTarget, timeLimitMinutes sql.NullInt64
	var timeStartStr, timeFinishStr sql.NullString
	var seasonName sql.NullString
	var competitionYear sql.NullInt64
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&c.ID,
		&c.CompetitionTitle,
		&c.CompetitionType,
		&c.ModelType,
		&c.ModelScale,
		&c.TrackName,
		&lapCountTarget,
		&timeLimitMinutes,
		&timeStartStr,
		&timeFinishStr,
		&c.Status,
		&seasonName,
		&competitionYear,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get competition: %w", err)
	}

	if lapCountTarget.Valid {
		lct := int(lapCountTarget.Int64)
		c.LapCountTarget = &lct
	}
	if timeLimitMinutes.Valid {
		tlm := int(timeLimitMinutes.Int64)
		c.TimeLimitMinutes = &tlm
	}
	if timeStartStr.Valid {
		if t, err := time.Parse(time.RFC3339, timeStartStr.String); err == nil {
			c.TimeStart = &t
		}
	}
	if timeFinishStr.Valid {
		if t, err := time.Parse(time.RFC3339, timeFinishStr.String); err == nil {
			c.TimeFinish = &t
		}
	}
	if seasonName.Valid {
		sn := seasonName.String
		c.SeasonName = &sn
	}
	if competitionYear.Valid {
		cy := int(competitionYear.Int64)
		c.CompetitionYear = &cy
	}
	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		c.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
		c.UpdatedAt = t
	}

	return &c, nil
}

// Create inserts a new competition
func (r *CompetitionRepository) Create(competition *models.Competition) error {
	now := time.Now().Format(time.RFC3339)

	var lapCountTarget sql.NullInt64
	if competition.LapCountTarget != nil {
		lapCountTarget = sql.NullInt64{Int64: int64(*competition.LapCountTarget), Valid: true}
	}

	var timeLimitMinutes sql.NullInt64
	if competition.TimeLimitMinutes != nil {
		timeLimitMinutes = sql.NullInt64{Int64: int64(*competition.TimeLimitMinutes), Valid: true}
	}

	var timeStartStr, timeFinishStr sql.NullString
	if competition.TimeStart != nil {
		timeStartStr = sql.NullString{String: competition.TimeStart.Format(time.RFC3339), Valid: true}
	}
	if competition.TimeFinish != nil {
		timeFinishStr = sql.NullString{String: competition.TimeFinish.Format(time.RFC3339), Valid: true}
	}

	var seasonName sql.NullString
	if competition.SeasonName != nil {
		seasonName = sql.NullString{String: *competition.SeasonName, Valid: true}
	}

	var competitionYear sql.NullInt64
	if competition.CompetitionYear != nil {
		competitionYear = sql.NullInt64{Int64: int64(*competition.CompetitionYear), Valid: true}
	}

	result, err := r.db.Exec(`
		INSERT INTO competitions (id, competition_title, competition_type, model_type, model_scale, 
		                          track_name, lap_count_target, time_limit_minutes, time_start, 
		                          time_finish, status, season_name, competition_year, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		competition.ID,
		competition.CompetitionTitle,
		competition.CompetitionType,
		competition.ModelType,
		competition.ModelScale,
		competition.TrackName,
		lapCountTarget,
		timeLimitMinutes,
		timeStartStr,
		timeFinishStr,
		competition.Status,
		seasonName,
		competitionYear,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create competition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating competition")
	}

	return nil
}

// Update updates an existing competition
func (r *CompetitionRepository) Update(competition *models.Competition) error {
	now := time.Now().Format(time.RFC3339)

	var lapCountTarget sql.NullInt64
	if competition.LapCountTarget != nil {
		lapCountTarget = sql.NullInt64{Int64: int64(*competition.LapCountTarget), Valid: true}
	}

	var timeLimitMinutes sql.NullInt64
	if competition.TimeLimitMinutes != nil {
		timeLimitMinutes = sql.NullInt64{Int64: int64(*competition.TimeLimitMinutes), Valid: true}
	}

	var timeStartStr, timeFinishStr sql.NullString
	if competition.TimeStart != nil {
		timeStartStr = sql.NullString{String: competition.TimeStart.Format(time.RFC3339), Valid: true}
	}
	if competition.TimeFinish != nil {
		timeFinishStr = sql.NullString{String: competition.TimeFinish.Format(time.RFC3339), Valid: true}
	}

	var seasonName sql.NullString
	if competition.SeasonName != nil {
		seasonName = sql.NullString{String: *competition.SeasonName, Valid: true}
	}

	var competitionYear sql.NullInt64
	if competition.CompetitionYear != nil {
		competitionYear = sql.NullInt64{Int64: int64(*competition.CompetitionYear), Valid: true}
	}

	result, err := r.db.Exec(`
		UPDATE competitions
		SET competition_title = ?, competition_type = ?, model_type = ?, model_scale = ?, 
		    track_name = ?, lap_count_target = ?, time_limit_minutes = ?, time_start = ?, 
		    time_finish = ?, status = ?, season_name = ?, competition_year = ?, updated_at = ?
		WHERE id = ?
	`,
		competition.CompetitionTitle,
		competition.CompetitionType,
		competition.ModelType,
		competition.ModelScale,
		competition.TrackName,
		lapCountTarget,
		timeLimitMinutes,
		timeStartStr,
		timeFinishStr,
		competition.Status,
		seasonName,
		competitionYear,
		now,
		competition.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update competition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competition not found")
	}

	return nil
}

// Delete removes a competition by ID
func (r *CompetitionRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM competitions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete competition: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("competition not found")
	}

	return nil
}

// Count returns total number of competitions
func (r *CompetitionRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM competitions`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count competitions: %w", err)
	}
	return count, nil
}

// GetByStatus returns competitions filtered by status
func (r *CompetitionRepository) GetByStatus(status string) ([]models.Competition, error) {
	rows, err := r.db.Query(`
		SELECT id, competition_title, competition_type, model_type, model_scale, track_name, 
		       lap_count_target, time_limit_minutes, time_start, time_finish, status, 
		       created_at, updated_at
		FROM competitions
		WHERE status = ?
		ORDER BY time_start DESC
	`, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query competitions by status: %w", err)
	}
	defer rows.Close()

	var competitions []models.Competition
	for rows.Next() {
		var c models.Competition
		var lapCountTarget, timeLimitMinutes sql.NullInt64
		var timeStartStr, timeFinishStr sql.NullString
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&c.ID,
			&c.CompetitionTitle,
			&c.CompetitionType,
			&c.ModelType,
			&c.ModelScale,
			&c.TrackName,
			&lapCountTarget,
			&timeLimitMinutes,
			&timeStartStr,
			&timeFinishStr,
			&c.Status,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan competition: %w", err)
		}

		if lapCountTarget.Valid {
			lct := int(lapCountTarget.Int64)
			c.LapCountTarget = &lct
		}
		if timeLimitMinutes.Valid {
			tlm := int(timeLimitMinutes.Int64)
			c.TimeLimitMinutes = &tlm
		}
		if timeStartStr.Valid {
			if t, err := time.Parse(time.RFC3339, timeStartStr.String); err == nil {
				c.TimeStart = &t
			}
		}
		if timeFinishStr.Valid {
			if t, err := time.Parse(time.RFC3339, timeFinishStr.String); err == nil {
				c.TimeFinish = &t
			}
		}
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			c.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			c.UpdatedAt = t
		}

		competitions = append(competitions, c)
	}

	return competitions, rows.Err()
}
