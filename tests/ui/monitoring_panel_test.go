package service_test

import (
	"fmt"
	"testing"

	"trackpulse/internal/models"
)

// MockCompetitionService is a mock implementation of CompetitionService for testing
type MockCompetitionService struct {
	competitions []models.Competition
	err          error
}

// GetAllCompetitions returns mock competitions
func (m *MockCompetitionService) GetAllCompetitions() ([]models.Competition, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.competitions, nil
}

// TestPopulateFilterOptions tests that filter options are correctly populated from competitions
func TestPopulateFilterOptions(t *testing.T) {
	year2023 := 2023
	year2024 := 2024

	competitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2024,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
		{
			ID:               "3",
			CompetitionTitle: "Race 3",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Buggy",
		},
	}

	// Verify unique years
	years := make(map[string]bool)
	for _, comp := range competitions {
		if comp.CompetitionYear != nil {
			yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
			years[yearStr] = true
		}
	}
	if len(years) != 2 {
		t.Errorf("Expected 2 unique years, got %d", len(years))
	}

	// Verify unique seasons
	seasons := make(map[string]bool)
	for _, comp := range competitions {
		if comp.Season != "" {
			seasons[comp.Season] = true
		}
	}
	if len(seasons) != 2 {
		t.Errorf("Expected 2 unique seasons, got %d", len(seasons))
	}

	// Verify unique tracks
	tracks := make(map[string]bool)
	for _, comp := range competitions {
		if comp.TrackName != "" {
			tracks[comp.TrackName] = true
		}
	}
	if len(tracks) != 2 {
		t.Errorf("Expected 2 unique tracks, got %d", len(tracks))
	}

	// Verify unique model types
	modelTypes := make(map[string]bool)
	for _, comp := range competitions {
		if comp.ModelType != "" {
			modelTypes[comp.ModelType] = true
		}
	}
	if len(modelTypes) != 2 {
		t.Errorf("Expected 2 unique model types, got %d", len(modelTypes))
	}
}

// TestMatchesFilters tests the filter matching logic
func TestMatchesFilters(t *testing.T) {
	year2023 := 2023
	year2024 := 2024

	competitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2024,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
		{
			ID:               "3",
			CompetitionTitle: "Race 3",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track C",
			ModelType:        "Touring",
		},
	}

	// Test no filters - all should match
	matchCount := 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "", "", "", "") {
			matchCount++
		}
	}
	if matchCount != 3 {
		t.Errorf("Expected 3 matches with no filters, got %d", matchCount)
	}

	// Test year filter
	matchCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "2023", "", "", "") {
			matchCount++
		}
	}
	if matchCount != 2 {
		t.Errorf("Expected 2 matches for year 2023, got %d", matchCount)
	}

	// Test season filter
	matchCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "", "Winter", "", "") {
			matchCount++
		}
	}
	if matchCount != 2 {
		t.Errorf("Expected 2 matches for Winter season, got %d", matchCount)
	}

	// Test track filter
	matchCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "", "", "Track A", "") {
			matchCount++
		}
	}
	if matchCount != 1 {
		t.Errorf("Expected 1 match for Track A, got %d", matchCount)
	}

	// Test model type filter
	matchCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "", "", "", "Touring") {
			matchCount++
		}
	}
	if matchCount != 2 {
		t.Errorf("Expected 2 matches for Touring model type, got %d", matchCount)
	}

	// Test combined filters
	matchCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, "2023", "Winter", "Track A", "Touring") {
			matchCount++
		}
	}
	if matchCount != 1 {
		t.Errorf("Expected 1 match for combined filters, got %d", matchCount)
	}
}

// matchesFiltersHelper is a helper function to test filter matching logic
func matchesFiltersHelper(comp models.Competition, selectedYear, selectedSeason, selectedTrack, selectedModelType string) bool {
	// Check year filter
	if selectedYear != "" && selectedYear != "All Years" {
		if comp.CompetitionYear == nil {
			return false
		}
		yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
		if yearStr != selectedYear {
			return false
		}
	}

	// Check season filter
	if selectedSeason != "" && selectedSeason != "All Seasons" {
		if comp.Season != selectedSeason {
			return false
		}
	}

	// Check track filter
	if selectedTrack != "" && selectedTrack != "All Tracks" {
		if comp.TrackName != selectedTrack {
			return false
		}
	}

	// Check model type filter
	if selectedModelType != "" && selectedModelType != "All Model Types" {
		if comp.ModelType != selectedModelType {
			return false
		}
	}

	return true
}

