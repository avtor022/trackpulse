package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RaceParticipantsRepository interface {
	GetAll() ([]models.RaceParticipant, error)
	GetByID(id string) (*models.RaceParticipant, error)
	Create(participant *models.RaceParticipant) error
	Update(participant *models.RaceParticipant) error
	Delete(id string) error
	GetByRaceID(raceID string) ([]models.RaceParticipant, error)
	GetByRacerModelID(racerModelID string) ([]models.RaceParticipant, error)
	GetParticipantWithDetails(participantID string) (*ParticipantWithDetails, error)
}

type ParticipantWithDetails struct {
	ID              string
	RaceID          string
	RaceTitle       string
	RacerModelID    string
	RacerNumber     int
	RacerFullName   string
	ModelBrand      string
	ModelName       string
	GridPosition    *int
	IsFinished      bool
	Disqualified    bool
	DNFReason       *string
	CreatedAt       string
	UpdatedAt       string
}

type raceParticipantsRepo struct {
	db *sql.DB
}

func NewRaceParticipantsRepository(db *sql.DB) RaceParticipantsRepository {
	return &raceParticipantsRepo{db: db}
}

func (r *raceParticipantsRepo) GetAll() ([]models.RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, 
						disqualified, dnf_reason, created_at, updated_at 
					FROM race_participants ORDER BY race_id, grid_position`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.RaceParticipant
	for rows.Next() {
		var participant models.RaceParticipant
		var gridPosition *int
		var dnfReason *string

		err := rows.Scan(
			&participant.ID,
			&participant.RaceID,
			&participant.RacerModelID,
			&gridPosition,
			&participant.IsFinished,
			&participant.Disqualified,
			&dnfReason,
			&participant.CreatedAt,
			&participant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		participant.GridPosition = gridPosition
		participant.DNFReason = dnfReason

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *raceParticipantsRepo) GetByID(id string) (*models.RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, 
						disqualified, dnf_reason, created_at, updated_at 
					FROM race_participants WHERE id = ?`

	var participant models.RaceParticipant
	var gridPosition *int
	var dnfReason *string

	err := r.db.QueryRow(query, id).Scan(
		&participant.ID,
		&participant.RaceID,
		&participant.RacerModelID,
		&gridPosition,
		&participant.IsFinished,
		&participant.Disqualified,
		&dnfReason,
		&participant.CreatedAt,
		&participant.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("race participant not found")
		}
		return nil, err
	}

	// Handle nullable fields
	participant.GridPosition = gridPosition
	participant.DNFReason = dnfReason

	return &participant, nil
}

func (r *raceParticipantsRepo) Create(participant *models.RaceParticipant) error {
	query := `INSERT INTO race_participants (id, race_id, racer_model_id, grid_position, 
																			is_finished, disqualified, dnf_reason, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var gridPosition *int
	var dnfReason *string
	
	if participant.GridPosition != nil {
		gridPosition = participant.GridPosition
	}
	if participant.DNFReason != nil {
		dnfReason = participant.DNFReason
	}

	_, err := r.db.Exec(query,
		participant.ID,
		participant.RaceID,
		participant.RacerModelID,
		gridPosition,
		participant.IsFinished,
		participant.Disqualified,
		dnfReason,
		participant.CreatedAt,
		participant.UpdatedAt,
	)
	return err
}

func (r *raceParticipantsRepo) Update(participant *models.RaceParticipant) error {
	query := `UPDATE race_participants SET race_id = ?, racer_model_id = ?, 
						grid_position = ?, is_finished = ?, disqualified = ?, 
						dnf_reason = ?, updated_at = ?
					WHERE id = ?`

	// Extract values for nullable fields
	var gridPosition *int
	var dnfReason *string
	
	if participant.GridPosition != nil {
		gridPosition = participant.GridPosition
	}
	if participant.DNFReason != nil {
		dnfReason = participant.DNFReason
	}

	_, err := r.db.Exec(query,
		participant.RaceID,
		participant.RacerModelID,
		gridPosition,
		participant.IsFinished,
		participant.Disqualified,
		dnfReason,
		participant.UpdatedAt,
		participant.ID,
	)
	return err
}

func (r *raceParticipantsRepo) Delete(id string) error {
	query := `DELETE FROM race_participants WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("race participant not found")
	}

	return nil
}

func (r *raceParticipantsRepo) GetByRaceID(raceID string) ([]models.RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, 
						disqualified, dnf_reason, created_at, updated_at 
					FROM race_participants WHERE race_id = ? ORDER BY grid_position`

	rows, err := r.db.Query(query, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.RaceParticipant
	for rows.Next() {
		var participant models.RaceParticipant
		var gridPosition *int
		var dnfReason *string

		err := rows.Scan(
			&participant.ID,
			&participant.RaceID,
			&participant.RacerModelID,
			&gridPosition,
			&participant.IsFinished,
			&participant.Disqualified,
			&dnfReason,
			&participant.CreatedAt,
			&participant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		participant.GridPosition = gridPosition
		participant.DNFReason = dnfReason

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *raceParticipantsRepo) GetByRacerModelID(racerModelID string) ([]models.RaceParticipant, error) {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, 
						disqualified, dnf_reason, created_at, updated_at 
					FROM race_participants WHERE racer_model_id = ? ORDER BY race_id`

	rows, err := r.db.Query(query, racerModelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.RaceParticipant
	for rows.Next() {
		var participant models.RaceParticipant
		var gridPosition *int
		var dnfReason *string

		err := rows.Scan(
			&participant.ID,
			&participant.RaceID,
			&participant.RacerModelID,
			&gridPosition,
			&participant.IsFinished,
			&participant.Disqualified,
			&dnfReason,
			&participant.CreatedAt,
			&participant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		participant.GridPosition = gridPosition
		participant.DNFReason = dnfReason

		participants = append(participants, participant)
	}

	return participants, nil
}

func (r *raceParticipantsRepo) GetParticipantWithDetails(participantID string) (*ParticipantWithDetails, error) {
	query := `
		SELECT 
			rp.id,
			rp.race_id,
			r.race_title,
			rp.racer_model_id,
			rc.racer_number,
			rc.full_name,
			rm.brand,
			rm.model_name,
			rp.grid_position,
			rp.is_finished,
			rp.disqualified,
			rp.dnf_reason,
			rp.created_at,
			rp.updated_at
		FROM race_participants rp
		JOIN races r ON rp.race_id = r.id
		JOIN racer_models rcm ON rp.racer_model_id = rcm.id
		JOIN racers rc ON rcm.racer_id = rc.id
		JOIN rc_models rm ON rcm.rc_model_id = rm.id
		WHERE rp.id = ?
	`

	var participant ParticipantWithDetails

	err := r.db.QueryRow(query, participantID).Scan(
		&participant.ID,
		&participant.RaceID,
		&participant.RaceTitle,
		&participant.RacerModelID,
		&participant.RacerNumber,
		&participant.RacerFullName,
		&participant.ModelBrand,
		&participant.ModelName,
		&participant.GridPosition,
		&participant.IsFinished,
		&participant.Disqualified,
		&participant.DNFReason,
		&participant.CreatedAt,
		&participant.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("race participant with details not found")
		}
		return nil, err
	}

	return &participant, nil
}