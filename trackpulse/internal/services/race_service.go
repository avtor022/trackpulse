package services

import (
	"time"
	"trackpulse/internal/repository"
	"trackpulse/internal/models"
)

type RaceService interface {
	StartRace(raceID string) error
	StopRace(raceID string) error
	GetLiveStandings(raceID string) ([]StandingsRow, error)
	CreateRace(race *models.Race) error
	UpdateRace(race *models.Race) error
	GetRaceByID(raceID string) (*models.Race, error)
	GetActiveRaces() ([]models.Race, error)
	AddParticipantToRace(raceID string, racerModelID string) error
	RemoveParticipantFromRace(raceID string, participantID string) error
	GetRaceParticipants(raceID string) ([]RaceParticipantInfo, error)
}

type StandingsRow struct {
	Position        int
	RaceParticipantID string
	RacerNumber     int
	RacerFullName   string
	ModelBrand      string
	ModelName       string
	NumberOfLaps    int
	BestLapTimeMs   int
	BestLapNumber   int
	LastLapTimeMs   int
	TotalRaceTimeMs int
	LastPassTime    *time.Time
	Status          string // "Racing", "Finished", "DNF", "DNS"
}

type RaceParticipantInfo struct {
	ID            string
	RaceID        string
	RacerModelID  string
	RacerNumber   int
	RacerFullName string
	ModelBrand    string
	ModelName     string
	GridPosition  *int
	IsFinished    bool
	Disqualified  bool
	DNFReason     *string
}

type raceService struct {
	repos *repository.Repositories
}

func NewRaceService(repos *repository.Repositories, cfg *config.Config) RaceService {
	return &raceService{
		repos: repos,
	}
}