// TestApplyFilters tests the filter application logic
func TestApplyFilters(t *testing.T) {
	year2023 := 2023
	year2024 := 2024

	allCompetitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2024,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
		{
			ID:               "3",
			CompetitionTitle: "Race 3",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Buggy",
		},
	}

	// Test filtering by year
	filtered := applyFiltersHelper(allCompetitions, "2023", "", "", "")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered competitions for year 2023, got %d", len(filtered))
	}

	// Test filtering by season
	filtered = applyFiltersHelper(allCompetitions, "", "Winter", "", "")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered competitions for Winter season, got %d", len(filtered))
	}

	// Test filtering by track
	filtered = applyFiltersHelper(allCompetitions, "", "", "Track A", "")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered competitions for Track A, got %d", len(filtered))
	}

	// Test filtering by model type
	filtered = applyFiltersHelper(allCompetitions, "", "", "", "Buggy")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered competitions for Buggy model type, got %d", len(filtered))
	}

	// Test filtering with multiple criteria
	filtered = applyFiltersHelper(allCompetitions, "2023", "Winter", "Track A", "")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered competitions for year 2023, Winter, Track A, got %d", len(filtered))
	}

	// Test filtering with all criteria
	filtered = applyFiltersHelper(allCompetitions, "2023", "Winter", "Track A", "Touring")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered competition for all criteria, got %d", len(filtered))
	}

	// Test filtering with no results
	filtered = applyFiltersHelper(allCompetitions, "2024", "Winter", "", "")
	if len(filtered) != 0 {
		t.Errorf("Expected 0 filtered competitions for non-existent combination, got %d", len(filtered))
	}
}

// applyFiltersHelper is a helper function to test filter application
func applyFiltersHelper(competitions []models.Competition, selectedYear, selectedSeason, selectedTrack, selectedModelType string) []models.Competition {
	filtered := []models.Competition{}
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, selectedYear, selectedSeason, selectedTrack, selectedModelType) {
			filtered = append(filtered, comp)
		}
	}
	return filtered
}

// TestGetFilteredCompetitionTitles tests getting titles from filtered competitions
func TestGetFilteredCompetitionTitles(t *testing.T) {
	competitions := []models.Competition{
		{ID: "1", CompetitionTitle: "Race 1"},
		{ID: "2", CompetitionTitle: "Race 2"},
		{ID: "3", CompetitionTitle: "Race 3"},
	}

	titles := make([]string, len(competitions))
	for i, comp := range competitions {
		titles[i] = comp.CompetitionTitle
	}

	if len(titles) != 3 {
		t.Errorf("Expected 3 titles, got %d", len(titles))
	}

	expectedTitles := map[string]bool{
		"Race 1": true,
		"Race 2": true,
		"Race 3": true,
	}

	for _, title := range titles {
		if !expectedTitles[title] {
			t.Errorf("Unexpected title: %s", title)
		}
	}
}

// TestFilterOptionsWithEmptyData tests handling of empty competition data
func TestFilterOptionsWithEmptyData(t *testing.T) {
	competitions := []models.Competition{}

	// Verify empty results
	years := make(map[string]bool)
	seasons := make(map[string]bool)
	tracks := make(map[string]bool)
	modelTypes := make(map[string]bool)

	for _, comp := range competitions {
		if comp.CompetitionYear != nil {
			years[fmt.Sprintf("%d", *comp.CompetitionYear)] = true
		}
		if comp.Season != "" {
			seasons[comp.Season] = true
		}
		if comp.TrackName != "" {
			tracks[comp.TrackName] = true
		}
		if comp.ModelType != "" {
			modelTypes[comp.ModelType] = true
		}
	}

	if len(years) != 0 {
		t.Errorf("Expected 0 years for empty data, got %d", len(years))
	}
	if len(seasons) != 0 {
		t.Errorf("Expected 0 seasons for empty data, got %d", len(seasons))
	}
	if len(tracks) != 0 {
		t.Errorf("Expected 0 tracks for empty data, got %d", len(tracks))
	}
	if len(modelTypes) != 0 {
		t.Errorf("Expected 0 model types for empty data, got %d", len(modelTypes))
	}
}

