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
	lapService            *service.LapService
	selectedCompetitionID string
	statusLabel           *widget.Label
	allCompetitions       []models.Competition
	competitionButton     *widget.Button
	startButton           *widget.Button
	stopButton            *widget.Button
	timerLabel            *widget.Label
	timer                 *Timer
	competitionFilter     *CompetitionFilter
	standings             []models.CompetitionStanding
	standingsTable        *widget.Table
	stopUpdates           chan struct{}
	timerCompetitionID    string
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, lapService *service.LapService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:         mainWindow,
		competitionService: competitionService,
		lapService:         lapService,
		stopUpdates:        make(chan struct{}),
	}

	p.content = p.createContent()
	if lapService != nil {
		lapService.SetErrorHandler(func(err error) {
			fyne.Do(func() {
				p.statusLabel.SetText(fmt.Sprintf(locale.T("monitoring.save_error"), err))
			})
		})
		go p.updateStandingsLoop()
		if active := lapService.GetActiveCompetition(); active != nil {
			p.onCompetitionSelected(active.ID)
		}
	}
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
		p.statusLabel,
		widget.NewSeparator(),
	)

	p.standingsTable = p.createStandingsTable()

	// Layout
	content := container.NewBorder(
		selectorContainer,
		nil,
		nil,
		nil,
		p.standingsTable,
	)

	// Load competitions
	p.refreshCompetitions()

	return content
}

func (p *MonitoringPanel) createStandingsTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return len(p.standings), 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			if id.Row < 0 || id.Row >= len(p.standings) {
				return
			}
			row := p.standings[id.Row]
			var text string
			switch id.Col {
			case 0:
				text = fmt.Sprintf("%d", id.Row+1)
			case 1:
				text = fmt.Sprintf("%d", row.CompetitorNumber)
			case 2:
				text = row.FullName
			case 3:
				text = fmt.Sprintf("%d", row.NumberOfLaps)
			case 4:
				text = "–"
				if row.BestLapTimeMs > 0 {
					text = fmt.Sprintf("%02d:%02d:%03d", row.BestLapTimeMs/60000, row.BestLapTimeMs/1000%60, row.BestLapTimeMs%1000)
				}
			}
			obj.(*widget.Label).SetText(text)
		},
	)
	table.ShowHeaderRow = true
	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		keys := []string{"monitoring.position", "monitoring.competitor_number", "monitoring.competitor_name", "monitoring.laps", "monitoring.best_lap"}
		if id.Col >= 0 && id.Col < len(keys) {
			obj.(*widget.Label).SetText(locale.T(keys[id.Col]))
		}
	}
	for i, width := range []float32{70, 150, 300, 120, 220} {
		table.SetColumnWidth(i, width)
	}
	return table
}

func (p *MonitoringPanel) refreshStandings() {
	if p.lapService == nil || p.selectedCompetitionID == "" {
		p.standings = nil
	} else {
		standings, err := p.lapService.GetStandings(p.selectedCompetitionID)
		if err != nil {
			p.statusLabel.SetText(locale.T("status.error_loading"))
			fmt.Printf("failed to load standings: %v\n", err)
			return
		}
		p.standings = standings
	}
	if p.standingsTable != nil {
		p.standingsTable.Refresh()
	}
}

func (p *MonitoringPanel) updateStandingsLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fyne.Do(func() { p.refreshStandings() })
		case <-p.stopUpdates:
			return
		}
	}
}

