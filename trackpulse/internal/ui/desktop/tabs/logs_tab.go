package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateRawScansTab(db *sql.DB) *fyne.Container {
	// Create table for displaying raw scans
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllRawScans(db))
			return count, 5 // Rows: count, Columns: 5 (Timestamp, Tag Value, Reader Type, COM Port, Processed)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			scans := getAllRawScans(db)
			if id.Row < len(scans) {
				scan := scans[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(scan.Timestamp)
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(scan.TagValue)
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(scan.ReaderType)
				case 3:
					if scan.COMPort != nil {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(*scan.COMPort)
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
					}
				case 4:
					if scan.IsProcessed {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("Yes")
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("No")
					}
				}
			}
		},
	)
	table.SetColumnWidth(0, 150) // Timestamp column
	table.SetColumnWidth(1, 150) // Tag Value column
	table.SetColumnWidth(2, 100) // Reader Type column
	table.SetColumnWidth(3, 100) // COM Port column
	table.SetColumnWidth(4, 80)  // Processed column

	// Create form for actions
	cleanOldLogsButton := widget.NewButton("Очистить логи старше 1 года", func() {
		// Clean old logs
		cleanOldLogs(db)
	})
	
	cleanAllLogsButton := widget.NewButton("Полная очистка логов", func() {
		// Clean all logs
		cleanAllLogs(db)
	})

	actionForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: cleanOldLogsButton},
			{Text: "", Widget: cleanAllLogsButton},
		},
	}

	topContainer := container.NewBorder(nil, actionForm, nil, nil, table)
	return topContainer
}

// Helper function to get all raw scans from DB
func getAllRawScans(db *sql.DB) []models.RawScan {
	scan := &models.RawScan{}
	scans, err := scan.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RawScan{}
	}
	return scans
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