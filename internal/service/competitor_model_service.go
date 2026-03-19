package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// CompetitorModelService handles business logic for competitor models (transponders)
type CompetitorModelService struct {
	repo            *repository.CompetitorModelRepository
	competitorRepo  *repository.CompetitorRepository
	modelRepo       *repository.RCModelRepository
}

// NewCompetitorModelService creates a new competitor model service
func NewCompetitorModelService(repo *repository.CompetitorModelRepository, competitorRepo *repository.CompetitorRepository, modelRepo *repository.RCModelRepository) *CompetitorModelService {
	return &CompetitorModelService{repo: repo, competitorRepo: competitorRepo, modelRepo: modelRepo}
}

// GetAllCompetitorModels returns all competitor models
func (s *CompetitorModelService) GetAllCompetitorModels() ([]models.CompetitorModel, error) {
	return s.repo.GetAll()
}

// GetCompetitorModelByID returns a competitor model by ID
func (s *CompetitorModelService) GetCompetitorModelByID(id string) (*models.CompetitorModel, error) {
	return s.repo.GetByID(id)
}

// CreateCompetitorModel creates a new competitor model with validation
func (s *CompetitorModelService) CreateCompetitorModel(cm *models.CompetitorModel) error {
	// Validate required fields
	if cm.CompetitorID == "" {
		return fmt.Errorf("competitor is required")
	}
	if cm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if cm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if competitor exists
	competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
	if err != nil {
		return fmt.Errorf("failed to get competitor: %w", err)
	}
	if competitor == nil {
		return fmt.Errorf("competitor not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(cm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number already exists
	existing, err := s.repo.GetByTransponderNumber(cm.TransponderNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing transponder: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("transponder number '%s' already exists", cm.TransponderNumber)
	}

	// Set default values
	if cm.TransponderType == "" {
		cm.TransponderType = "RFID"
	}
	cm.IsActive = true

	// Generate UUID
	cm.ID = uuid.New().String()

	return s.repo.Create(cm)
}

// UpdateCompetitorModel updates an existing competitor model with validation
func (s *CompetitorModelService) UpdateCompetitorModel(cm *models.CompetitorModel) error {
	// Validate required fields
	if cm.ID == "" {
		return fmt.Errorf("competitor model ID is required")
	}
	if cm.CompetitorID == "" {
		return fmt.Errorf("competitor is required")
	}
	if cm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if cm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if competitor model exists
	existing, err := s.repo.GetByID(cm.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing competitor model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("competitor model not found")
	}

	// Check if competitor exists
	competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
	if err != nil {
		return fmt.Errorf("failed to get competitor: %w", err)
	}
	if competitor == nil {
		return fmt.Errorf("competitor not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(cm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number is taken by another competitor model
	if existing.TransponderNumber != cm.TransponderNumber {
		transponderExists, err := s.repo.GetByTransponderNumber(cm.TransponderNumber)
		if err != nil {
			return fmt.Errorf("failed to check transponder number: %w", err)
		}
		if transponderExists != nil && transponderExists.ID != cm.ID {
			return fmt.Errorf("transponder number '%s' already exists", cm.TransponderNumber)
		}
	}

	return s.repo.Update(cm)
}

// DeleteCompetitorModel deletes a competitor model by ID
func (s *CompetitorModelService) DeleteCompetitorModel(id string) error {
	if id == "" {
		return fmt.Errorf("competitor model ID is required")
	}
	return s.repo.Delete(id)
}

// GetCompetitorModelCount returns total number of competitor models
func (s *CompetitorModelService) GetCompetitorModelCount() (int, error) {
	return s.repo.Count()
}
