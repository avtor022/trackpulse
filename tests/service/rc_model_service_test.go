package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RCModelBrandRepository interface for testing
type RCModelBrandRepository interface {
	GetAll() ([]models.RCModelBrand, error)
	GetByName(name string) (*models.RCModelBrand, error)
	Create(name string) (*models.RCModelBrand, error)
	GetOrCreate(name string) (*models.RCModelBrand, error)
	Delete(name string) error
}

// RCModelScaleRepository interface for testing
type RCModelScaleRepository interface {
	GetAll() ([]models.RCModelScale, error)
	GetByName(name string) (*models.RCModelScale, error)
	Create(name string) (*models.RCModelScale, error)
	GetOrCreate(name string) (*models.RCModelScale, error)
	Delete(name string) error
}

// RCModelTypeRepository interface for testing
type RCModelTypeRepository interface {
	GetAll() ([]models.RCModelType, error)
	GetByName(name string) (*models.RCModelType, error)
	Create(name string) (*models.RCModelType, error)
	GetOrCreate(name string) (*models.RCModelType, error)
	Delete(name string) error
}

// MockRCModelRepository is a mock implementation of RCModelRepository for testing
type MockRCModelRepository struct {
	models map[string]*models.RCModel
}

func NewMockRCModelRepository() *MockRCModelRepository {
	return &MockRCModelRepository{
		models: make(map[string]*models.RCModel),
	}
}

func (m *MockRCModelRepository) GetAll() ([]models.RCModel, error) {
	result := make([]models.RCModel, 0, len(m.models))
	for _, model := range m.models {
		result = append(result, *model)
	}
	return result, nil
}

func (m *MockRCModelRepository) GetByID(id string) (*models.RCModel, error) {
	if model, ok := m.models[id]; ok {
		return model, nil
	}
	return nil, nil
}

func (m *MockRCModelRepository) Create(model *models.RCModel) error {
	m.models[model.ID] = model
	return nil
}

func (m *MockRCModelRepository) Update(model *models.RCModel) error {
	if _, ok := m.models[model.ID]; !ok {
		return nil // Not found
	}
	m.models[model.ID] = model
	return nil
}

func (m *MockRCModelRepository) Delete(id string) error {
	delete(m.models, id)
	return nil
}

func (m *MockRCModelRepository) Count() (int, error) {
	return len(m.models), nil
}

func (m *MockRCModelRepository) GetAllModelNames() ([]string, error) {
	names := make([]string, 0, len(m.models))
	for _, model := range m.models {
		names = append(names, model.ModelName)
	}
	return names, nil
}

// MockRCModelBrandRepository is a mock implementation of RCModelBrandRepository for testing
type MockRCModelBrandRepository struct {
	brands map[string]*models.RCModelBrand
}

func NewMockRCModelBrandRepository() *MockRCModelBrandRepository {
	return &MockRCModelBrandRepository{
		brands: make(map[string]*models.RCModelBrand),
	}
}

func (m *MockRCModelBrandRepository) GetAll() ([]models.RCModelBrand, error) {
	result := make([]models.RCModelBrand, 0, len(m.brands))
	for _, brand := range m.brands {
		result = append(result, *brand)
	}
	return result, nil
}

func (m *MockRCModelBrandRepository) GetByName(name string) (*models.RCModelBrand, error) {
	if brand, ok := m.brands[name]; ok {
		return brand, nil
	}
	return nil, nil
}

func (m *MockRCModelBrandRepository) Create(name string) (*models.RCModelBrand, error) {
	id := uuid.New().String()
	brand := &models.RCModelBrand{
		ID:   id,
		Name: name,
	}
	m.brands[name] = brand
	return brand, nil
}

func (m *MockRCModelBrandRepository) GetOrCreate(name string) (*models.RCModelBrand, error) {
	if brand, ok := m.brands[name]; ok {
		return brand, nil
	}
	return m.Create(name)
}

func (m *MockRCModelBrandRepository) Delete(name string) error {
	delete(m.brands, name)
	return nil
}

// MockRCModelScaleRepository is a mock implementation of RCModelScaleRepository for testing
type MockRCModelScaleRepository struct {
	scales map[string]*models.RCModelScale
}

