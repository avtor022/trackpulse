package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// MockCompetitionRepository is a mock implementation of CompetitionRepository for testing
type MockCompetitionRepository struct {
	competitions map[string]*models.Competition
}

func NewMockCompetitionRepository() *MockCompetitionRepository {
	return &MockCompetitionRepository{
		competitions: make(map[string]*models.Competition),
	}
}

func (m *MockCompetitionRepository) GetAll() ([]models.Competition, error) {
	result := make([]models.Competition, 0, len(m.competitions))
	for _, c := range m.competitions {
		result = append(result, *c)
	}
	return result, nil
}

func (m *MockCompetitionRepository) GetByID(id string) (*models.Competition, error) {
	if c, ok := m.competitions[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCompetitionRepository) Create(competition *models.Competition) error {
	m.competitions[competition.ID] = competition
	return nil
}

func (m *MockCompetitionRepository) Update(competition *models.Competition) error {
	if _, ok := m.competitions[competition.ID]; !ok {
		return nil // Not found
	}
	m.competitions[competition.ID] = competition
	return nil
}

func (m *MockCompetitionRepository) Delete(id string) error {
	delete(m.competitions, id)
	return nil
}

func (m *MockCompetitionRepository) Count() (int, error) {
	return len(m.competitions), nil
}

func (m *MockCompetitionRepository) GetByStatus(status string) ([]models.Competition, error) {
	var result []models.Competition
	for _, c := range m.competitions {
		if c.Status == status {
			result = append(result, *c)
		}
	}
	return result, nil
}

// MockRCModelTypeRepository is a mock implementation of RCModelTypeRepository for testing
type MockRCModelTypeRepositoryForCompetition struct {
	types map[string]*models.RCModelType
}

func NewMockRCModelTypeRepositoryForCompetition() *MockRCModelTypeRepositoryForCompetition {
	return &MockRCModelTypeRepositoryForCompetition{
		types: make(map[string]*models.RCModelType),
	}
}

func (m *MockRCModelTypeRepositoryForCompetition) GetAll() ([]models.RCModelType, error) {
	result := make([]models.RCModelType, 0, len(m.types))
	for _, t := range m.types {
		result = append(result, *t)
	}
	return result, nil
}

func (m *MockRCModelTypeRepositoryForCompetition) GetByName(name string) (*models.RCModelType, error) {
	if t, ok := m.types[name]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockRCModelTypeRepositoryForCompetition) Create(name string) (*models.RCModelType, error) {
	id := uuid.New().String()
	t := &models.RCModelType{
		ID:   id,
		Name: name,
	}
	m.types[name] = t
	return t, nil
}

func (m *MockRCModelTypeRepositoryForCompetition) GetOrCreate(name string) (*models.RCModelType, error) {
	if t, ok := m.types[name]; ok {
		return t, nil
	}
	return m.Create(name)
}

func (m *MockRCModelTypeRepositoryForCompetition) Delete(name string) error {
	delete(m.types, name)
	return nil
}

// MockRCModelScaleRepositoryForCompetition is a mock implementation of RCModelScaleRepository for testing
type MockRCModelScaleRepositoryForCompetition struct {
	scales map[string]*models.RCModelScale
}

func NewMockRCModelScaleRepositoryForCompetition() *MockRCModelScaleRepositoryForCompetition {
	return &MockRCModelScaleRepositoryForCompetition{
		scales: make(map[string]*models.RCModelScale),
	}
}

func (m *MockRCModelScaleRepositoryForCompetition) GetAll() ([]models.RCModelScale, error) {
	result := make([]models.RCModelScale, 0, len(m.scales))
	for _, scale := range m.scales {
		result = append(result, *scale)
	}
	return result, nil
}

func (m *MockRCModelScaleRepositoryForCompetition) GetByName(name string) (*models.RCModelScale, error) {
	if scale, ok := m.scales[name]; ok {
		return scale, nil
	}
	return nil, nil
}

func (m *MockRCModelScaleRepositoryForCompetition) Create(name string) (*models.RCModelScale, error) {
	id := uuid.New().String()
	scale := &models.RCModelScale{
		ID:   id,
		Name: name,
	}
	m.scales[name] = scale
	return scale, nil
}

func (m *MockRCModelScaleRepositoryForCompetition) GetOrCreate(name string) (*models.RCModelScale, error) {
	if scale, ok := m.scales[name]; ok {
		return scale, nil
	}
	return m.Create(name)
}

func (m *MockRCModelScaleRepositoryForCompetition) Delete(name string) error {
	delete(m.scales, name)
	return nil
}

// TestCompetitionService_GetAllCompetitions tests the GetAllCompetitions method
func TestCompetitionService_GetAllCompetitions(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Add some test data
	timeStart := time.Now().Add(time.Hour)
	c1 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Summer Race 2025",
		CompetitionType:  "Time Attack",
		ModelType:        "Buggy",
		ModelScale:       "1/10",
		TrackName:        "Track A",
		Status:           "scheduled",
		TimeStart:        &timeStart,
	}
	c2 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Winter Championship",
		CompetitionType:  "Qualifying",
		ModelType:        "Truck",
		ModelScale:       "1/8",
		TrackName:        "Track B",
		Status:           "completed",
	}
	mockCompRepo.Create(c1)
	mockCompRepo.Create(c2)

	competitions, err := svc.GetAllCompetitions()
	if err != nil {
		t.Fatalf("GetAllCompetitions() error = %v", err)
	}
	if len(competitions) != 2 {
		t.Errorf("GetAllCompetitions() expected 2 competitions, got %d", len(competitions))
	}
}

