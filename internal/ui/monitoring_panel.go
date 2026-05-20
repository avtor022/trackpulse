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
	content                *fyne.Container
	mainWindow             fyne.Window
	competitionService     *service.CompetitionService
	participantService     *service.CompetitionParticipantService
	selectedCompetition    string
	selectedCompetitionID  string
	statusLabel            *widget.Label
	allCompetitions        []models.Competition
	competitionButton      *widget.Button
	startButton            *widget.Button
	stopButton             *widget.Button
	timerLabel             *widget.Label
	timer                  *Timer
	competitionFilter      *CompetitionFilter
	participantsTable      *widget.Table
	participantsContainer  *fyne.Container
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, participantService *service.CompetitionParticipantService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:         mainWindow,
		competitionService: competitionService,
		participantService: participantService,
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

	// Create participants table container (will be populated when competition is selected)
	p.participantsContainer = p.createParticipantsTable()

	// Main content area with participants table
	monitoringContent := container.NewBorder(nil, nil, nil, nil, p.participantsContainer)

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

// createParticipantsTable creates the participants registration table
func (p *MonitoringPanel) createParticipantsTable() *fyne.Container {
	// Headers for the table
	headers := []string{
		locale.T("monitoring.participants.transponder"),
		locale.T("monitoring.participants.number"),
		locale.T("monitoring.participants.name"),
		locale.T("monitoring.participants.model"),
		locale.T("monitoring.participants.scale"),
		locale.T("monitoring.participants.laps"),
		locale.T("monitoring.participants.best_lap"),
	}

	var data [][]string

	p.participantsTable = widget.NewTable(
		func() (int, int) {
			if len(data) == 0 {
				return 0, 0
			}
			return len(data) + 1, len(headers) // +1 for header row
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			if id.Row == 0 {
				// Header row
				label.Text = headers[id.Col]
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// Data row
				if id.Col < len(data[id.Row-1]) {
					label.Text = data[id.Row-1][id.Col]
				}
			}
		},
	)

	p.participantsTable.ShowHeaderRow = false
	p.participantsTable.ShowHeaderColumn = false

	return container.NewScroll(p.participantsTable)
}

// updateParticipantsTable updates the participants table with current data
func (p *MonitoringPanel) updateParticipantsTable() {
	if p.participantsTable == nil || p.selectedCompetitionID == "" {
		return
	}

	if p.participantService == nil {
		return
	}

	participantsData, err := p.participantService.GetParticipantRegistrationData(p.selectedCompetitionID)
	if err != nil {
		fmt.Println("ERROR loading participants data:", err)
		return
	}

	// Build table data
	var data [][]string
	for _, pd := range participantsData {
		transponderStatus := "✗"
		if pd.TransponderWorked {
			transponderStatus = "✓"
		}

		bestLapStr := "-"
		if pd.BestLapTimeMs > 0 {
			bestLapStr = formatDuration(pd.BestLapTimeMs)
		}

		row := []string{
			transponderStatus,
			fmt.Sprintf("%d", pd.CompetitorNumber),
			pd.FullName,
			pd.ModelName,
			pd.ModelScale,
			fmt.Sprintf("%d", pd.LapCount),
			bestLapStr,
		}
		data = append(data, row)
	}

	// Update table and refresh
	p.participantsTable.Refresh()
}

// formatDuration formats milliseconds to MM:SS.mmm format
func formatDuration(ms int) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	milliseconds := ms % 1000
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
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
			// Update participants table after competition selection
			p.updateParticipantsTable()
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
	// Also update participants table if a competition is selected
	if p.selectedCompetitionID != "" {
		p.updateParticipantsTable()
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