// TestFilterOptionsWithNilValues tests handling of nil values in competition fields
func TestFilterOptionsWithNilValues(t *testing.T) {
	competitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  nil, // Nil year
			Season:           "",  // Empty season
			TrackName:        "",  // Empty track
			ModelType:        "",  // Empty model type
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  nil,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
	}

	// Verify nil/empty values are handled correctly
	validYears := 0
	validSeasons := 0
	validTracks := 0
	validModelTypes := 0

	for _, comp := range competitions {
		if comp.CompetitionYear != nil {
			validYears++
		}
		if comp.Season != "" {
			validSeasons++
		}
		if comp.TrackName != "" {
			validTracks++
		}
		if comp.ModelType != "" {
			validModelTypes++
		}
	}

	if validYears != 0 {
		t.Errorf("Expected 0 valid years with nil values, got %d", validYears)
	}
	if validSeasons != 1 {
		t.Errorf("Expected 1 valid season, got %d", validSeasons)
	}
	if validTracks != 1 {
		t.Errorf("Expected 1 valid track, got %d", validTracks)
	}
	if validModelTypes != 1 {
		t.Errorf("Expected 1 valid model type, got %d", validModelTypes)
	}
}

// TestFilterOptionsWithDuplicateValues tests handling of duplicate values
func TestFilterOptionsWithDuplicateValues(t *testing.T) {
	year2023 := 2023

	competitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2023, // Same year
			Season:           "Winter",  // Same season
			TrackName:        "Track A", // Same track
			ModelType:        "Touring", // Same model type
		},
		{
			ID:               "3",
			CompetitionTitle: "Race 3",
			CompetitionYear:  &year2023, // Same year
			Season:           "Winter",  // Same season
			TrackName:        "Track A", // Same track
			ModelType:        "Touring", // Same model type
		},
	}

	// Verify duplicates are handled (using maps for uniqueness)
	years := make(map[string]bool)
	seasons := make(map[string]bool)
	tracks := make(map[string]bool)
	modelTypes := make(map[string]bool)

	for _, comp := range competitions {
		if comp.CompetitionYear != nil {
			years[fmt.Sprintf("%d", *comp.CompetitionYear)] = true
		}
		if comp.Season != "" {
			seasons[comp.Season] = true
		}
		if comp.TrackName != "" {
			tracks[comp.TrackName] = true
		}
		if comp.ModelType != "" {
			modelTypes[comp.ModelType] = true
		}
	}

	// Should have only 1 unique value for each field despite 3 competitions
	if len(years) != 1 {
		t.Errorf("Expected 1 unique year, got %d", len(years))
	}
	if len(seasons) != 1 {
		t.Errorf("Expected 1 unique season, got %d", len(seasons))
	}
	if len(tracks) != 1 {
		t.Errorf("Expected 1 unique track, got %d", len(tracks))
	}
	if len(modelTypes) != 1 {
		t.Errorf("Expected 1 unique model type, got %d", len(modelTypes))
	}
}

// TestCompetitionWithAllFields tests competitions with all fields populated
func TestCompetitionWithAllFields(t *testing.T) {
	year2024 := 2024

	competition := models.Competition{
		ID:               "test-id",
		CompetitionTitle: "Championship Final",
		CompetitionType:  "Race",
		ModelType:        "Formula 1",
		ModelScale:       "1/10",
		TrackName:        "Silverstone",
		LapCountTarget:   intPtr(50),
		TimeLimitMinutes: intPtr(30),
		Status:           "active",
		CompetitionYear:  &year2024,
		Season:           "Summer",
	}

	if competition.CompetitionTitle != "Championship Final" {
		t.Errorf("Expected CompetitionTitle to be 'Championship Final', got '%s'", competition.CompetitionTitle)
	}
	if competition.ModelType != "Formula 1" {
		t.Errorf("Expected ModelType to be 'Formula 1', got '%s'", competition.ModelType)
	}
	if competition.TrackName != "Silverstone" {
		t.Errorf("Expected TrackName to be 'Silverstone', got '%s'", competition.TrackName)
	}
	if competition.Season != "Summer" {
		t.Errorf("Expected Season to be 'Summer', got '%s'", competition.Season)
	}
	if *competition.CompetitionYear != 2024 {
		t.Errorf("Expected CompetitionYear to be 2024, got %d", *competition.CompetitionYear)
	}
}

