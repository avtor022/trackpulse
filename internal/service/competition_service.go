package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitionRepositoryInterface defines the interface for competition data access
type CompetitionRepositoryInterface interface {
	GetAll() ([]models.Competition, error)
	GetByID(id string) (*models.Competition, error)
	Create(competition *models.Competition) error
	Update(competition *models.Competition) error
	Delete(id string) error
	Count() (int, error)
	GetByStatus(status string) ([]models.Competition, error)
}

// RCModelTypeRepositoryInterface defines the interface for RC model type data access
type RCModelTypeRepositoryInterface interface {
	GetAll() ([]models.RCModelType, error)
	GetByName(name string) (*models.RCModelType, error)
	Create(name string) (*models.RCModelType, error)
	GetOrCreate(name string) (*models.RCModelType, error)
	Delete(name string) error
}

// RCModelScaleRepositoryInterface defines the interface for RC model scale data access
type RCModelScaleRepositoryInterface interface {
	GetAll() ([]models.RCModelScale, error)
	GetByName(name string) (*models.RCModelScale, error)
	Create(name string) (*models.RCModelScale, error)
	GetOrCreate(name string) (*models.RCModelScale, error)
	Delete(name string) error
}

// CompetitionService handles business logic for competitions
type CompetitionService struct {
	repo       CompetitionRepositoryInterface
	modelTypes RCModelTypeRepositoryInterface
	scales     RCModelScaleRepositoryInterface
	tracks     CompetitionTrackRepositoryInterface
}

// CompetitionTrackRepositoryInterface defines the interface for competition track data access
type CompetitionTrackRepositoryInterface interface {
	GetAll() ([]models.CompetitionTrack, error)
	GetByName(name string) (*models.CompetitionTrack, error)
	Create(name string) (*models.CompetitionTrack, error)
	GetOrCreate(name string) (*models.CompetitionTrack, error)
	Delete(name string) error
}

// NewCompetitionService creates a new competition service
func NewCompetitionService(repo CompetitionRepositoryInterface, modelTypes RCModelTypeRepositoryInterface, scales RCModelScaleRepositoryInterface, tracks CompetitionTrackRepositoryInterface) *CompetitionService {
	return &CompetitionService{repo: repo, modelTypes: modelTypes, scales: scales, tracks: tracks}
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
	return s.scales.GetAll()
}

// GetAllCompetitionTracks returns all competition tracks for competition selection
func (s *CompetitionService) GetAllCompetitionTracks() ([]models.CompetitionTrack, error) {
	return s.tracks.GetAll()
}

// AddCompetitionTrack adds a new competition track
func (s *CompetitionService) AddCompetitionTrack(name string) error {
	_, err := s.tracks.Create(name)
	return err
}

// DeleteCompetitionTrack deletes a competition track by name
func (s *CompetitionService) DeleteCompetitionTrack(name string) error {
	return s.tracks.Delete(name)
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
