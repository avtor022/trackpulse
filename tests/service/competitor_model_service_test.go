package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// CompetitorModelRepository interface for testing
type CompetitorModelRepository interface {
	GetAll() ([]models.CompetitorModel, error)
	GetByID(id string) (*models.CompetitorModel, error)
	GetByTransponderNumber(number string) (*models.CompetitorModel, error)
	Create(cm *models.CompetitorModel) error
	Update(cm *models.CompetitorModel) error
	Delete(id string) error
	Count() (int, error)
}

// CompetitorRepository interface for testing
type CompetitorRepository interface {
	GetAll() ([]models.Competitor, error)
	GetByID(id string) (*models.Competitor, error)
	GetByNumber(number int) (*models.Competitor, error)
	Create(competitor *models.Competitor) error
	Update(competitor *models.Competitor) error
	Delete(id string) error
	Count() (int, error)
}

// RCModelRepository interface for testing
type RCModelRepository interface {
	GetAll() ([]models.RCModel, error)
	GetByID(id string) (*models.RCModel, error)
	Create(model *models.RCModel) error
	Update(model *models.RCModel) error
	Delete(id string) error
	Count() (int, error)
}

// MockCompetitorModelRepository is a mock implementation of CompetitorModelRepository for testing
type MockCompetitorModelRepository struct {
	competitorModels    map[string]*models.CompetitorModel
	transponderIndex    map[string]*models.CompetitorModel
}

func NewMockCompetitorModelRepository() *MockCompetitorModelRepository {
	return &MockCompetitorModelRepository{
		competitorModels: make(map[string]*models.CompetitorModel),
		transponderIndex: make(map[string]*models.CompetitorModel),
	}
}

func (m *MockCompetitorModelRepository) GetAll() ([]models.CompetitorModel, error) {
	result := make([]models.CompetitorModel, 0, len(m.competitorModels))
	for _, cm := range m.competitorModels {
		result = append(result, *cm)
	}
	return result, nil
}

func (m *MockCompetitorModelRepository) GetByID(id string) (*models.CompetitorModel, error) {
	if cm, ok := m.competitorModels[id]; ok {
		return cm, nil
	}
	return nil, nil
}

func (m *MockCompetitorModelRepository) GetByTransponderNumber(number string) (*models.CompetitorModel, error) {
	if cm, ok := m.transponderIndex[number]; ok {
		return cm, nil
	}
	return nil, nil
}

func (m *MockCompetitorModelRepository) Create(cm *models.CompetitorModel) error {
	m.competitorModels[cm.ID] = cm
	m.transponderIndex[cm.TransponderNumber] = cm
	return nil
}

func (m *MockCompetitorModelRepository) Update(cm *models.CompetitorModel) error {
	if _, ok := m.competitorModels[cm.ID]; !ok {
		return nil // Not found
	}
	m.competitorModels[cm.ID] = cm
	m.transponderIndex[cm.TransponderNumber] = cm
	return nil
}

func (m *MockCompetitorModelRepository) Delete(id string) error {
	if cm, ok := m.competitorModels[id]; ok {
		delete(m.transponderIndex, cm.TransponderNumber)
		delete(m.competitorModels, id)
	}
	return nil
}

func (m *MockCompetitorModelRepository) Count() (int, error) {
	return len(m.competitorModels), nil
}

// MockCompetitorRepositoryForCM is a mock for CompetitorRepository used in CompetitorModelService tests
type MockCompetitorRepositoryForCM struct {
	competitors map[string]*models.Competitor
}

func NewMockCompetitorRepositoryForCM() *MockCompetitorRepositoryForCM {
	return &MockCompetitorRepositoryForCM{
		competitors: make(map[string]*models.Competitor),
	}
}

func (m *MockCompetitorRepositoryForCM) GetAll() ([]models.Competitor, error) {
	result := make([]models.Competitor, 0, len(m.competitors))
	for _, c := range m.competitors {
		result = append(result, *c)
	}
	return result, nil
}

