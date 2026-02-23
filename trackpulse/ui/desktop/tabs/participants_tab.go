package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewParticipantsTab creates the race participants management tab
func NewParticipantsTab(db *sql.DB) *container.TabItem {
	title := widget.NewLabel("Участники (Participants)")
	content := container.NewVBox(
		title,
		widget.NewLabel("Management of race participants"),
		// Additional widgets for participants management would go here
	)

	return container.NewTabItem("Participants", content)
}