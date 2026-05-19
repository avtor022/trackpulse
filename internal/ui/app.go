package ui

import (
	"database/sql"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"trackpulse/internal/locale"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
)

// App represents the main application UI
type App struct {
	fyneApp                fyne.App
	mainWindow             fyne.Window
	competitorService      *service.CompetitorService
	modelService           *service.RCModelService
	settingsService        *service.SettingsService
	competitorModelService *service.CompetitorModelService
	competitionService     *service.CompetitionService
	participantService     *service.CompetitionParticipantService
	lapService             *service.LapService
	rawScanRepo            *repository.RawScanRepository
	config                 *Config
	tabs                   *container.AppTabs
	monitoringPanel        *MonitoringPanel
	competitorPanel        *CompetitorPanel
	modelPanel             *ModelPanel
	competitorModelPanel   *CompetitorModelPanel
	competitionPanel       *CompetitionPanel
	participantPanel       *ParticipantPanel
	settingsPanel          *SettingsPanel
	logsPanel              *LogsPanel
	portScanner            *PortScanner
}

// GlobalApp holds a reference to the main app for locale change notifications
var GlobalApp *App

// Config holds UI configuration
type Config struct {
	Language string
	Title    string
}

// NewApp creates a new TrackPulse application
func NewApp(competitorService *service.CompetitorService, modelService *service.RCModelService, settingsService *service.SettingsService, competitorModelService *service.CompetitorModelService, competitionService *service.CompetitionService, participantService *service.CompetitionParticipantService, language string, db *sql.DB) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("TrackPulse")

	// Initialize repositories for lap service
	rawScanRepo := repository.NewRawScanRepository(db)
	competitorModelRepo := repository.NewCompetitorModelRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)
	participantRepo := repository.NewCompetitionParticipantRepository(db)
	competitionLapsRepo := repository.NewCompetitionLapsRepository(db)

	// Initialize lap service with buffered processing
	lapService := service.NewLapService(rawScanRepo, competitorModelRepo, competitionRepo, participantRepo, competitionLapsRepo)
	lapService.Start()

	appInstance := &App{
		fyneApp:                fyneApp,
		mainWindow:             mainWindow,
		competitorService:      competitorService,
		modelService:           modelService,
		settingsService:        settingsService,
		competitorModelService: competitorModelService,
		competitionService:     competitionService,
		participantService:     participantService,
		lapService:             lapService,
		rawScanRepo:            rawScanRepo,
		config: &Config{
			Language: language,
			Title:    "TrackPulse",
		},
	}

	// Set global reference for locale change notifications
	GlobalApp = appInstance

	return appInstance
}

// Run starts the application UI
func (a *App) Run() {
	a.mainWindow.SetContent(a.createMainContent())
	a.mainWindow.Resize(fyne.NewSize(1200, 800))
	a.mainWindow.ShowAndRun()
}

// createMainContent builds the main tabbed interface
func (a *App) createMainContent() *container.AppTabs {
	a.tabs = container.NewAppTabs(
		container.NewTabItem(locale.T("tab.monitoring"), a.createMonitoringTab()),
		container.NewTabItem(locale.T("tab.competitors"), a.createCompetitorsTab()),
		container.NewTabItem(locale.T("tab.models"), a.createModelsTab()),
		container.NewTabItem(locale.T("tab.transponders"), a.createTranspondersTab()),
		container.NewTabItem(locale.T("tab.competitions"), a.createCompetitionsTab()),
		container.NewTabItem(locale.T("tab.participants"), a.createParticipantsTab()),
		container.NewTabItem(locale.T("tab.logs"), a.createLogsTab()),
		container.NewTabItem(locale.T("tab.settings"), a.createSettingsTab()),
	)

	a.tabs.SetTabLocation(container.TabLocationTop)

	// Set up tab change listener to refresh monitoring panel when switching to it
	a.tabs.OnSelected = func(ti *container.TabItem) {
		if ti == a.tabs.Items[0] && a.monitoringPanel != nil {
			// Refresh monitoring panel data when switching to monitoring tab
			a.monitoringPanel.UpdateData()
		}
		if ti == a.tabs.Items[1] && a.competitorPanel != nil {
			// Refresh competitor panel data when switching to competitors tab
			a.competitorPanel.RefreshData()
		}
		if ti == a.tabs.Items[2] && a.modelPanel != nil {
			// Refresh model panel data when switching to models tab
			a.modelPanel.RefreshData()
		}
		if ti == a.tabs.Items[4] && a.competitionPanel != nil {
			// Refresh competition panel data when switching to competitions tab
			a.competitionPanel.Refresh()
		}
		if ti == a.tabs.Items[5] && a.participantPanel != nil {
			// Refresh participant panel data when switching to participants tab
			a.participantPanel.Refresh()
		}
	}

	return a.tabs
}

