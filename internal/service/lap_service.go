package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

// LapScan represents a scanned RFID tag with timing information
type LapScan struct {
	ID             string
	TagValue       string
	Timestamp      time.Time
	ReaderType     string
	COMPort        string
	SignalStrength *int
}

// LapService handles real-time lap processing from RFID scans
type LapService struct {
	rawScanRepo         *repository.RawScanRepository
	competitorModelRepo *repository.CompetitorModelRepository
	competitionRepo     *repository.CompetitionRepository
	participantRepo     *repository.CompetitionParticipantRepository
	competitionLapsRepo *repository.CompetitionLapsRepository

	// Buffered channel for async processing
	scanChannel chan LapScan

	// Active competition cache (in-memory)
	activeCompetition   *models.Competition
	activeCompetitionMu sync.RWMutex

	// Competitor models cache (transponder -> competitor_model_id)
	transponderCache map[string]string
	transponderMu    sync.RWMutex

	// Participant results cache (competition_participant_id -> current lap data)
	participantResults map[string]*ParticipantLapData
	resultsMu          sync.RWMutex

	// Reverse lookup cache for O(1) participant search (competitor_model_id -> participant_id)
	modelToParticipant map[string]string
	cacheMu            sync.RWMutex

	// Worker control
	stopWorker chan bool
	wg         sync.WaitGroup

	// Batch settings
	batchSize    int
	batchTimeout time.Duration
}

// ParticipantLapData holds real-time lap tracking for a participant
type ParticipantLapData struct {
	CompetitionParticipantID string
	CompetitionID            string
	CompetitorModelID        string
	StartTime                time.Time
	LastPassTime             time.Time
	LapCount                 int
	BestLapTimeMs            int
	BestLapNumber            int
	LastLapTimeMs            int
	TotalTimeMs              int
	LapTimes                 []int // in milliseconds
}

// NewLapService creates a new lap service with buffered processing
func NewLapService(
	rawScanRepo *repository.RawScanRepository,
	competitorModelRepo *repository.CompetitorModelRepository,
	competitionRepo *repository.CompetitionRepository,
	participantRepo *repository.CompetitionParticipantRepository,
	competitionLapsRepo *repository.CompetitionLapsRepository,
) *LapService {
	return &LapService{
		rawScanRepo:         rawScanRepo,
		competitorModelRepo: competitorModelRepo,
		competitionRepo:     competitionRepo,
		participantRepo:     participantRepo,
		competitionLapsRepo: competitionLapsRepo,
		scanChannel:         make(chan LapScan, 200), // Buffer for 200 scans
		transponderCache:    make(map[string]string),
		modelToParticipant:  make(map[string]string),
		participantResults:  make(map[string]*ParticipantLapData),
		stopWorker:          make(chan bool),
		batchSize:           50,                     // Write to DB every 50 scans
		batchTimeout:        100 * time.Millisecond, // Or every 100ms
	}
}

// Start begins the background worker for processing scans
func (s *LapService) Start() {
	s.wg.Add(1)
	go s.processingWorker()
}

// Stop gracefully shuts down the worker
func (s *LapService) Stop() {
	close(s.stopWorker)
	s.wg.Wait()
}

// SetActiveCompetition sets the currently active competition
func (s *LapService) SetActiveCompetition(comp *models.Competition) {
	s.activeCompetitionMu.Lock()
	defer s.activeCompetitionMu.Unlock()
	s.activeCompetition = comp

	// Reload participant cache for this competition
	if comp != nil && comp.Status == "in_progress" {
		s.loadParticipantCache(comp.ID)
	}
}

// GetActiveCompetition returns the currently active competition
func (s *LapService) GetActiveCompetition() *models.Competition {
	s.activeCompetitionMu.RLock()
	defer s.activeCompetitionMu.RUnlock()
	return s.activeCompetition
}

// loadParticipantCache loads all participants for a competition into memory
func (s *LapService) loadParticipantCache(competitionID string) {
	participants, err := s.participantRepo.GetByCompetitionID(competitionID)
	if err != nil {
		return
	}

	s.transponderMu.Lock()
	defer s.transponderMu.Unlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	// Clear and rebuild cache
	s.transponderCache = make(map[string]string)
	s.modelToParticipant = make(map[string]string)

	for _, p := range participants {
		cm, err := s.competitorModelRepo.GetByID(p.CompetitorModelID)
		if err != nil || cm == nil {
			continue
		}

		// Cache transponder -> competitor_model_id mapping
		s.transponderCache[cm.TransponderNumber] = cm.ID

		// Cache competitor_model_id -> participant_id for O(1) lookup
		s.modelToParticipant[cm.ID] = p.ID

		// Initialize lap data for this participant
		s.resultsMu.Lock()
		s.participantResults[p.ID] = &ParticipantLapData{
			CompetitionParticipantID: p.ID,
			CompetitionID:            competitionID,
			CompetitorModelID:        p.CompetitorModelID,
			StartTime:                time.Now(),
		}
		s.resultsMu.Unlock()
	}
}

// QueueScan adds a scan to the processing queue (non-blocking)
func (s *LapService) QueueScan(scan LapScan) bool {
	select {
	case s.scanChannel <- scan:
		return true
	default:
		// Channel full - drop scan (shouldn't happen with proper buffer sizing)
		return false
	}
}