func NewMockRCModelScaleRepository() *MockRCModelScaleRepository {
	return &MockRCModelScaleRepository{
		scales: make(map[string]*models.RCModelScale),
	}
}

func (m *MockRCModelScaleRepository) GetAll() ([]models.RCModelScale, error) {
	result := make([]models.RCModelScale, 0, len(m.scales))
	for _, scale := range m.scales {
		result = append(result, *scale)
	}
	return result, nil
}

func (m *MockRCModelScaleRepository) GetByName(name string) (*models.RCModelScale, error) {
	if scale, ok := m.scales[name]; ok {
		return scale, nil
	}
	return nil, nil
}

func (m *MockRCModelScaleRepository) Create(name string) (*models.RCModelScale, error) {
	id := uuid.New().String()
	scale := &models.RCModelScale{
		ID:   id,
		Name: name,
	}
	m.scales[name] = scale
	return scale, nil
}

func (m *MockRCModelScaleRepository) GetOrCreate(name string) (*models.RCModelScale, error) {
	if scale, ok := m.scales[name]; ok {
		return scale, nil
	}
	return m.Create(name)
}

func (m *MockRCModelScaleRepository) Delete(name string) error {
	delete(m.scales, name)
	return nil
}

// MockRCModelTypeRepository is a mock implementation of RCModelTypeRepository for testing
type MockRCModelTypeRepository struct {
	types map[string]*models.RCModelType
}

func NewMockRCModelTypeRepository() *MockRCModelTypeRepository {
	return &MockRCModelTypeRepository{
		types: make(map[string]*models.RCModelType),
	}
}

func (m *MockRCModelTypeRepository) GetAll() ([]models.RCModelType, error) {
	result := make([]models.RCModelType, 0, len(m.types))
	for _, t := range m.types {
		result = append(result, *t)
	}
	return result, nil
}

func (m *MockRCModelTypeRepository) GetByName(name string) (*models.RCModelType, error) {
	if t, ok := m.types[name]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockRCModelTypeRepository) Create(name string) (*models.RCModelType, error) {
	id := uuid.New().String()
	t := &models.RCModelType{
		ID:   id,
		Name: name,
	}
	m.types[name] = t
	return t, nil
}

func (m *MockRCModelTypeRepository) GetOrCreate(name string) (*models.RCModelType, error) {
	if t, ok := m.types[name]; ok {
		return t, nil
	}
	return m.Create(name)
}

func (m *MockRCModelTypeRepository) Delete(name string) error {
	delete(m.types, name)
	return nil
}

// NewRCModelServiceForTest creates a new RC model service for testing
func NewRCModelServiceForTest(modelRepo *MockRCModelRepository, brandRepo *MockRCModelBrandRepository, scaleRepo *MockRCModelScaleRepository, typeRepo *MockRCModelTypeRepository) *RCModelService {
	return &RCModelService{modelRepo: modelRepo, brandRepo: brandRepo, scaleRepo: scaleRepo, typeRepo: typeRepo}
}

// RCModelService handles business logic for RC models (for testing purposes)
type RCModelService struct {
	modelRepo   *MockRCModelRepository
	brandRepo   *MockRCModelBrandRepository
	scaleRepo   *MockRCModelScaleRepository
	typeRepo    *MockRCModelTypeRepository
}

// GetAllModels returns all RC models
func (s *RCModelService) GetAllModels() ([]models.RCModel, error) {
	return s.modelRepo.GetAll()
}

// GetModelByID returns an RC model by ID
func (s *RCModelService) GetModelByID(id string) (*models.RCModel, error) {
	return s.modelRepo.GetByID(id)
}

