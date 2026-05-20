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
	participantService    *service.CompetitionParticipantService
	competitorModelService *service.CompetitorModelService
	competitorService     *service.CompetitorService
	modelService          *service.RCModelService
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
	resultsTable          *widget.Table
	lapHistoryList        *widget.List
	currentLapData        map[string]*service.ParticipantLapData
	participants          []models.CompetitionParticipant
	competitorModels      map[string]models.CompetitorModel
	competitors           map[string]models.Competitor
	rcModels              map[string]models.RCModel
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, participantService *service.CompetitionParticipantService, competitorModelService *service.CompetitorModelService, competitorService *service.CompetitorService, modelService *service.RCModelService, lapService *service.LapService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:             mainWindow,
		competitionService:     competitionService,
		participantService:     participantService,
		competitorModelService: competitorModelService,
		competitorService:      competitorService,
		modelService:           modelService,
		lapService:             lapService,
		currentLapData:         make(map[string]*service.ParticipantLapData),
		competitorModels:       make(map[string]models.CompetitorModel),
		competitors:            make(map[string]models.Competitor),
		rcModels:               make(map[string]models.RCModel),
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

	// Create results table
	p.resultsTable = widget.NewTable(
		func() (int, int) {
			return len(p.participants), 8 // 8 columns: Position, Name, Model, Laps, Best Lap, Last Lap, Total Time, Gap
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Sample")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			if label, ok := obj.(*widget.Label); ok {
				if id.Row < len(p.participants) {
					p.updateTableCell(label, id)
				}
			}
		},
	)

	// Create lap history list
	p.lapHistoryList = widget.NewList(
		func() int {
			if p.currentLapData != nil {
				count := 0
				for _, data := range p.currentLapData {
					count += len(data.LapTimes)
				}
				return count
			}
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Lap history")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if label, ok := obj.(*widget.Label); ok {
				p.updateLapHistoryItem(label, id)
			}
		},
	)

	// Main content area with results table and lap history
	monitoringContent := container.NewHSplit(
		container.NewScroll(p.resultsTable),
		container.NewBorder(
			widget.NewLabel(locale.T("monitoring.lap_history")),
			nil,
			nil,
			nil,
			container.NewScroll(p.lapHistoryList),
		),
	)

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
			// Load participants and lap data for display
			p.loadCompetitionData(comp.ID)
			return
		}
	}
}

// loadCompetitionData loads participants and lap data for the selected competition
func (p *MonitoringPanel) loadCompetitionData(competitionID string) {
	// Load participants
	participants, err := p.participantService.GetParticipantsByCompetitionID(competitionID)
	if err != nil {
		fmt.Println("ERROR loading participants:", err)
		return
	}
	p.participants = participants

	// Load competitor models and build cache
	allModels, err := p.competitorModelService.GetAllCompetitorModels()
	if err != nil {
		fmt.Println("ERROR loading competitor models:", err)
		return
	}
	p.competitorModels = make(map[string]models.CompetitorModel)
	for _, cm := range allModels {
		p.competitorModels[cm.ID] = cm
	}

	// Load competitors
	allCompetitors, err := p.competitorService.GetAllCompetitors()
	if err != nil {
		fmt.Println("ERROR loading competitors:", err)
		return
	}
	p.competitors = make(map[string]models.Competitor)
	for _, c := range allCompetitors {
		p.competitors[c.ID] = c
	}

	// Load RC models
	allModelsRC, err := p.modelService.GetAllModels()
	if err != nil {
		fmt.Println("ERROR loading RC models:", err)
		return
	}
	p.rcModels = make(map[string]models.RCModel)
	for _, m := range allModelsRC {
		p.rcModels[m.ID] = m
	}

	// Get current lap data from LapService
	p.currentLapData = p.lapService.GetParticipantResults()

	// Refresh UI
	fyne.Do(func() {
		if p.resultsTable != nil {
			p.resultsTable.Refresh()
		}
		if p.lapHistoryList != nil {
			p.lapHistoryList.Refresh()
		}
	})
}

