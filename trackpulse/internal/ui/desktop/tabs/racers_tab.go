package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateRacersTab(db *sql.DB) *fyne.Container {
	// Create data binding for racers
	racersData := binding.BindStringList(
		&[]string{}, // Will be populated with racer names
	)

	// Create table for displaying racers
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllRacers(db))
			return count, 3 // Rows: count, Columns: 3 (Number, Full Name, Rating)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			racers := getAllRacers(db)
			if id.Row < len(racers) {
				racer := racers[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(racer.RacerNumber))
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(racer.FullName)
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(racer.Rating))
				}
			}
		},
	)
	table.SetColumnWidth(0, 80)  // Number column
	table.SetColumnWidth(1, 200) // Name column
	table.SetColumnWidth(2, 80)  // Rating column

	// Create form for adding/editing racers
	numberEntry := widget.NewEntry()
	numberEntry.PlaceHolder = "Number"
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Full Name"
	ratingEntry := widget.NewEntry()
	ratingEntry.PlaceHolder = "Rating"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Number", Widget: numberEntry},
			{Text: "Full Name", Widget: nameEntry},
			{Text: "Rating", Widget: ratingEntry},
		},
		OnSubmit: func() {
			// Add new racer to database
			racer := &models.Racer{
				// ID will be generated
				RacerNumber: parseInt(numberEntry.Text),
				FullName:    nameEntry.Text,
				Rating:      parseInt(ratingEntry.Text),
				// Other fields will be set by the model
			}
			
			// Call racer.Create(db) method
			// TODO: Implement proper error handling
		},
	}

	topContainer := container.NewBorder(nil, form, nil, nil, table)
	return topContainer
}

// Helper function to get all racers from DB
func getAllRacers(db *sql.DB) []models.Racer {
	racer := &models.Racer{}
	racers, err := racer.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.Racer{}
	}
	return racers
}

// Helper function to convert string to int
func parseInt(s string) int {
	// TODO: Implement proper parsing with error handling
	var i int
	// Using fmt.Sscanf or strconv.Atoi
	return i
}