// CreateModel creates a new RC model with validation
func (s *RCModelService) CreateModel(model *models.RCModel) error {
	// Validate required fields
	if model.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if model.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	if model.Scale == "" {
		return fmt.Errorf("scale is required")
	}
	if model.ModelType == "" {
		return fmt.Errorf("model type is required")
	}

	// Ensure brand exists in the brands table
	_, err := s.brandRepo.GetOrCreate(model.Brand)
	if err != nil {
		return fmt.Errorf("failed to ensure brand exists: %w", err)
	}

	// Ensure scale exists in the scales table
	_, err = s.scaleRepo.GetOrCreate(model.Scale)
	if err != nil {
		return fmt.Errorf("failed to ensure scale exists: %w", err)
	}

	// Ensure model type exists in the types table
	_, err = s.typeRepo.GetOrCreate(model.ModelType)
	if err != nil {
		return fmt.Errorf("failed to ensure model type exists: %w", err)
	}

	// Generate UUID
	model.ID = uuid.New().String()

	return s.modelRepo.Create(model)
}

// UpdateModel updates an existing RC model with validation
func (s *RCModelService) UpdateModel(model *models.RCModel) error {
	// Validate required fields
	if model.ID == "" {
		return fmt.Errorf("model ID is required")
	}
	if model.Brand == "" {
		return fmt.Errorf("brand is required")
	}
	if model.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	if model.Scale == "" {
		return fmt.Errorf("scale is required")
	}
	if model.ModelType == "" {
		return fmt.Errorf("model type is required")
	}

	// Check if model exists
	existing, err := s.modelRepo.GetByID(model.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing model: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("model not found")
	}

	// Ensure brand exists in the brands table
	_, err = s.brandRepo.GetOrCreate(model.Brand)
	if err != nil {
		return fmt.Errorf("failed to ensure brand exists: %w", err)
	}

	// Ensure scale exists in the scales table
	_, err = s.scaleRepo.GetOrCreate(model.Scale)
	if err != nil {
		return fmt.Errorf("failed to ensure scale exists: %w", err)
	}

	// Ensure model type exists in the types table
	_, err = s.typeRepo.GetOrCreate(model.ModelType)
	if err != nil {
		return fmt.Errorf("failed to ensure model type exists: %w", err)
	}

	return s.modelRepo.Update(model)
}

// DeleteModel deletes an RC model by ID
func (s *RCModelService) DeleteModel(id string) error {
	if id == "" {
		return fmt.Errorf("model ID is required")
	}
	return s.modelRepo.Delete(id)
}

// GetModelCount returns total number of RC models
func (s *RCModelService) GetModelCount() (int, error) {
	return s.modelRepo.Count()
}

// GetAllBrands returns all RC model brands
func (s *RCModelService) GetAllBrands() ([]models.RCModelBrand, error) {
	return s.brandRepo.GetAll()
}

// GetAllModelNames returns all unique model names
func (s *RCModelService) GetAllModelNames() ([]string, error) {
	return s.modelRepo.GetAllModelNames()
}

// AddBrand adds a new brand to the database
func (s *RCModelService) AddBrand(name string) error {
	if name == "" {
		return fmt.Errorf("brand name is required")
	}

	// Check if brand already exists
	existing, err := s.brandRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing brand: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("brand '%s' already exists", name)
	}

	// Create new brand
	_, err = s.brandRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create brand: %w", err)
	}

	return nil
}

// DeleteBrand deletes a brand from the database
func (s *RCModelService) DeleteBrand(name string) error {
	if name == "" {
		return fmt.Errorf("brand name is required")
	}

	// Check if brand is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.Brand == name {
			return fmt.Errorf("cannot delete brand '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete brand
	err = s.brandRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	return nil
}

// GetAllScales returns all RC model scales
func (s *RCModelService) GetAllScales() ([]models.RCModelScale, error) {
	return s.scaleRepo.GetAll()
}

// AddScale adds a new scale to the database
func (s *RCModelService) AddScale(name string) error {
	if name == "" {
		return fmt.Errorf("scale name is required")
	}

	// Check if scale already exists
	existing, err := s.scaleRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing scale: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("scale '%s' already exists", name)
	}

	// Create new scale
	_, err = s.scaleRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create scale: %w", err)
	}

	return nil
}

// DeleteScale deletes a scale from the database
func (s *RCModelService) DeleteScale(name string) error {
	if name == "" {
		return fmt.Errorf("scale name is required")
	}

	// Check if scale is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.Scale == name {
			return fmt.Errorf("cannot delete scale '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete scale
	err = s.scaleRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete scale: %w", err)
	}

	return nil
}

