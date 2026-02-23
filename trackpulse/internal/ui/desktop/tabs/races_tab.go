package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateRacesTab(db *sql.DB) *fyne.Container {
	// Create table for displaying races
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllRaces(db))
			return count, 4 // Rows: count, Columns: 4 (Title, Type, Status, Track)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			races := getAllRaces(db)
			if id.Row < len(races) {
				race := races[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(race.RaceTitle)
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(race.RaceType)
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(race.Status)
				case 3:
					if race.TrackName != nil {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(*race.TrackName)
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
					}
				}
			}
		},
	)
	table.SetColumnWidth(0, 200) // Title column
	table.SetColumnWidth(1, 100) // Type column
	table.SetColumnWidth(2, 100) // Status column
	table.SetColumnWidth(3, 150) // Track column

	// Create form for adding/editing races
	titleEntry := widget.NewEntry()
	titleEntry.PlaceHolder = "Race Title"
	typeSelect := widget.NewSelect([]string{"qualifying", "main", "final"}, func(s string) {})
	statusSelect := widget.NewSelect([]string{"scheduled", "active", "finished", "cancelled"}, func(s string) {})
	trackEntry := widget.NewEntry()
	trackEntry.PlaceHolder = "Track Name"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Type", Widget: typeSelect},
			{Text: "Status", Widget: statusSelect},
			{Text: "Track", Widget: trackEntry},
		},
		OnSubmit: func() {
			// Add new race to database
			race := &models.Race{
				// ID will be generated
				RaceTitle:  titleEntry.Text,
				RaceType:   typeSelect.Selected,
				Status:     statusSelect.Selected,
				TrackName:  &trackEntry.Text,
				// Other fields will be set by the model
			}
			
			// Call race.Create(db) method
			// TODO: Implement proper error handling
		},
	}

	topContainer := container.NewBorder(nil, form, nil, nil, table)
	return topContainer
}

// Helper function to get all races from DB
func getAllRaces(db *sql.DB) []models.Race {
	race := &models.Race{}
	races, err := race.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.Race{}
	}
	return races
}