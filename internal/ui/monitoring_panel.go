package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// MonitoringPanel represents the monitoring panel UI
type MonitoringPanel struct {
	content              *fyne.Container
	mainWindow           fyne.Window
	competitionService   *service.CompetitionService
	selectedCompetition  string
	selectedCompetitionID string
	statusLabel          *widget.Label
	allCompetitions      []models.Competition
	competitionButton    *widget.Button
	startButton          *widget.Button
	filteredCompetitions []models.Competition
	filteredYears        []string
	filteredSeasons      []string
	filteredTracks       []string
	filteredModelTypes   []string
	yearButton           *widget.Button
	seasonButton         *widget.Button
	trackButton          *widget.Button
	modelTypeButton      *widget.Button
	selectedYear         string
	selectedSeason       string
	selectedTrack        string
	selectedModelType    string
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:         mainWindow,
		competitionService: competitionService,
	}

	p.content = p.createContent()
	return p
}

// createContent builds the monitoring panel content
func (p *MonitoringPanel) createContent() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Button to open competition selection popup using reference_popup.go without add/delete buttons
	p.competitionButton = widget.NewButton(locale.T("form.competition.select"), func() {
		p.showCompetitionPopup()
	})

	// Start button - disabled until competition is selected
	p.startButton = widget.NewButton(locale.T("button.start"), func() {
		p.startMonitoring()
	})
	p.startButton.Disable()

	// Filter buttons using reference_popup.go without add/delete functionality
	p.yearButton = widget.NewButton(locale.T("filter.all_years"), func() {
		p.showYearFilterPopup()
	})
	p.seasonButton = widget.NewButton(locale.T("filter.all_seasons"), func() {
		p.showSeasonFilterPopup()
	})
	p.trackButton = widget.NewButton(locale.T("filter.all_tracks"), func() {
		p.showTrackFilterPopup()
	})
	p.modelTypeButton = widget.NewButton(locale.T("filter.all_model_types"), func() {
		p.showModelTypeFilterPopup()
	})

	// Filter container
	filterContainer := container.NewGridWithColumns(4,
		widget.NewLabel(locale.T("common.year")),
		widget.NewLabel(locale.T("common.season")),
		widget.NewLabel(locale.T("common.track")),
		widget.NewLabel(locale.T("common.model_type")),
		p.yearButton,
		p.seasonButton,
		p.trackButton,
		p.modelTypeButton,
	)

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewSeparator(),
		filterContainer,
		widget.NewSeparator(),
		container.NewHBox(p.competitionButton, p.startButton),
		widget.NewSeparator(),
	)

	// Main content area (placeholder for future monitoring widgets)
	monitoringContent := widget.NewLabel(locale.T("monitoring.placeholder"))
	monitoringContent.Alignment = fyne.TextAlignCenter

	// Layout
	content := container.NewBorder(
		selectorContainer,
		nil,
		nil,
		nil,
		monitoringContent,
	)

	// Load competitions
	p.refreshCompetitions()

	return content
}

// refreshCompetitions loads all competitions from the service
func (p *MonitoringPanel) refreshCompetitions() {
	if p.competitionService == nil {
		return
	}

	var err error
	p.allCompetitions, err = p.competitionService.GetAllCompetitions()
	if err != nil {
		fmt.Println("ERROR loading competitions:", err)
		p.statusLabel.SetText(locale.T("status.error_loading"))
		return
	}

	if len(p.allCompetitions) == 0 {
		p.statusLabel.SetText(locale.T("status.no_competitions"))
	} else {
		p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_competitions"), len(p.allCompetitions)))
	}

	// Populate filter options from loaded competitions
	p.populateFilterOptions()

	// Apply filters and update filtered list
	p.applyFilters()
}

