package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewAuditTab creates the audit log tab
func NewAuditTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Аудит (Audit Log)")
	content := container.NewVBox(
		title,
		widget.NewLabel("View audit log entries"),
		// Additional widgets for audit log display would go here
	)

	return container.NewTabItem("Audit Log", content)
}