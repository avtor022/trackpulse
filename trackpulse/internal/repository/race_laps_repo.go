package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RaceLapsRepository interface {
	GetAll() ([]models.RaceLap, error)
	GetByID(id string) (*models.RaceLap, error)
	Create(lap *models.RaceLap) error
	Update(lap *models.RaceLap) error
	Delete(id string) error
	GetByRaceParticipantID(raceParticipantID string) (*models.RaceLap, error)
	UpdateRaceResults(raceParticipantID string, numberOfLaps int, bestLapTimeMs int, bestLapNumber int, lastLapTimeMs int, lastPassTime *string, totalRaceTimeMs int) error
}

type raceLapsRepo struct {
	db *sql.DB
}

func NewRaceLapsRepository(db *sql.DB) RaceLapsRepository {
	return &raceLapsRepo{db: db}
}

func (r *raceLapsRepo) GetAll() ([]models.RaceLap, error) {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, 
						best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, 
						total_race_time_ms, created_at, updated_at 
					FROM race_laps ORDER BY number_of_laps DESC, best_lap_time_ms ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var raceLaps []models.RaceLap
	for rows.Next() {
		var lap models.RaceLap
		var timeFinish, lastPassTime *string

		err := rows.Scan(
			&lap.ID,
			&lap.RaceParticipantID,
			&lap.TimeStart,
			&timeFinish,
			&lap.NumberOfLaps,
			&lap.BestLapTimeMs,
			&lap.BestLapNumber,
			&lap.LastLapTimeMs,
			&lastPassTime,
			&lap.TotalRaceTimeMs,
			&lap.CreatedAt,
			&lap.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if timeFinish != nil {
			// Convert string to *time.Time
		}
		lap.TimeFinish = nil // Placeholder until we properly handle time conversion
		lap.LastPassTime = nil // Placeholder until we properly handle time conversion

		raceLaps = append(raceLaps, lap)
	}

	return raceLaps, nil
}

func (r *raceLapsRepo) GetByID(id string) (*models.RaceLap, error) {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, 
						best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, 
						total_race_time_ms, created_at, updated_at 
					FROM race_laps WHERE id = ?`

	var lap models.RaceLap
	var timeFinish, lastPassTime *string

	err := r.db.QueryRow(query, id).Scan(
		&lap.ID,
		&lap.RaceParticipantID,
		&lap.TimeStart,
		&timeFinish,
		&lap.NumberOfLaps,
		&lap.BestLapTimeMs,
		&lap.BestLapNumber,
		&lap.LastLapTimeMs,
		&lastPassTime,
		&lap.TotalRaceTimeMs,
		&lap.CreatedAt,
		&lap.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("race lap not found")
		}
		return nil, err
	}

	// Handle nullable fields
	if timeFinish != nil {
		// Convert string to *time.Time
	}
	lap.TimeFinish = nil // Placeholder until we properly handle time conversion
	lap.LastPassTime = nil // Placeholder until we properly handle time conversion

	return &lap, nil
}

func (r *raceLapsRepo) Create(lap *models.RaceLap) error {
	query := `INSERT INTO race_laps (id, race_participant_id, time_start, time_finish, 
																	number_of_laps, best_lap_time_ms, best_lap_number, 
																	last_lap_time_ms, last_pass_time, total_race_time_ms, 
																	created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var timeFinish, lastPassTime *string
	
	if lap.TimeFinish != nil {
		// Convert *time.Time to string
	}
	if lap.LastPassTime != nil {
		// Convert *time.Time to string
	}

	_, err := r.db.Exec(query,
		lap.ID,
		lap.RaceParticipantID,
		lap.TimeStart,
		timeFinish, // Pass nil or converted string
		lap.NumberOfLaps,
		lap.BestLapTimeMs,
		lap.BestLapNumber,
		lap.LastLapTimeMs,
		lastPassTime, // Pass nil or converted string
		lap.TotalRaceTimeMs,
		lap.CreatedAt,
		lap.UpdatedAt,
	)
	return err
}

func (r *raceLapsRepo) Update(lap *models.RaceLap) error {
	query := `UPDATE race_laps SET race_participant_id = ?, time_start = ?, 
						time_finish = ?, number_of_laps = ?, best_lap_time_ms = ?, 
						best_lap_number = ?, last_lap_time_ms = ?, last_pass_time = ?, 
						total_race_time_ms = ?, updated_at = ?
					WHERE id = ?`

	// Extract values for nullable fields
	var timeFinish, lastPassTime *string
	
	if lap.TimeFinish != nil {
		// Convert *time.Time to string
	}
	if lap.LastPassTime != nil {
		// Convert *time.Time to string
	}

	_, err := r.db.Exec(query,
		lap.RaceParticipantID,
		lap.TimeStart,
		timeFinish, // Pass nil or converted string
		lap.NumberOfLaps,
		lap.BestLapTimeMs,
		lap.BestLapNumber,
		lap.LastLapTimeMs,
		lastPassTime, // Pass nil or converted string
		lap.TotalRaceTimeMs,
		lap.UpdatedAt,
		lap.ID,
	)
	return err
}

func (r *raceLapsRepo) Delete(id string) error {
	query := `DELETE FROM race_laps WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("race lap not found")
	}

	return nil
}

func (r *raceLapsRepo) GetByRaceParticipantID(raceParticipantID string) (*models.RaceLap, error) {
	query := `SELECT id, race_participant_id, time_start, time_finish, number_of_laps, 
						best_lap_time_ms, best_lap_number, last_lap_time_ms, last_pass_time, 
						total_race_time_ms, created_at, updated_at 
					FROM race_laps WHERE race_participant_id = ?`

	var lap models.RaceLap
	var timeFinish, lastPassTime *string

	err := r.db.QueryRow(query, raceParticipantID).Scan(
		&lap.ID,
		&lap.RaceParticipantID,
		&lap.TimeStart,
		&timeFinish,
		&lap.NumberOfLaps,
		&lap.BestLapTimeMs,
		&lap.BestLapNumber,
		&lap.LastLapTimeMs,
		&lastPassTime,
		&lap.TotalRaceTimeMs,
		&lap.CreatedAt,
		&lap.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no record exists, we might want to create a default one
			// For now, just return the error
			return nil, errors.New("race lap not found")
		}
		return nil, err
	}

	// Handle nullable fields
	if timeFinish != nil {
		// Convert string to *time.Time
	}
	lap.TimeFinish = nil // Placeholder until we properly handle time conversion
	lap.LastPassTime = nil // Placeholder until we properly handle time conversion

	return &lap, nil
}

func (r *raceLapsRepo) UpdateRaceResults(raceParticipantID string, numberOfLaps int, bestLapTimeMs int, bestLapNumber int, lastLapTimeMs int, lastPassTime *string, totalRaceTimeMs int) error {
	// First check if a record already exists for this race participant
	queryCheck := `SELECT COUNT(*) FROM race_laps WHERE race_participant_id = ?`
	
	var count int
	err := r.db.QueryRow(queryCheck, raceParticipantID).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		// Update existing record
		query := `UPDATE race_laps SET number_of_laps = ?, best_lap_time_ms = ?, 
							best_lap_number = ?, last_lap_time_ms = ?, last_pass_time = ?, 
							total_race_time_ms = ?, updated_at = datetime('now')
						WHERE race_participant_id = ?`
		
		_, err := r.db.Exec(query,
			numberOfLaps,
			bestLapTimeMs,
			bestLapNumber,
			lastLapTimeMs,
			lastPassTime,
			totalRaceTimeMs,
			raceParticipantID,
		)
		return err
	} else {
		// Create new record
		id := generateUUID() // This would need to be implemented
		
		query := `INSERT INTO race_laps (id, race_participant_id, time_start, number_of_laps, 
																		best_lap_time_ms, best_lap_number, last_lap_time_ms, 
																		last_pass_time, total_race_time_ms, created_at, updated_at)
							VALUES (?, ?, datetime('now'), ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`
		
		_, err := r.db.Exec(query,
			id,
			raceParticipantID,
			numberOfLaps,
			bestLapTimeMs,
			bestLapNumber,
			lastLapTimeMs,
			lastPassTime,
			totalRaceTimeMs,
		)
		return err
	}
}

// Helper function to generate UUIDs (would need to import appropriate package)
func generateUUID() string {
	// Implementation would use a UUID generation library
	// For now, returning a placeholder
	return "placeholder-uuid"
}