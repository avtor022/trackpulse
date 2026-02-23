package desktop

import (
	"database/sql"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"trackpulse/internal/race"
	"trackpulse/internal/websocket"
	"trackpulse/ui/desktop/tabs"
)

// StartUI initializes and starts the desktop user interface
func StartUI(db *sql.DB, raceController *race.Controller, wsServer *websocket.Hub) {
	myApp := app.New()
	myWindow := myApp.NewWindow("TrackPulse - RC Race Timing System")
	myWindow.Resize(fyne.NewSize(1024, 768))

	// Create tab container
	racersTab := tabs.NewRacersTab(db)
	modelsTab := tabs.NewModelsTab(db)
	racerModelsTab := tabs.NewRacerModelsTab(db)
	racesTab := tabs.NewRacesTab(db)
	participantsTab := tabs.NewParticipantsTab(db)
	trackingTab := tabs.NewTrackingTab(db, raceController, wsServer)
	logsTab := tabs.NewLogsTab(db)
	auditTab := tabs.NewAuditTab(db)
	settingsTab := tabs.NewSettingsTab(db)

	tabContainer := container.NewAppTabs(
		container.NewTabItem("Гонщики (Racers)", racersTab),
		container.NewTabItem("Модели (RC Models)", modelsTab),
		container.NewTabItem("Связки (Racer Models)", racerModelsTab),
		container.NewTabItem("Заезды (Races)", racesTab),
		container.NewTabItem("Участники (Participants)", participantsTab),
		container.NewTabItem("Трекинг заезда (Race Tracking)", trackingTab),
		container.NewTabItem("Логи сканирований (Raw Scans)", logsTab),
		container.NewTabItem("Аудит (Audit Log)", auditTab),
		container.NewTabItem("Настройки (Settings)", settingsTab),
	)

	// Set initial tab
	tabContainer.SelectIndex(0)

	myWindow.SetContent(tabContainer)
	myWindow.ShowAndRun()
}