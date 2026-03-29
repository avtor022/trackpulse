package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// CompetitionService handles business logic for competitions
type CompetitionService struct {
	repo       *repository.CompetitionRepository
	modelTypes *repository.RCModelTypeRepository
	scaleRepo  *repository.RCModelScaleRepository
}

// NewCompetitionService creates a new competition service
func NewCompetitionService(repo *repository.CompetitionRepository, modelTypes *repository.RCModelTypeRepository, scaleRepo *repository.RCModelScaleRepository) *CompetitionService {
	return &CompetitionService{repo: repo, modelTypes: modelTypes, scaleRepo: scaleRepo}
}

// GetAllCompetitions returns all competitions
func (s *CompetitionService) GetAllCompetitions() ([]models.Competition, error) {
	return s.repo.GetAll()
}

// GetCompetitionByID returns a competition by ID
func (s *CompetitionService) GetCompetitionByID(id string) (*models.Competition, error) {
	return s.repo.GetByID(id)
}

// GetAllModelTypes returns all RC model types for competition selection
func (s *CompetitionService) GetAllModelTypes() ([]models.RCModelType, error) {
	return s.modelTypes.GetAll()
}

// GetAllModelScales returns all RC model scales for competition selection
func (s *CompetitionService) GetAllModelScales() ([]models.RCModelScale, error) {
	return s.scaleRepo.GetAll()
}

// CreateCompetition creates a new competition with validation
func (s *CompetitionService) CreateCompetition(competition *models.Competition) error {
	// Validate required fields
	if competition.CompetitionTitle == "" {
		return fmt.Errorf("competition title is required")
	}
	if competition.CompetitionType == "" {
		return fmt.Errorf("competition type is required")
	}

	// Generate UUID
	competition.ID = uuid.New().String()

	// Set default status if not provided
	if competition.Status == "" {
		competition.Status = "scheduled"
	}

	return s.repo.Create(competition)
}

// UpdateCompetition updates an existing competition with validation
func (s *CompetitionService) UpdateCompetition(competition *models.Competition) error {
	// Validate required fields
	if competition.ID == "" {
		return fmt.Errorf("competition ID is required")
	}
	if competition.CompetitionTitle == "" {
		return fmt.Errorf("competition title is required")
	}
	if competition.CompetitionType == "" {
		return fmt.Errorf("competition type is required")
	}

	// Check if competition exists
	existing, err := s.repo.GetByID(competition.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing competition: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("competition not found")
	}

	return s.repo.Update(competition)
}

// DeleteCompetition deletes a competition by ID
func (s *CompetitionService) DeleteCompetition(id string) error {
	if id == "" {
		return fmt.Errorf("competition ID is required")
	}
	return s.repo.Delete(id)
}

// GetCompetitionCount returns total number of competitions
func (s *CompetitionService) GetCompetitionCount() (int, error) {
	return s.repo.Count()
}

// GetCompetitionsByStatus returns competitions filtered by status
func (s *CompetitionService) GetCompetitionsByStatus(status string) ([]models.Competition, error) {
	return s.repo.GetByStatus(status)
}
