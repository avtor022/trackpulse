package service

import (
	"fmt"
	"log"
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
	scanChannel chan scanRequest
	errorMu     sync.RWMutex
	onError     func(error)

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
	queueMu    sync.RWMutex
	stopped    bool

	// Batch settings
	batchSize    int
	batchTimeout time.Duration
}

type scanRequest struct {
	scan LapScan
	done chan struct{}
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
		scanChannel:         make(chan scanRequest, 200), // Buffer for 200 scans
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
	s.queueMu.Lock()
	if s.stopped {
		s.queueMu.Unlock()
		s.wg.Wait()
		return
	}
	s.stopped = true
	close(s.stopWorker)
	s.queueMu.Unlock()
	s.wg.Wait()
}

// SetActiveCompetition sets the currently active competition
func (s *LapService) SetActiveCompetition(comp *models.Competition) error {
	s.activeCompetitionMu.Lock()
	defer s.activeCompetitionMu.Unlock()

	// Reload participant cache for this competition
	if comp != nil && comp.Status == "in_progress" {
		if err := s.loadParticipantCache(comp.ID); err != nil {
			return err
		}
	}
	s.activeCompetition = comp
	return nil
}

// GetActiveCompetition returns the currently active competition
func (s *LapService) GetActiveCompetition() *models.Competition {
	s.activeCompetitionMu.RLock()
	defer s.activeCompetitionMu.RUnlock()
	return s.activeCompetition
}

// loadParticipantCache loads all participants for a competition into memory
func (s *LapService) loadParticipantCache(competitionID string) error {
	participants, err := s.participantRepo.GetByCompetitionID(competitionID)
	if err != nil {
		return fmt.Errorf("failed to load participants: %w", err)
	}
	previous, err := s.competitionLapsRepo.GetAllByCompetitionID(competitionID)
	if err != nil {
		return fmt.Errorf("failed to load lap results: %w", err)
	}
	previousByParticipant := make(map[string]*models.CompetitionLaps, len(previous))
	for _, result := range previous {
		previousByParticipant[result.CompetitionParticipantID] = result
	}
	transponders := make(map[string]string)
	modelParticipants := make(map[string]string)
	results := make(map[string]*ParticipantLapData)
	for _, p := range participants {
		cm, err := s.competitorModelRepo.GetByID(p.CompetitorModelID)
		if err != nil {
			return fmt.Errorf("failed to load transponder: %w", err)
		}
		if cm == nil {
			return fmt.Errorf("transponder not found for participant %s", p.ID)
		}
		transponders[cm.TransponderNumber] = cm.ID
		modelParticipants[cm.ID] = p.ID
		data := &ParticipantLapData{CompetitionParticipantID: p.ID, CompetitionID: competitionID, CompetitorModelID: p.CompetitorModelID}
		if saved := previousByParticipant[p.ID]; saved != nil {
			data.StartTime = saved.TimeStart
			if saved.LastPassTime != nil {
				data.LastPassTime = *saved.LastPassTime
			}
			data.LapCount = saved.NumberOfLaps
			data.BestLapTimeMs = saved.BestLapTimeMs
			data.BestLapNumber = saved.BestLapNumber
			data.LastLapTimeMs = saved.LastLapTimeMs
			data.TotalTimeMs = saved.TotalCompetitionTimeMs
		}
		results[p.ID] = data
	}

	s.transponderMu.Lock()
	defer s.transponderMu.Unlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.resultsMu.Lock()
	defer s.resultsMu.Unlock()
	s.transponderCache = transponders
	s.modelToParticipant = modelParticipants
	s.participantResults = results
	return nil
}

// QueueScan adds a scan to the processing queue, applying backpressure when full.
func (s *LapService) QueueScan(scan LapScan) bool {
	s.queueMu.RLock()
	defer s.queueMu.RUnlock()
	if s.stopped {
		return false
	}
	select {
	case s.scanChannel <- scanRequest{scan: scan}:
		return true
	case <-s.stopWorker:
		return false
	}
}

// SetErrorHandler registers a handler for background persistence errors.
func (s *LapService) SetErrorHandler(handler func(error)) {
	s.errorMu.Lock()
	defer s.errorMu.Unlock()
	s.onError = handler
}

func (s *LapService) reportError(err error) {
	log.Printf("lap service: %v", err)
	s.errorMu.RLock()
	handler := s.onError
	s.errorMu.RUnlock()
	if handler != nil {
		handler(err)
	}
}

// Flush waits for scans already queued to finish processing.
func (s *LapService) Flush() {
	s.queueMu.RLock()
	defer s.queueMu.RUnlock()
	if s.stopped {
		return
	}
	done := make(chan struct{})
	s.scanChannel <- scanRequest{done: done}
	<-done
}

