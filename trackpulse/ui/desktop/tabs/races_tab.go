package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewRacesTab creates the races management tab
func NewRacesTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Заезды (Races)")
	content := container.NewVBox(
		title,
		widget.NewLabel("CRUD operations for races"),
		// Additional widgets for races management would go here
	)

	return container.NewTabItem("Races", content)
}