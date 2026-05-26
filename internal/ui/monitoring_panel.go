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
	content               *fyne.Container
	mainWindow            fyne.Window
	competitionService    *service.CompetitionService
	lapService            *service.LapService
	selectedCompetition   string
	selectedCompetitionID string
	statusLabel           *widget.Label
	allCompetitions       []models.Competition
	competitionButton     *widget.Button
	startButton           *widget.Button
	stopButton            *widget.Button
	timerLabel            *widget.Label
	timer                 *Timer
	competitionFilter     *CompetitionFilter
	monitoringTable       *widget.Table
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, lapService *service.LapService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:         mainWindow,
		competitionService: competitionService,
		lapService:         lapService,
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

	// Stop button - disabled by default, enabled only when competition status is in_progress
	p.stopButton = widget.NewButton(locale.T("button.stop"), func() {
		p.stopMonitoring()
	})
	p.stopButton.Disable()

	// Timer label - displays elapsed time during monitoring
	p.timerLabel = widget.NewLabel("00:00:00.00")
	p.timerLabel.Alignment = fyne.TextAlignCenter

	// Initialize timer
	p.timer = NewTimer(p.mainWindow, p.timerLabel)

	// Create competition filter component
	p.competitionFilter = NewCompetitionFilter(p.mainWindow, func() {
		// Filter changed callback - refresh competition popup if open
	})

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewSeparator(),
		p.competitionFilter.CreateContent(),
		widget.NewSeparator(),
		container.NewHBox(p.competitionButton, p.startButton, p.stopButton, p.timerLabel),
		widget.NewSeparator(),
	)

	// Initialize monitoring table (will be populated when competition is selected)
	p.monitoringTable = p.createMonitoringTable()

	// Layout
	content := container.NewBorder(
		selectorContainer,
		nil,
		nil,
		nil,
		p.monitoringTable,
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

	// Set competitions to filter component
	if p.competitionFilter != nil {
		p.competitionFilter.SetCompetitions(p.allCompetitions)
	}
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
				result := make([]ReferenceItem, len(p.competitionFilter.GetFilteredCompetitions()))
				for i, comp := range p.competitionFilter.GetFilteredCompetitions() {
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
		p.competitionFilter.GetFilteredCompetitionTitles(),
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
		p.timer.Stop()
		p.statusLabel.SetText(locale.T("status.ready"))
		if p.competitionButton != nil {
			p.competitionButton.SetText(locale.T("form.competition.select"))
		}
		if p.startButton != nil {
			p.startButton.Disable()
		}
		if p.stopButton != nil {
			p.stopButton.Disable()
		}
		if p.timerLabel != nil {
			p.timer.Reset()
		}
		p.selectedCompetitionID = ""
		return
	}

	// Stop and reset timer when switching to a different competition
	p.timer.Stop()

	// Find the selected competition
	for _, comp := range p.allCompetitions {
		if comp.CompetitionTitle == selected {
			p.selectedCompetitionID = comp.ID
			p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), comp.CompetitionTitle, comp.Status))
			if p.competitionButton != nil {
				p.competitionButton.SetText(comp.CompetitionTitle)
			}
			// Enable Start button only if competition status is scheduled
			if p.startButton != nil {
				if comp.Status == "scheduled" {
					p.startButton.Enable()
				} else {
					p.startButton.Disable()
				}
			}
			// Enable Stop button only if competition status is in_progress
			if p.stopButton != nil {
				if comp.Status == "in_progress" {
					p.stopButton.Enable()
					// Start timer when competition is in progress
					p.timer.Start(comp.TimeLimitMinutes, func() {
						p.stopMonitoring()
					})
				} else {
					p.stopButton.Disable()
				}
			}
			// Refresh monitoring table with participant data
			p.refreshMonitoringTable()
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
	// Refresh the selected competition state from DB
	if p.selectedCompetition != "" {
		p.onCompetitionSelected(p.selectedCompetition)
	}
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

	// Refresh competitions list to get updated status from DB
	p.refreshCompetitions()

	// Reload and update the selected competition data from DB
	p.onCompetitionSelected(p.selectedCompetition)

	// Show success message
	dialog.ShowInformation(locale.T("dialog.success"), locale.T("dialog.competition_started"), p.mainWindow)
}

// stopMonitoring stops the selected competition by changing its status to "finished"
func (p *MonitoringPanel) stopMonitoring() {
	if p.selectedCompetitionID == "" {
		dialog.ShowError(fmt.Errorf(locale.T("error.no_competition_selected")), p.mainWindow)
		return
	}

	err := p.competitionService.StopCompetition(p.selectedCompetitionID)
	if err != nil {
		dialog.ShowError(err, p.mainWindow)
		return
	}

	// Refresh competitions list to get updated status from DB
	p.refreshCompetitions()

	// Reload and update the selected competition data from DB
	p.onCompetitionSelected(p.selectedCompetition)

	// Show success message
	dialog.ShowInformation(locale.T("dialog.success"), locale.T("dialog.competition_stopped"), p.mainWindow)
}

// createMonitoringTable creates the table widget for displaying participant lap data
func (p *MonitoringPanel) createMonitoringTable() *widget.Table {
	// Table data storage
	var data []service.MonitoringParticipantData

	table := widget.NewTable(
		func() (int, int) {
			return len(data), 7 // rows, columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Sample")
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if id.Row < len(data) {
				row := data[id.Row]
				switch id.Col {
				case 0: // Transponder Active
					if row.TransponderActive {
						label.SetText("✓")
					} else {
						label.SetText("✗")
					}
				case 1: // Competitor Number
					label.SetText(fmt.Sprintf("%d", row.CompetitorNumber))
				case 2: // Full Name
					label.SetText(row.FullName)
				case 3: // Model (Scale + Name)
					label.SetText(fmt.Sprintf("%s %s", row.ModelScale, row.ModelName))
				case 4: // Laps
					label.SetText(fmt.Sprintf("%d", row.Laps))
				case 5: // Best Lap Time
					label.SetText(row.BestLapTimeFormatted)
				case 6: // Transponder Number
					label.SetText(row.TransponderNumber)
				}
			}
		},
	)

	// Set column widths
	table.SetColumnWidth(0, 80)  // Transponder Active
	table.SetColumnWidth(1, 80)  // Number
	table.SetColumnWidth(2, 200) // Full Name
	table.SetColumnWidth(3, 180) // Model
	table.SetColumnWidth(4, 80)  // Laps
	table.SetColumnWidth(5, 120) // Best Lap Time
	table.SetColumnWidth(6, 120) // Transponder Number

	// Store reference for updating
	p.monitoringTable = table

	return table
}

// refreshMonitoringTable updates the monitoring table with current data
func (p *MonitoringPanel) refreshMonitoringTable() {
	if p.lapService == nil || p.selectedCompetitionID == "" {
		return
	}

	data, err := p.lapService.GetMonitoringData(p.selectedCompetitionID)
	if err != nil {
		fmt.Println("ERROR loading monitoring data:", err)
		return
	}

	if p.monitoringTable != nil {
		p.monitoringTable.Refresh()
	}
}
