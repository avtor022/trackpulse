package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type LapHistoryRepository interface {
	GetAll() ([]models.LapHistory, error)
	GetByID(id string) (*models.LapHistory, error)
	Create(history *models.LapHistory) error
	Update(history *models.LapHistory) error
	Delete(id string) error
	GetByRaceParticipantID(raceParticipantID string) ([]models.LapHistory, error)
	GetByRaceParticipantAndLapNumber(raceParticipantID string, lapNumber int) (*models.LapHistory, error)
	GetBestLapForParticipant(raceParticipantID string) (*models.LapHistory, error)
	GetValidLapsForParticipant(raceParticipantID string) ([]models.LapHistory, error)
}

type lapHistoryRepo struct {
	db *sql.DB
}

func NewLapHistoryRepository(db *sql.DB) LapHistoryRepository {
	return &lapHistoryRepo{db: db}
}

func (r *lapHistoryRepo) GetAll() ([]models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history ORDER BY race_participant_id, lap_number`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.LapHistory
	for rows.Next() {
		var history models.LapHistory
		var invalidationReason *string

		err := rows.Scan(
			&history.ID,
			&history.RaceParticipantID,
			&history.LapNumber,
			&history.LapTimeMs,
			&history.StartTime,
			&history.EndTime,
			&history.IsValid,
			&invalidationReason,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		history.InvalidationReason = invalidationReason

		histories = append(histories, history)
	}

	return histories, nil
}

func (r *lapHistoryRepo) GetByID(id string) (*models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history WHERE id = ?`

	var history models.LapHistory
	var invalidationReason *string

	err := r.db.QueryRow(query, id).Scan(
		&history.ID,
		&history.RaceParticipantID,
		&history.LapNumber,
		&history.LapTimeMs,
		&history.StartTime,
		&history.EndTime,
		&history.IsValid,
		&invalidationReason,
		&history.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lap history not found")
		}
		return nil, err
	}

	// Handle nullable fields
	history.InvalidationReason = invalidationReason

	return &history, nil
}

func (r *lapHistoryRepo) Create(history *models.LapHistory) error {
	query := `INSERT INTO lap_history (id, race_participant_id, lap_number, lap_time_ms, 
																		start_time, end_time, is_valid, invalidation_reason, created_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var invalidationReason *string
	
	if history.InvalidationReason != nil {
		invalidationReason = history.InvalidationReason
	}

	_, err := r.db.Exec(query,
		history.ID,
		history.RaceParticipantID,
		history.LapNumber,
		history.LapTimeMs,
		history.StartTime,
		history.EndTime,
		history.IsValid,
		invalidationReason,
		history.CreatedAt,
	)
	return err
}

func (r *lapHistoryRepo) Update(history *models.LapHistory) error {
	query := `UPDATE lap_history SET race_participant_id = ?, lap_number = ?, 
						lap_time_ms = ?, start_time = ?, end_time = ?, is_valid = ?, 
						invalidation_reason = ?, created_at = ?
					WHERE id = ?`

	// Extract values for nullable fields
	var invalidationReason *string
	
	if history.InvalidationReason != nil {
		invalidationReason = history.InvalidationReason
	}

	_, err := r.db.Exec(query,
		history.RaceParticipantID,
		history.LapNumber,
		history.LapTimeMs,
		history.StartTime,
		history.EndTime,
		history.IsValid,
		invalidationReason,
		history.CreatedAt,
		history.ID,
	)
	return err
}

func (r *lapHistoryRepo) Delete(id string) error {
	query := `DELETE FROM lap_history WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("lap history not found")
	}

	return nil
}

func (r *lapHistoryRepo) GetByRaceParticipantID(raceParticipantID string) ([]models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history WHERE race_participant_id = ? 
					ORDER BY lap_number`

	rows, err := r.db.Query(query, raceParticipantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.LapHistory
	for rows.Next() {
		var history models.LapHistory
		var invalidationReason *string

		err := rows.Scan(
			&history.ID,
			&history.RaceParticipantID,
			&history.LapNumber,
			&history.LapTimeMs,
			&history.StartTime,
			&history.EndTime,
			&history.IsValid,
			&invalidationReason,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		history.InvalidationReason = invalidationReason

		histories = append(histories, history)
	}

	return histories, nil
}

func (r *lapHistoryRepo) GetByRaceParticipantAndLapNumber(raceParticipantID string, lapNumber int) (*models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history WHERE race_participant_id = ? AND lap_number = ?`

	var history models.LapHistory
	var invalidationReason *string

	err := r.db.QueryRow(query, raceParticipantID, lapNumber).Scan(
		&history.ID,
		&history.RaceParticipantID,
		&history.LapNumber,
		&history.LapTimeMs,
		&history.StartTime,
		&history.EndTime,
		&history.IsValid,
		&invalidationReason,
		&history.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lap history not found")
		}
		return nil, err
	}

	// Handle nullable fields
	history.InvalidationReason = invalidationReason

	return &history, nil
}

func (r *lapHistoryRepo) GetBestLapForParticipant(raceParticipantID string) (*models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history 
					WHERE race_participant_id = ? AND is_valid = 1
					ORDER BY lap_time_ms ASC LIMIT 1`

	var history models.LapHistory
	var invalidationReason *string

	err := r.db.QueryRow(query, raceParticipantID).Scan(
		&history.ID,
		&history.RaceParticipantID,
		&history.LapNumber,
		&history.LapTimeMs,
		&history.StartTime,
		&history.EndTime,
		&history.IsValid,
		&invalidationReason,
		&history.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no valid laps found for participant")
		}
		return nil, err
	}

	// Handle nullable fields
	history.InvalidationReason = invalidationReason

	return &history, nil
}

func (r *lapHistoryRepo) GetValidLapsForParticipant(raceParticipantID string) ([]models.LapHistory, error) {
	query := `SELECT id, race_participant_id, lap_number, lap_time_ms, start_time, 
						end_time, is_valid, invalidation_reason, created_at 
					FROM lap_history 
					WHERE race_participant_id = ? AND is_valid = 1
					ORDER BY lap_number`

	rows, err := r.db.Query(query, raceParticipantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.LapHistory
	for rows.Next() {
		var history models.LapHistory
		var invalidationReason *string

		err := rows.Scan(
			&history.ID,
			&history.RaceParticipantID,
			&history.LapNumber,
			&history.LapTimeMs,
			&history.StartTime,
			&history.EndTime,
			&history.IsValid,
			&invalidationReason,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		history.InvalidationReason = invalidationReason

		histories = append(histories, history)
	}

	return histories, nil
}