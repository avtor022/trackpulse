package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateAuditTab(db *sql.DB) *fyne.Container {
	// Create table for displaying audit logs
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllAuditLogs(db))
			return count, 5 // Rows: count, Columns: 5 (Timestamp, Action Type, Entity Type, User, Details)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			logs := getAllAuditLogs(db)
			if id.Row < len(logs) {
				log := logs[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(log.Timestamp)
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(log.ActionType)
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(log.EntityType)
				case 3:
					if log.UserName != nil {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(*log.UserName)
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
					}
				case 4:
					if log.Details != nil {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(*log.Details)
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
					}
				}
			}
		},
	)
	table.SetColumnWidth(0, 150) // Timestamp column
	table.SetColumnWidth(1, 120) // Action Type column
	table.SetColumnWidth(2, 120) // Entity Type column
	table.SetColumnWidth(3, 100) // User column
	table.SetColumnWidth(4, 300) // Details column

	// Create form for actions (read-only view)
	infoLabel := widget.NewLabel("Журнал аудита (только для чтения)")

	topContainer := container.NewBorder(nil, infoLabel, nil, nil, table)
	return topContainer
}

// Helper function to get all audit logs from DB
func getAllAuditLogs(db *sql.DB) []models.AuditLog {
	log := &models.AuditLog{}
	logs, err := log.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.AuditLog{}
	}
	return logs
}