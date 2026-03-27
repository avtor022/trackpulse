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

// CompetitorModelPanel represents the Transponders management panel
type CompetitorModelPanel struct {
	competitorModelService *service.CompetitorModelService
	competitorService      *service.CompetitorService
	modelService           *service.RCModelService
	content                *fyne.Container
	table                  *widget.Table
	statusLabel            *widget.Label
	window                 fyne.Window              // Reference to window for dialogs
	selectedID             string                   // ID of selected competitor model
	allCompetitorModels    []models.CompetitorModel // Cache of all competitor models
	allCompetitors         []models.Competitor      // Cache of all competitors
	allModels              []models.RCModel         // Cache of all RC models
	headers                []string                 // Localized table headers
}

// updateLocale updates all localized text in the panel
func (p *CompetitorModelPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("competitor.header.name"),
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
func (p *CompetitorModelPanel) Refresh() {
	p.updateLocale()
}

// NewCompetitorModelPanel creates a new transponder management panel
func NewCompetitorModelPanel(competitorModelService *service.CompetitorModelService, competitorService *service.CompetitorService, modelService *service.RCModelService, window fyne.Window) *CompetitorModelPanel {
	panel := &CompetitorModelPanel{
		competitorModelService: competitorModelService,
		competitorService:      competitorService,
		modelService:           modelService,
		window:                 window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the transponder panel UI
func (p *CompetitorModelPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying transponders
	p.table = p.createCompetitorModelTable()

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
func (p *CompetitorModelPanel) createToolbar() *widget.Toolbar {
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

// createCompetitorModelTable creates the data table for transponders
func (p *CompetitorModelPanel) createCompetitorModelTable() *widget.Table {
	// First load data
	p.refreshData()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allCompetitorModels) == 0 {
				return 0, 0
			}
			return len(p.allCompetitorModels), 8 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allCompetitorModels) {
				o.(*widget.Label).SetText("")
				return
			}
			rm := p.allCompetitorModels[i.Row]

			// Find competitor name
			competitorName := "-"
			for _, c := range p.allCompetitors {
				if c.ID == rm.CompetitorID {
					competitorName = c.FullName
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
				o.(*widget.Label).SetText(competitorName)
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
	table.SetColumnWidth(1, 200) // Competitor Name
	table.SetColumnWidth(2, 250) // Model Name
	table.SetColumnWidth(3, 120) // Transponder Number
	table.SetColumnWidth(4, 100) // Transponder Type
	table.SetColumnWidth(5, 80)  // Active
	table.SetColumnWidth(6, 150) // Created At
	table.SetColumnWidth(7, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allCompetitorModels) {
			p.selectedID = p.allCompetitorModels[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allCompetitorModels[id.Row].TransponderNumber))
		}
	}

	return table
}

// refreshData reloads the competitor model data
func (p *CompetitorModelPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allCompetitorModels, err = p.competitorModelService.GetAllCompetitorModels()
		if err != nil {
			fmt.Println("ERROR refreshing competitor models:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}

		// Load competitors
		p.allCompetitors, err = p.competitorService.GetAllCompetitors()
		if err != nil {
			fmt.Println("ERROR refreshing competitors:", err)
		}

		// Load models
		p.allModels, err = p.modelService.GetAllModels()
		if err != nil {
			fmt.Println("ERROR refreshing models:", err)
		}

		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allCompetitorModels) == 0 {
			p.statusLabel.SetText(locale.T("status.no_transponders"))
		} else {
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_transponders"), len(p.allCompetitorModels)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new transponder
func (p *CompetitorModelPanel) showCreateDialog() {
	p.showCompetitorModelDialog(locale.T("dialog.new_transponder.title"), nil)
}

// showEditDialog shows the dialog for editing an existing transponder
func (p *CompetitorModelPanel) showEditDialog() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competitor model in cache
	for _, rm := range p.allCompetitorModels {
		if rm.ID == p.selectedID {
			p.showCompetitorModelDialog(locale.T("dialog.edit.title"), &rm)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected transponder
func (p *CompetitorModelPanel) deleteSelected() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competitor model in cache
	var rmToDelete *models.CompetitorModel
	for i, rm := range p.allCompetitorModels {
		if rm.ID == p.selectedID {
			rmToDelete = &p.allCompetitorModels[i]
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
				if err := p.competitorModelService.DeleteCompetitorModel(rmToDelete.ID); err != nil {
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

// showCompetitorModelDialog shows a dialog for creating or editing a transponder
func (p *CompetitorModelPanel) showCompetitorModelDialog(title string, rm *models.CompetitorModel) {
	// Get all competitors
	allCompetitors, err := p.competitorService.GetAllCompetitors()
	if err != nil {
		fmt.Println("ERROR getting competitors:", err)
	}

	// Get all models
	allModels, err := p.modelService.GetAllModels()
	if err != nil {
		fmt.Println("ERROR getting models:", err)
	}

	// Build competitor options for select
	competitorOptions := make(map[string]string) // display -> ID
	var competitorDisplayNames []string
	for _, c := range allCompetitors {
		display := fmt.Sprintf("%s (#%d)", c.FullName, c.CompetitorNumber)
		competitorOptions[display] = c.ID
		competitorDisplayNames = append(competitorDisplayNames, display)
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

	// Create selects for competitor and model
	competitorSelect := widget.NewSelect(competitorDisplayNames, nil)
	competitorSelect.PlaceHolder = locale.T("common.select_one")

	modelSelect := widget.NewSelect(modelDisplayNames, nil)
	modelSelect.PlaceHolder = locale.T("common.select_one")

	if rm != nil {
		// Edit mode - populate fields
		transponderEntry.SetText(rm.TransponderNumber)
		transponderTypeEntry.SetText(rm.TransponderType)
		activeCheck.Checked = rm.IsActive

		// Select competitor
		for display, id := range competitorOptions {
			if id == rm.CompetitorID {
				competitorSelect.SetSelected(display)
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
		widget.NewFormItem(locale.T("form.transponder.competitor"), competitorSelect),
		widget.NewFormItem(locale.T("form.transponder.model"), modelSelect),
		widget.NewFormItem(locale.T("form.transponder.number"), transponderEntry),
		widget.NewFormItem(locale.T("form.transponder.type"), transponderTypeEntry),
		widget.NewFormItem(locale.T("form.transponder.active"), activeCheck),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Validate competitor selection
		if competitorSelect.Selected == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.competitor")), p.window)
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

		competitorID := competitorOptions[competitorSelect.Selected]
		modelID := modelOptions[modelSelect.Selected]

		var newRM *models.CompetitorModel
		if rm != nil {
			// Update existing
			newRM = rm
			newRM.CompetitorID = competitorID
			newRM.RCModelID = modelID
			newRM.TransponderNumber = transponderNumber
			newRM.TransponderType = transponderTypeEntry.Text
			newRM.IsActive = activeCheck.Checked

			if err := p.competitorModelService.UpdateCompetitorModel(newRM); err != nil {
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
			newRM = &models.CompetitorModel{
				CompetitorID:      competitorID,
				RCModelID:         modelID,
				TransponderNumber: transponderNumber,
				TransponderType:   transponderTypeEntry.Text,
				IsActive:          activeCheck.Checked,
			}

			if err := p.competitorModelService.CreateCompetitorModel(newRM); err != nil {
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
