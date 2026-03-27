package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"trackpulse/internal/locale"
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
	config                 *Config
	tabs                   *container.AppTabs
	competitorPanel        *CompetitorPanel
	modelPanel             *ModelPanel
	competitorModelPanel   *CompetitorModelPanel
	settingsPanel          *SettingsPanel
}

// Config holds UI configuration
type Config struct {
	Language string
	Title    string
}

// NewApp creates a new TrackPulse application
func NewApp(competitorService *service.CompetitorService, modelService *service.RCModelService, settingsService *service.SettingsService, competitorModelService *service.CompetitorModelService, language string) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("TrackPulse")

	return &App{
		fyneApp:                fyneApp,
		mainWindow:             mainWindow,
		competitorService:      competitorService,
		modelService:           modelService,
		settingsService:        settingsService,
		competitorModelService: competitorModelService,
		config: &Config{
			Language: language,
			Title:    "TrackPulse",
		},
	}
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
		container.NewTabItem(locale.T("tab.logs"), a.createLogsTab()),
		container.NewTabItem(locale.T("tab.settings"), a.createSettingsTab()),
	)

	a.tabs.SetTabLocation(container.TabLocationTop)
	return a.tabs
}

// createMonitoringTab creates the Live Monitoring tab
func (a *App) createMonitoringTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("app.welcome"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
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
	content := widget.NewLabel(locale.T("tab.competitions"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createLogsTab creates the Logs viewing tab
func (a *App) createLogsTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("tab.logs"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
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
			tab.Text = locale.T("tab.logs")
		case 6:
			tab.Text = locale.T("tab.settings")
		}
	}

	a.tabs.Refresh()

	// Refresh panels only if they have been created
	if a.competitorPanel != nil {
		a.competitorPanel.Refresh()
	}
	if a.modelPanel != nil {
		a.modelPanel.Refresh()
	}
	if a.competitorModelPanel != nil {
		a.competitorModelPanel.Refresh()
	}
	if a.settingsPanel != nil {
		a.settingsPanel.Refresh()
	}

	// Also update settings tab content
	a.tabs.Refresh()
}
