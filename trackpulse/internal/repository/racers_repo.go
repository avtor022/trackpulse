package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RacersRepository interface {
	GetAll() ([]models.Racer, error)
	GetByID(id string) (*models.Racer, error)
	Create(racer *models.Racer) error
	Update(racer *models.Racer) error
	Delete(id string) error
}

type racersRepo struct {
	db *sql.DB
}

func NewRacersRepository(db *sql.DB) RacersRepository {
	return &racersRepo{db: db}
}

func (r *racersRepo) GetAll() ([]models.Racer, error) {
	query := `SELECT id, racer_number, full_name, birthday, country, city, 
	                 rating, created_at, updated_at 
	          FROM racers ORDER BY racer_number`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var racers []models.Racer
	for rows.Next() {
		var racer models.Racer
		var birthday, country, city *string

		err := rows.Scan(
			&racer.ID,
			&racer.RacerNumber,
			&racer.FullName,
			&birthday,
			&country,
			&city,
			&racer.Rating,
			&racer.CreatedAt,
			&racer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		racer.Birthday = birthday
		racer.Country = country
		racer.City = city

		racers = append(racers, racer)
	}

	return racers, nil
}

func (r *racersRepo) GetByID(id string) (*models.Racer, error) {
	query := `SELECT id, racer_number, full_name, birthday, country, city, 
	                 rating, created_at, updated_at 
	          FROM racers WHERE id = ?`

	var racer models.Racer
	var birthday, country, city *string

	err := r.db.QueryRow(query, id).Scan(
		&racer.ID,
		&racer.RacerNumber,
		&racer.FullName,
		&birthday,
		&country,
		&city,
		&racer.Rating,
		&racer.CreatedAt,
		&racer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("racer not found")
		}
		return nil, err
	}

	// Handle nullable fields
	racer.Birthday = birthday
	racer.Country = country
	racer.City = city

	return &racer, nil
}

func (r *racersRepo) Create(racer *models.Racer) error {
	query := `INSERT INTO racers (id, racer_number, full_name, birthday, country, city, 
	                            rating, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		racer.ID,
		racer.RacerNumber,
		racer.FullName,
		racer.Birthday,
		racer.Country,
		racer.City,
		racer.Rating,
		racer.CreatedAt,
		racer.UpdatedAt,
	)
	return err
}

func (r *racersRepo) Update(racer *models.Racer) error {
	query := `UPDATE racers SET racer_number = ?, full_name = ?, birthday = ?, 
	                            country = ?, city = ?, rating = ?, updated_at = ? 
	          WHERE id = ?`

	_, err := r.db.Exec(query,
		racer.RacerNumber,
		racer.FullName,
		racer.Birthday,
		racer.Country,
		racer.City,
		racer.Rating,
		racer.UpdatedAt,
		racer.ID,
	)
	return err
}

func (r *racersRepo) Delete(id string) error {
	query := `DELETE FROM racers WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("racer not found")
	}

	return nil
}