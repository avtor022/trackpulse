package services

import (
	"sync"
	"time"
	"trackpulse/internal/repository"
	"trackpulse/internal/models"
)

type LapService interface {
	ProcessScan(tagValue string, readerType string, comPort string) error
	StartRace(raceID string) error
	StopRace(raceID string) error
	GetLiveStandings(raceID string) ([]StandingsRow, error)
	InitializeRace(raceID string) error
	FinishRace(raceID string) error
	GetRaceStatistics(raceID string) (*RaceStatistics, error)
	ValidateLap(raceParticipantID string, lapTimeMs int) error
	InvalidateLap(raceParticipantID string, lapNumber int, reason string) error
}

type RaceStatistics struct {
	TotalParticipants int
	CompletedLaps   int
	AverageLapTime  int
	FastestLapTime  int
	SlowestLapTime  int
	CurrentLeaders  []LeaderInfo
}

type LeaderInfo struct {
	RaceParticipantID string
	RacerNumber     int
	RacerFullName   string
	NumberOfLaps    int
	BestLapTimeMs   int
}

type lapService struct {
	repos       *repository.Repositories
	debounceMs  int
	lastScanTime map[string]time.Time
	mu          sync.Mutex
	activeRace  *string
}

func NewLapService(repos *repository.Repositories, cfg *Config) LapService {
	// Get debounce setting from system settings
	settingsRepo := repos.SystemSettings
	hardwareSettings, err := settingsRepo.GetHardwareSettings()
	var debounceMs int = 2000 // Default value
	if err == nil {
		if val, ok := hardwareSettings[models.SettingKeyDebounceMs]; ok {
			// Convert string to int
			// For now, using default if conversion fails
			debounceMs = 2000
		}
	}

	return &lapService{
		repos:        repos,
		debounceMs:   debounceMs,
		lastScanTime: make(map[string]time.Time),
		activeRace:   nil,
	}
}

func (s *lapService) ProcessScan(tagValue string, readerType string, comPort string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if there's an active race
	if s.activeRace == nil {
		// No active race, just log the scan
		return s.logRawScan(tagValue, readerType, comPort, nil)
	}

	// Check debounce
	if lastTime, exists := s.lastScanTime[tagValue]; exists {
		if time.Since(lastTime).Milliseconds() < int64(s.debounceMs) {
			// Still in debounce period, just log and return
			return s.logRawScan(tagValue, readerType, comPort, nil)
		}
	}
	s.lastScanTime[tagValue] = time.Now()

	// Find the racer model associated with this tag
	racerModel, err := s.repos.RacerModels.GetByTransponder(tagValue)
	if err != nil {
		// Unknown tag, log and return
		return s.logRawScan(tagValue, readerType, comPort, nil)
	}

	// Find the race participant for this racer model in the active race
	query := `
		SELECT id 
		FROM race_participants 
		WHERE race_id = ? AND racer_model_id = ?
	`
	var participantID string
	err = s.repos.DB.QueryRow(query, *s.activeRace, racerModel.ID).Scan(&participantID)
	if err != nil {
		// Racer is not participating in the active race, just log
		return s.logRawScan(tagValue, readerType, comPort, &racerModel.ID)
	}

	// Calculate lap time
	now := time.Now()
	var lapTimeMs int
	var lastPassTime *time.Time

	// Get the last pass time for this participant
	raceLap, err := s.repos.RaceLaps.GetByRaceParticipantID(participantID)
	if err == nil && raceLap.LastPassTime != nil {
		// Calculate lap time based on last pass
		lastPass, _ := time.Parse(time.RFC3339, *raceLap.LastPassTime) // Simplified parsing
		lapTimeMs = int(now.Sub(lastPass).Milliseconds())
		lastPassTime = &now
	} else {
		// First lap, no time calculation needed
		lapTimeMs = 0
		lastPassTime = &now
	}

	// Validate the lap time (e.g., ensure it's not too fast)
	if lapTimeMs > 0 && lapTimeMs < 1000 { // Less than 1 second, probably invalid
		// Log as invalid scan
		return s.logRawScan(tagValue, readerType, comPort, &racerModel.ID)
	}

	// Update race lap information
	newLapCount := raceLap.NumberOfLaps + 1
	
	// Update best lap time if this is better
	newBestLapTime := raceLap.BestLapTimeMs
	newBestLapNumber := raceLap.BestLapNumber
	if (lapTimeMs > 0 && newBestLapTime == 0) || (lapTimeMs > 0 && lapTimeMs < newBestLapTime) {
		newBestLapTime = lapTimeMs
		newBestLapNumber = newLapCount
	}

	// Update race_laps record
	err = s.repos.RaceLaps.UpdateRaceResults(
		participantID,
		newLapCount,
		newBestLapTime,
		newBestLapNumber,
		lapTimeMs,
		pointerToString(lastPassTime),
		int(now.Sub(raceLap.TimeStart).Milliseconds()),
	)
	if err != nil {
		return err
	}

	// Create lap history record
	lapHistory := &models.LapHistory{
		ID:                generateUUID(), // Would need to implement proper UUID generation
		RaceParticipantID: participantID,
		LapNumber:         newLapCount,
		LapTimeMs:         lapTimeMs,
		StartTime:         *raceLap.LastPassTime, // Previous last pass time
		EndTime:           now,
		IsValid:           true,
		CreatedAt:         now,
	}
	err = s.repos.LapHistory.Create(lapHistory)
	if err != nil {
		return err
	}

	// Log the processed scan
	return s.logRawScan(tagValue, readerType, comPort, &racerModel.ID)
}

