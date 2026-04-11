package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// MockCompetitorRepository is a mock implementation of CompetitorRepository for testing
type MockCompetitorRepository struct {
	competitors map[string]*models.Competitor
	numberIndex map[int]*models.Competitor
}

func NewMockCompetitorRepository() *MockCompetitorRepository {
	return &MockCompetitorRepository{
		competitors: make(map[string]*models.Competitor),
		numberIndex: make(map[int]*models.Competitor),
	}
}

func (m *MockCompetitorRepository) GetAll() ([]models.Competitor, error) {
	result := make([]models.Competitor, 0, len(m.competitors))
	for _, c := range m.competitors {
		result = append(result, *c)
	}
	return result, nil
}

func (m *MockCompetitorRepository) GetByID(id string) (*models.Competitor, error) {
	if c, ok := m.competitors[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCompetitorRepository) GetByNumber(number int) (*models.Competitor, error) {
	if c, ok := m.numberIndex[number]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCompetitorRepository) Create(competitor *models.Competitor) error {
	m.competitors[competitor.ID] = competitor
	m.numberIndex[competitor.CompetitorNumber] = competitor
	return nil
}

func (m *MockCompetitorRepository) Update(competitor *models.Competitor) error {
	if _, ok := m.competitors[competitor.ID]; !ok {
		return nil // Not found
	}
	m.competitors[competitor.ID] = competitor
	m.numberIndex[competitor.CompetitorNumber] = competitor
	return nil
}

func (m *MockCompetitorRepository) Delete(id string) error {
	if c, ok := m.competitors[id]; ok {
		delete(m.numberIndex, c.CompetitorNumber)
		delete(m.competitors, id)
	}
	return nil
}

func (m *MockCompetitorRepository) Count() (int, error) {
	return len(m.competitors), nil
}

func TestCompetitorService_GetAllCompetitors(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	// Add some test data
	birthday := time.Now().AddDate(-10, 0, 0)
	c1 := &models.Competitor{
		ID:               uuid.New().String(),
		CompetitorNumber: 1,
		FullName:         "John Doe",
		Country:          "USA",
		City:             "New York",
		Birthday:         &birthday,
		Rating:           1500,
	}
	c2 := &models.Competitor{
		ID:               uuid.New().String(),
		CompetitorNumber: 2,
		FullName:         "Jane Smith",
		Country:          "UK",
		City:             "London",
		Rating:           1600,
	}
	mockRepo.Create(c1)
	mockRepo.Create(c2)

	competitors, err := svc.GetAllCompetitors()
	if err != nil {
		t.Fatalf("GetAllCompetitors() error = %v", err)
	}
	if len(competitors) != 2 {
		t.Errorf("GetAllCompetitors() expected 2 competitors, got %d", len(competitors))
	}
}

func TestCompetitorService_CreateCompetitor_Success(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	birthday := time.Now().AddDate(-15, 0, 0)
	competitor := &models.Competitor{
		CompetitorNumber: 1,
		FullName:         "Test User",
		Country:          "USA",
		City:             "Boston",
		Birthday:         &birthday,
		Rating:           1200,
	}

	err := svc.CreateCompetitor(competitor)
	if err != nil {
		t.Fatalf("CreateCompetitor() error = %v", err)
	}
	if competitor.ID == "" {
		t.Error("CreateCompetitor() should generate UUID")
	}

	// Verify it was added
	retrieved, _ := mockRepo.GetByID(competitor.ID)
	if retrieved == nil || retrieved.FullName != "Test User" {
		t.Error("CreateCompetitor() failed to store competitor")
	}
}

func TestCompetitorService_CreateCompetitor_EmptyName(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	competitor := &models.Competitor{
		CompetitorNumber: 1,
		FullName:         "",
		Country:          "USA",
	}

	err := svc.CreateCompetitor(competitor)
	if err == nil {
		t.Error("CreateCompetitor() expected error for empty name, got nil")
	}
}

func TestCompetitorService_CreateCompetitor_InvalidNumber(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	competitor := &models.Competitor{
		CompetitorNumber: 0,
		FullName:         "Test User",
	}

	err := svc.CreateCompetitor(competitor)
	if err == nil {
		t.Error("CreateCompetitor() expected error for invalid number, got nil")
	}

	competitor.CompetitorNumber = -1
	err = svc.CreateCompetitor(competitor)
	if err == nil {
		t.Error("CreateCompetitor() expected error for negative number, got nil")
	}
}

func TestCompetitorService_CreateCompetitor_DuplicateNumber(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	// Create first competitor
	c1 := &models.Competitor{
		CompetitorNumber: 5,
		FullName:         "First User",
	}
	mockRepo.Create(c1)

	// Try to create second with same number
	c2 := &models.Competitor{
		CompetitorNumber: 5,
		FullName:         "Second User",
	}

	err := svc.CreateCompetitor(c2)
	if err == nil {
		t.Error("CreateCompetitor() expected error for duplicate number, got nil")
	}
}

func TestCompetitorService_UpdateCompetitor_Success(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	// Create initial competitor
	birthday := time.Now().AddDate(-20, 0, 0)
	c := &models.Competitor{
		ID:               uuid.New().String(),
		CompetitorNumber: 10,
		FullName:         "Original Name",
		Country:          "USA",
		City:             "Chicago",
		Birthday:         &birthday,
		Rating:           1000,
	}
	mockRepo.Create(c)

	// Update
	c.FullName = "Updated Name"
	c.Rating = 1500
	err := svc.UpdateCompetitor(c)
	if err != nil {
		t.Fatalf("UpdateCompetitor() error = %v", err)
	}

	// Verify update
	retrieved, _ := mockRepo.GetByID(c.ID)
	if retrieved == nil || retrieved.FullName != "Updated Name" || retrieved.Rating != 1500 {
		t.Error("UpdateCompetitor() failed to update competitor")
	}
}

func TestCompetitorService_UpdateCompetitor_EmptyID(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	c := &models.Competitor{
		ID:               "",
		CompetitorNumber: 1,
		FullName:         "Test",
	}

	err := svc.UpdateCompetitor(c)
	if err == nil {
		t.Error("UpdateCompetitor() expected error for empty ID, got nil")
	}
}

func TestCompetitorService_UpdateCompetitor_NotFound(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	c := &models.Competitor{
		ID:               uuid.New().String(),
		CompetitorNumber: 999,
		FullName:         "Non-existent",
	}

	err := svc.UpdateCompetitor(c)
	if err == nil {
		t.Error("UpdateCompetitor() expected error for non-existent competitor, got nil")
	}
}

func TestCompetitorService_DeleteCompetitor_Success(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	// Create competitor
	c := &models.Competitor{
		ID:               uuid.New().String(),
		CompetitorNumber: 7,
		FullName:         "To Delete",
	}
	mockRepo.Create(c)

	// Delete
	err := svc.DeleteCompetitor(c.ID)
	if err != nil {
		t.Fatalf("DeleteCompetitor() error = %v", err)
	}

	// Verify deletion
	retrieved, _ := mockRepo.GetByID(c.ID)
	if retrieved != nil {
		t.Error("DeleteCompetitor() failed to delete competitor")
	}
}

func TestCompetitorService_DeleteCompetitor_EmptyID(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	err := svc.DeleteCompetitor("")
	if err == nil {
		t.Error("DeleteCompetitor() expected error for empty ID, got nil")
	}
}

func TestCompetitorService_GetCompetitorCount(t *testing.T) {
	mockRepo := NewMockCompetitorRepository()
	_ = mockRepo
	svc := service.NewCompetitorService(mockRepo)

	// Empty count
	count, err := svc.GetCompetitorCount()
	if err != nil {
		t.Fatalf("GetCompetitorCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetCompetitorCount() expected 0, got %d", count)
	}

	// Add competitors
	for i := 1; i <= 5; i++ {
		c := &models.Competitor{
			ID:               uuid.New().String(),
			CompetitorNumber: i,
			FullName:         "User",
		}
		mockRepo.Create(c)
	}

	count, err = svc.GetCompetitorCount()
	if err != nil {
		t.Fatalf("GetCompetitorCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("GetCompetitorCount() expected 5, got %d", count)
	}
}
