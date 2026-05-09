package service

import (
	"fmt"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitionParticipantRepositoryInterface defines the interface for competition participant data access
type CompetitionParticipantRepositoryInterface interface {
	GetAll() ([]models.CompetitionParticipant, error)
	GetByID(id string) (*models.CompetitionParticipant, error)
	GetByCompetitionID(competitionID string) ([]models.CompetitionParticipant, error)
	Create(participant *models.CompetitionParticipant) error
	Update(participant *models.CompetitionParticipant) error
	Delete(id string) error
	DeleteByCompetitionAndCompetitorModel(competitionID, competitorModelID string) error
	Exists(competitionID, competitorModelID string) (bool, error)
	Count() (int, error)
}

// CompetitorModelServiceInterface defines the interface for competitor model data access
type CompetitorModelServiceInterface interface {
	GetAllCompetitorModels() ([]models.CompetitorModel, error)
	GetCompetitorModelByID(id string) (*models.CompetitorModel, error)
}

// CompetitionServiceInterface defines the interface for competition data access
type CompetitionServiceInterface interface {
	GetAllCompetitions() ([]models.Competition, error)
	GetCompetitionByID(id string) (*models.Competition, error)
}

// CompetitionParticipantService handles business logic for competition participants
type CompetitionParticipantService struct {
	repo                 CompetitionParticipantRepositoryInterface
	competitorModelService CompetitorModelServiceInterface
	competitionService     CompetitionServiceInterface
}

// NewCompetitionParticipantService creates a new competition participant service
func NewCompetitionParticipantService(repo CompetitionParticipantRepositoryInterface, competitorModelService CompetitorModelServiceInterface, competitionService CompetitionServiceInterface) *CompetitionParticipantService {
	return &CompetitionParticipantService{
		repo:                 repo,
		competitorModelService: competitorModelService,
		competitionService:     competitionService,
	}
}

// GetAllParticipants returns all competition participants
func (s *CompetitionParticipantService) GetAllParticipants() ([]models.CompetitionParticipant, error) {
	return s.repo.GetAll()
}

// GetParticipantByID returns a participant by ID
func (s *CompetitionParticipantService) GetParticipantByID(id string) (*models.CompetitionParticipant, error) {
	return s.repo.GetByID(id)
}

// GetParticipantsByCompetitionID returns all participants for a specific competition
func (s *CompetitionParticipantService) GetParticipantsByCompetitionID(competitionID string) ([]models.CompetitionParticipant, error) {
	return s.repo.GetByCompetitionID(competitionID)
}

// AddParticipant adds a new participant to a competition
func (s *CompetitionParticipantService) AddParticipant(participant *models.CompetitionParticipant) error {
	// Validate required fields
	if participant.CompetitionID == "" {
		return fmt.Errorf("competition ID is required")
	}
	if participant.CompetitorModelID == "" {
		return fmt.Errorf("competitor model ID is required")
	}

	// Check if competition exists
	competition, err := s.competitionService.GetCompetitionByID(participant.CompetitionID)
	if err != nil {
		return fmt.Errorf("failed to get competition: %w", err)
	}
	if competition == nil {
		return fmt.Errorf("competition not found")
	}

	// Check if competitor model exists
	competitorModel, err := s.competitorModelService.GetCompetitorModelByID(participant.CompetitorModelID)
	if err != nil {
		return fmt.Errorf("failed to get competitor model: %w", err)
	}
	if competitorModel == nil {
		return fmt.Errorf("competitor model not found")
	}

	// Check if participant already exists
	exists, err := s.repo.Exists(participant.CompetitionID, participant.CompetitorModelID)
	if err != nil {
		return fmt.Errorf("failed to check existing participant: %w", err)
	}
	if exists {
		return fmt.Errorf("participant already exists for this competition")
	}

	// Generate UUID
	participant.ID = uuid.New().String()

	// Set default values
	if participant.GridPosition == nil {
		defaultPos := 0
		participant.GridPosition = &defaultPos
	}
	participant.IsFinished = false
	participant.Disqualified = false

	return s.repo.Create(participant)
}

// AddParticipantsBulk adds multiple participants to a competition
// Returns list of successfully added participant IDs and list of errors
func (s *CompetitionParticipantService) AddParticipantsBulk(competitionID string, competitorModelIDs []string) ([]string, []error) {
	var addedIDs []string
	var errors []error

	// Check if competition exists
	competition, err := s.competitionService.GetCompetitionByID(competitionID)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get competition: %w", err)}
	}
	if competition == nil {
		return nil, []error{fmt.Errorf("competition not found")}
	}

	for _, cmID := range competitorModelIDs {
		// Check if competitor model exists
		competitorModel, err := s.competitorModelService.GetCompetitorModelByID(cmID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get competitor model %s: %w", cmID, err))
			continue
		}
		if competitorModel == nil {
			errors = append(errors, fmt.Errorf("competitor model %s not found", cmID))
			continue
		}

		// Check if participant already exists
		exists, err := s.repo.Exists(competitionID, cmID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to check existing participant for %s: %w", cmID, err))
			continue
		}
		if exists {
			errors = append(errors, fmt.Errorf("participant with model %s already exists for this competition", cmID))
			continue
		}

		// Create participant
		participant := &models.CompetitionParticipant{
			ID:                uuid.New().String(),
			CompetitionID:     competitionID,
			CompetitorModelID: cmID,
			GridPosition:      nil,
			IsFinished:        false,
			Disqualified:      false,
			DNFReason:         "",
		}

		if err := s.repo.Create(participant); err != nil {
			errors = append(errors, fmt.Errorf("failed to create participant for %s: %w", cmID, err))
			continue
		}

		addedIDs = append(addedIDs, participant.ID)
	}

	return addedIDs, errors
}

// RemoveParticipant removes a participant from a competition
func (s *CompetitionParticipantService) RemoveParticipant(id string) error {
	if id == "" {
		return fmt.Errorf("participant ID is required")
	}
	return s.repo.Delete(id)
}

// RemoveParticipantByCompetitionAndModel removes a participant by competition and competitor model IDs
func (s *CompetitionParticipantService) RemoveParticipantByCompetitionAndModel(competitionID, competitorModelID string) error {
	if competitionID == "" || competitorModelID == "" {
		return fmt.Errorf("competition ID and competitor model ID are required")
	}
	return s.repo.DeleteByCompetitionAndCompetitorModel(competitionID, competitorModelID)
}

// GetParticipantCount returns total number of participants
func (s *CompetitionParticipantService) GetParticipantCount() (int, error) {
	return s.repo.Count()
}
