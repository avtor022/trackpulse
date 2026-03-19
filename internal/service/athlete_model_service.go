package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// AthleteModelService handles business logic for athlete models (transponders)
type AthleteModelService struct {
	repo       *repository.AthleteModelRepository
	athleteRepo  *repository.AthleteRepository
	modelRepo  *repository.RCModelRepository
}

// NewAthleteModelService creates a new athlete model service
func NewAthleteModelService(repo *repository.AthleteModelRepository, athleteRepo *repository.AthleteRepository, modelRepo *repository.RCModelRepository) *AthleteModelService {
	return &AthleteModelService{repo: repo, athleteRepo: athleteRepo, modelRepo: modelRepo}
}

// GetAllAthleteModels returns all athlete models
func (s *AthleteModelService) GetAllAthleteModels() ([]models.AthleteModel, error) {
	return s.repo.GetAll()
}

// GetAthleteModelByID returns an athlete model by ID
func (s *AthleteModelService) GetAthleteModelByID(id string) (*models.AthleteModel, error) {
	return s.repo.GetByID(id)
}

// CreateAthleteModel creates a new athlete model with validation
func (s *AthleteModelService) CreateAthleteModel(am *models.AthleteModel) error {
	// Validate required fields
	if am.AthleteID == "" {
		return fmt.Errorf("athlete is required")
	}
	if am.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if am.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if athlete exists
	athlete, err := s.athleteRepo.GetByID(am.AthleteID)
	if err != nil {
		return fmt.Errorf("failed to get athlete: %w", err)
	}
	if athlete == nil {
		return fmt.Errorf("athlete not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(am.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number already exists
	existing, err := s.repo.GetByTransponderNumber(am.TransponderNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing transponder: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("transponder number '%s' already exists", am.TransponderNumber)
	}

	// Set default values
	if am.TransponderType == "" {
		am.TransponderType = "RFID"
	}
	am.IsActive = true

	// Generate UUID
	am.ID = uuid.New().String()

	return s.repo.Create(am)
}

// UpdateAthleteModel updates an existing athlete model with validation
func (s *AthleteModelService) UpdateAthleteModel(am *models.AthleteModel) error {
	// Validate required fields
	if am.ID == "" {
		return fmt.Errorf("athlete model ID is required")
	}
	if am.AthleteID == "" {
		return fmt.Errorf("athlete is required")
	}
	if am.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if am.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if athlete model exists
	existing, err := s.repo.GetByID(am.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing athlete model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("athlete model not found")
	}

	// Check if athlete exists
	athlete, err := s.athleteRepo.GetByID(am.AthleteID)
	if err != nil {
		return fmt.Errorf("failed to get athlete: %w", err)
	}
	if athlete == nil {
		return fmt.Errorf("athlete not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(am.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number is taken by another athlete model
	if existing.TransponderNumber != am.TransponderNumber {
		transponderExists, err := s.repo.GetByTransponderNumber(am.TransponderNumber)
		if err != nil {
			return fmt.Errorf("failed to check transponder number: %w", err)
		}
		if transponderExists != nil && transponderExists.ID != am.ID {
			return fmt.Errorf("transponder number '%s' already exists", am.TransponderNumber)
		}
	}

	return s.repo.Update(am)
}

// DeleteAthleteModel deletes an athlete model by ID
func (s *AthleteModelService) DeleteAthleteModel(id string) error {
	if id == "" {
		return fmt.Errorf("athlete model ID is required")
	}
	return s.repo.Delete(id)
}

// GetAthleteModelCount returns total number of athlete models
func (s *AthleteModelService) GetAthleteModelCount() (int, error) {
	return s.repo.Count()
}