// GetAllModelTypes returns all RC model types
func (s *RCModelService) GetAllModelTypes() ([]models.RCModelType, error) {
	return s.typeRepo.GetAll()
}

// AddModelType adds a new model type to the database
func (s *RCModelService) AddModelType(name string) error {
	if name == "" {
		return fmt.Errorf("model type name is required")
	}

	// Check if model type already exists
	existing, err := s.typeRepo.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to check existing model type: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("model type '%s' already exists", name)
	}

	// Create new model type
	_, err = s.typeRepo.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create model type: %w", err)
	}

	return nil
}

// DeleteModelType deletes a model type from the database
func (s *RCModelService) DeleteModelType(name string) error {
	if name == "" {
		return fmt.Errorf("model type name is required")
	}

	// Check if model type is used by any models
	models, err := s.modelRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check models: %w", err)
	}

	for _, model := range models {
		if model.ModelType == name {
			return fmt.Errorf("cannot delete model type '%s': it is used by model '%s'", name, model.ModelName)
		}
	}

	// Delete model type
	err = s.typeRepo.Delete(name)
	if err != nil {
		return fmt.Errorf("failed to delete model type: %w", err)
	}

	return nil
}

// TestGetAllModels tests the GetAllModels method
func TestRCModelService_GetAllModels(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add some test data
	model1 := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	model2 := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Arrma",
		ModelName: "Kraton",
		Scale:     "1/8",
		ModelType: "Truck",
	}
	mockModelRepo.Create(model1)
	mockModelRepo.Create(model2)

	rcModels, err := svc.GetAllModels()
	if err != nil {
		t.Fatalf("GetAllModels() error = %v", err)
	}
	if len(rcModels) != 2 {
		t.Errorf("GetAllModels() expected 2 models, got %d", len(rcModels))
	}
}

// TestGetModelByID tests the GetModelByID method
func TestRCModelService_GetModelByID_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Team Associated",
		ModelName: "RC10",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	retrieved, err := svc.GetModelByID(modelID)
	if err != nil {
		t.Fatalf("GetModelByID() error = %v", err)
	}
	if retrieved == nil || retrieved.ModelName != "RC10" {
		t.Error("GetModelByID() failed to retrieve model")
	}
}

func TestRCModelService_GetModelByID_NotFound(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	retrieved, err := svc.GetModelByID(uuid.New().String())
	if err != nil {
		t.Fatalf("GetModelByID() error = %v", err)
	}
	if retrieved != nil {
		t.Error("GetModelByID() should return nil for non-existent model")
	}
}

// TestCreateModel tests the CreateModel method
func TestRCModelService_CreateModel_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		Brand:     "Xray",
		ModelName: "XB8",
		Scale:     "1/8",
		ModelType: "Buggy",
	}

	err := svc.CreateModel(model)
	if err != nil {
		t.Fatalf("CreateModel() error = %v", err)
	}
	if model.ID == "" {
		t.Error("CreateModel() should generate UUID")
	}

	// Verify it was added
	retrieved, _ := mockModelRepo.GetByID(model.ID)
	if retrieved == nil || retrieved.ModelName != "XB8" {
		t.Error("CreateModel() failed to store model")
	}
}

func TestRCModelService_CreateModel_EmptyBrand(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		Brand:     "",
		ModelName: "XB8",
		Scale:     "1/8",
		ModelType: "Buggy",
	}

	err := svc.CreateModel(model)
	if err == nil {
		t.Error("CreateModel() expected error for empty brand, got nil")
	}
}

func TestRCModelService_CreateModel_EmptyModelName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		Brand:     "Xray",
		ModelName: "",
		Scale:     "1/8",
		ModelType: "Buggy",
	}

	err := svc.CreateModel(model)
	if err == nil {
		t.Error("CreateModel() expected error for empty model name, got nil")
	}
}

func TestRCModelService_CreateModel_EmptyScale(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		Brand:     "Xray",
		ModelName: "XB8",
		Scale:     "",
		ModelType: "Buggy",
	}

	err := svc.CreateModel(model)
	if err == nil {
		t.Error("CreateModel() expected error for empty scale, got nil")
	}
}

