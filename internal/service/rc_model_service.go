package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// RCModelService handles business logic for RC models
type RCModelService struct {
	repo *repository.RCModelRepository
}

// NewRCModelService creates a new RC model service
func NewRCModelService(repo *repository.RCModelRepository) *RCModelService {
	return &RCModelService{repo: repo}
}

// GetAllModels returns all RC models
func (s *RCModelService) GetAllModels() ([]models.RCModel, error) {
	return s.repo.GetAll()
}

// GetModelByID returns an RC model by ID
func (s *RCModelService) GetModelByID(id string) (*models.RCModel, error) {
	return s.repo.GetByID(id)
}

// CreateModel creates a new RC model with validation
func (s *RCModelService) CreateModel(model *models.RCModel) error {
	// Validate required fields
	if model.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if model.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	if model.Scale == "" {
		return fmt.Errorf("scale is required")
	}
	if model.ModelType == "" {
		return fmt.Errorf("model type is required")
	}

	// Generate UUID
	model.ID = uuid.New().String()

	return s.repo.Create(model)
}

// UpdateModel updates an existing RC model with validation
func (s *RCModelService) UpdateModel(model *models.RCModel) error {
	// Validate required fields
	if model.ID == "" {
		return fmt.Errorf("model ID is required")
	}
	if model.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if model.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	if model.Scale == "" {
		return fmt.Errorf("scale is required")
	}
	if model.ModelType == "" {
		return fmt.Errorf("model type is required")
	}

	// Check if model exists
	existing, err := s.repo.GetByID(model.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("model not found")
	}

	return s.repo.Update(model)
}

// DeleteModel deletes an RC model by ID
func (s *RCModelService) DeleteModel(id string) error {
	if id == "" {
		return fmt.Errorf("model ID is required")
	}
	return s.repo.Delete(id)
}

// GetModelCount returns total number of RC models
func (s *RCModelService) GetModelCount() (int, error) {
	return s.repo.Count()
}

// GetUniqueBrands returns a list of unique brand names from all models
func (s *RCModelService) GetUniqueBrands() ([]string, error) {
	return s.repo.GetUniqueBrands()
}
