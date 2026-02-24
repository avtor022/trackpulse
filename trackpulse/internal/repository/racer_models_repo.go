package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RacerModelsRepository interface {
	GetAll() ([]models.RacerModel, error)
	GetByID(id string) (*models.RacerModel, error)
	Create(racerModel *models.RacerModel) error
	Update(racerModel *models.RacerModel) error
	Delete(id string) error
	GetByTransponder(transponderNumber string) (*models.RacerModel, error)
	GetByRacerID(racerID string) ([]models.RacerModel, error)
	GetByRCModelID(rcModelID string) ([]models.RacerModel, error)
}

type racerModelsRepo struct {
	db *sql.DB
}

func NewRacerModelsRepository(db *sql.DB) RacerModelsRepository {
	return &racerModelsRepo{db: db}
}

func (r *racerModelsRepo) GetAll() ([]models.RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, 
						is_active, created_at, updated_at 
					FROM racer_models ORDER BY transponder_number`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racerModels []models.RacerModel
	for rows.Next() {
		var racerModel models.RacerModel

		err := rows.Scan(
			&racerModel.ID,
			&racerModel.RacerID,
			&racerModel.RCModelID,
			&racerModel.TransponderNumber,
			&racerModel.TransponderType,
			&racerModel.IsActive,
			&racerModel.CreatedAt,
			&racerModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		racerModels = append(racerModels, racerModel)
	}

	return racerModels, nil
}

func (r *racerModelsRepo) GetByID(id string) (*models.RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, 
						is_active, created_at, updated_at 
					FROM racer_models WHERE id = ?`

	var racerModel models.RacerModel

	err := r.db.QueryRow(query, id).Scan(
		&racerModel.ID,
		&racerModel.RacerID,
		&racerModel.RCModelID,
		&racerModel.TransponderNumber,
		&racerModel.TransponderType,
		&racerModel.IsActive,
		&racerModel.CreatedAt,
		&racerModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("racer model not found")
		}
		return nil, err
	}

	return &racerModel, nil
}

func (r *racerModelsRepo) Create(racerModel *models.RacerModel) error {
	query := `INSERT INTO racer_models (id, racer_id, rc_model_id, transponder_number, 
																		transponder_type, is_active, created_at, updated_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		racerModel.ID,
		racerModel.RacerID,
		racerModel.RCModelID,
		racerModel.TransponderNumber,
		racerModel.TransponderType,
		racerModel.IsActive,
		racerModel.CreatedAt,
		racerModel.UpdatedAt,
	)
	return err
}

func (r *racerModelsRepo) Update(racerModel *models.RacerModel) error {
	query := `UPDATE racer_models SET racer_id = ?, rc_model_id = ?, 
						transponder_number = ?, transponder_type = ?, is_active = ?, 
						updated_at = ?
					WHERE id = ?`

	_, err := r.db.Exec(query,
		racerModel.RacerID,
		racerModel.RCModelID,
		racerModel.TransponderNumber,
		racerModel.TransponderType,
		racerModel.IsActive,
		racerModel.UpdatedAt,
		racerModel.ID,
	)
	return err
}

func (r *racerModelsRepo) Delete(id string) error {
	query := `DELETE FROM racer_models WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("racer model not found")
	}

	return nil
}

func (r *racerModelsRepo) GetByTransponder(transponderNumber string) (*models.RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, 
						is_active, created_at, updated_at 
					FROM racer_models WHERE transponder_number = ? AND is_active = 1`

	var racerModel models.RacerModel

	err := r.db.QueryRow(query, transponderNumber).Scan(
		&racerModel.ID,
		&racerModel.RacerID,
		&racerModel.RCModelID,
		&racerModel.TransponderNumber,
		&racerModel.TransponderType,
		&racerModel.IsActive,
		&racerModel.CreatedAt,
		&racerModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("racer model not found")
		}
		return nil, err
	}

	return &racerModel, nil
}

func (r *racerModelsRepo) GetByRacerID(racerID string) ([]models.RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, 
						is_active, created_at, updated_at 
					FROM racer_models WHERE racer_id = ? ORDER BY transponder_number`

	rows, err := r.db.Query(query, racerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racerModels []models.RacerModel
	for rows.Next() {
		var racerModel models.RacerModel

		err := rows.Scan(
			&racerModel.ID,
			&racerModel.RacerID,
			&racerModel.RCModelID,
			&racerModel.TransponderNumber,
			&racerModel.TransponderType,
			&racerModel.IsActive,
			&racerModel.CreatedAt,
			&racerModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		racerModels = append(racerModels, racerModel)
	}

	return racerModels, nil
}

func (r *racerModelsRepo) GetByRCModelID(rcModelID string) ([]models.RacerModel, error) {
	query := `SELECT id, racer_id, rc_model_id, transponder_number, transponder_type, 
						is_active, created_at, updated_at 
					FROM racer_models WHERE rc_model_id = ? ORDER BY transponder_number`

	rows, err := r.db.Query(query, rcModelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racerModels []models.RacerModel
	for rows.Next() {
		var racerModel models.RacerModel

		err := rows.Scan(
			&racerModel.ID,
			&racerModel.RacerID,
			&racerModel.RCModelID,
			&racerModel.TransponderNumber,
			&racerModel.TransponderType,
			&racerModel.IsActive,
			&racerModel.CreatedAt,
			&racerModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		racerModels = append(racerModels, racerModel)
	}

	return racerModels, nil
}