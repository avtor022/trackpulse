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
	content             *fyne.Container
	mainWindow          fyne.Window
	competitionService  *service.CompetitionService
	selectedCompetition string
	statusLabel         *widget.Label
	allCompetitions     []models.Competition
	competitionButton   *widget.Button
	yearFilterSelect    *widget.Select
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

	// Year filter dropdown
	years := p.getAvailableYears()
	yearOptions := append([]string{locale.T("common.all_years")}, years...)
	p.yearFilterSelect = widget.NewSelect(yearOptions, func(selected string) {
		p.onYearFilterChanged(selected)
	})
	if len(years) == 0 {
		p.yearFilterSelect.SetSelected(locale.T("common.all_years"))
	} else {
		p.yearFilterSelect.SetSelected(locale.T("common.all_years"))
	}

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(locale.T("common.filter_year"), p.yearFilterSelect),
		),
		p.competitionButton,
	)

	// Main content area (placeholder for future monitoring widgets)
	monitoringContent := widget.NewLabel(locale.T("monitoring.placeholder"))
	monitoringContent.Alignment = fyne.TextAlignCenter

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			selectorContainer,
			widget.NewSeparator(),
		),
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
}

// showCompetitionPopup shows the competition selection popup without add/delete buttons
func (p *MonitoringPanel) showCompetitionPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Apply year filter to get filtered competitions
	selectedYear := p.yearFilterSelect.Selected
	filteredCompetitions := p.filterCompetitionsByYear(selectedYear)

	// Convert competitions to ReferenceItem slice
	items := make([]ReferenceItem, len(filteredCompetitions))
	for i, comp := range filteredCompetitions {
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
				result := make([]ReferenceItem, len(filteredCompetitions))
				for i, comp := range filteredCompetitions {
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
		p.getFilteredCompetitionTitles(filteredCompetitions),
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

// getFilteredCompetitionTitles returns a slice of competition titles from filtered list
func (p *MonitoringPanel) getFilteredCompetitionTitles(competitions []models.Competition) []string {
	titles := make([]string, len(competitions))
	for i, comp := range competitions {
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

// getAvailableYears extracts unique years from competitions
func (p *MonitoringPanel) getAvailableYears() []string {
	yearMap := make(map[string]bool)
	for _, comp := range p.allCompetitions {
		if comp.TimeStart != nil {
			year := fmt.Sprintf("%d", comp.TimeStart.Year())
			yearMap[year] = true
		}
	}

	years := make([]string, 0, len(yearMap))
	for year := range yearMap {
		years = append(years, year)
	}

	// Sort years in descending order
	for i := 0; i < len(years); i++ {
		for j := i + 1; j < len(years); j++ {
			if years[i] < years[j] {
				years[i], years[j] = years[j], years[i]
			}
		}
	}

	return years
}

// filterCompetitionsByYear filters competitions by selected year
func (p *MonitoringPanel) filterCompetitionsByYear(year string) []models.Competition {
	if year == locale.T("common.all_years") || year == "" {
		return p.allCompetitions
	}

	var filtered []models.Competition
	for _, comp := range p.allCompetitions {
		if comp.TimeStart != nil {
			compYear := fmt.Sprintf("%d", comp.TimeStart.Year())
			if compYear == year {
				filtered = append(filtered, comp)
			}
		}
	}
	return filtered
}

// onYearFilterChanged handles year filter selection change
func (p *MonitoringPanel) onYearFilterChanged(selected string) {
	// Refresh the popup with filtered competitions if it's open
	// The actual filtering will be applied when showCompetitionPopup is called
	p.refreshCompetitions()
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
