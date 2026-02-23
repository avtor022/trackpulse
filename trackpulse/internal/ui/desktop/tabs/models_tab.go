package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateModelsTab(db *sql.DB) *fyne.Container {
	// Create table for displaying RC models
	table := widget.NewTable(
		func() (int, int) { 
			count := len(getAllModels(db))
			return count, 4 // Rows: count, Columns: 4 (Brand, Model Name, Scale, Type)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			models := getAllModels(db)
			if id.Row < len(models) {
				model := models[id.Row]
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
				
				switch id.Col {
				case 0:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(model.Brand)
				case 1:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(model.ModelName)
				case 2:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(model.Scale)
				case 3:
					obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(model.ModelType)
				}
			}
		},
	)
	table.SetColumnWidth(0, 120) // Brand column
	table.SetColumnWidth(1, 150) // Model Name column
	table.SetColumnWidth(2, 80)  // Scale column
	table.SetColumnWidth(3, 100) // Type column

	// Create form for adding/editing models
	brandEntry := widget.NewEntry()
	brandEntry.PlaceHolder = "Brand"
	modelNameEntry := widget.NewEntry()
	modelNameEntry.PlaceHolder = "Model Name"
	scaleEntry := widget.NewEntry()
	scaleEntry.PlaceHolder = "Scale"
	typeEntry := widget.NewEntry()
	typeEntry.PlaceHolder = "Type"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Brand", Widget: brandEntry},
			{Text: "Model Name", Widget: modelNameEntry},
			{Text: "Scale", Widget: scaleEntry},
			{Text: "Type", Widget: typeEntry},
		},
		OnSubmit: func() {
			// Add new model to database
			rcModel := &models.RCModel{
				// ID will be generated
				Brand:     brandEntry.Text,
				ModelName: modelNameEntry.Text,
				Scale:     scaleEntry.Text,
				ModelType: typeEntry.Text,
				// Other fields will be set by the model
			}
			
			// Call rcModel.Create(db) method
			// TODO: Implement proper error handling
		},
	}

	topContainer := container.NewBorder(nil, form, nil, nil, table)
	return topContainer
}

// Helper function to get all models from DB
func getAllModels(db *sql.DB) []models.RCModel {
	model := &models.RCModel{}
	models, err := model.GetAll(db)
	if err != nil {
		// Handle error appropriately
		return []models.RCModel{}
	}
	return models
}