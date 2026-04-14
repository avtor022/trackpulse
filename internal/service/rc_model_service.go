package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// RCModelService handles business logic for RC models
type RCModelService struct {
	modelRepo   *repository.RCModelRepository
	brandRepo   *repository.RCModelBrandRepository
	scaleRepo   *repository.RCModelScaleRepository
	typeRepo    *repository.RCModelTypeRepository
	trackRepo   *repository.RCModelTrackRepository
}

// NewRCModelService creates a new RC model service
func NewRCModelService(modelRepo *repository.RCModelRepository, brandRepo *repository.RCModelBrandRepository, scaleRepo *repository.RCModelScaleRepository, typeRepo *repository.RCModelTypeRepository, trackRepo *repository.RCModelTrackRepository) *RCModelService {
	return &RCModelService{modelRepo: modelRepo, brandRepo: brandRepo, scaleRepo: scaleRepo, typeRepo: typeRepo, trackRepo: trackRepo}
}

// GetAllModels returns all RC models
func (s *RCModelService) GetAllModels() ([]models.RCModel, error) {
	return s.modelRepo.GetAll()
}

// GetModelByID returns an RC model by ID
func (s *RCModelService) GetModelByID(id string) (*models.RCModel, error) {
	return s.modelRepo.GetByID(id)
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

	// Ensure brand exists in the brands table
	_, err := s.brandRepo.GetOrCreate(model.Brand)
	if err != nil {
		return fmt.Errorf("failed to ensure brand exists: %w", err)
	}

	// Ensure scale exists in the scales table
	_, err = s.scaleRepo.GetOrCreate(model.Scale)
	if err != nil {
		return fmt.Errorf("failed to ensure scale exists: %w", err)
	}

	// Ensure model type exists in the types table
	_, err = s.typeRepo.GetOrCreate(model.ModelType)
	if err != nil {
		return fmt.Errorf("failed to ensure model type exists: %w", err)
	}

	// Generate UUID
	model.ID = uuid.New().String()

	return s.modelRepo.Create(model)
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
	existing, err := s.modelRepo.GetByID(model.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("model not found")
	}

	// Ensure brand exists in the brands table
	_, err = s.brandRepo.GetOrCreate(model.Brand)
	if err != nil {
		return fmt.Errorf("failed to ensure brand exists: %w", err)
	}

	// Ensure scale exists in the scales table
	_, err = s.scaleRepo.GetOrCreate(model.Scale)
	if err != nil {
		return fmt.Errorf("failed to ensure scale exists: %w", err)
	}

	// Ensure model type exists in the types table
	_, err = s.typeRepo.GetOrCreate(model.ModelType)
	if err != nil {
		return fmt.Errorf("failed to ensure model type exists: %w", err)
	}

	return s.modelRepo.Update(model)
}

// DeleteModel deletes an RC model by ID
func (s *RCModelService) DeleteModel(id string) error {
	if id == "" {
		return fmt.Errorf("model ID is required")
	}
	return s.modelRepo.Delete(id)
}

// GetModelCount returns total number of RC models
func (s *RCModelService) GetModelCount() (int, error) {
	return s.modelRepo.Count()
}

// GetAllBrands returns all RC model brands
func (s *RCModelService) GetAllBrands() ([]models.RCModelBrand, error) {
	return s.brandRepo.GetAll()
}

// GetAllModelNames returns all unique model names
func (s *RCModelService) GetAllModelNames() ([]string, error) {
	return s.modelRepo.GetAllModelNames()
}

// AddBrand adds a new brand to the database
func (s *RCModelService) AddBrand(name string) error {
	if name == "" {
		return fmt.Errorf("brand name is required")
	}

	// Check if brand already exists
	existing, err := s.brandRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing brand: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("brand '%s' already exists", name)
	}

	// Create new brand
	_, err = s.brandRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}

	return nil
}

// DeleteBrand deletes a brand from the database
func (s *RCModelService) DeleteBrand(name string) error {
	if name == "" {
		return fmt.Errorf("brand name is required")
	}

	// Check if brand is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.Brand == name {
			return fmt.Errorf("cannot delete brand '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete brand
	err = s.brandRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}

// GetAllScales returns all RC model scales
func (s *RCModelService) GetAllScales() ([]models.RCModelScale, error) {
	return s.scaleRepo.GetAll()
}

// AddScale adds a new scale to the database
func (s *RCModelService) AddScale(name string) error {
	if name == "" {
		return fmt.Errorf("scale name is required")
	}

	// Check if scale already exists
	existing, err := s.scaleRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing scale: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("scale '%s' already exists", name)
	}

	// Create new scale
	_, err = s.scaleRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create scale: %w", err)
	}

	return nil
}

// DeleteScale deletes a scale from the database
func (s *RCModelService) DeleteScale(name string) error {
	if name == "" {
		return fmt.Errorf("scale name is required")
	}

	// Check if scale is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.Scale == name {
			return fmt.Errorf("cannot delete scale '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete scale
	err = s.scaleRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete scale: %w", err)
	}

	return nil
}

// GetAllModelTypes returns all RC model types
func (s *RCModelService) GetAllModelTypes() ([]models.RCModelType, error) {
	return s.typeRepo.GetAll()
}

// AddModelType adds a new model type to the database
func (s *RCModelService) AddModelType(name string) error {
	if name == "" {
		return fmt.Errorf("model type name is required")
	}

	// Check if model type already exists
	existing, err := s.typeRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing model type: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("model type '%s' already exists", name)
	}

	// Create new model type
	_, err = s.typeRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create model type: %w", err)
	}

	return nil
}

// DeleteModelType deletes a model type from the database
func (s *RCModelService) DeleteModelType(name string) error {
	if name == "" {
		return fmt.Errorf("model type name is required")
	}

	// Check if model type is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.ModelType == name {
			return fmt.Errorf("cannot delete model type '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete model type
	err = s.typeRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete model type: %w", err)
	}

	return nil
}

// GetAllTracks returns all RC model tracks
func (s *RCModelService) GetAllTracks() ([]models.RCModelTrack, error) {
	return s.trackRepo.GetAll()
}

// AddTrack adds a new track to the database
func (s *RCModelService) AddTrack(name string) error {
	if name == "" {
		return fmt.Errorf("track name is required")
	}

	// Check if track already exists
	existing, err := s.trackRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing track: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("track '%s' already exists", name)
	}

	// Create new track
	_, err = s.trackRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}

	return nil
}

// DeleteTrack deletes a track from the database
func (s *RCModelService) DeleteTrack(name string) error {
	if name == "" {
		return fmt.Errorf("track name is required")
	}

	// Delete track
	err := s.trackRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	return nil
}