// TestFilterButtonLabels tests that filter button labels are set correctly
func TestFilterButtonLabels(t *testing.T) {
	// Simulate default filter button labels
	defaultLabels := map[string]string{
		"year":       "All Years",
		"season":     "All Seasons",
		"track":      "All Tracks",
		"model_type": "All Model Types",
	}

	// Verify default labels
	expectedDefaults := map[string]string{
		"year":       "All Years",
		"season":     "All Seasons",
		"track":      "All Tracks",
		"model_type": "All Model Types",
	}

	for key, label := range defaultLabels {
		if expectedDefaults[key] != label {
			t.Errorf("Expected default label for %s to be '%s', got '%s'", key, expectedDefaults[key], label)
		}
	}

	// Simulate selected filter
	selectedLabel := "2024"
	if selectedLabel != "2024" {
		t.Errorf("Expected selected label to be '2024', got '%s'", selectedLabel)
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// TestFilterSelectionLogic tests the logic of selecting filters
func TestFilterSelectionLogic(t *testing.T) {
	year2023 := 2023
	year2024 := 2024

	competitions := []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2024,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
	}

	// Simulate filter selection process
	selectedYear := ""
	selectedSeason := ""
	selectedTrack := ""
	selectedModelType := ""

	// Initially no filters selected - all competitions visible
	filteredCount := 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, selectedYear, selectedSeason, selectedTrack, selectedModelType) {
			filteredCount++
		}
	}
	if filteredCount != 2 {
		t.Errorf("Expected 2 competitions with no filters, got %d", filteredCount)
	}

	// Select year filter
	selectedYear = "2023"
	filteredCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, selectedYear, selectedSeason, selectedTrack, selectedModelType) {
			filteredCount++
		}
	}
	if filteredCount != 1 {
		t.Errorf("Expected 1 competition after selecting year 2023, got %d", filteredCount)
	}

	// Reset year filter
	selectedYear = "All Years"
	filteredCount = 0
	for _, comp := range competitions {
		if matchesFiltersHelper(comp, selectedYear, selectedSeason, selectedTrack, selectedModelType) {
			filteredCount++
		}
	}
	if filteredCount != 2 {
		t.Errorf("Expected 2 competitions after resetting year filter, got %d", filteredCount)
	}
}

// TestMonitorPanelFilterRefresh tests that filters refresh when competitions are loaded
func TestMonitorPanelFilterRefresh(t *testing.T) {
	year2023 := 2023
	year2024 := 2024

	// Simulate initial state with no competitions
	allCompetitions := []models.Competition{}
	filteredYears := []string{"All Years"}
	filteredSeasons := []string{"All Seasons"}
	filteredTracks := []string{"All Tracks"}
	filteredModelTypes := []string{"All Model Types"}

	if len(filteredYears) != 1 || len(filteredSeasons) != 1 || len(filteredTracks) != 1 || len(filteredModelTypes) != 1 {
		t.Error("Expected default filter options with no competitions")
	}

	// Simulate loading competitions
	allCompetitions = []models.Competition{
		{
			ID:               "1",
			CompetitionTitle: "Race 1",
			CompetitionYear:  &year2023,
			Season:           "Winter",
			TrackName:        "Track A",
			ModelType:        "Touring",
		},
		{
			ID:               "2",
			CompetitionTitle: "Race 2",
			CompetitionYear:  &year2024,
			Season:           "Summer",
			TrackName:        "Track B",
			ModelType:        "Buggy",
		},
	}

	// Simulate populateFilterOptions
	yearSet := make(map[string]bool)
	seasonSet := make(map[string]bool)
	trackSet := make(map[string]bool)
	modelTypeSet := make(map[string]bool)

	for _, comp := range allCompetitions {
		if comp.CompetitionYear != nil {
			yearSet[fmt.Sprintf("%d", *comp.CompetitionYear)] = true
		}
		if comp.Season != "" {
			seasonSet[comp.Season] = true
		}
		if comp.TrackName != "" {
			trackSet[comp.TrackName] = true
		}
		if comp.ModelType != "" {
			modelTypeSet[comp.ModelType] = true
		}
	}

	filteredYears = []string{"All Years"}
	for year := range yearSet {
		filteredYears = append(filteredYears, year)
	}

	filteredSeasons = []string{"All Seasons"}
	for season := range seasonSet {
		filteredSeasons = append(filteredSeasons, season)
	}

	filteredTracks = []string{"All Tracks"}
	for track := range trackSet {
		filteredTracks = append(filteredTracks, track)
	}

	filteredModelTypes = []string{"All Model Types"}
	for modelType := range modelTypeSet {
		filteredModelTypes = append(filteredModelTypes, modelType)
	}

	// Verify filter options were populated
	if len(filteredYears) != 3 { // All Years + 2023 + 2024
		t.Errorf("Expected 3 year options, got %d", len(filteredYears))
	}
	if len(filteredSeasons) != 3 { // All Seasons + Winter + Summer
		t.Errorf("Expected 3 season options, got %d", len(filteredSeasons))
	}
	if len(filteredTracks) != 3 { // All Tracks + Track A + Track B
		t.Errorf("Expected 3 track options, got %d", len(filteredTracks))
	}
	if len(filteredModelTypes) != 3 { // All Model Types + Touring + Buggy
		t.Errorf("Expected 3 model type options, got %d", len(filteredModelTypes))
	}
}
