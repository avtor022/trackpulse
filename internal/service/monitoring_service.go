package service

import (
	"trackpulse/internal/models"
)

// MonitoringParticipantData represents complete participant data for monitoring display
type MonitoringParticipantData struct {
	TransponderActive        bool   // Transponder workability flag (true after first pass)
	CompetitorNumber         int    // Participant number
	FullName                 string // Full name of participant
	ModelName                string // RC model name
	ModelScale               string // RC model scale
	LapCount                 int    // Number of laps completed
	BestLapTimeMs            int    // Best lap time in milliseconds
	TotalCompetitionTimeMs   int    // Total competition time in milliseconds
	LastLapTimeMs            int    // Last lap time in milliseconds
}

// MonitoringService handles business logic for monitoring display
type MonitoringService struct {
	participantRepo     CompetitionParticipantRepositoryInterface
	competitorModelRepo CompetitorModelRepositoryInterface
	competitorRepo      CompetitorRepositoryInterface
	rcModelRepo         RCModelRepositoryInterface
	competitionLapsRepo CompetitionLapsRepositoryInterface
	lapService          *LapService
}

// CompetitorModelRepositoryInterface defines interface for competitor model data access
type CompetitorModelRepositoryInterface interface {
	GetByID(id string) (*models.CompetitorModel, error)
}

// RCModelRepositoryInterface defines interface for RC model data access
type RCModelRepositoryInterface interface {
	GetByID(id string) (*models.RCModel, error)
}

// CompetitionLapsRepositoryInterface defines interface for competition laps data access
type CompetitionLapsRepositoryInterface interface {
	GetByParticipantID(participantID string) (*models.CompetitionLaps, error)
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(
	participantRepo CompetitionParticipantRepositoryInterface,
	competitorModelRepo CompetitorModelRepositoryInterface,
	competitorRepo CompetitorRepositoryInterface,
	rcModelRepo RCModelRepositoryInterface,
	competitionLapsRepo CompetitionLapsRepositoryInterface,
	lapService *LapService,
) *MonitoringService {
	return &MonitoringService{
		participantRepo:     participantRepo,
		competitorModelRepo: competitorModelRepo,
		competitorRepo:      competitorRepo,
		rcModelRepo:         rcModelRepo,
		competitionLapsRepo: competitionLapsRepo,
		lapService:          lapService,
	}
}

// GetParticipantsForCompetition returns complete participant data for a competition
// This includes both persisted data from DB and real-time data from LapService cache
func (s *MonitoringService) GetParticipantsForCompetition(competitionID string) ([]MonitoringParticipantData, error) {
	// Get all participants for this competition
	participants, err := s.participantRepo.GetByCompetitionID(competitionID)
	if err != nil {
		return nil, err
	}

	var result []MonitoringParticipantData

	for _, p := range participants {
		data := MonitoringParticipantData{}

		// Get competitor model (transponder info)
		cm, err := s.competitorModelRepo.GetByID(p.CompetitorModelID)
		if err != nil || cm == nil {
			continue
		}

		// Get competitor info
		competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
		if err != nil || competitor == nil {
			continue
		}

		// Get RC model info
		rcModel, err := s.rcModelRepo.GetByID(cm.RCModelID)
		if err != nil || rcModel == nil {
			continue
		}

		// Set basic info
		data.CompetitorNumber = competitor.CompetitorNumber
		data.FullName = competitor.FullName
		data.ModelName = rcModel.ModelName
		data.ModelScale = rcModel.Scale

		// Check if transponder is active (has made at least one pass)
		// First check real-time cache from LapService
		if s.lapService != nil {
			results := s.lapService.GetParticipantResults()
			if lapData, exists := results[p.ID]; exists {
				// Transponder is active if we have any lap data
				data.TransponderActive = lapData.LapCount > 0
				data.LapCount = lapData.LapCount
				data.BestLapTimeMs = lapData.BestLapTimeMs
				data.LastLapTimeMs = lapData.LastLapTimeMs
				data.TotalCompetitionTimeMs = lapData.TotalTimeMs
			}
		}

		// If no real-time data, try to get persisted data from DB
		if data.LapCount == 0 {
			laps, err := s.competitionLapsRepo.GetByParticipantID(p.ID)
			if err == nil && laps != nil {
				data.TransponderActive = laps.NumberOfLaps > 0
				data.LapCount = laps.NumberOfLaps
				data.BestLapTimeMs = laps.BestLapTimeMs
				data.LastLapTimeMs = laps.LastLapTimeMs
				data.TotalCompetitionTimeMs = laps.TotalCompetitionTimeMs
			}
		}

		result = append(result, data)
	}

	return result, nil
}