func TestRCModelService_CreateModel_EmptyModelType(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		Brand:     "Xray",
		ModelName: "XB8",
		Scale:     "1/8",
		ModelType: "",
	}

	err := svc.CreateModel(model)
	if err == nil {
		t.Error("CreateModel() expected error for empty model type, got nil")
	}
}

// TestUpdateModel tests the UpdateModel method
func TestRCModelService_UpdateModel_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Yokomo",
		ModelName: "YZ-2",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	// Update
	model.ModelName = "YZ-2 SR"
	model.MotorType = "Brushless"
	err := svc.UpdateModel(model)
	if err != nil {
		t.Fatalf("UpdateModel() error = %v", err)
	}

	// Verify update
	retrieved, _ := mockModelRepo.GetByID(model.ID)
	if retrieved == nil || retrieved.ModelName != "YZ-2 SR" || retrieved.MotorType != "Brushless" {
		t.Error("UpdateModel() failed to update model")
	}
}

func TestRCModelService_UpdateModel_EmptyID(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		ID:        "",
		Brand:     "Yokomo",
		ModelName: "YZ-2",
		Scale:     "1/10",
		ModelType: "Buggy",
	}

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for empty ID, got nil")
	}
}

func TestRCModelService_UpdateModel_EmptyBrand(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "",
		ModelName: "YZ-2",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for empty brand, got nil")
	}
}

func TestRCModelService_UpdateModel_EmptyModelName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Yokomo",
		ModelName: "",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for empty model name, got nil")
	}
}

func TestRCModelService_UpdateModel_EmptyScale(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Yokomo",
		ModelName: "YZ-2",
		Scale:     "",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for empty scale, got nil")
	}
}

func TestRCModelService_UpdateModel_EmptyModelType(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Yokomo",
		ModelName: "YZ-2",
		Scale:     "1/10",
		ModelType: "",
	}
	mockModelRepo.Create(model)

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for empty model type, got nil")
	}
}

func TestRCModelService_UpdateModel_NotFound(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	model := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Yokomo",
		ModelName: "YZ-2",
		Scale:     "1/10",
		ModelType: "Buggy",
	}

	err := svc.UpdateModel(model)
	if err == nil {
		t.Error("UpdateModel() expected error for non-existent model, got nil")
	}
}

// TestDeleteModel tests the DeleteModel method
func TestRCModelService_DeleteModel_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	modelID := uuid.New().String()
	model := &models.RCModel{
		ID:        modelID,
		Brand:     "Losi",
		ModelName: "TEN-T",
		Scale:     "1/10",
		ModelType: "Truck",
	}
	mockModelRepo.Create(model)

	err := svc.DeleteModel(modelID)
	if err != nil {
		t.Fatalf("DeleteModel() error = %v", err)
	}

	// Verify deletion
	retrieved, _ := mockModelRepo.GetByID(modelID)
	if retrieved != nil {
		t.Error("DeleteModel() failed to delete model")
	}
}

func TestRCModelService_DeleteModel_EmptyID(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.DeleteModel("")
	if err == nil {
		t.Error("DeleteModel() expected error for empty ID, got nil")
	}
}

// TestGetModelCount tests the GetModelCount method
func TestRCModelService_GetModelCount(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Empty count
	count, err := svc.GetModelCount()
	if err != nil {
		t.Fatalf("GetModelCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetModelCount() expected 0, got %d", count)
	}

	// Add models
	for i := 1; i <= 5; i++ {
		model := &models.RCModel{
			ID:        uuid.New().String(),
			Brand:     "Brand",
			ModelName: "Model",
			Scale:     "1/10",
			ModelType: "Buggy",
		}
		mockModelRepo.Create(model)
	}

	count, err = svc.GetModelCount()
	if err != nil {
		t.Fatalf("GetModelCount() error = %v", err)
	}
	if count != 5 {
		t.Errorf("GetModelCount() expected 5, got %d", count)
	}
}

