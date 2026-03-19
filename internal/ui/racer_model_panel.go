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

// RacerModelPanel represents the Transponders management panel
type RacerModelPanel struct {
	racerModelService *service.RacerModelService
	racerService      *service.RacerService
	modelService      *service.RCModelService
	content           *fyne.Container
	table             *widget.Table
	statusLabel       *widget.Label
	window            fyne.Window          // Reference to window for dialogs
	selectedID        string               // ID of selected racer model
	allRacerModels    []models.RacerModel  // Cache of all racer models
	allRacers         []models.Racer       // Cache of all racers
	allModels         []models.RCModel     // Cache of all RC models
	headers           []string             // Localized table headers
}

// updateLocale updates all localized text in the panel
func (p *RacerModelPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("racer.header.name"),
		locale.T("model.header.name"),
		locale.T("transponder.header.number"),
		locale.T("transponder.header.type"),
		locale.T("transponder.header.active"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *RacerModelPanel) Refresh() {
	p.updateLocale()
}

// NewRacerModelPanel creates a new transponder management panel
func NewRacerModelPanel(racerModelService *service.RacerModelService, racerService *service.RacerService, modelService *service.RCModelService, window fyne.Window) *RacerModelPanel {
	panel := &RacerModelPanel{
		racerModelService: racerModelService,
		racerService:      racerService,
		modelService:      modelService,
		window:            window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the transponder panel UI
func (p *RacerModelPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying transponders
	p.table = p.createRacerModelTable()

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
func (p *RacerModelPanel) createToolbar() *widget.Toolbar {
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

// createRacerModelTable creates the data table for transponders
func (p *RacerModelPanel) createRacerModelTable() *widget.Table {
	// First load data
	p.refreshData()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allRacerModels) == 0 {
				return 0, 0
			}
			return len(p.allRacerModels), 8 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allRacerModels) {
				o.(*widget.Label).SetText("")
				return
			}
			rm := p.allRacerModels[i.Row]

			// Find racer name
			racerName := "-"
			for _, r := range p.allRacers {
				if r.ID == rm.RacerID {
					racerName = r.FullName
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
				o.(*widget.Label).SetText(racerName)
			case 2:
				o.(*widget.Label).SetText(modelName)
			case 3:
				o.(*widget.Label).SetText(rm.TransponderNumber)
			case 4:
				o.(*widget.Label).SetText(rm.TransponderType)
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
	table.SetColumnWidth(3, 120) // Transponder Number
	table.SetColumnWidth(4, 100) // Transponder Type
	table.SetColumnWidth(5, 80)  // Active
	table.SetColumnWidth(6, 150) // Created At
	table.SetColumnWidth(7, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allRacerModels) {
			p.selectedID = p.allRacerModels[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allRacerModels[id.Row].TransponderNumber))
		}
	}

	return table
}

// refreshData reloads the racer model data
func (p *RacerModelPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allRacerModels, err = p.racerModelService.GetAllRacerModels()
		if err != nil {
			fmt.Println("ERROR refreshing racer models:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}

		// Load racers
		p.allRacers, err = p.racerService.GetAllRacers()
		if err != nil {
			fmt.Println("ERROR refreshing racers:", err)
		}

		// Load models
		p.allModels, err = p.modelService.GetAllModels()
		if err != nil {
			fmt.Println("ERROR refreshing models:", err)
		}

		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allRacerModels) == 0 {
			p.statusLabel.SetText("No transponders found")
		} else {
			p.statusLabel.SetText(fmt.Sprintf("Loaded %d transponders", len(p.allRacerModels)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new transponder
func (p *RacerModelPanel) showCreateDialog() {
	p.showRacerModelDialog("Create New Transponder", nil)
}

// showEditDialog shows the dialog for editing an existing transponder
func (p *RacerModelPanel) showEditDialog() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected racer model in cache
	for _, rm := range p.allRacerModels {
		if rm.ID == p.selectedID {
			p.showRacerModelDialog(locale.T("dialog.edit.title"), &rm)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected transponder
func (p *RacerModelPanel) deleteSelected() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected racer model in cache
	var rmToDelete *models.RacerModel
	for i, rm := range p.allRacerModels {
		if rm.ID == p.selectedID {
			rmToDelete = &p.allRacerModels[i]
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
		fmt.Sprintf(locale.T("dialog.delete.message"), rmToDelete.TransponderNumber),
		func(confirmed bool) {
			if confirmed {
				if err := p.racerModelService.DeleteRacerModel(rmToDelete.ID); err != nil {
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

// showRacerModelDialog shows a dialog for creating or editing a transponder
func (p *RacerModelPanel) showRacerModelDialog(title string, rm *models.RacerModel) {
	// Get all racers
	allRacers, err := p.racerService.GetAllRacers()
	if err != nil {
		fmt.Println("ERROR getting racers:", err)
	}

	// Get all models
	allModels, err := p.modelService.GetAllModels()
	if err != nil {
		fmt.Println("ERROR getting models:", err)
	}

	// Build racer options for select
	racerOptions := make(map[string]string) // display -> ID
	var racerDisplayNames []string
	for _, r := range allRacers {
		display := fmt.Sprintf("%s (#%d)", r.FullName, r.RacerNumber)
		racerOptions[display] = r.ID
		racerDisplayNames = append(racerDisplayNames, display)
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
	transponderEntry := widget.NewEntry()
	transponderEntry.SetPlaceHolder(locale.T("form.transponder.number_placeholder"))

	transponderTypeEntry := widget.NewEntry()
	transponderTypeEntry.SetText("RFID")

	activeCheck := widget.NewCheck(locale.T("form.transponder.active"), nil)
	activeCheck.Checked = true

	// Create selects for racer and model
	racerSelect := widget.NewSelect(racerDisplayNames, nil)
	racerSelect.PlaceHolder = locale.T("form.transponder.select_racer")

	modelSelect := widget.NewSelect(modelDisplayNames, nil)
	modelSelect.PlaceHolder = locale.T("form.transponder.select_model")

	if rm != nil {
		// Edit mode - populate fields
		transponderEntry.SetText(rm.TransponderNumber)
		transponderTypeEntry.SetText(rm.TransponderType)
		activeCheck.Checked = rm.IsActive

		// Select racer
		for display, id := range racerOptions {
			if id == rm.RacerID {
				racerSelect.SetSelected(display)
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
		widget.NewFormItem(locale.T("form.transponder.racer"), racerSelect),
		widget.NewFormItem(locale.T("form.transponder.model"), modelSelect),
		widget.NewFormItem(locale.T("form.transponder.number"), transponderEntry),
		widget.NewFormItem(locale.T("form.transponder.type"), transponderTypeEntry),
		widget.NewFormItem(locale.T("form.transponder.active"), activeCheck),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Validate racer selection
		if racerSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.racer")), p.window)
			return
		}

		// Validate model selection
		if modelSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.model")), p.window)
			return
		}

		// Validate transponder number
		transponderNumber := transponderEntry.Text
		if transponderNumber == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.transponder_number")), p.window)
			return
		}

		racerID := racerOptions[racerSelect.Selected]
		modelID := modelOptions[modelSelect.Selected]

		var newRM *models.RacerModel
		if rm != nil {
			// Update existing
			newRM = rm
			newRM.RacerID = racerID
			newRM.RCModelID = modelID
			newRM.TransponderNumber = transponderNumber
			newRM.TransponderType = transponderTypeEntry.Text
			newRM.IsActive = activeCheck.Checked

			if err := p.racerModelService.UpdateRacerModel(newRM); err != nil {
				fmt.Println("ERROR updating transponder:", err)
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
			newRM = &models.RacerModel{
				RacerID:           racerID,
				RCModelID:         modelID,
				TransponderNumber: transponderNumber,
				TransponderType:   transponderTypeEntry.Text,
				IsActive:          activeCheck.Checked,
			}

			if err := p.racerModelService.CreateRacerModel(newRM); err != nil {
				fmt.Println("ERROR creating transponder:", err)
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
