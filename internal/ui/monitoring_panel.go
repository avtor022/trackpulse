package ui

import (
	"fmt"
	"time"

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
	monitoringService     *service.MonitoringService
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
	participantsTable     *widget.Table
	participantData       []service.MonitoringParticipantData
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, monitoringService *service.MonitoringService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:        mainWindow,
		competitionService: competitionService,
		monitoringService:  monitoringService,
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

	// Main content area with participants table
	p.participantsTable = widget.NewTable(
		func() (int, int) {
			if len(p.participantData) == 0 {
				return 0, 0
			}
			return len(p.participantData), 7
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.participantData) {
				o.(*widget.Label).SetText("")
				return
			}
			data := p.participantData[i.Row]

			switch i.Col {
			case 0: // Transponder Active
				status := "❌"
				if data.TransponderActive {
					status = "✅"
				}
				o.(*widget.Label).SetText(status)
			case 1: // Competitor Number
				o.(*widget.Label).SetText(fmt.Sprintf("%d", data.CompetitorNumber))
			case 2: // Full Name
				o.(*widget.Label).SetText(data.FullName)
			case 3: // Model
				o.(*widget.Label).SetText(fmt.Sprintf("%s (%s)", data.ModelName, data.ModelScale))
			case 4: // Lap Count
				o.(*widget.Label).SetText(fmt.Sprintf("%d", data.LapCount))
			case 5: // Best Lap Time
				if data.BestLapTimeMs > 0 {
					o.(*widget.Label).SetText(formatDuration(time.Duration(data.BestLapTimeMs) * time.Millisecond))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6: // Total Time
				if data.TotalCompetitionTimeMs > 0 {
					o.(*widget.Label).SetText(formatDuration(time.Duration(data.TotalCompetitionTimeMs) * time.Millisecond))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		},
	)
	p.participantsTable.ShowHeaderRow = true
	p.participantsTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("Header", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}
	p.participantsTable.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		headers := []string{
			"RFID",
			locale.T("competitor.header.number"),
			locale.T("competitor.header.name"),
			locale.T("participants.table.model"),
			locale.T("competition.header.laps"),
			locale.T("monitoring.best_lap"),
			locale.T("monitoring.total_time"),
		}
		if id.Col >= 0 && id.Col < len(headers) {
			o.(*widget.Label).SetText(headers[id.Col])
			o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
		}
	}

	// Set column widths
	p.participantsTable.SetColumnWidth(0, 60)  // RFID status
	p.participantsTable.SetColumnWidth(1, 80)  // Number
	p.participantsTable.SetColumnWidth(2, 200) // Name
	p.participantsTable.SetColumnWidth(3, 200) // Model
	p.participantsTable.SetColumnWidth(4, 80)  // Laps
	p.participantsTable.SetColumnWidth(5, 120) // Best Lap
	p.participantsTable.SetColumnWidth(6, 120) // Total Time

	monitoringContent := container.NewScroll(p.participantsTable)
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
		p.clearParticipantsTable()
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
			// Load participant data for monitoring display immediately after selection
			p.loadParticipantData()
			return
		}
	}
}

// loadParticipantData loads participant data for the selected competition
func (p *MonitoringPanel) loadParticipantData() {
	if p.selectedCompetitionID == "" || p.monitoringService == nil {
		return
	}

	data, err := p.monitoringService.GetParticipantsForCompetition(p.selectedCompetitionID)
	if err != nil {
		fmt.Println("ERROR loading participant data:", err)
		return
	}

	p.participantData = data
	if p.participantsTable != nil {
		p.participantsTable.Refresh()
	}
}

// clearParticipantsTable clears the participants table
func (p *MonitoringPanel) clearParticipantsTable() {
	p.participantData = nil
	if p.participantsTable != nil {
		p.participantsTable.Refresh()
	}
}

// formatDuration formats a duration into MM:SS.ms format
func formatDuration(d time.Duration) string {
	ms := d.Milliseconds()
	minutes := ms / 60000
	seconds := (ms % 60000) / 1000
	millis := ms % 1000
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, millis)
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
