package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateRacerModelsTab(db *sql.DB) *fyne.Container {
	// Create table for displaying racer-model associations
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllRacerModels(db))
			return count, 3 // Rows: count, Columns: 3 (Racer, Model, Transponder)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			racerModels := getAllRacerModels(db)
			if id.Row < len(racerModels) {
				racerModel := racerModels[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(getRacerName(db, racerModel.RacerID))
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(getModelName(db, racerModel.RCModelID))
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(racerModel.TransponderNumber)
				}
			}
		},
	)
	table.SetColumnWidth(0, 200) // Racer column
	table.SetColumnWidth(1, 200) // Model column
	table.SetColumnWidth(2, 150) // Transponder column

	// Create form for adding/editing racer-model associations
	racerSelect := widget.NewSelect([]string{}, func(s string) {})
	modelSelect := widget.NewSelect([]string{}, func(s string) {})
	transponderEntry := widget.NewEntry()
	transponderEntry.PlaceHolder = "Transponder Number"

	// Populate racer select
	populateRacerSelect(db, racerSelect)

	// Populate model select
	populateModelSelect(db, modelSelect)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Racer", Widget: racerSelect},
			{Text: "Model", Widget: modelSelect},
			{Text: "Transponder", Widget: transponderEntry},
		},
		OnSubmit: func() {
			// Add new racer-model association to database
			racerModel := &models.RacerModel{
				// ID will be generated
				RacerID:           racerSelect.Selected, // This needs to be the ID, not name
				RCModelID:         modelSelect.Selected, // This needs to be the ID, not name
				TransponderNumber: transponderEntry.Text,
				TransponderType:   "RFID",
				IsActive:          true,
				// Other fields will be set by the model
			}
			
			// Call racerModel.Create(db) method
			// TODO: Implement proper error handling
		},
	}

	topContainer := container.NewBorder(nil, form, nil, nil, table)
	return topContainer
}

// Helper function to get all racer-model associations from DB
func getAllRacerModels(db *sql.DB) []models.RacerModel {
	racerModel := &models.RacerModel{}
	racerModels, err := racerModel.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RacerModel{}
	}
	return racerModels
}

// Helper function to get racer name by ID
func getRacerName(db *sql.DB, racerID string) string {
	racer := &models.Racer{}
	err := racer.GetByID(db, racerID)
	if err != nil {
		return "Unknown Racer"
	}
	return racer.FullName
}

// Helper function to get model name by ID
func getModelName(db *sql.DB, modelID string) string {
	model := &models.RCModel{}
	err := model.GetByID(db, modelID)
	if err != nil {
		return "Unknown Model"
	}
	return model.Brand + " " + model.ModelName
}

// Helper function to populate racer select dropdown
func populateRacerSelect(db *sql.DB, selectWidget *widget.Select) {
	racers := getAllRacers(db)
	options := make([]string, len(racers))
	idMap := make(map[string]string) // Map name to ID
	
	for i, racer := range racers {
		options[i] = racer.FullName
		idMap[racer.FullName] = racer.ID
	}
	
	selectWidget.Options = options
	selectWidget.OnChanged = func(s string) {
		// Store the selected ID for later use
		// This would require extending the Select functionality or using a custom widget
	}
}

// Helper function to populate model select dropdown
func populateModelSelect(db *sql.DB, selectWidget *widget.Select) {
	models := getAllModels(db)
	options := make([]string, len(models))
	idMap := make(map[string]string) // Map name to ID
	
	for i, model := range models {
		options[i] = model.Brand + " " + model.ModelName
		idMap[model.Brand + " " + model.ModelName] = model.ID
	}
	
	selectWidget.Options = options
	selectWidget.OnChanged = func(s string) {
		// Store the selected ID for later use
		// This would require extending the Select functionality or using a custom widget
	}
}

// Helper function to get all racers from DB (from racers_tab.go)
func getAllRacers(db *sql.DB) []models.Racer {
	racer := &models.Racer{}
	racers, err := racer.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.Racer{}
	}
	return racers
}

// Helper function to get all models from DB (from models_tab.go)
func getAllModels(db *sql.DB) []models.RCModel {
	model := &models.RCModel{}
	models, err := model.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RCModel{}
	}
	return models
}