// createMonitoringTab creates the Live Monitoring tab
func (a *App) createMonitoringTab() fyne.CanvasObject {
	a.monitoringPanel = NewMonitoringPanel(a.competitionService, a.mainWindow)
	return a.monitoringPanel.content
}

// createCompetitorsTab creates the Competitors management tab
func (a *App) createCompetitorsTab() fyne.CanvasObject {
	a.competitorPanel = NewCompetitorPanel(a.competitorService, a.mainWindow)
	return a.competitorPanel.content
}

// createModelsTab creates the Models management tab
func (a *App) createModelsTab() fyne.CanvasObject {
	a.modelPanel = NewModelPanel(a.modelService, a.mainWindow)
	return a.modelPanel.content
}

// createTranspondersTab creates the Transponders management tab
func (a *App) createTranspondersTab() fyne.CanvasObject {
	a.competitorModelPanel = NewCompetitorModelPanel(a.competitorModelService, a.competitorService, a.modelService, a.mainWindow)
	return a.competitorModelPanel.content
}

// createCompetitionsTab creates the Competitions management tab
func (a *App) createCompetitionsTab() fyne.CanvasObject {
	a.competitionPanel = NewCompetitionPanel(a.competitionService, a.mainWindow)
	return a.competitionPanel.content
}

// createParticipantsTab creates the Participants binding tab
func (a *App) createParticipantsTab() fyne.CanvasObject {
	a.participantPanel = NewParticipantPanel(a.competitionService, a.participantService, a.competitorModelService, a.modelService, a.mainWindow)
	return a.participantPanel.content
}

// createLogsTab creates the Logs viewing tab
func (a *App) createLogsTab() fyne.CanvasObject {
	a.logsPanel = NewLogsPanel(a.mainWindow)
	
	// Initialize port scanner with lap service
	a.portScanner = NewPortScanner(a.logsPanel, a.lapService)
	
	// Create container with logs and port scanner
	return container.NewBorder(
		a.portScanner.BuildUI(),
		nil,
		nil,
		nil,
		a.logsPanel.content,
	)
}

// createSettingsTab creates the Settings tab
func (a *App) createSettingsTab() fyne.CanvasObject {
	a.settingsPanel = NewSettingsPanel(a.settingsService, a.config, a.mainWindow)
	return a.settingsPanel.content
}

// refreshUI updates all UI elements with new locale strings
func (a *App) refreshUI() {
	// Refresh tab titles
	for i, tab := range a.tabs.Items {
		switch i {
		case 0:
			tab.Text = locale.T("tab.monitoring")
		case 1:
			tab.Text = locale.T("tab.competitors")
		case 2:
			tab.Text = locale.T("tab.models")
		case 3:
			tab.Text = locale.T("tab.transponders")
		case 4:
			tab.Text = locale.T("tab.competitions")
		case 5:
			tab.Text = locale.T("tab.participants")
		case 6:
			tab.Text = locale.T("tab.logs")
		case 7:
			tab.Text = locale.T("tab.settings")
		}
	}

	a.tabs.Refresh()

	// Refresh panels only if they have been created
	if a.monitoringPanel != nil {
		a.monitoringPanel.Refresh()
	}
	if a.competitorPanel != nil {
		a.competitorPanel.Refresh()
	}
	if a.modelPanel != nil {
		a.modelPanel.Refresh()
	}
	if a.competitorModelPanel != nil {
		a.competitorModelPanel.Refresh()
	}
	if a.competitionPanel != nil {
		a.competitionPanel.Refresh()
	}
	if a.participantPanel != nil {
		a.participantPanel.Refresh()
	}
	if a.settingsPanel != nil {
		a.settingsPanel.Refresh()
	}
	if a.logsPanel != nil {
		a.logsPanel.Refresh()
	}
	// Also update settings tab content
	a.tabs.Refresh()
}
