package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// AthleteModelPanel represents the AthleteModels management panel
type AthleteModelPanel struct {
	athleteModelService *service.AthleteModelService
	athleteService      *service.AthleteService
	modelService      *service.RCModelService
	content           *fyne.Container
	table             *widget.Table
	statusLabel       *widget.Label
	window            fyne.Window          // Reference to window for dialogs
	selectedID        string               // ID of selected athlete model
	allAthleteModels    []models.AthleteModel  // Cache of all athlete models
	allAthletes         []models.Athlete       // Cache of all athletes
	allModels         []models.RCModel     // Cache of all RC models
	headers           []string             // Localized table headers
}

// updateLocale updates all localized text in the panel
func (p *AthleteModelPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("athlete.header.name"),
		locale.T("model.header.name"),
		locale.T("athletemodel.header.number"),
		locale.T("athletemodel.header.type"),
		locale.T("athletemodel.header.active"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *AthleteModelPanel) Refresh() {
	p.updateLocale()
}

// NewAthleteModelPanel creates a new AthleteModel management panel
func NewAthleteModelPanel(athleteModelService *service.AthleteModelService, athleteService *service.AthleteService, modelService *service.RCModelService, window fyne.Window) *AthleteModelPanel {
	panel := &AthleteModelPanel{
		athleteModelService: athleteModelService,
		athleteService:      athleteService,
		modelService:      modelService,
		window:            window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the AthleteModel panel UI
func (p *AthleteModelPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying AthleteModels
	p.table = p.createAthleteModelTable()

	// Layout
	content := container.NewBorder(
		container.NewHBox(toolbar, p.statusLabel), // Top
		nil,     // Bottom
		nil,     // Left
		nil,     // Right
		p.table, // Content
	)

	p.content = content
	p.refreshData()

	return content
}

// createToolbar creates the action toolbar
func (p *AthleteModelPanel) createToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			p.showCreateDialog()
		}),
		widget.NewToolbarAction(theme.ContentRedoIcon(), func() {
			p.showEditDialog()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			p.deleteSelected()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			p.refreshData()
		}),
	)
}

// createAthleteModelTable creates the data table for AthleteModels
func (p *AthleteModelPanel) createAthleteModelTable() *widget.Table {
	// First load data
	p.refreshData()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allAthleteModels) == 0 {
				return 0, 0
			}
			return len(p.allAthleteModels), 8 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allAthleteModels) {
				o.(*widget.Label).SetText("")
				return
			}
			rm := p.allAthleteModels[i.Row]

			// Find athlete name
			athleteName := "-"
			for _, r := range p.allAthletes {
				if r.ID == rm.AthleteID {
					athleteName = r.FullName
					break
				}
			}

			// Find model name
			modelName := "-"
			for _, m := range p.allModels {
				if m.ID == rm.RCModelID {
					modelName = m.Brand + " " + m.ModelName
					break
				}
			}

			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(rm.ID)
			case 1:
				o.(*widget.Label).SetText(athleteName)
			case 2:
				o.(*widget.Label).SetText(modelName)
			case 3:
				o.(*widget.Label).SetText(rm.AthleteModelNumber)
			case 4:
				o.(*widget.Label).SetText(rm.AthleteModelType)
			case 5:
				if rm.IsActive {
					o.(*widget.Label).SetText("✓")
				} else {
					o.(*widget.Label).SetText("✗")
				}
			case 6:
				if !rm.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(rm.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 7:
				if !rm.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(rm.UpdatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
			// Ensure text truncation with ellipsis to prevent overflow
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		},
	)

	// Initialize headers
	p.updateLocale()

	table.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("Header")
		label.Truncation = fyne.TextTruncateEllipsis
		return label
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(p.headers) {
			o.(*widget.Label).SetText(p.headers[id.Col])
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		}
	}

	// Enable header row display
	table.ShowHeaderRow = true

	// Set column widths for better visibility
	table.SetColumnWidth(0, 280) // ID
	table.SetColumnWidth(1, 200) // Racer Name
	table.SetColumnWidth(2, 250) // Model Name
	table.SetColumnWidth(3, 120) // AthleteModel Number
	table.SetColumnWidth(4, 100) // AthleteModel Type
	table.SetColumnWidth(5, 80)  // Active
	table.SetColumnWidth(6, 150) // Created At
	table.SetColumnWidth(7, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allAthleteModels) {
			p.selectedID = p.allAthleteModels[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allAthleteModels[id.Row].AthleteModelNumber))
		}
	}

	return table
}