// TestCompetitionService_GetCompetitionByID_Success tests successful retrieval by ID
func TestCompetitionService_GetCompetitionByID_Success(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	compID := uuid.New().String()
	competition := &models.Competition{
		ID:               compID,
		CompetitionTitle: "Test Competition",
		CompetitionType:  "Race",
		ModelType:        "Buggy",
		ModelScale:       "1/10",
		TrackName:        "Main Track",
		Status:           "scheduled",
	}
	mockCompRepo.Create(competition)

	retrieved, err := svc.GetCompetitionByID(compID)
	if err != nil {
		t.Fatalf("GetCompetitionByID() error = %v", err)
	}
	if retrieved == nil || retrieved.CompetitionTitle != "Test Competition" {
		t.Error("GetCompetitionByID() failed to retrieve competition")
	}
}

// TestCompetitionService_GetCompetitionByID_NotFound tests retrieval of non-existent competition
func TestCompetitionService_GetCompetitionByID_NotFound(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	retrieved, err := svc.GetCompetitionByID(uuid.New().String())
	if err != nil {
		t.Fatalf("GetCompetitionByID() error = %v", err)
	}
	if retrieved != nil {
		t.Error("GetCompetitionByID() should return nil for non-existent competition")
	}
}

// TestCompetitionService_CreateCompetition_Success tests successful creation
func TestCompetitionService_CreateCompetition_Success(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	competition := &models.Competition{
		CompetitionTitle: "New Competition",
		CompetitionType:  "Time Attack",
		ModelType:        "Buggy",
		ModelScale:       "1/10",
		TrackName:        "Test Track",
	}

	err := svc.CreateCompetition(competition)
	if err != nil {
		t.Fatalf("CreateCompetition() error = %v", err)
	}
	if competition.ID == "" {
		t.Error("CreateCompetition() should generate UUID")
	}
	if competition.Status != "scheduled" {
		t.Errorf("CreateCompetition() should set default status to 'scheduled', got '%s'", competition.Status)
	}

	// Verify it was added
	retrieved, _ := mockCompRepo.GetByID(competition.ID)
	if retrieved == nil || retrieved.CompetitionTitle != "New Competition" {
		t.Error("CreateCompetition() failed to store competition")
	}
}

