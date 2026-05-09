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
	selectedCompetition   string
	selectedCompetitionID string
	statusLabel           *widget.Label
	allCompetitions       []models.Competition
	competitionButton     *widget.Button
	startButton           *widget.Button
	stopButton            *widget.Button
	timerLabel            *widget.Label
	timerTicker           *time.Ticker
	timerStop             chan struct{}
	filteredCompetitions  []models.Competition
	filteredYears         []string
	filteredSeasons       []string
	filteredTracks        []string
	filteredModelTypes    []string
	yearButton            *widget.Button
	seasonButton          *widget.Button
	trackButton           *widget.Button
	modelTypeButton       *widget.Button
	selectedYear          string
	selectedSeason        string
	selectedTrack         string
	selectedModelType     string
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
	p.timerLabel = widget.NewLabel("")
	p.timerLabel.Alignment = fyne.TextAlignCenter

	// Filter buttons using reference_popup.go without add/delete functionality
	p.yearButton = widget.NewButton(locale.T("filter.all_years"), func() {
		p.showYearFilterPopup()
	})
	p.seasonButton = widget.NewButton(locale.T("filter.all_seasons"), func() {
		p.showSeasonFilterPopup()
	})
	p.trackButton = widget.NewButton(locale.T("filter.all_tracks"), func() {
		p.showTrackFilterPopup()
	})
	p.modelTypeButton = widget.NewButton(locale.T("filter.all_model_types"), func() {
		p.showModelTypeFilterPopup()
	})

	// Filter container
	filterContainer := container.NewGridWithColumns(4,
		widget.NewLabel(locale.T("common.year")),
		widget.NewLabel(locale.T("common.season")),
		widget.NewLabel(locale.T("common.track")),
		widget.NewLabel(locale.T("common.model_type")),
		p.yearButton,
		p.seasonButton,
		p.trackButton,
		p.modelTypeButton,
	)

	// Selector container
	selectorContainer := container.NewVBox(
		widget.NewSeparator(),
		filterContainer,
		widget.NewSeparator(),
		container.NewHBox(p.competitionButton, p.startButton, p.stopButton, p.timerLabel),
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

	// Populate filter options from loaded competitions
	p.populateFilterOptions()

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

	p.filteredYears = []string{locale.T("filter.all_years")}
	for year := range yearSet {
		p.filteredYears = append(p.filteredYears, year)
	}

	p.filteredSeasons = []string{locale.T("filter.all_seasons")}
	for season := range seasonSet {
		p.filteredSeasons = append(p.filteredSeasons, season)
	}

	p.filteredTracks = []string{locale.T("filter.all_tracks")}
	for track := range trackSet {
		p.filteredTracks = append(p.filteredTracks, track)
	}

	p.filteredModelTypes = []string{locale.T("filter.all_model_types")}
	for modelType := range modelTypeSet {
		p.filteredModelTypes = append(p.filteredModelTypes, modelType)
	}
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
	if p.selectedYear != "" && p.selectedYear != locale.T("filter.all_years") {
		if comp.CompetitionYear == nil {
			return false
		}
		yearStr := fmt.Sprintf("%d", *comp.CompetitionYear)
		if yearStr != p.selectedYear {
			return false
		}
	}

	// Check season filter
	if p.selectedSeason != "" && p.selectedSeason != locale.T("filter.all_seasons") {
		if comp.Season != p.selectedSeason {
			return false
		}
	}

	// Check track filter
	if p.selectedTrack != "" && p.selectedTrack != locale.T("filter.all_tracks") {
		if comp.TrackName != p.selectedTrack {
			return false
		}
	}

	// Check model type filter
	if p.selectedModelType != "" && p.selectedModelType != locale.T("filter.all_model_types") {
		if comp.ModelType != p.selectedModelType {
			return false
		}
	}

	return true
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

// showYearFilterPopup shows the year filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showYearFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredYears))
	copy(items, p.filteredYears)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
		ReferencePopupConfig{
			Title:          "common.year",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				p.selectedYear = selected
				p.yearButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedYear = selected
			p.yearButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showSeasonFilterPopup shows the season filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showSeasonFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredSeasons))
	copy(items, p.filteredSeasons)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
		ReferencePopupConfig{
			Title:          "common.season",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				p.selectedSeason = selected
				p.seasonButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedSeason = selected
			p.seasonButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showTrackFilterPopup shows the track filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showTrackFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredTracks))
	copy(items, p.filteredTracks)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
		ReferencePopupConfig{
			Title:          "common.track",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				p.selectedTrack = selected
				p.trackButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedTrack = selected
			p.trackButton.SetText(selected)
			p.applyFilters()
		},
		func(opts []string) {},
	)
	popupManager.ShowPopupWithoutAddDelete(mainDialog, &currentDialog, func(d dialog.Dialog) {
		currentDialog = d
	})
}