// ProcessScan processes a single RFID scan and records lap if valid
func (s *LapService) ProcessScan(scan LapScan) {
	// Check if we have an active competition
	s.activeCompetitionMu.RLock()
	activeComp := s.activeCompetition
	s.activeCompetitionMu.RUnlock()

	if activeComp == nil || activeComp.Status != "in_progress" {
		return // No active race
	}

	// Look up competitor model from cache
	s.transponderMu.RLock()
	competitorModelID, exists := s.transponderCache[scan.TagValue]
	s.transponderMu.RUnlock()

	if !exists {
		return // Unknown transponder
	}

	// Find participant for this competitor model using O(1) lookup
	s.cacheMu.RLock()
	participantID, exists := s.modelToParticipant[competitorModelID]
	s.cacheMu.RUnlock()

	if !exists {
		return // Participant not found
	}

	// Calculate lap
	s.recordLap(participantID, scan.Timestamp)
}

// recordLap records a lap for a participant
func (s *LapService) recordLap(participantID string, timestamp time.Time) {
	s.resultsMu.Lock()
	defer s.resultsMu.Unlock()

	data, exists := s.participantResults[participantID]
	if !exists {
		return
	}

	// Calculate lap time
	lapTimeMs := 0
	if data.LastPassTime.IsZero() {
		// First lap - start timer
		data.StartTime = timestamp
		data.LastPassTime = timestamp
		data.LapCount = 1
		data.TotalTimeMs = 0
		
		// Mark transponder as worked on first pass
		s.markTransponderWorked(participantID)
	} else {
		// Subsequent lap
		lapTimeMs = int(timestamp.Sub(data.LastPassTime).Milliseconds())
		data.LastPassTime = timestamp
		data.LapCount++
		data.TotalTimeMs = int(timestamp.Sub(data.StartTime).Milliseconds())
		data.LastLapTimeMs = lapTimeMs
		data.LapTimes = append(data.LapTimes, lapTimeMs)

		// Update best lap
		if data.BestLapTimeMs == 0 || lapTimeMs < data.BestLapTimeMs {
			data.BestLapTimeMs = lapTimeMs
			data.BestLapNumber = data.LapCount
		}
	}
}

// markTransponderWorked sets transponder_worked = true for the participant
func (s *LapService) markTransponderWorked(participantID string) {
	participant, err := s.participantRepo.GetByID(participantID)
	if err != nil || participant == nil {
		return
	}
	
	// Only update if not already marked
	if !participant.TransponderWorked {
		participant.TransponderWorked = true
		s.participantRepo.Update(participant)
	}
}

// GetParticipantResults returns current results for all participants
func (s *LapService) GetParticipantResults() map[string]*ParticipantLapData {
	s.resultsMu.RLock()
	defer s.resultsMu.RUnlock()

	results := make(map[string]*ParticipantLapData)
	for k, v := range s.participantResults {
		results[k] = v
	}
	return results
}

// PersistResults writes current results to database
func (s *LapService) PersistResults() error {
	s.resultsMu.RLock()
	defer s.resultsMu.RUnlock()

	for participantID, data := range s.participantResults {
		laps := &models.CompetitionLaps{
			ID:                       uuid.New().String(),
			CompetitionParticipantID: participantID,
			TimeStart:                data.StartTime,
			NumberOfLaps:             data.LapCount,
			BestLapTimeMs:            data.BestLapTimeMs,
			BestLapNumber:            data.BestLapNumber,
			LastLapTimeMs:            data.LastLapTimeMs,
			LastPassTime:             &data.LastPassTime,
			TotalCompetitionTimeMs:   data.TotalTimeMs,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
		}

		if err := s.competitionLapsRepo.Upsert(laps); err != nil {
			return fmt.Errorf("failed to persist results for participant %s: %w", participantID, err)
		}
	}

	return nil
}

// processingWorker is the background goroutine that processes scans in batches
func (s *LapService) processingWorker() {
	defer s.wg.Done()

	batch := make([]*models.RawScan, 0, s.batchSize)
	ticker := time.NewTicker(s.batchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopWorker:
			// Flush remaining batch before exit
			if len(batch) > 0 {
				s.flushBatch(batch)
			}
			return

		case scan := <-s.scanChannel:
			// Create RawScan record
			rawScan := &models.RawScan{
				ID:             uuid.New().String(),
				Timestamp:      scan.Timestamp,
				TagValue:       scan.TagValue,
				ReaderType:     scan.ReaderType,
				COMPort:        scan.COMPort,
				SignalStrength: scan.SignalStrength,
				IsProcessed:    false,
			}

			batch = append(batch, rawScan)

			// Process lap immediately (don't wait for batch)
			s.ProcessScan(scan)

			// Flush batch if full
			if len(batch) >= s.batchSize {
				s.flushBatch(batch)
				batch = make([]*models.RawScan, 0, s.batchSize)
			}

		case <-ticker.C:
			// Timeout - flush current batch
			if len(batch) > 0 {
				s.flushBatch(batch)
				batch = make([]*models.RawScan, 0, s.batchSize)
			}
		}
	}
}

// flushBatch writes a batch of raw scans to the database
func (s *LapService) flushBatch(batch []*models.RawScan) {
	if len(batch) == 0 {
		return
	}

	err := s.rawScanRepo.CreateBulk(batch)
	if err != nil {
		// Log error but don't block processing
		// In production, you'd use a proper logger
		return
	}
}

// RefreshTransponderCache reloads the transponder cache from database
func (s *LapService) RefreshTransponderCache() {
	allModels, err := s.competitorModelRepo.GetAll()
	if err != nil {
		return
	}

	s.transponderMu.Lock()
	defer s.transponderMu.Unlock()

	s.transponderCache = make(map[string]string)
	for _, cm := range allModels {
		s.transponderCache[cm.TransponderNumber] = cm.ID
	}
}

// GetTransponderForModel returns the competitor model ID for a transponder
func (s *LapService) GetTransponderForModel(transponderNumber string) (string, bool) {
	s.transponderMu.RLock()
	defer s.transponderMu.RUnlock()

	id, exists := s.transponderCache[transponderNumber]
	return id, exists
}