// TestCompetitionService_CreateCompetition_EmptyTitle tests validation for empty title
func TestCompetitionService_CreateCompetition_EmptyTitle(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	competition := &models.Competition{
		CompetitionTitle: "",
		CompetitionType:  "Race",
	}

	err := svc.CreateCompetition(competition)
	if err == nil {
		t.Error("CreateCompetition() expected error for empty title, got nil")
	}
}

// TestCompetitionService_CreateCompetition_EmptyType tests validation for empty type
func TestCompetitionService_CreateCompetition_EmptyType(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	competition := &models.Competition{
		CompetitionTitle: "Valid Title",
		CompetitionType:  "",
	}

	err := svc.CreateCompetition(competition)
	if err == nil {
		t.Error("CreateCompetition() expected error for empty type, got nil")
	}
}

// TestCompetitionService_UpdateCompetition_Success tests successful update
func TestCompetitionService_UpdateCompetition_Success(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Create initial competition
	compID := uuid.New().String()
	timeStart := time.Now().Add(time.Hour)
	c := &models.Competition{
		ID:                compID,
		CompetitionTitle:  "Original Title",
		CompetitionType:   "Qualifying",
		ModelType:         "Buggy",
		ModelScale:        "1/10",
		TrackName:         "Track 1",
		Status:            "scheduled",
		TimeStart:         &timeStart,
		LapCountTarget:    intPtr(10),
		TimeLimitMinutes:  intPtr(5),
	}
	mockCompRepo.Create(c)

	// Update
	c.CompetitionTitle = "Updated Title"
	c.Status = "running"
	err := svc.UpdateCompetition(c)
	if err != nil {
		t.Fatalf("UpdateCompetition() error = %v", err)
	}

	// Verify update
	retrieved, _ := mockCompRepo.GetByID(c.ID)
	if retrieved == nil || retrieved.CompetitionTitle != "Updated Title" || retrieved.Status != "running" {
		t.Error("UpdateCompetition() failed to update competition")
	}
}

// TestCompetitionService_UpdateCompetition_EmptyID tests validation for empty ID
func TestCompetitionService_UpdateCompetition_EmptyID(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	c := &models.Competition{
		ID:               "",
		CompetitionTitle: "Test",
		CompetitionType:  "Race",
	}

	err := svc.UpdateCompetition(c)
	if err == nil {
		t.Error("UpdateCompetition() expected error for empty ID, got nil")
	}
}

// TestCompetitionService_UpdateCompetition_EmptyTitle tests validation for empty title on update
func TestCompetitionService_UpdateCompetition_EmptyTitle(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	c := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "",
		CompetitionType:  "Race",
	}
	mockCompRepo.Create(c)

	err := svc.UpdateCompetition(c)
	if err == nil {
		t.Error("UpdateCompetition() expected error for empty title, got nil")
	}
}

// TestCompetitionService_UpdateCompetition_EmptyType tests validation for empty type on update
func TestCompetitionService_UpdateCompetition_EmptyType(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	c := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Test Title",
		CompetitionType:  "",
	}
	mockCompRepo.Create(c)

	err := svc.UpdateCompetition(c)
	if err == nil {
		t.Error("UpdateCompetition() expected error for empty type, got nil")
	}
}

// TestCompetitionService_UpdateCompetition_NotFound tests update of non-existent competition
func TestCompetitionService_UpdateCompetition_NotFound(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	c := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Non-existent",
		CompetitionType:  "Race",
	}

	err := svc.UpdateCompetition(c)
	if err == nil {
		t.Error("UpdateCompetition() expected error for non-existent competition, got nil")
	}
}

// TestCompetitionService_DeleteCompetition_Success tests successful deletion
func TestCompetitionService_DeleteCompetition_Success(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Create competition
	c := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "To Delete",
		CompetitionType:  "Race",
	}
	mockCompRepo.Create(c)

	// Delete
	err := svc.DeleteCompetition(c.ID)
	if err != nil {
		t.Fatalf("DeleteCompetition() error = %v", err)
	}

	// Verify deletion
	retrieved, _ := mockCompRepo.GetByID(c.ID)
	if retrieved != nil {
		t.Error("DeleteCompetition() failed to delete competition")
	}
}