func (s *raceService) StartRace(raceID string) error {
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

func (s *raceService) StopRace(raceID string) error {
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

func (s *raceService) GetLiveStandings(raceID string) ([]StandingsRow, error) {
	// Get all participants for the race with their details
	query := `
		SELECT 
			rp.id,
			rc.racer_number,
			rc.full_name,
			rm.brand,
			rm.model_name,
			rl.number_of_laps,
			rl.best_lap_time_ms,
			rl.best_lap_number,
			rl.last_lap_time_ms,
			rl.total_race_time_ms,
			rl.last_pass_time,
			rp.is_finished,
			rp.disqualified,
			rp.dnf_reason
		FROM race_participants rp
		JOIN racer_models rcm ON rp.racer_model_id = rcm.id
		JOIN racers rc ON rcm.racer_id = rc.id
		JOIN rc_models rm ON rcm.rc_model_id = rm.id
		LEFT JOIN race_laps rl ON rp.id = rl.race_participant_id
		WHERE rp.race_id = ?
		ORDER BY rl.number_of_laps DESC, rl.best_lap_time_ms ASC
	`

	// Since we don't have direct access to the DB here, we'll need to get the data differently
	// Get participants first
	participants, err := s.repos.RaceParticipants.GetByRaceID(raceID)
	if err != nil {
		return nil, err
	}

	var standings []StandingsRow
	for i, participant := range participants {
		// Get participant details
		details, err := s.repos.RaceParticipants.GetParticipantWithDetails(participant.ID)
		if err != nil {
			continue // Skip if we can't get details
		}

		// Get race lap info
		raceLap, err := s.repos.RaceLaps.GetByRaceParticipantID(participant.ID)
		if err != nil {
			// If no race lap record exists, create a default entry
			raceLap = &models.RaceLap{
				NumberOfLaps:    0,
				BestLapTimeMs:   0,
				BestLapNumber:   0,
				LastLapTimeMs:   0,
				TotalRaceTimeMs: 0,
			}
		}

		// Determine status
		status := "Racing"
		if participant.IsFinished {
			status = "Finished"
		} else if participant.Disqualified {
			status = "DSQ"
		} else if participant.DNFReason != nil && *participant.DNFReason != "" {
			status = "DNF"
		}

		standing := StandingsRow{
			Position:        i + 1,
			RaceParticipantID: participant.ID,
			RacerNumber:     details.RacerNumber,
			RacerFullName:   details.RacerFullName,
			ModelBrand:      details.ModelBrand,
			ModelName:       details.ModelName,
			NumberOfLaps:    raceLap.NumberOfLaps,
			BestLapTimeMs:   raceLap.BestLapTimeMs,
			BestLapNumber:   raceLap.BestLapNumber,
			LastLapTimeMs:   raceLap.LastLapTimeMs,
			TotalRaceTimeMs: raceLap.TotalRaceTimeMs,
			LastPassTime:    nil, // Would need to extract from raceLap.LastPassTime string
			Status:          status,
		}

		standings = append(standings, standing)
	}

	return standings, nil
}

func (s *raceService) CreateRace(race *models.Race) error {
	if race.ID == "" {
		race.ID = generateUUID() // Would need to implement proper UUID generation
	}
	
	if race.Status == "" {
		race.Status = "scheduled"
	}
	
	now := time.Now().UTC()
	if race.CreatedAt.IsZero() {
		race.CreatedAt = now
	}
	race.UpdatedAt = now

	return s.repos.Races.Create(race)
}

func (s *raceService) UpdateRace(race *models.Race) error {
	race.UpdatedAt = time.Now().UTC()
	return s.repos.Races.Update(race)
}

func (s *raceService) GetRaceByID(raceID string) (*models.Race, error) {
	return s.repos.Races.GetByID(raceID)
}

func (s *raceService) GetActiveRaces() ([]models.Race, error) {
	return s.repos.Races.GetActiveRaces()
}

func (s *raceService) AddParticipantToRace(raceID string, racerModelID string) error {
	// Check if racer model exists
	_, err := s.repos.RacerModels.GetByID(racerModelID)
	if err != nil {
		return err
	}

	// Check if race exists
	_, err = s.repos.Races.GetByID(raceID)
	if err != nil {
		return err
	}

	// Check if participant already exists in this race
	participants, err := s.repos.RaceParticipants.GetByRaceID(raceID)
	if err != nil {
		return err
	}

	for _, p := range participants {
		if p.RacerModelID == racerModelID {
			return nil // Already exists, no error
		}
	}

	// Create new participant
	participant := &models.RaceParticipant{
		ID:             generateUUID(), // Would need to implement proper UUID generation
		RaceID:         raceID,
		RacerModelID:   racerModelID,
		IsFinished:     false,
		Disqualified:   false,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err = s.repos.RaceParticipants.Create(participant)
	if err != nil {
		return err
	}

	// Create corresponding race_laps record
	raceLap := &models.RaceLap{
		ID:                generateUUID(), // Would need to implement proper UUID generation
		RaceParticipantID: participant.ID,
		NumberOfLaps:      0,
		BestLapTimeMs:     0,
		BestLapNumber:     0,
		LastLapTimeMs:     0,
		TotalRaceTimeMs:   0,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	return s.repos.RaceLaps.Create(raceLap)
}

func (s *raceService) RemoveParticipantFromRace(raceID string, participantID string) error {
	// Verify that the participant belongs to the specified race
	participant, err := s.repos.RaceParticipants.GetByID(participantID)
	if err != nil {
		return err
	}

	if participant.RaceID != raceID {
		return nil // Participant doesn't belong to this race, nothing to remove
	}

	// Remove the participant (this will also remove related lap records due to cascade)
	return s.repos.RaceParticipants.Delete(participantID)
}

func (s *raceService) GetRaceParticipants(raceID string) ([]RaceParticipantInfo, error) {
	participants, err := s.repos.RaceParticipants.GetByRaceID(raceID)
	if err != nil {
		return nil, err
	}

	var participantInfos []RaceParticipantInfo
	for _, p := range participants {
		details, err := s.repos.RaceParticipants.GetParticipantWithDetails(p.ID)
		if err != nil {
			continue // Skip if we can't get details
		}

		info := RaceParticipantInfo{
			ID:            p.ID,
			RaceID:        p.RaceID,
			RacerModelID:  p.RacerModelID,
			RacerNumber:   details.RacerNumber,
			RacerFullName: details.RacerFullName,
			ModelBrand:    details.ModelBrand,
			ModelName:     details.ModelName,
			GridPosition:  p.GridPosition,
			IsFinished:    p.IsFinished,
			Disqualified:  p.Disqualified,
			DNFReason:     p.DNFReason,
		}

		participantInfos = append(participantInfos, info)
	}

	return participantInfos, nil
}

// Helper function to generate UUIDs (would need to import appropriate package)
func generateUUID() string {
	// Implementation would use a UUID generation library
	// For now, returning a placeholder
	return "placeholder-uuid"
}