package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// RacerService handles business logic for racers
type RacerService struct {
	repo *repository.RacerRepository
}

// NewRacerService creates a new racer service
func NewRacerService(repo *repository.RacerRepository) *RacerService {
	return &RacerService{repo: repo}
}

// GetAllRacers returns all racers
func (s *RacerService) GetAllRacers() ([]models.Racer, error) {
	return s.repo.GetAll()
}

// GetRacerByID returns a racer by ID
func (s *RacerService) GetRacerByID(id string) (*models.Racer, error) {
	return s.repo.GetByID(id)
}

// CreateRacer creates a new racer with validation
func (s *RacerService) CreateRacer(racer *models.Racer) error {
	// Validate required fields
	if racer.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if racer.RacerNumber <= 0 {
		return fmt.Errorf("racer number must be positive")
	}

	// Check if racer number already exists
	existing, err := s.repo.GetByNumber(racer.RacerNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing racer: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("racer number %d already exists", racer.RacerNumber)
	}

	// Generate UUID
	racer.ID = uuid.New().String()

	return s.repo.Create(racer)
}

// UpdateRacer updates an existing racer with validation
func (s *RacerService) UpdateRacer(racer *models.Racer) error {
	// Validate required fields
	if racer.ID == "" {
		return fmt.Errorf("racer ID is required")
	}
	if racer.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if racer.RacerNumber <= 0 {
		return fmt.Errorf("racer number must be positive")
	}

	// Check if racer exists
	existing, err := s.repo.GetByID(racer.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing racer: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("racer not found")
	}

	// Check if racer number is taken by another racer
	if existing.RacerNumber != racer.RacerNumber {
		numberExists, err := s.repo.GetByNumber(racer.RacerNumber)
		if err != nil {
			return fmt.Errorf("failed to check racer number: %w", err)
		}
		if numberExists != nil && numberExists.ID != racer.ID {
			return fmt.Errorf("racer number %d already exists", racer.RacerNumber)
		}
	}

	return s.repo.Update(racer)
}

// DeleteRacer deletes a racer by ID
func (s *RacerService) DeleteRacer(id string) error {
	if id == "" {
		return fmt.Errorf("racer ID is required")
	}
	return s.repo.Delete(id)
}

// GetRacerCount returns total number of racers
func (s *RacerService) GetRacerCount() (int, error) {
	return s.repo.Count()
}
