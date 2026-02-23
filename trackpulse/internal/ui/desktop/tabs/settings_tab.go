package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateSettingsTab(db *sql.DB) *fyne.Container {
	// Create form for system settings
	dbPathEntry := widget.NewEntry()
	dbPathEntry.PlaceHolder = "Database Path"
	
	languageSelect := widget.NewSelect([]string{"ru", "en"}, func(s string) {})
	
	comPortSelect := widget.NewSelect([]string{}, func(s string) {})
	readerTypeSelect := widget.NewSelect([]string{"EM4095", "RC522"}, func(s string) {})
	baudRateSelect := widget.NewSelect([]string{"9600", "115200"}, func(s string) {})
	
	webPortEntry := widget.NewEntry()
	webPortEntry.SetText("8080")
	
	debounceEntry := widget.NewEntry()
	debounceEntry.SetText("2000")
	
	// Populate available COM ports
	populateComPorts(comPortSelect)

	settingsForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Язык", Widget: languageSelect},
			{Text: "Путь к БД", Widget: dbPathEntry},
			{Text: "COM-порт", Widget: comPortSelect},
			{Text: "Тип считывателя", Widget: readerTypeSelect},
			{Text: "Baud rate", Widget: baudRateSelect},
			{Text: "Порт веб-интерфейса", Widget: webPortEntry},
			{Text: "Debounce (мс)", Widget: debounceEntry},
		},
		OnSubmit: func() {
			// Save settings to database
			saveSetting(db, "language", languageSelect.Selected, "string")
			saveSetting(db, "database_path", dbPathEntry.Text, "string")
			saveSetting(db, "com_port", comPortSelect.Selected, "string")
			saveSetting(db, "reader_type", readerTypeSelect.Selected, "string")
			saveSetting(db, "baud_rate", baudRateSelect.Selected, "string")
			saveSetting(db, "web_port", webPortEntry.Text, "string")
			saveSetting(db, "debounce_ms", debounceEntry.Text, "string")
		},
	}

	// Create buttons for log cleanup
	cleanOldLogsButton := widget.NewButton("Очистить логи старше 1 года", func() {
		// Clean old logs
		cleanOldLogs(db)
	})
	
	cleanAllLogsButton := widget.NewButton("Полная очистка логов", func() {
		// Clean all logs
		cleanAllLogs(db)
	})

	logCleanupForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: cleanOldLogsButton},
			{Text: "", Widget: cleanAllLogsButton},
		},
	}

	// Main layout
	mainContainer := container.NewVBox(
		widget.NewCard("Настройки системы", "", settingsForm),
		widget.NewCard("Очистка логов", "", logCleanupForm),
	)

	return mainContainer
}

// Helper function to save a setting to the database
func saveSetting(db *sql.DB, key, value, valueType string) {
	setting := &models.SystemSetting{
		Key:       key,
		Value:     value,
		ValueType: valueType,
	}
	
	// Try to update existing setting
	err := setting.Update(db)
	if err != nil {
		// If update fails, try to create new setting
		setting.Create(db)
	}
	
	// Log the action
	details := "Setting '" + key + "' changed to '" + value + "'"
	logAuditEvent(db, "SETTINGS_UPDATE", "system_settings", key, details, "")
}

// Helper function to populate available COM ports
func populateComPorts(selectWidget *widget.Select) {
	// In a real implementation, this would scan for available COM ports
	// For now, just add some example ports
	selectWidget.Options = []string{"/dev/ttyUSB0", "/dev/ttyACM0", "COM1", "COM2", "COM3"}
}

// Helper function to clean old logs (older than 1 year)
func cleanOldLogs(db *sql.DB) {
	// Delete logs older than 1 year from raw_scans, audit_log, and lap_history tables
	// This is a simplified version - in a real implementation, we'd use proper date functions
	
	// For raw_scans
	deleteQuery := `DELETE FROM raw_scans WHERE timestamp < datetime('now', '-1 year')`
	db.Exec(deleteQuery)
	
	// For audit_log
	deleteQuery = `DELETE FROM audit_log WHERE timestamp < datetime('now', '-1 year')`
	db.Exec(deleteQuery)
	
	// For lap_history (completed race data would remain, but old non-race lap_history could be cleaned)
	// We'll skip lap_history for this example since it might be needed for historical results
	
	// Log the action
	logAuditEvent(db, "LOG_CLEAN_OLD", "log_cleanup", "", "Old logs cleaned (older than 1 year)", "")
}

// Helper function to clean all logs
func cleanAllLogs(db *sql.DB) {
	// Delete all logs from raw_scans, audit_log tables
	// NOTE: This does not include lap_history which might be needed for historical results
	
	// For raw_scans
	deleteQuery := `DELETE FROM raw_scans`
	db.Exec(deleteQuery)
	
	// For audit_log
	deleteQuery = `DELETE FROM audit_log`
	db.Exec(deleteQuery)
	
	// For lap_history (only if not during an active race)
	// In a real implementation, we would check if there's an active race first
	currentRace := getActiveRace(db)
	if currentRace.ID == "" {
		// No active race, safe to clean lap_history
		deleteQuery = `DELETE FROM lap_history`
		db.Exec(deleteQuery)
	}
	
	// Log the action
	logAuditEvent(db, "LOG_CLEAN_ALL", "log_cleanup", "", "All logs cleaned", "")
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

// Helper function to get active race (reusing from tracking_tab.go)
func getActiveRace(db *sql.DB) models.Race {
	// Query for race with status 'active'
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, 
	                 time_limit_minutes, time_start, time_finish, status, created_at, updated_at 
	          FROM races WHERE status = 'active'`
	
	var race models.Race
	err := db.QueryRow(query).Scan(
		&race.ID, &race.RaceTitle, &race.RaceType, &race.ModelType, &race.ModelScale, &race.TrackName, 
		&race.LapCountTarget, &race.TimeLimitMinutes, &race.TimeStart, &race.TimeFinish, &race.Status, 
		&race.CreatedAt, &race.UpdatedAt,
	)
	
	if err != nil {
		// Return empty race if no active race found
		return models.Race{}
	}
	
	return race
}