// refreshData reloads the athlete model data
func (p *AthleteModelPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allAthleteModels, err = p.athleteModelService.GetAllAthleteModels()
		if err != nil {
			fmt.Println("ERROR refreshing athlete models:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}

		// Load athletes
		p.allAthletes, err = p.athleteService.GetAllAthletes()
		if err != nil {
			fmt.Println("ERROR refreshing athletes:", err)
		}

		// Load models
		p.allModels, err = p.modelService.GetAllModels()
		if err != nil {
			fmt.Println("ERROR refreshing models:", err)
		}

		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allAthleteModels) == 0 {
			p.statusLabel.SetText(locale.T("status.no_athletemodels"))
		} else {
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_athletemodels"), len(p.allAthleteModels)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new AthleteModel
func (p *AthleteModelPanel) showCreateDialog() {
	p.showRacerModelDialog(locale.T("dialog.new_athletemodel.title"), nil)
}

// showEditDialog shows the dialog for editing an existing AthleteModel
func (p *AthleteModelPanel) showEditDialog() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected athlete model in cache
	for _, rm := range p.allAthleteModels {
		if rm.ID == p.selectedID {
			p.showRacerModelDialog(locale.T("dialog.edit.title"), &rm)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected AthleteModel
func (p *AthleteModelPanel) deleteSelected() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected athlete model in cache
	var rmToDelete *models.AthleteModel
	for i, rm := range p.allAthleteModels {
		if rm.ID == p.selectedID {
			rmToDelete = &p.allAthleteModels[i]
			break
		}
	}

	if rmToDelete == nil {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
		return
	}

	// Show confirmation dialog
	dialog.ShowConfirm(
		locale.T("dialog.delete.title"),
		fmt.Sprintf(locale.T("dialog.delete.message"), rmToDelete.AthleteModelNumber),
		func(confirmed bool) {
			if confirmed {
				if err := p.athleteModelService.DeleteAthleteModel(rmToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText(locale.T("status.delete_failed") + ": " + err.Error())
				} else {
					p.refreshData()
					p.selectedID = ""
					p.statusLabel.SetText(locale.T("status.deleted_success"))
				}
			}
		},
		p.window,
	)
}

// showRacerModelDialog shows a dialog for creating or editing a AthleteModel
func (p *AthleteModelPanel) showRacerModelDialog(title string, rm *models.AthleteModel) {
	// Get all athletes
	allAthletes, err := p.athleteService.GetAllAthletes()
	if err != nil {
		fmt.Println("ERROR getting athletes:", err)
	}

	// Get all models
	allModels, err := p.modelService.GetAllModels()
	if err != nil {
		fmt.Println("ERROR getting models:", err)
	}

	// Build athlete options for select
	athleteOptions := make(map[string]string) // display -> ID
	var athleteDisplayNames []string
	for _, r := range allAthletes {
		display := fmt.Sprintf("%s (#%d)", r.FullName, r.RacerNumber)
		athleteOptions[display] = r.ID
		athleteDisplayNames = append(athleteDisplayNames, display)
	}

	// Build model options for select
	modelOptions := make(map[string]string) // display -> ID
	var modelDisplayNames []string
	for _, m := range allModels {
		display := fmt.Sprintf("%s %s (%s)", m.Brand, m.ModelName, m.Scale)
		modelOptions[display] = m.ID
		modelDisplayNames = append(modelDisplayNames, display)
	}

	// Create form fields
	AthleteModelEntry := widget.NewEntry()
	AthleteModelEntry.SetPlaceHolder(locale.T("form.AthleteModel.number_placeholder"))

	AthleteModelTypeEntry := widget.NewEntry()
	AthleteModelTypeEntry.SetText("RFID")

	activeCheck := widget.NewCheck(locale.T("form.AthleteModel.active"), nil)
	activeCheck.Checked = true

	// Create selects for athlete and model
	athleteSelect := widget.NewSelect(athleteDisplayNames, nil)
	athleteSelect.PlaceHolder = locale.T("form.AthleteModel.select_athlete")

	modelSelect := widget.NewSelect(modelDisplayNames, nil)
	modelSelect.PlaceHolder = locale.T("form.AthleteModel.select_model")

	if rm != nil {
		// Edit mode - populate fields
		AthleteModelEntry.SetText(rm.AthleteModelNumber)
		AthleteModelTypeEntry.SetText(rm.AthleteModelType)
		activeCheck.Checked = rm.IsActive

		// Select athlete
		for display, id := range athleteOptions {
			if id == rm.AthleteID {
				athleteSelect.SetSelected(display)
				break
			}
		}

		// Select model
		for display, id := range modelOptions {
			if id == rm.RCModelID {
				modelSelect.SetSelected(display)
				break
			}
		}
	}

	// Create form with localized labels
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.AthleteModel.athlete"), athleteSelect),
		widget.NewFormItem(locale.T("form.AthleteModel.model"), modelSelect),
		widget.NewFormItem(locale.T("form.AthleteModel.number"), AthleteModelEntry),
		widget.NewFormItem(locale.T("form.AthleteModel.type"), AthleteModelTypeEntry),
		widget.NewFormItem(locale.T("form.AthleteModel.active"), activeCheck),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Validate athlete selection
		if athleteSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.athlete")), p.window)
			return
		}

		// Validate model selection
		if modelSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.model")), p.window)
			return
		}

		// Validate AthleteModel number
		AthleteModelNumber := AthleteModelEntry.Text
		if AthleteModelNumber == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.AthleteModel_number")), p.window)
			return
		}

		athleteID := athleteOptions[athleteSelect.Selected]
		modelID := modelOptions[modelSelect.Selected]

		var newRM *models.AthleteModel
		if rm != nil {
			// Update existing
			newRM = rm
			newRM.AthleteID = athleteID
			newRM.RCModelID = modelID
			newRM.AthleteModelNumber = AthleteModelNumber
			newRM.AthleteModelType = AthleteModelTypeEntry.Text
			newRM.IsActive = activeCheck.Checked

			if err := p.athleteModelService.UpdateAthleteModel(newRM); err != nil {
				fmt.Println("ERROR updating AthleteModel:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText(locale.T("status.updated_success"))

			// Close dialog and refresh data in main thread
			d.Hide()
			fyne.Do(func() {
				p.refreshData()
			})
		} else {
			// Create new
			newRM = &models.AthleteModel{
				AthleteID:        athleteID,
				RCModelID:        modelID,
				AthleteModelNumber: AthleteModelNumber,
				AthleteModelType:   AthleteModelTypeEntry.Text,
				IsActive:          activeCheck.Checked,
			}

			if err := p.athleteModelService.CreateAthleteModel(newRM); err != nil {
				fmt.Println("ERROR creating AthleteModel:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText(locale.T("status.created_success"))

			// Close dialog and refresh data in main thread
			d.Hide()
			fyne.Do(func() {
				p.refreshData()
			})
		}
	})

	// Create cancel button
	cancelBtn := widget.NewButton(locale.T("common.cancel"), func() {
		p.statusLabel.SetText(locale.T("status.operation_cancelled"))
		d.Hide()
	})

	// Set dialog buttons
	d.SetButtons([]fyne.CanvasObject{cancelBtn, saveBtn})

	// Show dialog first
	d.Show()

	// Set dialog size to 50% of parent window
	parentSize := p.window.Canvas().Size()

	// Calculate 50% of parent width for dialog
	dialogWidth := parentSize.Width * 0.5
	if dialogWidth < 600 {
		dialogWidth = 600 // Minimum width
	}

	// Calculate dialog height (reasonable portion of parent)
	dialogHeight := parentSize.Height * 0.6
	if dialogHeight < 400 {
		dialogHeight = 400 // Minimum height
	}

	// Resize the dialog window
	d.Resize(fyne.NewSize(dialogWidth, dialogHeight))
}
