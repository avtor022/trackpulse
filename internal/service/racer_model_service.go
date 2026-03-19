package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// RacerModelService handles business logic for racer models (transponders)
type RacerModelService struct {
	repo       *repository.RacerModelRepository
	racerRepo  *repository.RacerRepository
	modelRepo  *repository.RCModelRepository
}

// NewRacerModelService creates a new racer model service
func NewRacerModelService(repo *repository.RacerModelRepository, racerRepo *repository.RacerRepository, modelRepo *repository.RCModelRepository) *RacerModelService {
	return &RacerModelService{repo: repo, racerRepo: racerRepo, modelRepo: modelRepo}
}

// GetAllRacerModels returns all racer models
func (s *RacerModelService) GetAllRacerModels() ([]models.RacerModel, error) {
	return s.repo.GetAll()
}

// GetRacerModelByID returns a racer model by ID
func (s *RacerModelService) GetRacerModelByID(id string) (*models.RacerModel, error) {
	return s.repo.GetByID(id)
}

// CreateRacerModel creates a new racer model with validation
func (s *RacerModelService) CreateRacerModel(rm *models.RacerModel) error {
	// Validate required fields
	if rm.RacerID == "" {
		return fmt.Errorf("racer is required")
	}
	if rm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if rm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if racer exists
	racer, err := s.racerRepo.GetByID(rm.RacerID)
	if err != nil {
		return fmt.Errorf("failed to get racer: %w", err)
	}
	if racer == nil {
		return fmt.Errorf("racer not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(rm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number already exists
	existing, err := s.repo.GetByTransponderNumber(rm.TransponderNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing transponder: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("transponder number '%s' already exists", rm.TransponderNumber)
	}

	// Set default values
	if rm.TransponderType == "" {
		rm.TransponderType = "RFID"
	}
	rm.IsActive = true

	// Generate UUID
	rm.ID = uuid.New().String()

	return s.repo.Create(rm)
}

// UpdateRacerModel updates an existing racer model with validation
func (s *RacerModelService) UpdateRacerModel(rm *models.RacerModel) error {
	// Validate required fields
	if rm.ID == "" {
		return fmt.Errorf("racer model ID is required")
	}
	if rm.RacerID == "" {
		return fmt.Errorf("racer is required")
	}
	if rm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if rm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if racer model exists
	existing, err := s.repo.GetByID(rm.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing racer model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("racer model not found")
	}

	// Check if racer exists
	racer, err := s.racerRepo.GetByID(rm.RacerID)
	if err != nil {
		return fmt.Errorf("failed to get racer: %w", err)
	}
	if racer == nil {
		return fmt.Errorf("racer not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(rm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number is taken by another racer model
	if existing.TransponderNumber != rm.TransponderNumber {
		transponderExists, err := s.repo.GetByTransponderNumber(rm.TransponderNumber)
		if err != nil {
			return fmt.Errorf("failed to check transponder number: %w", err)
		}
		if transponderExists != nil && transponderExists.ID != rm.ID {
			return fmt.Errorf("transponder number '%s' already exists", rm.TransponderNumber)
		}
	}

	return s.repo.Update(rm)
}

// DeleteRacerModel deletes a racer model by ID
func (s *RacerModelService) DeleteRacerModel(id string) error {
	if id == "" {
		return fmt.Errorf("racer model ID is required")
	}
	return s.repo.Delete(id)
}

// GetRacerModelCount returns total number of racer models
func (s *RacerModelService) GetRacerModelCount() (int, error) {
	return s.repo.Count()
}
