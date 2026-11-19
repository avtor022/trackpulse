package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
)

// CompetitionFilter represents a reusable competition filter component
type CompetitionFilter struct {
	mainWindow           fyne.Window
	selectedYear         string
	selectedSeason       string
	selectedTrack        string
	selectedModelType    string
	filteredYears        []string
	filteredSeasons      []string
	filteredTracks       []string
	filteredModelTypes   []string
	yearButton           *widget.Button
	seasonButton         *widget.Button
	trackButton          *widget.Button
	modelTypeButton      *widget.Button
	labels               []*widget.Label
	onFilterChanged      func()
	allCompetitions      []models.Competition
	filteredCompetitions []models.Competition
}

// NewCompetitionFilter creates a new competition filter component
func NewCompetitionFilter(mainWindow fyne.Window, onFilterChanged func()) *CompetitionFilter {
	cf := &CompetitionFilter{
		mainWindow:      mainWindow,
		onFilterChanged: onFilterChanged,
	}
	return cf
}

// SetCompetitions sets the competitions to filter and populates filter options
func (cf *CompetitionFilter) SetCompetitions(competitions []models.Competition) {
	cf.allCompetitions = competitions
	cf.populateFilterOptions()
	cf.applyFilters()
}

// CreateContent creates the filter UI container
func (cf *CompetitionFilter) CreateContent() *fyne.Container {
	cf.yearButton = widget.NewButton(locale.T("filter.all_years"), func() {
		cf.showYearFilterPopup()
	})
	cf.seasonButton = widget.NewButton(locale.T("filter.all_seasons"), func() {
		cf.showSeasonFilterPopup()
	})
	cf.trackButton = widget.NewButton(locale.T("filter.all_tracks"), func() {
		cf.showTrackFilterPopup()
	})
	cf.modelTypeButton = widget.NewButton(locale.T("filter.all_model_types"), func() {
		cf.showModelTypeFilterPopup()
	})

	cf.labels = []*widget.Label{
		widget.NewLabel(locale.T("common.year")),
		widget.NewLabel(locale.T("common.season")),
		widget.NewLabel(locale.T("common.track")),
		widget.NewLabel(locale.T("common.model_type")),
	}
	filterContainer := container.NewGridWithColumns(4,
		cf.labels[0], cf.labels[1], cf.labels[2], cf.labels[3],
		cf.yearButton,
		cf.seasonButton,
		cf.trackButton,
		cf.modelTypeButton,
	)

	return container.NewVBox(
		widget.NewSeparator(),
		filterContainer,
		widget.NewSeparator(),
	)
}

