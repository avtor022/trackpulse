package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewSettingsTab creates the settings tab
func NewSettingsTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Настройки (Settings)")
	content := container.NewVBox(
		title,
		widget.NewLabel("Application settings and configuration"),
		// Additional widgets for settings configuration would go here
	)

	return container.NewTabItem("Settings", content)
}