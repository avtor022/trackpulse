package ui

import (
	"trackpulse/internal/config"
	"trackpulse/internal/services"
)

// Application represents the main application
type Application struct {
	services *services.Services
	config   *config.Config
	// GUI framework specific implementation would go here
	// For example, if using Fyne: window *fyne.Window
}

// NewApplication creates a new application instance
func NewApplication(svc *services.Services, cfg *config.Config) *Application {
	return &Application{
		services: svc,
		config:   cfg,
	}
}

// Run starts the application
func (a *Application) Run() {
	// Initialize the GUI framework
	// Show the main window
	// Handle events
	// This is a placeholder implementation
	a.showMainWindow()
}

// showMainWindow shows the main application window
func (a *Application) showMainWindow() {
	// Implementation would depend on the chosen GUI framework
	// For example, with Fyne:
	/*
		app := app.New()
		window := app.NewWindow("TrackPulse")
		window.Resize(fyne.NewSize(1024, 768))

		// Check if we need to show login dialog first
		if a.config.RequireLogin {
			loginWin := windows.NewLoginWindow(a.services.AuthService, a.translator)
			if !loginWin.Show() {
				// User cancelled login, exit the app
				return
			}
		}

		// Create and show main window
		mainWin := windows.NewMainWindow(a.services, a.translator)
		mainWin.Show(window)

		window.ShowAndRun()
	*/
	
	// For now, just print that the app is running
	println("TrackPulse application started...")
}