// TestCompetitionService_DeleteCompetition_EmptyID tests validation for empty ID on delete
func TestCompetitionService_DeleteCompetition_EmptyID(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	err := svc.DeleteCompetition("")
	if err == nil {
		t.Error("DeleteCompetition() expected error for empty ID, got nil")
	}
}

// TestCompetitionService_GetCompetitionCount tests the count functionality
func TestCompetitionService_GetCompetitionCount(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Empty count
	count, err := svc.GetCompetitionCount()
	if err != nil {
		t.Fatalf("GetCompetitionCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetCompetitionCount() expected 0, got %d", count)
	}

	// Add competitions
	for i := 1; i <= 5; i++ {
		c := &models.Competition{
			ID:               uuid.New().String(),
			CompetitionTitle: "Competition",
			CompetitionType:  "Race",
		}
		mockCompRepo.Create(c)
	}

	count, err = svc.GetCompetitionCount()
	if err != nil {
		t.Fatalf("GetCompetitionCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("GetCompetitionCount() expected 5, got %d", count)
	}
}

// TestCompetitionService_GetCompetitionsByStatus tests filtering by status
func TestCompetitionService_GetCompetitionsByStatus(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Add competitions with different statuses
	c1 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Scheduled Race",
		CompetitionType:  "Race",
		Status:           "scheduled",
	}
	c2 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Running Race",
		CompetitionType:  "Race",
		Status:           "running",
	}
	c3 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Completed Race",
		CompetitionType:  "Race",
		Status:           "completed",
	}
	c4 := &models.Competition{
		ID:               uuid.New().String(),
		CompetitionTitle: "Another Scheduled",
		CompetitionType:  "Qualifying",
		Status:           "scheduled",
	}
	mockCompRepo.Create(c1)
	mockCompRepo.Create(c2)
	mockCompRepo.Create(c3)
	mockCompRepo.Create(c4)

	// Test filtering by "scheduled" status
	scheduled, err := svc.GetCompetitionsByStatus("scheduled")
	if err != nil {
		t.Fatalf("GetCompetitionsByStatus() error = %v", err)
	}
	if len(scheduled) != 2 {
		t.Errorf("GetCompetitionsByStatus('scheduled') expected 2, got %d", len(scheduled))
	}

	// Test filtering by "running" status
	running, err := svc.GetCompetitionsByStatus("running")
	if err != nil {
		t.Fatalf("GetCompetitionsByStatus() error = %v", err)
	}
	if len(running) != 1 {
		t.Errorf("GetCompetitionsByStatus('running') expected 1, got %d", len(running))
	}

	// Test filtering by non-existent status
	empty, err := svc.GetCompetitionsByStatus("cancelled")
	if err != nil {
		t.Fatalf("GetCompetitionsByStatus() error = %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("GetCompetitionsByStatus('cancelled') expected 0, got %d", len(empty))
	}
}

// TestCompetitionService_GetAllModelTypes tests getting all model types
func TestCompetitionService_GetAllModelTypes(t *testing.T) {
	mockCompRepo := NewMockCompetitionRepository()
	mockTypeRepo := NewMockRCModelTypeRepositoryForCompetition()
	svc := service.NewCompetitionService(mockCompRepo, mockTypeRepo, NewMockRCModelScaleRepositoryForCompetition())

	// Add some model types
	mockTypeRepo.Create("Buggy")
	mockTypeRepo.Create("Truck")
	mockTypeRepo.Create("Touring Car")

	modelTypes, err := svc.GetAllModelTypes()
	if err != nil {
		t.Fatalf("GetAllModelTypes() error = %v", err)
	}
	if len(modelTypes) != 3 {
		t.Errorf("GetAllModelTypes() expected 3 model types, got %d", len(modelTypes))
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