func (m *MockCompetitorRepositoryForCM) GetByID(id string) (*models.Competitor, error) {
	if c, ok := m.competitors[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCompetitorRepositoryForCM) GetByNumber(number int) (*models.Competitor, error) {
	for _, c := range m.competitors {
		if c.CompetitorNumber == number {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockCompetitorRepositoryForCM) Create(competitor *models.Competitor) error {
	m.competitors[competitor.ID] = competitor
	return nil
}

func (m *MockCompetitorRepositoryForCM) Update(competitor *models.Competitor) error {
	if _, ok := m.competitors[competitor.ID]; !ok {
		return nil
	}
	m.competitors[competitor.ID] = competitor
	return nil
}

func (m *MockCompetitorRepositoryForCM) Delete(id string) error {
	delete(m.competitors, id)
	return nil
}

func (m *MockCompetitorRepositoryForCM) Count() (int, error) {
	return len(m.competitors), nil
}

// MockRCModelRepositoryForCM is a mock for RCModelRepository used in CompetitorModelService tests
type MockRCModelRepositoryForCM struct {
	models map[string]*models.RCModel
}

func NewMockRCModelRepositoryForCM() *MockRCModelRepositoryForCM {
	return &MockRCModelRepositoryForCM{
		models: make(map[string]*models.RCModel),
	}
}

func (m *MockRCModelRepositoryForCM) GetAll() ([]models.RCModel, error) {
	result := make([]models.RCModel, 0, len(m.models))
	for _, model := range m.models {
		result = append(result, *model)
	}
	return result, nil
}

func (m *MockRCModelRepositoryForCM) GetByID(id string) (*models.RCModel, error) {
	if model, ok := m.models[id]; ok {
		return model, nil
	}
	return nil, nil
}

func (m *MockRCModelRepositoryForCM) Create(model *models.RCModel) error {
	m.models[model.ID] = model
	return nil
}

func (m *MockRCModelRepositoryForCM) Update(model *models.RCModel) error {
	if _, ok := m.models[model.ID]; !ok {
		return nil
	}
	m.models[model.ID] = model
	return nil
}

func (m *MockRCModelRepositoryForCM) Delete(id string) error {
	delete(m.models, id)
	return nil
}

func (m *MockRCModelRepositoryForCM) Count() (int, error) {
	return len(m.models), nil
}

// NewCompetitorModelService creates a new competitor model service for testing
func NewCompetitorModelService(repo CompetitorModelRepository, competitorRepo CompetitorRepository, modelRepo RCModelRepository) *CompetitorModelService {
	return &CompetitorModelService{repo: repo, competitorRepo: competitorRepo, modelRepo: modelRepo}
}

// CompetitorModelService handles business logic for competitor models (transponders)
type CompetitorModelService struct {
	repo           CompetitorModelRepository
	competitorRepo CompetitorRepository
	modelRepo      RCModelRepository
}

// GetAllCompetitorModels returns all competitor models
func (s *CompetitorModelService) GetAllCompetitorModels() ([]models.CompetitorModel, error) {
	return s.repo.GetAll()
}

// GetCompetitorModelByID returns a competitor model by ID
func (s *CompetitorModelService) GetCompetitorModelByID(id string) (*models.CompetitorModel, error) {
	return s.repo.GetByID(id)
}

// CreateCompetitorModel creates a new competitor model with validation
func (s *CompetitorModelService) CreateCompetitorModel(cm *models.CompetitorModel) error {
	// Validate required fields
	if cm.CompetitorID == "" {
		return fmt.Errorf("competitor is required")
	}
	if cm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if cm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if competitor exists
	competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
	if err != nil {
		return fmt.Errorf("failed to get competitor: %w", err)
	}
	if competitor == nil {
		return fmt.Errorf("competitor not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(cm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number already exists
	existing, err := s.repo.GetByTransponderNumber(cm.TransponderNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing transponder: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("transponder number '%s' already exists", cm.TransponderNumber)
	}

	// Set default values
	if cm.TransponderType == "" {
		cm.TransponderType = "RFID"
	}
	cm.IsActive = true

	// Generate UUID
	cm.ID = uuid.New().String()

	return s.repo.Create(cm)
}

// UpdateCompetitorModel updates an existing competitor model with validation
func (s *CompetitorModelService) UpdateCompetitorModel(cm *models.CompetitorModel) error {
	// Validate required fields
	if cm.ID == "" {
		return fmt.Errorf("competitor model ID is required")
	}
	if cm.CompetitorID == "" {
		return fmt.Errorf("competitor is required")
	}
	if cm.RCModelID == "" {
		return fmt.Errorf("RC model is required")
	}
	if cm.TransponderNumber == "" {
		return fmt.Errorf("transponder number is required")
	}

	// Check if competitor model exists
	existing, err := s.repo.GetByID(cm.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing competitor model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("competitor model not found")
	}

	// Check if competitor exists
	competitor, err := s.competitorRepo.GetByID(cm.CompetitorID)
	if err != nil {
		return fmt.Errorf("failed to get competitor: %w", err)
	}
	if competitor == nil {
		return fmt.Errorf("competitor not found")
	}

	// Check if RC model exists
	rcModel, err := s.modelRepo.GetByID(cm.RCModelID)
	if err != nil {
		return fmt.Errorf("failed to get RC model: %w", err)
	}
	if rcModel == nil {
		return fmt.Errorf("RC model not found")
	}

	// Check if transponder number is taken by another competitor model
	if existing.TransponderNumber != cm.TransponderNumber {
		transponderExists, err := s.repo.GetByTransponderNumber(cm.TransponderNumber)
		if err != nil {
			return fmt.Errorf("failed to check transponder number: %w", err)
		}
		if transponderExists != nil && transponderExists.ID != cm.ID {
			return fmt.Errorf("transponder number '%s' already exists", cm.TransponderNumber)
		}
	}

	return s.repo.Update(cm)
}

// DeleteCompetitorModel deletes a competitor model by ID
func (s *CompetitorModelService) DeleteCompetitorModel(id string) error {
	if id == "" {
		return fmt.Errorf("competitor model ID is required")
	}
	return s.repo.Delete(id)
}

// GetCompetitorModelCount returns total number of competitor models
func (s *CompetitorModelService) GetCompetitorModelCount() (int, error) {
	return s.repo.Count()
}

func TestCompetitorModelService_GetAllCompetitorModels(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Add some test data
	cm1 := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP001",
		TransponderType:   "RFID",
		IsActive:          true,
	}
	cm2 := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP002",
		TransponderType:   "NFC",
		IsActive:          true,
	}
	mockRepo.Create(cm1)
	mockRepo.Create(cm2)

	competitorModels, err := svc.GetAllCompetitorModels()
	if err != nil {
		t.Fatalf("GetAllCompetitorModels() error = %v", err)
	}
	if len(competitorModels) != 2 {
		t.Errorf("GetAllCompetitorModels() expected 2 competitor models, got %d", len(competitorModels))
	}
}

func TestCompetitorModelService_CreateCompetitorModel_Success(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create prerequisite competitor and RC model
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID:               competitorID,
		CompetitorNumber: 1,
		FullName:         "Test Pilot",
		Country:          "USA",
		City:             "Miami",
	}
	mockCompetitorRepo.Create(competitor)

	rcModelID := uuid.New().String()
	rcModel := &models.RCModel{
		ID:        rcModelID,
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(rcModel)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      competitorID,
		RCModelID:         rcModelID,
		TransponderNumber: "TP100",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err != nil {
		t.Fatalf("CreateCompetitorModel() error = %v", err)
	}
	if competitorModel.ID == "" {
		t.Error("CreateCompetitorModel() should generate UUID")
	}
	if competitorModel.TransponderType != "RFID" {
		t.Errorf("CreateCompetitorModel() should default TransponderType to RFID, got %s", competitorModel.TransponderType)
	}
	if !competitorModel.IsActive {
		t.Error("CreateCompetitorModel() should set IsActive to true")
	}

	// Verify it was added
	retrieved, _ := mockRepo.GetByID(competitorModel.ID)
	if retrieved == nil || retrieved.TransponderNumber != "TP100" {
		t.Error("CreateCompetitorModel() failed to store competitor model")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_EmptyCompetitorID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      "",
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP100",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for empty competitor ID, got nil")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_EmptyRCModelID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      uuid.New().String(),
		RCModelID:         "",
		TransponderNumber: "TP100",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for empty RC model ID, got nil")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_EmptyTransponderNumber(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for empty transponder number, got nil")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_CompetitorNotFound(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create RC model but not competitor
	rcModelID := uuid.New().String()
	rcModel := &models.RCModel{
		ID:        rcModelID,
		Brand:     "Arrma",
		ModelName: "Kraton",
	}
	mockModelRepo.Create(rcModel)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      uuid.New().String(),
		RCModelID:         rcModelID,
		TransponderNumber: "TP100",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for non-existent competitor, got nil")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_RCModelNotFound(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create competitor but not RC model
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID:       competitorID,
		FullName: "Test Pilot",
	}
	mockCompetitorRepo.Create(competitor)

	competitorModel := &models.CompetitorModel{
		CompetitorID:      competitorID,
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP100",
	}

	err := svc.CreateCompetitorModel(competitorModel)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for non-existent RC model, got nil")
	}
}

func TestCompetitorModelService_CreateCompetitorModel_DuplicateTransponderNumber(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create prerequisite entities
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID:       competitorID,
		FullName: "Test Pilot",
	}
	mockCompetitorRepo.Create(competitor)

	rcModelID1 := uuid.New().String()
	rcModel1 := &models.RCModel{
		ID:        rcModelID1,
		Brand:     "Team Associated",
		ModelName: "RC10",
	}
	mockModelRepo.Create(rcModel1)

	rcModelID2 := uuid.New().String()
	rcModel2 := &models.RCModel{
		ID:        rcModelID2,
		Brand:     "Yokomo",
		ModelName: "YZ-2",
	}
	mockModelRepo.Create(rcModel2)

	// Create first competitor model
	cm1 := &models.CompetitorModel{
		CompetitorID:      competitorID,
		RCModelID:         rcModelID1,
		TransponderNumber: "TP200",
	}
	mockRepo.Create(cm1)

	// Try to create second with same transponder number
	cm2 := &models.CompetitorModel{
		CompetitorID:      competitorID,
		RCModelID:         rcModelID2,
		TransponderNumber: "TP200",
	}

	err := svc.CreateCompetitorModel(cm2)
	if err == nil {
		t.Error("CreateCompetitorModel() expected error for duplicate transponder number, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_Success(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create prerequisite entities
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID:       competitorID,
		FullName: "Test Pilot",
	}
	mockCompetitorRepo.Create(competitor)

	rcModelID := uuid.New().String()
	rcModel := &models.RCModel{
		ID:        rcModelID,
		Brand:     "Xray",
		ModelName: "XB8",
	}
	mockModelRepo.Create(rcModel)

	// Create initial competitor model
	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      competitorID,
		RCModelID:         rcModelID,
		TransponderNumber: "TP300",
		TransponderType:   "RFID",
		IsActive:          true,
	}
	mockRepo.Create(cm)

	// Update
	cm.TransponderNumber = "TP301"
	cm.TransponderType = "NFC"
	cm.IsActive = false
	err := svc.UpdateCompetitorModel(cm)
	if err != nil {
		t.Fatalf("UpdateCompetitorModel() error = %v", err)
	}

	// Verify update
	retrieved, _ := mockRepo.GetByID(cm.ID)
	if retrieved == nil || retrieved.TransponderNumber != "TP301" || retrieved.TransponderType != "NFC" || retrieved.IsActive {
		t.Error("UpdateCompetitorModel() failed to update competitor model")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_EmptyID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	cm := &models.CompetitorModel{
		ID:                "",
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP400",
	}

	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for empty ID, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_EmptyCompetitorID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      "",
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP400",
	}

	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for empty competitor ID, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_EmptyRCModelID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         "",
		TransponderNumber: "TP400",
	}

	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for empty RC model ID, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_EmptyTransponderNumber(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "",
	}

	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for empty transponder number, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_NotFound(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP500",
	}

	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for non-existent competitor model, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_CompetitorNotFound(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create RC model but not competitor
	rcModelID := uuid.New().String()
	rcModel := &models.RCModel{
		ID: rcModelID,
	}
	mockModelRepo.Create(rcModel)

	// Create competitor model in repo (simulating existing record)
	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         rcModelID,
		TransponderNumber: "TP600",
	}
	mockRepo.Create(cm)

	// Try to update with non-existent competitor
	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for non-existent competitor, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_RCModelNotFound(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create competitor but not RC model
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID: competitorID,
	}
	mockCompetitorRepo.Create(competitor)

	// Create competitor model in repo (simulating existing record)
	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      competitorID,
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP700",
	}
	mockRepo.Create(cm)

	// Try to update with non-existent RC model
	err := svc.UpdateCompetitorModel(cm)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for non-existent RC model, got nil")
	}
}

