package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/service"
)

// App represents the main application UI
type App struct {
	fyneApp         fyne.App
	mainWindow      fyne.Window
	racerService    *service.RacerService
	modelService    *service.RCModelService
	settingsService *service.SettingsService
	racerModelService *service.RacerModelService
	config          *Config
	tabs            *container.AppTabs
	racerPanel      *RacerPanel
	modelPanel      *ModelPanel
	racerModelPanel *RacerModelPanel
}

// Config holds UI configuration
type Config struct {
	Language string
	Title    string
}

// NewApp creates a new TrackPulse application
func NewApp(racerService *service.RacerService, modelService *service.RCModelService, settingsService *service.SettingsService, racerModelService *service.RacerModelService, language string) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("TrackPulse")

	return &App{
		fyneApp:           fyneApp,
		mainWindow:        mainWindow,
		racerService:      racerService,
		modelService:      modelService,
		settingsService:   settingsService,
		racerModelService: racerModelService,
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
		container.NewTabItem(locale.T("tab.racers"), a.createRacersTab()),
		container.NewTabItem(locale.T("tab.models"), a.createModelsTab()),
		container.NewTabItem(locale.T("tab.transponders"), a.createTranspondersTab()),
		container.NewTabItem(locale.T("tab.races"), a.createRacesTab()),
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

// createRacersTab creates the Racers management tab
func (a *App) createRacersTab() fyne.CanvasObject {
	a.racerPanel = NewRacerPanel(a.racerService, a.mainWindow)
	return a.racerPanel.content
}

// createModelsTab creates the Models management tab
func (a *App) createModelsTab() fyne.CanvasObject {
	a.modelPanel = NewModelPanel(a.modelService, a.mainWindow)
	return a.modelPanel.content
}

// createTranspondersTab creates the Transponders management tab
func (a *App) createTranspondersTab() fyne.CanvasObject {
	a.racerModelPanel = NewRacerModelPanel(a.racerModelService, a.racerService, a.modelService, a.mainWindow)
	return a.racerModelPanel.content
}

// createRacesTab creates the Races management tab
func (a *App) createRacesTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("tab.races"))
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
	// Create language selector
	languageLabel := widget.NewLabel(locale.T("settings.language"))

	// Build options for language select
	options := make([]string, 0, len(locale.SupportedLocales))
	for _, name := range locale.SupportedLocales {
		options = append(options, name)
	}

	// Find current language name
	currentName := "English"
	for code, name := range locale.SupportedLocales {
		if code == a.config.Language {
			currentName = name
			break
		}
	}

	// Create select without callback first
	languageSelect := widget.NewSelect(options, nil)
	
	// Set initial value without triggering callback
	languageSelect.SetSelected(currentName)
	
	// Now set the callback for future changes
	languageSelect.OnChanged = func(selected string) {
		// Find the code for the selected language
		var selectedCode string
		for code, name := range locale.SupportedLocales {
			if name == selected {
				selectedCode = code
				break
			}
		}

		if selectedCode != "" {
			locale.SetLocale(selectedCode)
			a.config.Language = selectedCode
			
			// Save to database
			if a.settingsService != nil {
				err := a.settingsService.SetLocale(selectedCode)
				if err != nil {
					// Log error but continue with UI update
					fmt.Printf("Failed to save locale: %v\n", err)
				}
			}
			
			a.refreshUI()
		}
	}

	// Create settings form
	form := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(languageLabel, languageSelect),
		widget.NewSeparator(),
	)

	return container.NewPadded(form)
}

// refreshUI updates all UI elements with new locale strings
func (a *App) refreshUI() {
	// Refresh tab titles
	for i, tab := range a.tabs.Items {
		switch i {
		case 0:
			tab.Text = locale.T("tab.monitoring")
		case 1:
			tab.Text = locale.T("tab.racers")
		case 2:
			tab.Text = locale.T("tab.models")
		case 3:
			tab.Text = locale.T("tab.transponders")
		case 4:
			tab.Text = locale.T("tab.races")
		case 5:
			tab.Text = locale.T("tab.logs")
		case 6:
			tab.Text = locale.T("tab.settings")
		}
	}

	a.tabs.Refresh()

	// Refresh panels only if they have been created
	if a.racerPanel != nil {
		a.racerPanel.Refresh()
	}
	if a.modelPanel != nil {
		a.modelPanel.Refresh()
	}
	if a.racerModelPanel != nil {
		a.racerModelPanel.Refresh()
	}
	
	// Also update settings tab content
	a.tabs.Refresh()
}
