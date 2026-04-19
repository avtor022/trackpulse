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

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewLabel(locale.T("form.competition.title")),
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

	// Convert competitions to ReferenceItem slice
	items := make([]ReferenceItem, len(p.allCompetitions))
	for i, comp := range p.allCompetitions {
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
				result := make([]ReferenceItem, len(p.allCompetitions))
				for i, comp := range p.allCompetitions {
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
		p.getCompetitionTitles(),
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
