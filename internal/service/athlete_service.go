package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// AthleteService handles business logic for athletes
type AthleteService struct {
	repo *repository.AthleteRepository
}

// NewAthleteService creates a new athlete service
func NewAthleteService(repo *repository.AthleteRepository) *AthleteService {
	return &AthleteService{repo: repo}
}

// GetAllAthletes returns all athletes
func (s *AthleteService) GetAllAthletes() ([]models.Athlete, error) {
	return s.repo.GetAll()
}

// GetAthleteByID returns an athlete by ID
func (s *AthleteService) GetAthleteByID(id string) (*models.Athlete, error) {
	return s.repo.GetByID(id)
}

// CreateAthlete creates a new athlete with validation
func (s *AthleteService) CreateAthlete(athlete *models.Athlete) error {
	// Validate required fields
	if athlete.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if athlete.RacerNumber <= 0 {
		return fmt.Errorf("racer number must be positive")
	}

	// Check if racer number already exists
	existing, err := s.repo.GetByNumber(athlete.RacerNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing athlete: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("racer number %d already exists", athlete.RacerNumber)
	}

	// Generate UUID
	athlete.ID = uuid.New().String()

	return s.repo.Create(athlete)
}

// UpdateAthlete updates an existing athlete with validation
func (s *AthleteService) UpdateAthlete(athlete *models.Athlete) error {
	// Validate required fields
	if athlete.ID == "" {
		return fmt.Errorf("athlete ID is required")
	}
	if athlete.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if athlete.RacerNumber <= 0 {
		return fmt.Errorf("racer number must be positive")
	}

	// Check if athlete exists
	existing, err := s.repo.GetByID(athlete.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing athlete: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("athlete not found")
	}

	// Check if racer number is taken by another athlete
	if existing.RacerNumber != athlete.RacerNumber {
		numberExists, err := s.repo.GetByNumber(athlete.RacerNumber)
		if err != nil {
			return fmt.Errorf("failed to check racer number: %w", err)
		}
		if numberExists != nil && numberExists.ID != athlete.ID {
			return fmt.Errorf("racer number %d already exists", athlete.RacerNumber)
		}
	}

	return s.repo.Update(athlete)
}

// DeleteAthlete deletes an athlete by ID
func (s *AthleteService) DeleteAthlete(id string) error {
	if id == "" {
		return fmt.Errorf("athlete ID is required")
	}
	return s.repo.Delete(id)
}

// GetAthleteCount returns total number of athletes
func (s *AthleteService) GetAthleteCount() (int, error) {
	return s.repo.Count()
}
