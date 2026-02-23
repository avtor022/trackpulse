package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewModelsTab creates the RC models management tab
func NewModelsTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Модели (RC Models)")
	content := container.NewVBox(
		title,
		widget.NewLabel("CRUD operations for RC models"),
		// Additional widgets for RC models management would go here
	)

	return container.NewTabItem("RC Models", content)
}