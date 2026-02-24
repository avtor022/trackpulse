package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RacesRepository interface {
	GetAll() ([]models.Race, error)
	GetByID(id string) (*models.Race, error)
	Create(race *models.Race) error
	Update(race *models.Race) error
	Delete(id string) error
	GetActiveRaces() ([]models.Race, error)
	GetRacesByStatus(status string) ([]models.Race, error)
}

type racesRepo struct {
	db *sql.DB
}

func NewRacesRepository(db *sql.DB) RacesRepository {
	return &racesRepo{db: db}
}

func (r *racesRepo) GetAll() ([]models.Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, 
						lap_count_target, time_limit_minutes, time_start, time_finish, status, 
						created_at, updated_at 
					FROM races ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []models.Race
	for rows.Next() {
		var race models.Race
		var modelType, modelScale, trackName *string
		var lapCountTarget, timeLimitMinutes *int
		var timeStart, timeFinish *string // Will convert to *time.Time later

		err := rows.Scan(
			&race.ID,
			&race.RaceTitle,
			&race.RaceType,
			&modelType,
			&modelScale,
			&trackName,
			&lapCountTarget,
			&timeLimitMinutes,
			&timeStart,
			&timeFinish,
			&race.Status,
			&race.CreatedAt,
			&race.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		race.ModelType = modelType
		race.ModelScale = modelScale
		race.TrackName = trackName
		if lapCountTarget != nil {
			race.LapCountTarget = lapCountTarget
		}
		if timeLimitMinutes != nil {
			race.TimeLimitMinutes = timeLimitMinutes
		}

		races = append(races, race)
	}

	return races, nil
}

func (r *racesRepo) GetByID(id string) (*models.Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, 
						lap_count_target, time_limit_minutes, time_start, time_finish, status, 
						created_at, updated_at 
					FROM races WHERE id = ?`

	var race models.Race
	var modelType, modelScale, trackName *string
	var lapCountTarget, timeLimitMinutes *int
	var timeStart, timeFinish *string // Will convert to *time.Time later

	err := r.db.QueryRow(query, id).Scan(
		&race.ID,
		&race.RaceTitle,
		&race.RaceType,
		&modelType,
		&modelScale,
		&trackName,
		&lapCountTarget,
		&timeLimitMinutes,
		&timeStart,
		&timeFinish,
		&race.Status,
		&race.CreatedAt,
		&race.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("race not found")
		}
		return nil, err
	}

	// Handle nullable fields
	race.ModelType = modelType
	race.ModelScale = modelScale
	race.TrackName = trackName
	if lapCountTarget != nil {
		race.LapCountTarget = lapCountTarget
	}
	if timeLimitMinutes != nil {
		race.TimeLimitMinutes = timeLimitMinutes
	}

	return &race, nil
}

func (r *racesRepo) Create(race *models.Race) error {
	query := `INSERT INTO races (id, race_title, race_type, model_type, model_scale, 
																track_name, lap_count_target, time_limit_minutes, time_start, 
																time_finish, status, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var modelType, modelScale, trackName *string
	var lapCountTarget, timeLimitMinutes *int
	
	if race.ModelType != nil {
		modelType = race.ModelType
	}
	if race.ModelScale != nil {
		modelScale = race.ModelScale
	}
	if race.TrackName != nil {
		trackName = race.TrackName
	}
	if race.LapCountTarget != nil {
		lapCountTarget = race.LapCountTarget
	}
	if race.TimeLimitMinutes != nil {
		timeLimitMinutes = race.TimeLimitMinutes
	}

	_, err := r.db.Exec(query,
		race.ID,
		race.RaceTitle,
		race.RaceType,
		modelType,
		modelScale,
		trackName,
		lapCountTarget,
		timeLimitMinutes,
		race.TimeStart,
		race.TimeFinish,
		race.Status,
		race.CreatedAt,
		race.UpdatedAt,
	)
	return err
}

func (r *racesRepo) Update(race *models.Race) error {
	query := `UPDATE races SET race_title = ?, race_type = ?, model_type = ?, 
						model_scale = ?, track_name = ?, lap_count_target = ?, 
						time_limit_minutes = ?, time_start = ?, time_finish = ?, 
						status = ?, updated_at = ?
					WHERE id = ?`

	// Extract values for nullable fields
	var modelType, modelScale, trackName *string
	var lapCountTarget, timeLimitMinutes *int
	
	if race.ModelType != nil {
		modelType = race.ModelType
	}
	if race.ModelScale != nil {
		modelScale = race.ModelScale
	}
	if race.TrackName != nil {
		trackName = race.TrackName
	}
	if race.LapCountTarget != nil {
		lapCountTarget = race.LapCountTarget
	}
	if race.TimeLimitMinutes != nil {
		timeLimitMinutes = race.TimeLimitMinutes
	}

	_, err := r.db.Exec(query,
		race.RaceTitle,
		race.RaceType,
		modelType,
		modelScale,
		trackName,
		lapCountTarget,
		timeLimitMinutes,
		race.TimeStart,
		race.TimeFinish,
		race.Status,
		race.UpdatedAt,
		race.ID,
	)
	return err
}

func (r *racesRepo) Delete(id string) error {
	query := `DELETE FROM races WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("race not found")
	}

	return nil
}

func (r *racesRepo) GetActiveRaces() ([]models.Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, 
						lap_count_target, time_limit_minutes, time_start, time_finish, status, 
						created_at, updated_at 
					FROM races WHERE status = 'active' ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []models.Race
	for rows.Next() {
		var race models.Race
		var modelType, modelScale, trackName *string
		var lapCountTarget, timeLimitMinutes *int
		var timeStart, timeFinish *string // Will convert to *time.Time later

		err := rows.Scan(
			&race.ID,
			&race.RaceTitle,
			&race.RaceType,
			&modelType,
			&modelScale,
			&trackName,
			&lapCountTarget,
			&timeLimitMinutes,
			&timeStart,
			&timeFinish,
			&race.Status,
			&race.CreatedAt,
			&race.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		race.ModelType = modelType
		race.ModelScale = modelScale
		race.TrackName = trackName
		if lapCountTarget != nil {
			race.LapCountTarget = lapCountTarget
		}
		if timeLimitMinutes != nil {
			race.TimeLimitMinutes = timeLimitMinutes
		}

		races = append(races, race)
	}

	return races, nil
}

func (r *racesRepo) GetRacesByStatus(status string) ([]models.Race, error) {
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, 
						lap_count_target, time_limit_minutes, time_start, time_finish, status, 
						created_at, updated_at 
					FROM races WHERE status = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []models.Race
	for rows.Next() {
		var race models.Race
		var modelType, modelScale, trackName *string
		var lapCountTarget, timeLimitMinutes *int
		var timeStart, timeFinish *string // Will convert to *time.Time later

		err := rows.Scan(
			&race.ID,
			&race.RaceTitle,
			&race.RaceType,
			&modelType,
			&modelScale,
			&trackName,
			&lapCountTarget,
			&timeLimitMinutes,
			&timeStart,
			&timeFinish,
			&race.Status,
			&race.CreatedAt,
			&race.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		race.ModelType = modelType
		race.ModelScale = modelScale
		race.TrackName = trackName
		if lapCountTarget != nil {
			race.LapCountTarget = lapCountTarget
		}
		if timeLimitMinutes != nil {
			race.TimeLimitMinutes = timeLimitMinutes
		}

		races = append(races, race)
	}

	return races, nil
}