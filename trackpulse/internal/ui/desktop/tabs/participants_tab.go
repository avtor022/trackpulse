package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateParticipantsTab(db *sql.DB) *fyne.Container {
	// Create table for displaying race participants
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllParticipants(db))
			return count, 4 // Rows: count, Columns: 4 (Race, Racer, Model, Grid Position)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			participants := getAllParticipants(db)
			if id.Row < len(participants) {
				participant := participants[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(getRaceTitle(db, participant.RaceID))
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(getRacerNameForParticipant(db, participant.RacerModelID))
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(getModelNameForParticipant(db, participant.RacerModelID))
				case 3:
					if participant.GridPosition != nil {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(*participant.GridPosition))
					} else {
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
					}
				}
			}
		},
	)
	table.SetColumnWidth(0, 150) // Race column
	table.SetColumnWidth(1, 150) // Racer column
	table.SetColumnWidth(2, 150) // Model column
	table.SetColumnWidth(3, 100) // Grid Position column

	// Create form for adding/editing participants
	raceSelect := widget.NewSelect([]string{}, func(s string) {})
	racerModelSelect := widget.NewSelect([]string{}, func(s string) {})
	gridPositionEntry := widget.NewEntry()
	gridPositionEntry.PlaceHolder = "Grid Position"

	// Populate race select
	populateRaceSelect(db, raceSelect)

	// Populate racer-model select
	populateRacerModelSelect(db, racerModelSelect)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Race", Widget: raceSelect},
			{Text: "Racer-Model", Widget: racerModelSelect},
			{Text: "Grid Position", Widget: gridPositionEntry},
		},
		OnSubmit: func() {
			// Parse grid position
			var gridPos *int
			if gridPositionEntry.Text != "" {
				pos := parseInt(gridPositionEntry.Text)
				gridPos = &pos
			}

			// Add new participant to database
			participant := &models.RaceParticipant{
				// ID will be generated
				RaceID:         raceSelect.Selected, // This needs to be the ID, not name
				RacerModelID:   racerModelSelect.Selected, // This needs to be the ID, not name
				GridPosition:   gridPos,
				IsFinished:     false,
				Disqualified:   false,
				// Other fields will be set by the model
			}
			
			// Call participant.Create(db) method
			// TODO: Implement proper error handling
		},
	}

	topContainer := container.NewBorder(nil, form, nil, nil, table)
	return topContainer
}

// Helper function to get all participants from DB
func getAllParticipants(db *sql.DB) []models.RaceParticipant {
	participant := &models.RaceParticipant{}
	participants, err := participant.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RaceParticipant{}
	}
	return participants
}

// Helper function to get race title by ID
func getRaceTitle(db *sql.DB, raceID string) string {
	race := &models.Race{}
	err := race.GetByID(db, raceID)
	if err != nil {
		return "Unknown Race"
	}
	return race.RaceTitle
}

// Helper function to get racer name by racer-model ID
func getRacerNameForParticipant(db *sql.DB, racerModelID string) string {
	racerModel := &models.RacerModel{}
	err := racerModel.GetByID(db, racerModelID)
	if err != nil {
		return "Unknown Racer"
	}
	
	racer := &models.Racer{}
	err = racer.GetByID(db, racerModel.RacerID)
	if err != nil {
		return "Unknown Racer"
	}
	
	return racer.FullName
}

// Helper function to get model name by racer-model ID
func getModelNameForParticipant(db *sql.DB, racerModelID string) string {
	racerModel := &models.RacerModel{}
	err := racerModel.GetByID(db, racerModelID)
	if err != nil {
		return "Unknown Model"
	}
	
	model := &models.RCModel{}
	err = model.GetByID(db, racerModel.RCModelID)
	if err != nil {
		return "Unknown Model"
	}
	
	return model.Brand + " " + model.ModelName
}

// Helper function to populate race select dropdown
func populateRaceSelect(db *sql.DB, selectWidget *widget.Select) {
	races := getAllRaces(db)
	options := make([]string, len(races))
	idMap := make(map[string]string) // Map name to ID
	
	for i, race := range races {
		options[i] = race.RaceTitle
		idMap[race.RaceTitle] = race.ID
	}
	
	selectWidget.Options = options
	selectWidget.OnChanged = func(s string) {
		// Store the selected ID for later use
		// This would require extending the Select functionality or using a custom widget
	}
}

// Helper function to populate racer-model select dropdown
func populateRacerModelSelect(db *sql.DB, selectWidget *widget.Select) {
	racerModels := getAllRacerModels(db)
	options := make([]string, len(racerModels))
	idMap := make(map[string]string) // Map name to ID
	
	for i, racerModel := range racerModels {
		racerName := getRacerName(db, racerModel.RacerID)
		modelName := getModelName(db, racerModel.RCModelID)
		option := racerName + " - " + modelName
		options[i] = option
		idMap[option] = racerModel.ID
	}
	
	selectWidget.Options = options
	selectWidget.OnChanged = func(s string) {
		// Store the selected ID for later use
		// This would require extending the Select functionality or using a custom widget
	}
}

// Helper function to get all races from DB (from races_tab.go)
func getAllRaces(db *sql.DB) []models.Race {
	race := &models.Race{}
	races, err := race.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.Race{}
	}
	return races
}

// Helper function to get all racer-models from DB (from racer_models_tab.go)
func getAllRacerModels(db *sql.DB) []models.RacerModel {
	racerModel := &models.RacerModel{}
	racerModels, err := racerModel.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RacerModel{}
	}
	return racerModels
}

// Helper function to get racer name by ID (from racer_models_tab.go)
func getRacerName(db *sql.DB, racerID string) string {
	racer := &models.Racer{}
	err := racer.GetByID(db, racerID)
	if err != nil {
		return "Unknown Racer"
	}
	return racer.FullName
}

// Helper function to get model name by ID (from racer_models_tab.go)
func getModelName(db *sql.DB, modelID string) string {
	model := &models.RCModel{}
	err := model.GetByID(db, modelID)
	if err != nil {
		return "Unknown Model"
	}
	return model.Brand + " " + model.ModelName
}

// Helper function to convert string to int (from racers_tab.go)
func parseInt(s string) int {
	// TODO: Implement proper parsing with error handling
	var i int
	// Using fmt.Sscanf or strconv.Atoi
	return i
}