func (s *lapService) StartRace(raceID string) error {
	// Set the active race
	s.activeRace = &raceID

	// Initialize race if needed (could be done in race service)
	return s.InitializeRace(raceID)
}

func (s *lapService) StopRace(raceID string) error {
	// Clear the active race
	s.activeRace = nil

	// Finish the race
	return s.FinishRace(raceID)
}

func (s *lapService) GetLiveStandings(raceID string) ([]StandingsRow, error) {
	// This would be implemented similarly to in race_service but could include real-time updates
	service := NewRaceService(s.repos, &Config{})
	return service.GetLiveStandings(raceID)
}

func (s *lapService) InitializeRace(raceID string) error {
	// Update race status to active and set start time
	race, err := s.repos.Races.GetByID(raceID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	race.Status = "active"
	race.TimeStart = &now
	race.UpdatedAt = now

	err = s.repos.Races.Update(race)
	if err != nil {
		return err
	}

	// Initialize race_laps records for all participants if they don't exist
	participants, err := s.repos.RaceParticipants.GetByRaceID(raceID)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		// Check if race_laps record exists for this participant
		_, err := s.repos.RaceLaps.GetByRaceParticipantID(participant.ID)
		if err != nil {
			// Record doesn't exist, create a new one
			raceLap := &models.RaceLap{
				ID:                generateUUID(), // Would need to implement proper UUID generation
				RaceParticipantID: participant.ID,
				TimeStart:         now,
				NumberOfLaps:      0,
				BestLapTimeMs:     0,
				BestLapNumber:     0,
				LastLapTimeMs:     0,
				TotalRaceTimeMs:   0,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			err = s.repos.RaceLaps.Create(raceLap)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *lapService) FinishRace(raceID string) error {
	// Update race status to finished and set finish time
	race, err := s.repos.Races.GetByID(raceID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	race.Status = "finished"
	race.TimeFinish = &now
	race.UpdatedAt = now

	return s.repos.Races.Update(race)
}

func (s *lapService) GetRaceStatistics(raceID string) (*RaceStatistics, error) {
	// Count total participants
	participants, err := s.repos.RaceParticipants.GetByRaceID(raceID)
	if err != nil {
		return nil, err
	}

	// Get all lap history for this race
	var allLaps []models.LapHistory
	for _, participant := range participants {
		participantLaps, err := s.repos.LapHistory.GetByRaceParticipantID(participant.ID)
		if err != nil {
			continue // Skip if we can't get laps for this participant
		}
		allLaps = append(allLaps, participantLaps...)
	}

	// Calculate statistics
	stats := &RaceStatistics{
		TotalParticipants: len(participants),
		CompletedLaps:   len(allLaps),
	}

	if len(allLaps) > 0 {
		var totalTime int
		var fastestTime int = 0
		var slowestTime int = 0

		for i, lap := range allLaps {
			totalTime += lap.LapTimeMs
			
			if fastestTime == 0 || lap.LapTimeMs < fastestTime {
				fastestTime = lap.LapTimeMs
			}
			
			if slowestTime == 0 || lap.LapTimeMs > slowestTime {
				slowestTime = lap.LapTimeMs
			}
		}

		stats.AverageLapTime = totalTime / len(allLaps)
		stats.FastestLapTime = fastestTime
		stats.SlowestLapTime = slowestTime
	}

	// Get current leaders (those with most laps, then fastest time)
	standings, err := s.GetLiveStandings(raceID)
	if err == nil {
		for _, standing := range standings {
			if len(stats.CurrentLeaders) < 5 { // Top 5 leaders
				stats.CurrentLeaders = append(stats.CurrentLeaders, LeaderInfo{
					RaceParticipantID: standing.RaceParticipantID,
					RacerNumber:     standing.RacerNumber,
					RacerFullName:   standing.RacerFullName,
					NumberOfLaps:    standing.NumberOfLaps,
					BestLapTimeMs:   standing.BestLapTimeMs,
				})
			} else {
				break
			}
		}
	}

	return stats, nil
}

func (s *lapService) ValidateLap(raceParticipantID string, lapTimeMs int) error {
	// Could implement validation rules here
	// For example: check if lap time is physically possible
	// For now, just return nil (valid)
	return nil
}

func (s *lapService) InvalidateLap(raceParticipantID string, lapNumber int, reason string) error {
	// Get the specific lap
	lap, err := s.repos.LapHistory.GetByRaceParticipantAndLapNumber(raceParticipantID, lapNumber)
	if err != nil {
		return err
	}

	// Mark as invalid
	lap.IsValid = false
	lap.InvalidationReason = &reason

	err = s.repos.LapHistory.Update(lap)
	if err != nil {
		return err
	}

	// Recalculate race results for this participant
	return s.recalculateRaceResults(raceParticipantID)
}

func (s *lapService) recalculateRaceResults(raceParticipantID string) error {
	// Get all valid laps for this participant
	validLaps, err := s.repos.LapHistory.GetValidLapsForParticipant(raceParticipantID)
	if err != nil {
		return err
	}

	if len(validLaps) == 0 {
		// No valid laps left, reset to zero
		return s.repos.RaceLaps.UpdateRaceResults(raceParticipantID, 0, 0, 0, 0, nil, 0)
	}

	// Calculate new stats
	var totalLaps int = len(validLaps)
	var bestLapTime int = 0
	var bestLapNumber int = 0
	var totalTime int = 0

	for i, lap := range validLaps {
		totalTime += lap.LapTimeMs
		
		if bestLapTime == 0 || lap.LapTimeMs < bestLapTime {
			bestLapTime = lap.LapTimeMs
			bestLapNumber = lap.LapNumber
		}
	}

	// Get the last pass time from the most recent valid lap
	var lastPassTimeStr *string
	if len(validLaps) > 0 {
		lastLap := validLaps[len(validLaps)-1]
		timeStr := lastLap.EndTime.Format(time.RFC3339)
		lastPassTimeStr = &timeStr
	}

	// Update race results
	return s.repos.RaceLaps.UpdateRaceResults(
		raceParticipantID,
		totalLaps,
		bestLapTime,
		bestLapNumber,
		validLaps[len(validLaps)-1].LapTimeMs, // Last lap time
		lastPassTimeStr,
		totalTime,
	)
}

func (s *lapService) logRawScan(tagValue string, readerType string, comPort string, racerModelID *string) error {
	scan := &models.RawScan{
		ID:                 generateUUID(), // Would need to implement proper UUID generation
		Timestamp:          time.Now(),
		TagValue:           tagValue,
		ReaderType:         readerType,
		IsProcessed:        racerModelID != nil,
		CreatedAt:          time.Now(),
	}
	
	if comPort != "" {
		scan.COMPort = &comPort
	}
	
	if racerModelID != nil {
		scan.LinkedRacerModelID = racerModelID
	}

	return s.repos.RawScans.Create(scan)
}

// Helper function to convert *time.Time to *string
func pointerToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.Format(time.RFC3339)
	return &str
}