func TestCompetitorModelService_UpdateCompetitorModel_DuplicateTransponderNumber(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create prerequisite entities
	competitorID := uuid.New().String()
	competitor := &models.Competitor{
		ID: competitorID,
	}
	mockCompetitorRepo.Create(competitor)

	rcModelID1 := uuid.New().String()
	rcModel1 := &models.RCModel{
		ID: rcModelID1,
	}
	mockModelRepo.Create(rcModel1)

	rcModelID2 := uuid.New().String()
	rcModel2 := &models.RCModel{
		ID: rcModelID2,
	}
	mockModelRepo.Create(rcModel2)

	// Create two competitor models with different transponder numbers
	cm1 := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      competitorID,
		RCModelID:         rcModelID1,
		TransponderNumber: "TP800",
	}
	mockRepo.Create(cm1)

	cm2 := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      competitorID,
		RCModelID:         rcModelID2,
		TransponderNumber: "TP801",
	}
	mockRepo.Create(cm2)

	// Try to update cm2 to use cm1's transponder number (create a new struct to avoid modifying the stored one)
	cm2Update := &models.CompetitorModel{
		ID:                cm2.ID,
		CompetitorID:      competitorID,
		RCModelID:         rcModelID2,
		TransponderNumber: "TP800", // Try to use cm1's transponder number
	}
	err := svc.UpdateCompetitorModel(cm2Update)
	if err == nil {
		t.Error("UpdateCompetitorModel() expected error for duplicate transponder number, got nil")
	}
}

