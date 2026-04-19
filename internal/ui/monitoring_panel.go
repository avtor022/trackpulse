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
	content            *fyne.Container
	mainWindow         fyne.Window
	competitionService *service.CompetitionService
	competitionSelect  *widget.Select
	statusLabel        *widget.Label
	allCompetitions    []models.Competition
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

	// Competition selector
	p.competitionSelect = widget.NewSelect([]string{}, func(selected string) {
		p.onCompetitionSelected(selected)
	})
	p.competitionSelect.PlaceHolder = locale.T("form.competition.type_placeholder")

	// Button to open competition selection popup
	selectBtn := widget.NewButton(locale.T("common.select_one"), func() {
		p.showCompetitionPopup()
	})

	// Competition info display
	infoLabel := widget.NewLabel("")
	infoLabel.Wrapping = fyne.TextWrapWord

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewLabel(locale.T("form.competition.title")),
		container.NewHBox(selectBtn, p.competitionSelect),
		infoLabel,
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

	// Update select options
	options := make([]string, len(p.allCompetitions))
	for i, comp := range p.allCompetitions {
		options[i] = comp.CompetitionTitle
	}
	p.competitionSelect.Options = options
	p.competitionSelect.Refresh()

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
			Title:          "common.select_one",
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
				p.competitionSelect.SetSelected(selected)
				p.onCompetitionSelected(selected)
			},
			UpdateOptions: func(opts []string) {
				p.competitionSelect.Options = opts
			},
		},
		p.getCompetitionTitles(),
		"",
		func(selected string) {
			p.competitionSelect.SetSelected(selected)
			p.onCompetitionSelected(selected)
		},
		func(opts []string) {
			p.competitionSelect.Options = opts
		},
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
		return
	}

	// Find the selected competition
	for _, comp := range p.allCompetitions {
		if comp.CompetitionTitle == selected {
			p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), comp.CompetitionTitle, comp.Status))
			return
		}
	}
}

// Refresh updates the panel with new locale strings
func (p *MonitoringPanel) Refresh() {
	p.content = p.createContent()
}
