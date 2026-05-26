package service

import (
"fmt"
"trackpulse/internal/models"
"trackpulse/internal/repository"
)

// MonitoringService provides combined data for monitoring display
type MonitoringService struct {
participantRepo     *repository.CompetitionParticipantRepository
competitorModelRepo *repository.CompetitorModelRepository
competitorRepo      *repository.CompetitorRepository
rcModelRepo         *repository.RCModelRepository
competitionLapsRepo *repository.CompetitionLapsRepository
lapService          *LapService
}

// ParticipantMonitoringData represents complete participant data for monitoring
type ParticipantMonitoringData struct {
CompetitorNumber    int
FullName            string
ModelName           string
ModelScale          string
LapCount            int
BestLapTimeMs       int
TransponderVerified bool
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(
participantRepo *repository.CompetitionParticipantRepository,
competitorModelRepo *repository.CompetitorModelRepository,
competitorRepo *repository.CompetitorRepository,
rcModelRepo *repository.RCModelRepository,
competitionLapsRepo *repository.CompetitionLapsRepository,
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

// GetParticipantsMonitoringData returns monitoring data for all participants in a competition
func (s *MonitoringService) GetParticipantsMonitoringData(competitionID string) ([]ParticipantMonitoringData, error) {
// Get all participants for the competition
participants, err := s.participantRepo.GetByCompetitionID(competitionID)
if err != nil {
return nil, fmt.Errorf("failed to get participants: %w", err)
}

var result []ParticipantMonitoringData

for _, p := range participants {
data := ParticipantMonitoringData{}

// Get competitor model
cm, err := s.competitorModelRepo.GetByID(p.CompetitorModelID)
if err != nil || cm == nil {
continue
}

// Get competitor info
competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
if err != nil || competitor == nil {
continue
}
data.CompetitorNumber = competitor.CompetitorNumber
data.FullName = competitor.FullName

// Get RC model info
rcModel, err := s.rcModelRepo.GetByID(cm.RCModelID)
if err != nil || rcModel == nil {
continue
}
data.ModelName = rcModel.ModelName
data.ModelScale = rcModel.Scale

// Try to get live data from lap service first
lapData := s.getLiveLapData(p.ID)
if lapData != nil {
data.LapCount = lapData.LapCount
data.BestLapTimeMs = lapData.BestLapTimeMs
data.TransponderVerified = lapData.TransponderVerified
} else {
// Fallback to database
laps, err := s.competitionLapsRepo.GetByParticipantID(p.ID)
if err == nil && laps != nil {
data.LapCount = laps.NumberOfLaps
data.BestLapTimeMs = laps.BestLapTimeMs
// Transponder verified if we have any laps
data.TransponderVerified = laps.NumberOfLaps > 0
}
}

result = append(result, data)
}

return result, nil
}

// getLiveLapData gets live lap data from the lap service cache
func (s *MonitoringService) getLiveLapData(participantID string) *ParticipantLapData {
if s.lapService == nil {
return nil
}

results := s.lapService.GetParticipantResults()
if data, exists := results[participantID]; exists {
return data
}
return nil
}