func TestCompetitorModelService_DeleteCompetitorModel_Success(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Create competitor model
	cm := &models.CompetitorModel{
		ID:                uuid.New().String(),
		CompetitorID:      uuid.New().String(),
		RCModelID:         uuid.New().String(),
		TransponderNumber: "TP900",
	}
	mockRepo.Create(cm)

	// Delete
	err := svc.DeleteCompetitorModel(cm.ID)
	if err != nil {
		t.Fatalf("DeleteCompetitorModel() error = %v", err)
	}

	// Verify deletion
	retrieved, _ := mockRepo.GetByID(cm.ID)
	if retrieved != nil {
		t.Error("DeleteCompetitorModel() failed to delete competitor model")
	}
}

func TestCompetitorModelService_DeleteCompetitorModel_EmptyID(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	err := svc.DeleteCompetitorModel("")
	if err == nil {
		t.Error("DeleteCompetitorModel() expected error for empty ID, got nil")
	}
}

func TestCompetitorModelService_GetCompetitorModelCount(t *testing.T) {
	mockRepo := NewMockCompetitorModelRepository()
	mockCompetitorRepo := NewMockCompetitorRepositoryForCM()
	mockModelRepo := NewMockRCModelRepositoryForCM()
	
	svc := NewCompetitorModelService(mockRepo, mockCompetitorRepo, mockModelRepo)

	// Empty count
	count, err := svc.GetCompetitorModelCount()
	if err != nil {
		t.Fatalf("GetCompetitorModelCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetCompetitorModelCount() expected 0, got %d", count)
	}

	// Add competitor models
	for i := 1; i <= 5; i++ {
		cm := &models.CompetitorModel{
			ID:                uuid.New().String(),
			CompetitorID:      uuid.New().String(),
			RCModelID:         uuid.New().String(),
			TransponderNumber: time.Now().Format(time.RFC3339Nano),
		}
		mockRepo.Create(cm)
	}

	count, err = svc.GetCompetitorModelCount()
	if err != nil {
		t.Fatalf("GetCompetitorModelCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("GetCompetitorModelCount() expected 5, got %d", count)
	}
}
