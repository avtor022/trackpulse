package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewRacerModelsTab creates the racer models management tab
func NewRacerModelsTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Связки (Racer Models)")
	content := container.NewVBox(
		title,
		widget.NewLabel("CRUD operations for linking racers to their models"),
		// Additional widgets for racer models management would go here
	)

	return container.NewTabItem("Racer Models", content)
}