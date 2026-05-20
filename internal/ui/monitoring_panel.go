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
	content                  *fyne.Container
	mainWindow               fyne.Window
	competitionService       *service.CompetitionService
	participantService       *service.CompetitionParticipantService
	lapService               *service.LapService
	competitorModelService   *service.CompetitorModelService
	competitorService        *service.CompetitorService
	rcModelService           *service.RCModelService
	selectedCompetition      string
	selectedCompetitionID    string
	statusLabel              *widget.Label
	allCompetitions          []models.Competition
	competitionButton        *widget.Button
	startButton              *widget.Button
	stopButton               *widget.Button
	timerLabel               *widget.Label
	timer                    *Timer
	competitionFilter        *CompetitionFilter
	participantsTable        *widget.Table
	currentParticipants      []models.CompetitionParticipant
	lapDataMap               map[string]*service.ParticipantLapData
}

// NewMonitoringPanel creates a new monitoring panel
func NewMonitoringPanel(competitionService *service.CompetitionService, participantService *service.CompetitionParticipantService, lapService *service.LapService, competitorModelService *service.CompetitorModelService, competitorService *service.CompetitorService, rcModelService *service.RCModelService, mainWindow fyne.Window) *MonitoringPanel {
	p := &MonitoringPanel{
		mainWindow:             mainWindow,
		competitionService:     competitionService,
		participantService:     participantService,
		lapService:             lapService,
		competitorModelService: competitorModelService,
		competitorService:      competitorService,
		rcModelService:         rcModelService,
		lapDataMap:             make(map[string]*service.ParticipantLapData),
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

	// Create participants table for displaying registered participants and lap data
	p.participantsTable = widget.NewTable(
		func() (int, int) {
			return len(p.currentParticipants), 7
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.currentParticipants) {
				o.(*widget.Label).SetText("")
				return
			}
			participant := p.currentParticipants[i.Row]
			
			// Get lap data from lapService if available
			var lapData *service.ParticipantLapData
			if p.lapDataMap != nil {
				lapData = p.lapDataMap[participant.ID]
			}
			
			switch i.Col {
			case 0: // Position
				pos := i.Row + 1
				o.(*widget.Label).SetText(fmt.Sprintf("%d", pos))
			case 1: // Competitor Name
				o.(*widget.Label).SetText(p.getCompetitorName(participant.CompetitorModelID))
			case 2: // Model
				o.(*widget.Label).SetText(p.getRCModelName(participant.CompetitorModelID))
			case 3: // Transponder
				o.(*widget.Label).SetText(p.getTransponderNumber(participant.CompetitorModelID))
			case 4: // Laps
				if lapData != nil {
					o.(*widget.Label).SetText(fmt.Sprintf("%d", lapData.LapCount))
				} else {
					o.(*widget.Label).SetText("0")
				}
			case 5: // Best Lap
				if lapData != nil && lapData.BestLapTimeMs > 0 {
					o.(*widget.Label).SetText(p.formatLapTime(lapData.BestLapTimeMs))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6: // Last Lap
				if lapData != nil && lapData.LastLapTimeMs > 0 {
					o.(*widget.Label).SetText(p.formatLapTime(lapData.LastLapTimeMs))
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
			locale.T("monitoring.table.position"),
			locale.T("monitoring.table.competitor"),
			locale.T("monitoring.table.model"),
			locale.T("monitoring.table.transponder"),
			locale.T("monitoring.table.laps"),
			locale.T("monitoring.table.best_lap"),
			locale.T("monitoring.table.last_lap"),
		}
		if id.Col >= 0 && id.Col < len(headers) {
			o.(*widget.Label).SetText(headers[id.Col])
			o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
		}
	}
	
	// Set column widths
	p.participantsTable.SetColumnWidth(0, 60)  // Position
	p.participantsTable.SetColumnWidth(1, 200) // Competitor
	p.participantsTable.SetColumnWidth(2, 150) // Model
	p.participantsTable.SetColumnWidth(3, 120) // Transponder
	p.participantsTable.SetColumnWidth(4, 80)  // Laps
	p.participantsTable.SetColumnWidth(5, 100) // Best Lap
	p.participantsTable.SetColumnWidth(6, 100) // Last Lap

	// Main content area with participants table
	monitoringContent := container.NewBorder(
		widget.NewLabelWithStyle(locale.T("monitoring.participants_title"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewScroll(p.participantsTable),
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
		p.currentParticipants = nil
		p.lapDataMap = make(map[string]*service.ParticipantLapData)
		if p.participantsTable != nil {
			p.participantsTable.Refresh()
		}
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
			
			// Load participants for the selected competition
			p.loadParticipantsForCompetition(comp.ID)
			
			return
		}
	}
}

// loadParticipantsForCompetition loads participants for the selected competition
func (p *MonitoringPanel) loadParticipantsForCompetition(competitionID string) {
	if p.participantService == nil {
		return
	}
	
	participants, err := p.participantService.GetParticipantsByCompetitionID(competitionID)
	if err != nil {
		fmt.Println("ERROR loading participants:", err)
		return
	}
	
	p.currentParticipants = participants
	
	// Initialize lap data map for all participants
	p.lapDataMap = make(map[string]*service.ParticipantLapData)
	for _, participant := range participants {
		p.lapDataMap[participant.ID] = &service.ParticipantLapData{
			CompetitionParticipantID: participant.ID,
			CompetitionID:            competitionID,
			CompetitorModelID:        participant.CompetitorModelID,
		}
	}
	
	// Update UI with participants list
	if p.participantsTable != nil {
		p.participantsTable.Refresh()
	}
}

// refreshLapData updates lap data from lapService
func (p *MonitoringPanel) refreshLapData() {
	if p.lapService == nil || p.selectedCompetitionID == "" {
		return
	}
	
	// Get current lap data from lap service
	lapData := p.lapService.GetParticipantResults()
	
	// Update our local map
	for participantID, data := range lapData {
		if _, exists := p.lapDataMap[participantID]; exists {
			p.lapDataMap[participantID] = data
		}
	}
	
	// Refresh table to show updated data
	fyne.Do(func() {
		if p.participantsTable != nil {
			p.participantsTable.Refresh()
		}
	})
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
	// Start periodic lap data refresh if competition is in progress
	if p.selectedCompetitionID != "" {
		p.startLapDataRefresh()
	}
}

// startLapDataRefresh starts periodic refresh of lap data
func (p *MonitoringPanel) startLapDataRefresh() {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for range ticker.C {
			if p.selectedCompetitionID == "" {
				return
			}
			p.refreshLapData()
		}
	}()
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

// getCompetitorName returns the competitor name for a competitor model ID
func (p *MonitoringPanel) getCompetitorName(cmID string) string {
	if p.competitorModelService == nil {
		return cmID[:8]
	}
	
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return cmID[:8]
	}
	
	if p.competitorService == nil {
		return cm.TransponderNumber
	}
	
	competitor, err := p.competitorService.GetCompetitorByID(cm.CompetitorID)
	if err != nil || competitor == nil {
		return cm.TransponderNumber
	}
	
	return fmt.Sprintf("#%d %s", competitor.CompetitorNumber, competitor.FullName)
}

// getRCModelName returns the RC model name for a competitor model ID
func (p *MonitoringPanel) getRCModelName(cmID string) string {
	if p.competitorModelService == nil || p.rcModelService == nil {
		return "Model"
	}
	
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return "Model"
	}
	
	model, err := p.rcModelService.GetModelByID(cm.RCModelID)
	if err != nil || model == nil {
		return "Model"
	}
	
	return fmt.Sprintf("%s %s", model.Brand, model.ModelName)
}

// getTransponderNumber returns the transponder number for a competitor model ID
func (p *MonitoringPanel) getTransponderNumber(cmID string) string {
	if p.competitorModelService == nil {
		return "TP" + cmID[:4]
	}
	
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return "TP" + cmID[:4]
	}
	
	return cm.TransponderNumber
}

// formatLapTime formats lap time in milliseconds to MM:SS.mmm format
func (p *MonitoringPanel) formatLapTime(lapTimeMs int) string {
	if lapTimeMs <= 0 {
		return "-"
	}
	
	totalSeconds := lapTimeMs / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	milliseconds := lapTimeMs % 1000
	
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
}