// Close stops the background refresh loop.
func (p *MonitoringPanel) Close() {
	close(p.stopUpdates)
	p.timer.Stop()
	if p.lapService != nil {
		p.lapService.SetErrorHandler(nil)
	}
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
	items := container.NewVBox()
	popup := dialog.NewCustomWithoutButtons(locale.T("form.competition.select"), container.NewVScroll(items), p.mainWindow)
	for _, comp := range p.competitionFilter.GetFilteredCompetitions() {
		id := comp.ID
		button := widget.NewButton(comp.CompetitionTitle, func() {
			popup.Hide()
			p.onCompetitionSelected(id)
		})
		items.Add(button)
	}
	popup.SetButtons([]fyne.CanvasObject{widget.NewButton(locale.T("common.close"), popup.Hide)})
	popup.Resize(fyne.NewSize(500, 400))
	popup.Show()
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
func (p *MonitoringPanel) onCompetitionSelected(selectedID string) {
	if selectedID == "" {
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
		p.refreshStandings()
		return
	}

	// Find the selected competition
	for _, comp := range p.allCompetitions {
		if comp.ID == selectedID {
			p.selectedCompetitionID = comp.ID
			p.statusLabel.SetText(fmt.Sprintf("%s: %s (%s)", locale.T("common.selected"), comp.CompetitionTitle, comp.Status))
			if p.competitionButton != nil {
				p.competitionButton.SetText(comp.CompetitionTitle)
			}
			// Enable Start button only if competition status is scheduled
			if p.startButton != nil {
				if comp.Status == "scheduled" && (p.lapService == nil || p.lapService.GetActiveCompetition() == nil) {
					p.startButton.Enable()
				} else {
					p.startButton.Disable()
				}
			}
			// Enable Stop button only if competition status is in_progress
			if p.stopButton != nil {
				if comp.Status == "in_progress" {
					if p.lapService != nil && p.lapService.GetActiveCompetition() == nil {
						if err := p.lapService.SetActiveCompetition(&comp); err != nil {
							dialog.ShowError(err, p.mainWindow)
							p.stopButton.Disable()
							return
						}
					}
					active := p.lapService != nil && p.lapService.GetActiveCompetition() != nil && p.lapService.GetActiveCompetition().ID == comp.ID
					if !active {
						p.stopButton.Disable()
						p.refreshStandings()
						return
					}
					p.stopButton.Enable()
					// Start timer when competition is in progress
					if p.timerCompetitionID != comp.ID {
						p.timerCompetitionID = comp.ID
						start := time.Now()
						if comp.TimeStart != nil {
							start = *comp.TimeStart
						}
						id := comp.ID
						p.timer.Start(start, comp.TimeLimitMinutes, func() bool { return p.stopCompetition(id) })
					}
				} else {
					p.stopButton.Disable()
				}
			}
			p.refreshStandings()
			return
		}
	}
}

// Refresh updates the panel with new locale strings
func (p *MonitoringPanel) Refresh() {
	p.startButton.SetText(locale.T("button.start"))
	p.stopButton.SetText(locale.T("button.stop"))
	if p.selectedCompetitionID == "" {
		p.competitionButton.SetText(locale.T("form.competition.select"))
	}
	p.standingsTable.Refresh()
}

// UpdateData reloads competition data and refreshes filter options
func (p *MonitoringPanel) UpdateData() {
	p.refreshCompetitions()
	// Refresh the selected competition state from DB
	if p.selectedCompetitionID != "" {
		p.onCompetitionSelected(p.selectedCompetitionID)
	}
}

// startMonitoring starts monitoring for the selected competition
func (p *MonitoringPanel) startMonitoring() {
	if p.selectedCompetitionID == "" {
		dialog.ShowError(fmt.Errorf(locale.T("error.no_competition_selected")), p.mainWindow)
		return
	}
	if p.lapService != nil && p.lapService.GetActiveCompetition() != nil {
		dialog.ShowError(fmt.Errorf("%s", locale.T("error.competition_in_progress")), p.mainWindow)
		return
	}

	err := p.competitionService.StartCompetition(p.selectedCompetitionID)
	if err != nil {
		dialog.ShowError(err, p.mainWindow)
		return
	}
	comp, err := p.competitionService.GetCompetitionByID(p.selectedCompetitionID)
	if err != nil || comp == nil {
		p.competitionService.SetCompetitionScheduled(p.selectedCompetitionID)
		if err == nil {
			err = fmt.Errorf("competition not found")
		}
		dialog.ShowError(err, p.mainWindow)
		return
	}
	if p.lapService != nil {
		if err := p.lapService.SetActiveCompetition(comp); err != nil {
			p.competitionService.SetCompetitionScheduled(p.selectedCompetitionID)
			dialog.ShowError(err, p.mainWindow)
			return
		}
	}

	// Refresh competitions list to get updated status from DB
	p.refreshCompetitions()

	// Reload and update the selected competition data from DB
	p.onCompetitionSelected(p.selectedCompetitionID)

	// Show success message
	dialog.ShowInformation(locale.T("dialog.success"), locale.T("dialog.competition_started"), p.mainWindow)
}

// stopMonitoring stops the selected competition by changing its status to "finished"
func (p *MonitoringPanel) stopMonitoring() {
	if p.selectedCompetitionID == "" {
		dialog.ShowError(fmt.Errorf(locale.T("error.no_competition_selected")), p.mainWindow)
		return
	}

	p.stopCompetition(p.selectedCompetitionID)
}

func (p *MonitoringPanel) stopCompetition(id string) bool {
	if p.lapService != nil {
		active := p.lapService.GetActiveCompetition()
		if active == nil || active.ID != id {
			return false
		}
		p.lapService.Flush()
		if err := p.lapService.PersistResults(); err != nil {
			dialog.ShowError(err, p.mainWindow)
			return false
		}
		p.lapService.SetActiveCompetition(nil)
	}
	err := p.competitionService.StopCompetition(id)
	if err != nil {
		if p.lapService != nil {
			if comp, loadErr := p.competitionService.GetCompetitionByID(id); loadErr == nil && comp != nil {
				p.lapService.SetActiveCompetition(comp)
			}
		}
		dialog.ShowError(err, p.mainWindow)
		return false
	}
	p.timer.Stop()
	p.timerCompetitionID = ""

	// Refresh competitions list to get updated status from DB
	p.refreshCompetitions()

	// Reload and update the selected competition data from DB
	p.onCompetitionSelected(p.selectedCompetitionID)

	// Show success message
	dialog.ShowInformation(locale.T("dialog.success"), locale.T("dialog.competition_stopped"), p.mainWindow)
	return true
}