// updateTableCell updates a cell in the results table
func (p *MonitoringPanel) updateTableCell(label *widget.Label, id widget.TableCellID) {
	if id.Row >= len(p.participants) {
		return
	}

	participant := p.participants[id.Row]
	cm, cmOk := p.competitorModels[participant.CompetitorModelID]
	
	var competitor models.Competitor
	var rcModel models.RCModel
	
	if cmOk {
		competitor = p.competitors[cm.CompetitorID]
		rcModel = p.rcModels[cm.RCModelID]
	}

	lapData, hasLapData := p.currentLapData[participant.ID]

	switch id.Col {
	case 0: // Position
		label.SetText(fmt.Sprintf("%d", id.Row+1))
	case 1: // Competitor Name
		if cmOk && competitor.ID != "" {
			label.SetText(competitor.FullName)
		} else {
			label.SetText("-")
		}
	case 2: // Model
		if cmOk && rcModel.ID != "" {
			label.SetText(fmt.Sprintf("%s %s", rcModel.Brand, rcModel.ModelName))
		} else {
			label.SetText("-")
		}
	case 3: // Laps
		if hasLapData {
			label.SetText(fmt.Sprintf("%d", lapData.LapCount))
		} else {
			label.SetText("0")
		}
	case 4: // Best Lap
		if hasLapData && lapData.BestLapTimeMs > 0 {
			label.SetText(formatLapTime(lapData.BestLapTimeMs))
		} else {
			label.SetText("-")
		}
	case 5: // Last Lap
		if hasLapData && lapData.LastLapTimeMs > 0 {
			label.SetText(formatLapTime(lapData.LastLapTimeMs))
		} else {
			label.SetText("-")
		}
	case 6: // Total Time
		if hasLapData {
			label.SetText(formatLapTime(lapData.TotalTimeMs))
		} else {
			label.SetText("00:00.000")
		}
	case 7: // Gap to Leader
		if id.Row == 0 || !hasLapData {
			label.SetText("-")
		} else {
			// Calculate gap to leader (first participant)
			if len(p.participants) > 0 {
				leaderID := p.participants[0].ID
				leaderData, hasLeader := p.currentLapData[leaderID]
				if hasLeader && hasLapData && leaderData.LapCount > 0 {
					gap := lapData.TotalTimeMs - leaderData.TotalTimeMs
					if gap > 0 {
						label.SetText("+" + formatLapTime(gap))
					} else {
						label.SetText("-")
					}
				} else {
					label.SetText("-")
				}
			} else {
				label.SetText("-")
			}
		}
	}
}

// updateLapHistoryItem updates a lap history list item
func (p *MonitoringPanel) updateLapHistoryItem(label *widget.Label, id widget.ListItemID) {
	// Collect all laps from all participants
	type LapEntry struct {
		CompetitorName string
		LapNumber      int
		LapTime        int
		Time           time.Time
	}
	
	var allLaps []LapEntry
	
	for participantID, data := range p.currentLapData {
		// Find participant
		var participantName string
		for _, part := range p.participants {
			if part.ID == participantID {
				if cm, ok := p.competitorModels[part.CompetitorModelID]; ok {
					if c, ok2 := p.competitors[cm.CompetitorID]; ok2 {
						participantName = c.FullName
					}
				}
				break
			}
		}
		
		// Add all laps for this participant
		for i, lapTime := range data.LapTimes {
			allLaps = append(allLaps, LapEntry{
				CompetitorName: participantName,
				LapNumber:      i + 1,
				LapTime:        lapTime,
			})
		}
	}
	
	if id < len(allLaps) {
		lap := allLaps[id]
		label.SetText(fmt.Sprintf("%s - Lap %d: %s", lap.CompetitorName, lap.LapNumber, formatLapTime(lap.LapTime)))
	} else {
		label.SetText("")
	}
}

// formatLapTime formats milliseconds into MM:SS.mmm format
func formatLapTime(ms int) string {
	duration := time.Duration(ms) * time.Millisecond
	minutes := int(duration.Minutes())
	seconds := duration.Seconds() - float64(minutes*60)
	return fmt.Sprintf("%d:%06.3f", minutes, seconds)
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
