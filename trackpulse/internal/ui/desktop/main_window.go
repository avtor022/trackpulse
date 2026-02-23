package desktop

import (
	"database/sql"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/ui/desktop/tabs"
)

type App struct {
	fyneApp    fyne.App
	mainWindow fyne.Window
	db         *sql.DB
}

func NewApp(db *sql.DB) *App {
	myApp := app.New()
	mainWindow := myApp.NewWindow("TrackPulse - RC Race Timing System")
	mainWindow.Resize(fyne.NewSize(1024, 768))

	return &App{
		fyneApp:    myApp,
		mainWindow: mainWindow,
		db:         db,
	}
}

func (a *App) Run() error {
	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Гонщики", tabs.CreateRacersTab(a.db)),
		container.NewTabItem("Модели", tabs.CreateModelsTab(a.db)),
		container.NewTabItem("Связки", tabs.CreateRacerModelsTab(a.db)),
		container.NewTabItem("Заезды", tabs.CreateRacesTab(a.db)),
		container.NewTabItem("Участники", tabs.CreateParticipantsTab(a.db)),
		container.NewTabItem("Трекинг заезда", tabs.CreateRaceTrackingTab(a.db)),
		container.NewTabItem("Логи сканирований", tabs.CreateRawScansTab(a.db)),
		container.NewTabItem("Аудит", tabs.CreateAuditTab(a.db)),
		container.NewTabItem("Статус Веб", tabs.CreateWebStatusTab(a.db)),
		container.NewTabItem("Настройки", tabs.CreateSettingsTab(a.db)),
	)

	a.mainWindow.SetContent(tabs)
	a.mainWindow.ShowAndRun()

	return nil
}