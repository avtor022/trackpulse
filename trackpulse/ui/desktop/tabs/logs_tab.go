package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewLogsTab creates the raw scans log tab
func NewLogsTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Логи сканирований (Raw Scans)")
	content := container.NewVBox(
		title,
		widget.NewLabel("View and manage raw RFID scan logs"),
		// Additional widgets for logs management would go here
	)

	return container.NewTabItem("Raw Scans", content)
}