// TestGetAllBrands tests the GetAllBrands method
func TestRCModelService_GetAllBrands(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add some brands
	mockBrandRepo.Create("Traxxas")
	mockBrandRepo.Create("Arrma")
	mockBrandRepo.Create("Team Associated")

	brands, err := svc.GetAllBrands()
	if err != nil {
		t.Fatalf("GetAllBrands() error = %v", err)
	}
	if len(brands) != 3 {
		t.Errorf("GetAllBrands() expected 3 brands, got %d", len(brands))
	}
}

// TestGetAllModelNames tests the GetAllModelNames method
func TestRCModelService_GetAllModelNames(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add some models
	model1 := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	model2 := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Arrma",
		ModelName: "Kraton",
		Scale:     "1/8",
		ModelType: "Truck",
	}
	mockModelRepo.Create(model1)
	mockModelRepo.Create(model2)

	names, err := svc.GetAllModelNames()
	if err != nil {
		t.Fatalf("GetAllModelNames() error = %v", err)
	}
	if len(names) != 2 {
		t.Errorf("GetAllModelNames() expected 2 names, got %d", len(names))
	}
}

// TestAddBrand tests the AddBrand method
func TestRCModelService_AddBrand_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddBrand("HPI")
	if err != nil {
		t.Fatalf("AddBrand() error = %v", err)
	}

	// Verify brand was added
	brand, _ := mockBrandRepo.GetByName("HPI")
	if brand == nil {
		t.Error("AddBrand() failed to add brand")
	}
}

func TestRCModelService_AddBrand_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddBrand("")
	if err == nil {
		t.Error("AddBrand() expected error for empty name, got nil")
	}
}

func TestRCModelService_AddBrand_Duplicate(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add brand first time
	mockBrandRepo.Create("HPI")

	// Try to add again
	err := svc.AddBrand("HPI")
	if err == nil {
		t.Error("AddBrand() expected error for duplicate brand, got nil")
	}
}

// TestDeleteBrand tests the DeleteBrand method
func TestRCModelService_DeleteBrand_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add brand
	mockBrandRepo.Create("HPI")

	err := svc.DeleteBrand("HPI")
	if err != nil {
		t.Fatalf("DeleteBrand() error = %v", err)
	}

	// Verify brand was deleted
	brand, _ := mockBrandRepo.GetByName("HPI")
	if brand != nil {
		t.Error("DeleteBrand() failed to delete brand")
	}
}

func TestRCModelService_DeleteBrand_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.DeleteBrand("")
	if err == nil {
		t.Error("DeleteBrand() expected error for empty name, got nil")
	}
}

func TestRCModelService_DeleteBrand_UsedByModel(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add brand
	mockBrandRepo.Create("Traxxas")

	// Add model using this brand
	model := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	// Try to delete brand
	err := svc.DeleteBrand("Traxxas")
	if err == nil {
		t.Error("DeleteBrand() expected error when brand is used by model, got nil")
	}
}

// TestGetAllScales tests the GetAllScales method
func TestRCModelService_GetAllScales(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add some scales
	mockScaleRepo.Create("1/10")
	mockScaleRepo.Create("1/8")
	mockScaleRepo.Create("1/5")

	scales, err := svc.GetAllScales()
	if err != nil {
		t.Fatalf("GetAllScales() error = %v", err)
	}
	if len(scales) != 3 {
		t.Errorf("GetAllScales() expected 3 scales, got %d", len(scales))
	}
}

// TestAddScale tests the AddScale method
func TestRCModelService_AddScale_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddScale("1/12")
	if err != nil {
		t.Fatalf("AddScale() error = %v", err)
	}

	// Verify scale was added
	scale, _ := mockScaleRepo.GetByName("1/12")
	if scale == nil {
		t.Error("AddScale() failed to add scale")
	}
}

func TestRCModelService_AddScale_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddScale("")
	if err == nil {
		t.Error("AddScale() expected error for empty name, got nil")
	}
}

func TestRCModelService_AddScale_Duplicate(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add scale first time
	mockScaleRepo.Create("1/10")

	// Try to add again
	err := svc.AddScale("1/10")
	if err == nil {
		t.Error("AddScale() expected error for duplicate scale, got nil")
	}
}

