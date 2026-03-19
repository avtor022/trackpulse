package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// CompetitorService handles business logic for competitors
type CompetitorService struct {
	repo *repository.CompetitorRepository
}

// NewCompetitorService creates a new competitor service
func NewCompetitorService(repo *repository.CompetitorRepository) *CompetitorService {
	return &CompetitorService{repo: repo}
}

// GetAllCompetitors returns all competitors
func (s *CompetitorService) GetAllCompetitors() ([]models.Competitor, error) {
	return s.repo.GetAll()
}

// GetCompetitorByID returns a competitor by ID
func (s *CompetitorService) GetCompetitorByID(id string) (*models.Competitor, error) {
	return s.repo.GetByID(id)
}

// CreateCompetitor creates a new competitor with validation
func (s *CompetitorService) CreateCompetitor(competitor *models.Competitor) error {
	// Validate required fields
	if competitor.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if competitor.CompetitorNumber <= 0 {
		return fmt.Errorf("competitor number must be positive")
	}

	// Check if competitor number already exists
	existing, err := s.repo.GetByNumber(competitor.CompetitorNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing competitor: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("competitor number %d already exists", competitor.CompetitorNumber)
	}

	// Generate UUID
	competitor.ID = uuid.New().String()

	return s.repo.Create(competitor)
}

// UpdateCompetitor updates an existing competitor with validation
func (s *CompetitorService) UpdateCompetitor(competitor *models.Competitor) error {
	// Validate required fields
	if competitor.ID == "" {
		return fmt.Errorf("competitor ID is required")
	}
	if competitor.FullName == "" {
		return fmt.Errorf("full name is required")
	}
	if competitor.CompetitorNumber <= 0 {
		return fmt.Errorf("competitor number must be positive")
	}

	// Check if competitor exists
	existing, err := s.repo.GetByID(competitor.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing competitor: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("competitor not found")
	}

	// Check if competitor number is taken by another competitor
	if existing.CompetitorNumber != competitor.CompetitorNumber {
		numberExists, err := s.repo.GetByNumber(competitor.CompetitorNumber)
		if err != nil {
			return fmt.Errorf("failed to check competitor number: %w", err)
		}
		if numberExists != nil && numberExists.ID != competitor.ID {
			return fmt.Errorf("competitor number %d already exists", competitor.CompetitorNumber)
		}
	}

	return s.repo.Update(competitor)
}

// DeleteCompetitor deletes a competitor by ID
func (s *CompetitorService) DeleteCompetitor(id string) error {
	if id == "" {
		return fmt.Errorf("competitor ID is required")
	}
	return s.repo.Delete(id)
}

// GetCompetitorCount returns total number of competitors
func (s *CompetitorService) GetCompetitorCount() (int, error) {
	return s.repo.Count()
}
