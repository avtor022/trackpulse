package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewRacersTab creates the racers management tab
func NewRacersTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Гонщики (Racers)")
	content := container.NewVBox(
		title,
		widget.NewLabel("CRUD operations for racers"),
		// Additional widgets for racers management would go here
	)

	return container.NewTabItem("Racers", content)
}