// populateFilterOptions populates filter dropdown options from competitions
func (p *MonitoringPanel) populateFilterOptions() {
	yearSet := make(map[string]bool)
	seasonSet := make(map[string]bool)
	trackSet := make(map[string]bool)
	modelTypeSet := make(map[string]bool)

	for _, comp := range p.allCompetitions {
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

	p.filteredYears = []string{locale.T("filter.all_years")}
	for year := range yearSet {
		p.filteredYears = append(p.filteredYears, year)
	}

	p.filteredSeasons = []string{locale.T("filter.all_seasons")}
	for season := range seasonSet {
		p.filteredSeasons = append(p.filteredSeasons, season)
	}

	p.filteredTracks = []string{locale.T("filter.all_tracks")}
	for track := range trackSet {
		p.filteredTracks = append(p.filteredTracks, track)
	}

	p.filteredModelTypes = []string{locale.T("filter.all_model_types")}
	for modelType := range modelTypeSet {
		p.filteredModelTypes = append(p.filteredModelTypes, modelType)
	}
}

// applyFilters applies selected filters to the competition list
func (p *MonitoringPanel) applyFilters() {
	p.filteredCompetitions = []models.Competition{}

	for _, comp := range p.allCompetitions {
		if p.matchesFilters(comp) {
			p.filteredCompetitions = append(p.filteredCompetitions, comp)
		}
	}
}

// matchesFilters checks if a competition matches the current filter selections
func (p *MonitoringPanel) matchesFilters(comp models.Competition) bool {
	// Check year filter
	if p.selectedYear != "" && p.selectedYear != locale.T("filter.all_years") {
		if comp.CompetitionYear == nil {
			return false
		}
		yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
		if yearStr != p.selectedYear {
			return false
		}
	}

	// Check season filter
	if p.selectedSeason != "" && p.selectedSeason != locale.T("filter.all_seasons") {
		if comp.Season != p.selectedSeason {
			return false
		}
	}

	// Check track filter
	if p.selectedTrack != "" && p.selectedTrack != locale.T("filter.all_tracks") {
		if comp.TrackName != p.selectedTrack {
			return false
		}
	}

	// Check model type filter
	if p.selectedModelType != "" && p.selectedModelType != locale.T("filter.all_model_types") {
		if comp.ModelType != p.selectedModelType {
			return false
		}
	}

	return true
}

// showCompetitionPopup shows the competition selection popup without add/delete buttons
func (p *MonitoringPanel) showCompetitionPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	popupManager := NewReferencePopupManager(
		p.mainWindow,
		ReferencePopupConfig{
			Title:          "form.competition.title",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(p.filteredCompetitions))
				for i, comp := range p.filteredCompetitions {
					result[i] = ReferenceItem{Name: comp.CompetitionTitle}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				p.selectedCompetition = selected
				p.onCompetitionSelected(selected)
			},
			UpdateOptions: func(opts []string) {},
		},
		p.getFilteredCompetitionTitles(),
		"",
		func(selected string) {
			p.selectedCompetition = selected
			p.onCompetitionSelected(selected)
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showYearFilterPopup shows the year filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showYearFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredYears))
	copy(items, p.filteredYears)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
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
				p.selectedYear = selected
				p.yearButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedYear = selected
			p.yearButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showSeasonFilterPopup shows the season filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showSeasonFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredSeasons))
	copy(items, p.filteredSeasons)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
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
				p.selectedSeason = selected
				p.seasonButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedSeason = selected
			p.seasonButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showTrackFilterPopup shows the track filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showTrackFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredTracks))
	copy(items, p.filteredTracks)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
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
				p.selectedTrack = selected
				p.trackButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedTrack = selected
			p.trackButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showModelTypeFilterPopup shows the model type filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showModelTypeFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredModelTypes))
	copy(items, p.filteredModelTypes)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
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
				p.selectedModelType = selected
				p.modelTypeButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedModelType = selected
			p.modelTypeButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// getFilteredCompetitionTitles returns a slice of filtered competition titles
func (p *MonitoringPanel) getFilteredCompetitionTitles() []string {
	titles := make([]string, len(p.filteredCompetitions))
	for i, comp := range p.filteredCompetitions {
		titles[i] = comp.CompetitionTitle
	}
	return titles
}

// getCompetitionTitles returns a slice of competition titles
func (p *MonitoringPanel) getCompetitionTitles() []string {
	titles := make([]string, len(p.allCompetitions))
	for i, comp := range p.allCompetitions {
		titles[i] = comp.CompetitionTitle
	}
	return titles
}

// onCompetitionSelected handles competition selection
func (p *MonitoringPanel) onCompetitionSelected(selected string) {
	if selected == "" {
		p.statusLabel.SetText(locale.T("status.ready"))
		if p.competitionButton != nil {
			p.competitionButton.SetText(locale.T("form.competition.select"))
		}
		if p.startButton != nil {
			p.startButton.Disable()
		}
		p.selectedCompetitionID = ""
		return
	}

	// Find the selected competition
	for _, comp := range p.allCompetitions {
		if comp.CompetitionTitle == selected {
			p.selectedCompetitionID = comp.ID
			p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), comp.CompetitionTitle, comp.Status))
			if p.competitionButton != nil {
				p.competitionButton.SetText(comp.CompetitionTitle)
			}
			if p.startButton != nil {
				p.startButton.Enable()
			}
			return
		}
	}
}

// Refresh updates the panel with new locale strings
func (p *MonitoringPanel) Refresh() {
	p.content = p.createContent()
}

// UpdateData reloads competition data and refreshes filter options
func (p *MonitoringPanel) UpdateData() {
	p.refreshCompetitions()
}

// startMonitoring starts monitoring for the selected competition
func (p *MonitoringPanel) startMonitoring() {
	if p.selectedCompetitionID == "" {
		dialog.ShowError(fmt.Errorf(locale.T("error.no_competition_selected")), p.mainWindow)
		return
	}

	err := p.competitionService.StartCompetition(p.selectedCompetitionID)
	if err != nil {
		dialog.ShowError(err, p.mainWindow)
		return
	}

	// Update status label to show the new status
	p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), p.selectedCompetition, locale.T("status.in_progress")))
	
	// Refresh competitions list to get updated status
	p.refreshCompetitions()
	
	// Show success message
	dialog.ShowInformation(locale.T("dialog.success"), locale.T("dialog.competition_started"), p.mainWindow)
}
