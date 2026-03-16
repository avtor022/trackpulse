package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/service"
)

// App represents the main application UI
type App struct {
	fyneApp      fyne.App
	mainWindow   fyne.Window
	racerService *service.RacerService
	modelService *service.RCModelService
	config       *Config
}

// Config holds UI configuration
type Config struct {
	Language string
	Title    string
}

// NewApp creates a new TrackPulse application
func NewApp(racerService *service.RacerService, modelService *service.RCModelService, language string) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("TrackPulse")

	return &App{
		fyneApp:      fyneApp,
		mainWindow:   mainWindow,
		racerService: racerService,
		modelService: modelService,
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
	tabs := container.NewAppTabs(
		container.NewTabItem("Monitoring", a.createMonitoringTab()),
		container.NewTabItem("Racers", a.createRacersTab()),
		container.NewTabItem("Models", a.createModelsTab()),
		container.NewTabItem("Transponders", a.createTranspondersTab()),
		container.NewTabItem("Races", a.createRacesTab()),
		container.NewTabItem("Logs", a.createLogsTab()),
		container.NewTabItem("Settings", a.createSettingsTab()),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}

// createMonitoringTab creates the Live Monitoring tab
func (a *App) createMonitoringTab() fyne.CanvasObject {
	content := widget.NewLabel("Live Race Monitoring - Coming Soon")
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createRacersTab creates the Racers management tab
func (a *App) createRacersTab() fyne.CanvasObject {
	return NewRacerPanel(a.racerService, a.mainWindow)
}

// createModelsTab creates the Models management tab
func (a *App) createModelsTab() fyne.CanvasObject {
	return NewModelPanel(a.modelService, a.mainWindow)
}

// createTranspondersTab creates the Transponders management tab
func (a *App) createTranspondersTab() fyne.CanvasObject {
	content := widget.NewLabel("Transponder Assignment - Coming Soon")
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createRacesTab creates the Races management tab
func (a *App) createRacesTab() fyne.CanvasObject {
	content := widget.NewLabel("Race Management - Coming Soon")
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createLogsTab creates the Logs viewing tab
func (a *App) createLogsTab() fyne.CanvasObject {
	content := widget.NewLabel("System Logs - Coming Soon")
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createSettingsTab creates the Settings tab
func (a *App) createSettingsTab() fyne.CanvasObject {
	// Create language selector
	languageLabel := widget.NewLabel(locale.T("settings.language"))
	
	// Build options for language select
	options := make([]string, 0, len(locale.SupportedLocales))
	for code, name := range locale.SupportedLocales {
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
	
	languageSelect := widget.NewSelect(options, func(selected string) {
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
			// TODO: Refresh UI with new language
		}
	})
	
	// Set default selection
	languageSelect.SetSelected(currentName)
	
	// Create settings form
	form := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(languageLabel, languageSelect),
		widget.NewSeparator(),
	)
	
	return container.NewPadded(form)
}
