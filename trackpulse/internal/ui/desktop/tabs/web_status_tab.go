package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateWebStatusTab(db *sql.DB) *fyne.Container {
	// Create labels to display web server status information
	statusLabel := widget.NewLabel("Статус: ")
	urlLabel := widget.NewLabel("URL доступа: ")
	portLabel := widget.NewLabel("Порт: ")
	clientsLabel := widget.NewLabel("Активные клиенты: ")
	uptimeLabel := widget.NewLabel("Время работы: ")
	requestsLabel := widget.NewLabel("Запросов в минуту: ")

	// Create buttons to control the web server
	startButton := widget.NewButton("Включить", func() {
		// Start the web server
		startWebServer(db)
	})
	stopButton := widget.NewButton("Выключить", func() {
		// Stop the web server
		stopWebServer(db)
	})

	// Create status indicators
	serverStatusLabel := widget.NewLabel("🔴 Выключен")
	serverStatusLabel.Importance = widget.LowImportance

	// Layout for server controls
	controlsContainer := container.NewHBox(startButton, stopButton)

	// Layout for status information
	statusContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("Статус", "", serverStatusLabel),
			widget.NewCard("Порт", "", portLabel),
		),
		container.NewGridWithColumns(2,
			widget.NewCard("URL доступа", "", urlLabel),
			widget.NewCard("Активные клиенты", "", clientsLabel),
		),
		container.NewGridWithColumns(2,
			widget.NewCard("Время работы", "", uptimeLabel),
			widget.NewCard("Запросов в минуту", "", requestsLabel),
		),
	)

	// Main layout
	mainContainer := container.NewBorder(
		container.NewVBox(
			widget.NewCard("Управление веб-сервером", "", controlsContainer),
			widget.NewCard("Информация о подключении", "", statusContainer),
		),
		nil, nil, nil,
	)

	return mainContainer
}

// Helper function to start the web server
func startWebServer(db *sql.DB) {
	// In a real implementation, this would start the HTTP and WebSocket servers
	// For now, just update the status in the database
	
	// Update system setting to indicate web server is enabled
	setting := &models.SystemSetting{
		Key:       "web_interface_enabled",
		Value:     "true",
		ValueType: "boolean",
	}
	
	// Try to update existing setting
	err := setting.Update(db)
	if err != nil {
		// If update fails, try to create new setting
		setting.Create(db)
	}
	
	// Log the action
	logAuditEvent(db, "WEB_SERVER_START", "web_server", "", "Web server started", "")
}

// Helper function to stop the web server
func stopWebServer(db *sql.DB) {
	// In a real implementation, this would stop the HTTP and WebSocket servers
	// For now, just update the status in the database
	
	// Update system setting to indicate web server is disabled
	setting := &models.SystemSetting{
		Key:       "web_interface_enabled",
		Value:     "false",
		ValueType: "boolean",
	}
	
	// Try to update existing setting
	err := setting.Update(db)
	if err != nil {
		// If update fails, try to create new setting
		setting.Create(db)
	}
	
	// Log the action
	logAuditEvent(db, "WEB_SERVER_STOP", "web_server", "", "Web server stopped", "")
}

// Helper function to log audit events
func logAuditEvent(db *sql.DB, actionType, entityType, entityID, details, userName string) {
	auditLog := &models.AuditLog{
		// ID will be generated
		ActionType: actionType,
		EntityType: entityType,
		EntityID:   &entityID,
		Details:    &details,
		UserName:   &userName,
		// Other fields will be set by the model
	}
	
	// Call auditLog.Create(db) method
	// TODO: Implement proper error handling
}

// Additional helper functions to get web server status information
func getWebServerStatus(db *sql.DB) string {
	setting := &models.SystemSetting{}
	err := setting.GetByKey(db, "web_interface_enabled")
	if err != nil {
		return "stopped" // Default to stopped if setting doesn't exist
	}
	
	if setting.Value == "true" {
		return "running"
	}
	
	return "stopped"
}

func getWebServerPort(db *sql.DB) string {
	setting := &models.SystemSetting{}
	err := setting.GetByKey(db, "web_server_port")
	if err != nil {
		return "8080" // Default port
	}
	
	return setting.Value
}

func getActiveClientsCount(db *sql.DB) int {
	// In a real implementation, this would come from the WebSocket hub
	// For now, returning a placeholder value
	return 0
}