// showModelTypeFilterPopup shows the model type filter popup using reference_popup.go without add/delete
func (p *MonitoringPanel) showModelTypeFilterPopup() {
	var currentDialog dialog.Dialog
	var mainDialog dialog.Dialog

	// Refresh filter options from competitions before showing popup
	p.populateFilterOptions()

	items := make([]string, len(p.filteredModelTypes))
	copy(items, p.filteredModelTypes)

	popupManager := NewReferencePopupManager(
		p.mainWindow,
		ReferencePopupConfig{
			Title:          "common.model_type",
			AddTitle:       "",
			AddLabel:       "",
			AddPlaceholder: "",
			DeleteMessage:  "",
			NewErrorExists: "",
			EnterNameInfo:  "",
			GetAllFunc: func() ([]ReferenceItem, error) {
				result := make([]ReferenceItem, len(items))
				for i, item := range items {
					result[i] = ReferenceItem{Name: item}
				}
				return result, nil
			},
			AddFunc:    func(name string) error { return nil },
			DeleteFunc: func(name string) error { return nil },
			OnItemSelected: func(selected string) {
				p.selectedModelType = selected
				p.modelTypeButton.SetText(selected)
				p.applyFilters()
			},
			UpdateOptions: func(opts []string) {},
		},
		items,
		"",
		func(selected string) {
			p.selectedModelType = selected
			p.modelTypeButton.SetText(selected)
			p.applyFilters()
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
		p.stopTimer()
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
			p.timerLabel.SetText("")
		}
		p.selectedCompetitionID = ""
		return
	}

	// Stop and reset timer when switching to a different competition
	p.stopTimer()

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
					p.startTimer(comp.TimeLimitMinutes)
				} else {
					p.stopButton.Disable()
				}
			}
			return
		}
	}
}

// Refresh updates the panel with new locale strings
func (p *MonitoringPanel) Refresh() {
	p.content = p.createContent()
}

// startTimer starts the timer for monitoring elapsed time
func (p *MonitoringPanel) startTimer(timeLimitMinutes *int) {
	// Stop any existing timer first
	p.stopTimer()

	startTime := time.Now()
	var limitReached bool

	p.timerStop = make(chan struct{})
	p.timerTicker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-p.timerTicker.C:
				elapsed := time.Since(startTime)

				// Check if time limit is reached
				if timeLimitMinutes != nil && !limitReached {
					limitDuration := time.Duration(*timeLimitMinutes) * time.Minute
					if elapsed >= limitDuration {
						limitReached = true
						// Timer reached limit, stop it and notify on main thread
						fyne.Do(func() {
							p.stopTimer()
							p.stopMonitoring()
							dialog.ShowInformation(locale.T("dialog.info"), locale.T("dialog.time_limit_reached"), p.mainWindow)
						})
						return
					}
				}

				// Update timer label on main thread
				finalElapsed := elapsed
				fyne.Do(func() {
					if p.timerLabel != nil {
						p.timerLabel.SetText(p.formatDuration(finalElapsed))
					}
				})
			case <-p.timerStop:
				return
			}
		}
	}()
}

// stopTimer stops the running timer
func (p *MonitoringPanel) stopTimer() {
	if p.timerTicker != nil {
		p.timerTicker.Stop()
		p.timerTicker = nil
	}
	if p.timerStop != nil {
		close(p.timerStop)
		p.timerStop = nil
	}
	// Do not clear timer label - keep showing the last elapsed time
	// Timer label is only cleared when switching to a different competition
}

// formatDuration formats a duration as MM:SS or HH:MM:SS
func (p *MonitoringPanel) formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
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