// populateFilterOptions populates filter dropdown options from competitions
func (cf *CompetitionFilter) populateFilterOptions() {
	yearSet := make(map[string]bool)
	seasonSet := make(map[string]bool)
	trackSet := make(map[string]bool)
	modelTypeSet := make(map[string]bool)

	for _, comp := range cf.allCompetitions {
		if comp.CompetitionYear != nil {
			yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
			yearSet[yearStr] = true
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

	cf.filteredYears = []string{locale.T("filter.all_years")}
	for year := range yearSet {
		cf.filteredYears = append(cf.filteredYears, year)
	}

	cf.filteredSeasons = []string{locale.T("filter.all_seasons")}
	for season := range seasonSet {
		cf.filteredSeasons = append(cf.filteredSeasons, season)
	}

	cf.filteredTracks = []string{locale.T("filter.all_tracks")}
	for track := range trackSet {
		cf.filteredTracks = append(cf.filteredTracks, track)
	}

	cf.filteredModelTypes = []string{locale.T("filter.all_model_types")}
	for modelType := range modelTypeSet {
		cf.filteredModelTypes = append(cf.filteredModelTypes, modelType)
	}
}

// applyFilters applies selected filters to the competition list
func (cf *CompetitionFilter) applyFilters() {
	cf.filteredCompetitions = []models.Competition{}

	for _, comp := range cf.allCompetitions {
		if cf.matchesFilters(comp) {
			cf.filteredCompetitions = append(cf.filteredCompetitions, comp)
		}
	}

	if cf.onFilterChanged != nil {
		cf.onFilterChanged()
	}
}

// matchesFilters checks if a competition matches the current filter selections
func (cf *CompetitionFilter) matchesFilters(comp models.Competition) bool {
	// Check year filter
	if cf.selectedYear != "" {
		if comp.CompetitionYear == nil {
			return false
		}
		yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
		if yearStr != cf.selectedYear {
			return false
		}
	}

	// Check season filter
	if cf.selectedSeason != "" {
		if comp.Season != cf.selectedSeason {
			return false
		}
	}

	// Check track filter
	if cf.selectedTrack != "" {
		if comp.TrackName != cf.selectedTrack {
			return false
		}
	}

	// Check model type filter
	if cf.selectedModelType != "" {
		if comp.ModelType != cf.selectedModelType {
			return false
		}
	}

	return true
}

// GetFilteredCompetitions returns the filtered list of competitions
func (cf *CompetitionFilter) GetFilteredCompetitions() []models.Competition {
	return cf.filteredCompetitions
}

// GetFilteredCompetitionTitles returns titles of filtered competitions
func (cf *CompetitionFilter) GetFilteredCompetitionTitles() []string {
	titles := make([]string, len(cf.filteredCompetitions))
	for i, comp := range cf.filteredCompetitions {
		titles[i] = comp.CompetitionTitle
	}
	return titles
}

// RefreshButtons updates button text to reflect current filter selections
func (cf *CompetitionFilter) RefreshButtons() {
	if cf.yearButton != nil {
		cf.yearButton.SetText(filterLabel(cf.selectedYear, "filter.all_years"))
	}
	if cf.seasonButton != nil {
		cf.seasonButton.SetText(filterLabel(cf.selectedSeason, "filter.all_seasons"))
	}
	if cf.trackButton != nil {
		cf.trackButton.SetText(filterLabel(cf.selectedTrack, "filter.all_tracks"))
	}
	if cf.modelTypeButton != nil {
		cf.modelTypeButton.SetText(filterLabel(cf.selectedModelType, "filter.all_model_types"))
	}
}

func filterLabel(value, allKey string) string {
	if value == "" {
		return locale.T(allKey)
	}
	return value
}

func filterValue(value, allKey string) string {
	if value == locale.T(allKey) {
		return ""
	}
	return value
}

// RefreshLocale updates the filter labels without changing its selection.
func (cf *CompetitionFilter) RefreshLocale() {
	keys := []string{"common.year", "common.season", "common.track", "common.model_type"}
	for i, label := range cf.labels {
		label.SetText(locale.T(keys[i]))
	}
	cf.RefreshButtons()
}

// showYearFilterPopup shows the year filter popup
func (cf *CompetitionFilter) showYearFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	cf.populateFilterOptions()

	items := make([]string, len(cf.filteredYears))
	copy(items, cf.filteredYears)

	popupManager := NewReferencePopupManager(
		cf.mainWindow,
		ReferencePopupConfig{
			Title:          "common.year",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				cf.selectedYear = filterValue(selected, "filter.all_years")
				cf.RefreshButtons()
				cf.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			cf.selectedYear = filterValue(selected, "filter.all_years")
			cf.RefreshButtons()
			cf.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showSeasonFilterPopup shows the season filter popup
func (cf *CompetitionFilter) showSeasonFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	cf.populateFilterOptions()

	items := make([]string, len(cf.filteredSeasons))
	copy(items, cf.filteredSeasons)

	popupManager := NewReferencePopupManager(
		cf.mainWindow,
		ReferencePopupConfig{
			Title:          "common.season",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				cf.selectedSeason = filterValue(selected, "filter.all_seasons")
				cf.RefreshButtons()
				cf.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			cf.selectedSeason = filterValue(selected, "filter.all_seasons")
			cf.RefreshButtons()
			cf.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showTrackFilterPopup shows the track filter popup
func (cf *CompetitionFilter) showTrackFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	cf.populateFilterOptions()

	items := make([]string, len(cf.filteredTracks))
	copy(items, cf.filteredTracks)

	popupManager := NewReferencePopupManager(
		cf.mainWindow,
		ReferencePopupConfig{
			Title:          "common.track",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				cf.selectedTrack = filterValue(selected, "filter.all_tracks")
				cf.RefreshButtons()
				cf.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			cf.selectedTrack = filterValue(selected, "filter.all_tracks")
			cf.RefreshButtons()
			cf.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showModelTypeFilterPopup shows the model type filter popup
func (cf *CompetitionFilter) showModelTypeFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	cf.populateFilterOptions()

	items := make([]string, len(cf.filteredModelTypes))
	copy(items, cf.filteredModelTypes)

	popupManager := NewReferencePopupManager(
		cf.mainWindow,
		ReferencePopupConfig{
			Title:          "common.model_type",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				cf.selectedModelType = filterValue(selected, "filter.all_model_types")
				cf.RefreshButtons()
				cf.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			cf.selectedModelType = filterValue(selected, "filter.all_model_types")
			cf.RefreshButtons()
			cf.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}