// TestDeleteScale tests the DeleteScale method
func TestRCModelService_DeleteScale_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add scale
	mockScaleRepo.Create("1/12")

	err := svc.DeleteScale("1/12")
	if err != nil {
		t.Fatalf("DeleteScale() error = %v", err)
	}

	// Verify scale was deleted
	scale, _ := mockScaleRepo.GetByName("1/12")
	if scale != nil {
		t.Error("DeleteScale() failed to delete scale")
	}
}

func TestRCModelService_DeleteScale_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.DeleteScale("")
	if err == nil {
		t.Error("DeleteScale() expected error for empty name, got nil")
	}
}

func TestRCModelService_DeleteScale_UsedByModel(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add scale
	mockScaleRepo.Create("1/10")

	// Add model using this scale
	model := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	// Try to delete scale
	err := svc.DeleteScale("1/10")
	if err == nil {
		t.Error("DeleteScale() expected error when scale is used by model, got nil")
	}
}

// TestGetAllModelTypes tests the GetAllModelTypes method
func TestRCModelService_GetAllModelTypes(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add some types
	mockTypeRepo.Create("Buggy")
	mockTypeRepo.Create("Truck")
	mockTypeRepo.Create("Touring Car")

	types, err := svc.GetAllModelTypes()
	if err != nil {
		t.Fatalf("GetAllModelTypes() error = %v", err)
	}
	if len(types) != 3 {
		t.Errorf("GetAllModelTypes() expected 3 types, got %d", len(types))
	}
}

// TestAddModelType tests the AddModelType method
func TestRCModelService_AddModelType_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddModelType("Monster Truck")
	if err != nil {
		t.Fatalf("AddModelType() error = %v", err)
	}

	// Verify type was added
	modelType, _ := mockTypeRepo.GetByName("Monster Truck")
	if modelType == nil {
		t.Error("AddModelType() failed to add model type")
	}
}

func TestRCModelService_AddModelType_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.AddModelType("")
	if err == nil {
		t.Error("AddModelType() expected error for empty name, got nil")
	}
}

func TestRCModelService_AddModelType_Duplicate(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add type first time
	mockTypeRepo.Create("Buggy")

	// Try to add again
	err := svc.AddModelType("Buggy")
	if err == nil {
		t.Error("AddModelType() expected error for duplicate type, got nil")
	}
}

// TestDeleteModelType tests the DeleteModelType method
func TestRCModelService_DeleteModelType_Success(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add type
	mockTypeRepo.Create("Monster Truck")

	err := svc.DeleteModelType("Monster Truck")
	if err != nil {
		t.Fatalf("DeleteModelType() error = %v", err)
	}

	// Verify type was deleted
	modelType, _ := mockTypeRepo.GetByName("Monster Truck")
	if modelType != nil {
		t.Error("DeleteModelType() failed to delete model type")
	}
}

func TestRCModelService_DeleteModelType_EmptyName(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	err := svc.DeleteModelType("")
	if err == nil {
		t.Error("DeleteModelType() expected error for empty name, got nil")
	}
}

func TestRCModelService_DeleteModelType_UsedByModel(t *testing.T) {
	mockModelRepo := NewMockRCModelRepository()
	mockBrandRepo := NewMockRCModelBrandRepository()
	mockScaleRepo := NewMockRCModelScaleRepository()
	mockTypeRepo := NewMockRCModelTypeRepository()

	svc := NewRCModelServiceForTest(mockModelRepo, mockBrandRepo, mockScaleRepo, mockTypeRepo)

	// Add type
	mockTypeRepo.Create("Buggy")

	// Add model using this type
	model := &models.RCModel{
		ID:        uuid.New().String(),
		Brand:     "Traxxas",
		ModelName: "Slash",
		Scale:     "1/10",
		ModelType: "Buggy",
	}
	mockModelRepo.Create(model)

	// Try to delete type
	err := svc.DeleteModelType("Buggy")
	if err == nil {
		t.Error("DeleteModelType() expected error when type is used by model, got nil")
	}
}
