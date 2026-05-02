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
	statusLabel          *widget.Label
	allCompetitions      []models.Competition
	competitionButton    *widget.Button
	filteredCompetitions []models.Competition
	yearSelect           *widget.Select
	seasonSelect         *widget.Select
	trackSelect          *widget.Select
	modelTypeSelect      *widget.Select
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

	// Filter selects
	p.yearSelect = widget.NewSelect([]string{locale.T("filter.all_years")}, func(value string) {
		p.applyFilters()
	})
	p.seasonSelect = widget.NewSelect([]string{locale.T("filter.all_seasons")}, func(value string) {
		p.applyFilters()
	})
	p.trackSelect = widget.NewSelect([]string{locale.T("filter.all_tracks")}, func(value string) {
		p.applyFilters()
	})
	p.modelTypeSelect = widget.NewSelect([]string{locale.T("filter.all_model_types")}, func(value string) {
		p.applyFilters()
	})

	// Populate filter options
	p.populateFilterOptions()

	// Filter container
	filterContainer := container.NewGridWithColumns(4,
		widget.NewLabel(locale.T("common.year")),
		widget.NewLabel(locale.T("common.season")),
		widget.NewLabel(locale.T("common.track")),
		widget.NewLabel(locale.T("common.model_type")),
		p.yearSelect,
		p.seasonSelect,
		p.trackSelect,
		p.modelTypeSelect,
	)

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewSeparator(),
		filterContainer,
		widget.NewSeparator(),
		container.NewHBox(p.competitionButton),
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

	yearOptions := []string{locale.T("filter.all_years")}
	for year := range yearSet {
		yearOptions = append(yearOptions, year)
	}
	p.yearSelect.Options = yearOptions

	seasonOptions := []string{locale.T("filter.all_seasons")}
	for season := range seasonSet {
		seasonOptions = append(seasonOptions, season)
	}
	p.seasonSelect.Options = seasonOptions

	trackOptions := []string{locale.T("filter.all_tracks")}
	for track := range trackSet {
		trackOptions = append(trackOptions, track)
	}
	p.trackSelect.Options = trackOptions

	modelTypeOptions := []string{locale.T("filter.all_model_types")}
	for modelType := range modelTypeSet {
		modelTypeOptions = append(modelTypeOptions, modelType)
	}
	p.modelTypeSelect.Options = modelTypeOptions
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
	if p.yearSelect.Selected != "" && p.yearSelect.Selected != locale.T("filter.all_years") {
		if comp.CompetitionYear == nil {
			return false
		}
		yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
		if yearStr != p.yearSelect.Selected {
			return false
		}
	}

	// Check season filter
	if p.seasonSelect.Selected != "" && p.seasonSelect.Selected != locale.T("filter.all_seasons") {
		if comp.Season != p.seasonSelect.Selected {
			return false
		}
	}

	// Check track filter
	if p.trackSelect.Selected != "" && p.trackSelect.Selected != locale.T("filter.all_tracks") {
		if comp.TrackName != p.trackSelect.Selected {
			return false
		}
	}

	// Check model type filter
	if p.modelTypeSelect.Selected != "" && p.modelTypeSelect.Selected != locale.T("filter.all_model_types") {
		if comp.ModelType != p.modelTypeSelect.Selected {
			return false
		}
	}

	return true
}

// showCompetitionPopup shows the competition selection popup without add/delete buttons
func (p *MonitoringPanel) showCompetitionPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Convert filtered competitions to ReferenceItem slice
	items := make([]ReferenceItem, len(p.filteredCompetitions))
	for i, comp := range p.filteredCompetitions {
		items[i] = ReferenceItem{Name: comp.CompetitionTitle}
	}

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
		return
	}

	// Find the selected competition
	for _, comp := range p.allCompetitions {
		if comp.CompetitionTitle == selected {
			p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), comp.CompetitionTitle, comp.Status))
			if p.competitionButton != nil {
				p.competitionButton.SetText(comp.CompetitionTitle)
			}
			return
		}
	}
}

// Refresh updates the panel with new locale strings
func (p *MonitoringPanel) Refresh() {
	p.content = p.createContent()
}