// ProcessScan processes a single RFID scan and records lap if valid
func (s *LapService) ProcessScan(scan LapScan) {
	err := s.processScan(scan)
	if err != nil {
		s.reportError(err)
	}
}

func (s *LapService) processScan(scan LapScan) error {
	// Check if we have an active competition
	s.activeCompetitionMu.RLock()
	defer s.activeCompetitionMu.RUnlock()
	activeComp := s.activeCompetition

	if activeComp == nil || activeComp.Status != "in_progress" {
		return nil // No active race
	}
	if activeComp.TimeStart != nil && scan.Timestamp.Before(*activeComp.TimeStart) {
		return nil
	}

	// Look up competitor model from cache
	s.transponderMu.RLock()
	competitorModelID, exists := s.transponderCache[scan.TagValue]
	s.transponderMu.RUnlock()

	if !exists {
		return nil // Unknown transponder
	}

	// Find participant for this competitor model using O(1) lookup
	s.cacheMu.RLock()
	participantID, exists := s.modelToParticipant[competitorModelID]
	s.cacheMu.RUnlock()

	if !exists {
		return nil // Participant not found
	}

	// Calculate lap
	if err := s.recordLap(participantID, scan.Timestamp, activeComp.MinLapTimeMs); err != nil {
		return fmt.Errorf("failed to persist lap result: %w", err)
	}
	return nil
}

// recordLap records a lap for a participant
func (s *LapService) recordLap(participantID string, timestamp time.Time, minLapTimeMs int) error {
	s.resultsMu.Lock()
	defer s.resultsMu.Unlock()

	data, exists := s.participantResults[participantID]
	if !exists {
		return nil
	}

	// Calculate lap time
	lapTimeMs := 0
	if data.LastPassTime.IsZero() {
		// First lap - start timer
		data.StartTime = timestamp
		data.LastPassTime = timestamp
		data.LapCount = 0
		data.TotalTimeMs = 0
	} else {
		// Subsequent lap
		lapTimeMs = int(timestamp.Sub(data.LastPassTime).Milliseconds())
		if lapTimeMs <= 0 || lapTimeMs < minLapTimeMs {
			return nil
		}
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
	return s.persistParticipant(data)
}

// GetParticipantResults returns current results for all participants
func (s *LapService) GetParticipantResults() map[string]*ParticipantLapData {
	s.resultsMu.RLock()
	defer s.resultsMu.RUnlock()

	results := make(map[string]*ParticipantLapData)
	for k, v := range s.participantResults {
		copy := *v
		copy.LapTimes = append([]int(nil), v.LapTimes...)
		results[k] = &copy
	}
	return results
}

// GetStandings returns the participants and their persisted race results.
func (s *LapService) GetStandings(competitionID string) ([]models.CompetitionStanding, error) {
	return s.competitionLapsRepo.GetStandings(competitionID)
}

// PersistResults writes current results to database
func (s *LapService) PersistResults() error {
	s.resultsMu.RLock()
	defer s.resultsMu.RUnlock()

	for participantID, data := range s.participantResults {
		if data.LastPassTime.IsZero() {
			continue
		}
		if err := s.persistParticipant(data); err != nil {
			return fmt.Errorf("failed to persist results for participant %s: %w", participantID, err)
		}
	}

	return nil
}

func (s *LapService) persistParticipant(data *ParticipantLapData) error {
	laps := &models.CompetitionLaps{
		ID:                       uuid.New().String(),
		CompetitionParticipantID: data.CompetitionParticipantID,
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

	return s.competitionLapsRepo.Upsert(laps)
}

// processingWorker is the background goroutine that processes scans in batches
func (s *LapService) processingWorker() {
	defer s.wg.Done()

	batch := make([]*models.RawScan, 0, s.batchSize)
	ticker := time.NewTicker(s.batchTimeout)
	defer ticker.Stop()
	process := func(request scanRequest) {
		if request.done != nil {
			if len(batch) > 0 {
				s.flushBatch(batch)
				batch = make([]*models.RawScan, 0, s.batchSize)
			}
			close(request.done)
			return
		}
		scan := request.scan
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
		s.ProcessScan(scan)
		if len(batch) >= s.batchSize {
			s.flushBatch(batch)
			batch = make([]*models.RawScan, 0, s.batchSize)
		}
	}

	for {
		select {
		case <-s.stopWorker:
			for {
				select {
				case request := <-s.scanChannel:
					process(request)
				default:
					s.flushBatch(batch)
					return
				}
			}

		case request := <-s.scanChannel:
			process(request)

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
		s.reportError(fmt.Errorf("failed to save raw scans: %w", err))
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
