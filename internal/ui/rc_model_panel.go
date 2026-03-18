package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// ModelPanel represents the RC Models management panel
type ModelPanel struct {
	modelService    *service.RCModelService
	content         *fyne.Container
	table           *widget.Table
	statusLabel     *widget.Label
	window          fyne.Window      // Reference to window for dialogs
	selectedModelID string           // ID of selected model
	allModels       []models.RCModel // Cache of all models
	headers         []string         // Localized table headers
	brandSelect     *widget.Select   // Reference to brand select widget for locale updates
}

// updateLocale updates all localized text in the panel
func (p *ModelPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("model.header.brand"),
		locale.T("model.header.name"),
		locale.T("model.header.scale"),
		locale.T("model.header.type"),
		locale.T("model.header.motor"),
		locale.T("model.header.drive"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	// Update brand select placeholder if it exists
	if p.brandSelect != nil {
		p.brandSelect.PlaceHolder = locale.T("common.select_one")
	}

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *ModelPanel) Refresh() {
	p.updateLocale()
}

// NewModelPanel creates a new RC model management panel
func NewModelPanel(modelService *service.RCModelService, window fyne.Window) *ModelPanel {
	panel := &ModelPanel{
		modelService: modelService,
		window:       window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the model panel UI
func (p *ModelPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying models
	p.table = p.createModelTable()

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
func (p *ModelPanel) createToolbar() *widget.Toolbar {
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

// createModelTable creates the data table for RC models
func (p *ModelPanel) createModelTable() *widget.Table {
	// First load data
	p.allModels, _ = p.modelService.GetAllModels()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allModels) == 0 {
				return 0, 0
			}
			return len(p.allModels), 9 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allModels) {
				o.(*widget.Label).SetText("")
				return
			}
			model := p.allModels[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(model.ID)
			case 1:
				o.(*widget.Label).SetText(model.Brand)
			case 2:
				o.(*widget.Label).SetText(model.ModelName)
			case 3:
				o.(*widget.Label).SetText(model.Scale)
			case 4:
				o.(*widget.Label).SetText(model.ModelType)
			case 5:
				if model.MotorType != "" {
					o.(*widget.Label).SetText(model.MotorType)
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6:
				if model.DriveType != "" {
					o.(*widget.Label).SetText(model.DriveType)
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 7:
				if !model.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(model.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 8:
				if !model.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(model.UpdatedAt.Format("2006-01-02 15:04:05"))
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
	table.SetColumnWidth(1, 150) // Brand
	table.SetColumnWidth(2, 200) // Model Name
	table.SetColumnWidth(3, 80)  // Scale
	table.SetColumnWidth(4, 120) // Model Type
	table.SetColumnWidth(5, 100) // Motor Type
	table.SetColumnWidth(6, 80)  // Drive Type
	table.SetColumnWidth(7, 150) // Created At
	table.SetColumnWidth(8, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allModels) {
			p.selectedModelID = p.allModels[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s %s", p.allModels[id.Row].Brand, p.allModels[id.Row].ModelName))
		}
	}

	return table
}

// refreshData reloads the model data
func (p *ModelPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allModels, err = p.modelService.GetAllModels()
		if err != nil {
			fmt.Println("ERROR refreshing data:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}
		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allModels) == 0 {
			p.statusLabel.SetText("No models found")
		} else {
			p.statusLabel.SetText(fmt.Sprintf("Loaded %d models", len(p.allModels)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new model
func (p *ModelPanel) showCreateDialog() {
	p.showModelDialog("Create New RC Model", nil)
}

// showEditDialog shows the dialog for editing an existing model
func (p *ModelPanel) showEditDialog() {
	if p.selectedModelID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected model in cache
	for _, model := range p.allModels {
		if model.ID == p.selectedModelID {
			p.showModelDialog(locale.T("dialog.edit.title"), &model)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected model
func (p *ModelPanel) deleteSelected() {
	if p.selectedModelID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected model in cache
	var modelToDelete *models.RCModel
	for i, model := range p.allModels {
		if model.ID == p.selectedModelID {
			modelToDelete = &p.allModels[i]
			break
		}
	}

	if modelToDelete == nil {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
		return
	}

	// Show confirmation dialog
	dialog.ShowConfirm(
		locale.T("dialog.delete.title"),
		fmt.Sprintf(locale.T("dialog.delete.message"), modelToDelete.Brand+" "+modelToDelete.ModelName),
		func(confirmed bool) {
			if confirmed {
				if err := p.modelService.DeleteModel(modelToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText(locale.T("status.delete_failed") + ": " + err.Error())
				} else {
					p.refreshData()
					p.selectedModelID = ""
					p.statusLabel.SetText(locale.T("status.deleted_success"))
				}
			}
		},
		p.window,
	)
}

// showModelDialog shows a dialog for creating or editing a model
func (p *ModelPanel) showModelDialog(title string, model *models.RCModel) {
	// Get all brands from the separate reference table
	allBrands, err := p.modelService.GetAllBrands()
	if err != nil {
		fmt.Println("ERROR getting brands:", err)
		// Continue working even if brands could not be retrieved
	}

	// Extract brand names
	var existingBrands []string
	for _, brand := range allBrands {
		existingBrands = append(existingBrands, brand.Name)
	}

	// Create widget for brand selection - use Select with option to add new
	var brandSelect *widget.Select

	// Add option to create new brand
	newBrandOption := "+ " + locale.T("common.add") + " " + strings.TrimSuffix(locale.T("form.model.brand"), ":")
	selectOptions := append(existingBrands, newBrandOption)

	var mainDialog dialog.Dialog

	brandSelect = widget.NewSelect(selectOptions, func(selected string) {
		if selected == newBrandOption {
			// Show dialog to add new brand
			if mainDialog != nil {
				mainDialog.Hide() // Hide main dialog
			}

			newBrandEntry := widget.NewEntry()
			newBrandEntry.SetPlaceHolder(locale.T("dialog.add_brand.placeholder"))

			// Create label and input field vertically for better width
			label := widget.NewLabel(locale.T("dialog.add_brand.label"))
			entryContainer := container.NewVBox(label, newBrandEntry)

			newBrandDialog := dialog.NewCustomWithoutButtons(locale.T("dialog.add_brand.title"), entryContainer, p.window)

			cancelBtn := widget.NewButton(locale.T("common.cancel"), func() {
				newBrandDialog.Hide()
				// Return to main dialog
				if mainDialog != nil {
					mainDialog.Show()
				}
				brandSelect.SetSelected("")
			})

			saveBtn := widget.NewButton(locale.T("common.save"), func() {
				newBrandName := strings.TrimSpace(newBrandEntry.Text)
				if newBrandName == "" {
					dialog.ShowError(fmt.Errorf(locale.T("info.enter_brand_name")), p.window)
					return
				}

				// Check if brand already exists
				for _, b := range existingBrands {
					if strings.EqualFold(b, newBrandName) {
						dialog.ShowError(fmt.Errorf(locale.T("dialog.new_brand.error_exists"), newBrandName), p.window)
						return
					}
				}

				// Add new brand to reference table
				if err := p.modelService.AddBrand(newBrandName); err != nil {
					dialog.ShowError(err, p.window)
					return
				}

				// Update brand list
				existingBrands = append(existingBrands, newBrandName)
				selectOptions = append(existingBrands, newBrandOption)
				brandSelect.Options = selectOptions

				newBrandDialog.Hide()
				// Return to main dialog with new brand selected
				if mainDialog != nil {
					mainDialog.Show()
				}
				brandSelect.SetSelected(newBrandName)
			})

			newBrandDialog.SetButtons([]fyne.CanvasObject{cancelBtn, saveBtn})

			// Increase dialog width for new brand
			parentSize := p.window.Canvas().Size()
			dialogWidth := parentSize.Width * 0.6
			if dialogWidth < 700 {
				dialogWidth = 700
			}
			// Use standard content height, do not change it
			newBrandDialog.Resize(fyne.NewSize(dialogWidth, newBrandDialog.MinSize().Height))

			newBrandDialog.Show()
		}
	})

	// Set placeholder text for the select widget (localized)
	brandSelect.PlaceHolder = locale.T("common.select_one")

	// Store reference to brandSelect for locale updates
	p.brandSelect = brandSelect

	if model != nil && model.Brand != "" {
		brandSelect.SetSelected(model.Brand)
	}

	brandWidget := brandSelect

	// Create widget for model name - simple text entry
	modelNameEntry := widget.NewEntry()
	modelNameEntry.SetPlaceHolder(locale.T("form.model.name_placeholder"))
	modelNameEntry.Resize(fyne.NewSize(250, modelNameEntry.MinSize().Height))

	if model != nil && model.ModelName != "" {
		modelNameEntry.SetText(model.ModelName)
	}

	modelNameWidget := modelNameEntry

	scaleEntry := widget.NewEntry()
	scaleEntry.SetPlaceHolder(locale.T("form.model.scale_placeholder"))

	modelTypeEntry := widget.NewEntry()
	modelTypeEntry.SetPlaceHolder(locale.T("form.model.type_placeholder"))

	motorTypeEntry := widget.NewEntry()
	motorTypeEntry.SetPlaceHolder(locale.T("form.model.motor_placeholder"))

	driveTypeEntry := widget.NewEntry()
	driveTypeEntry.SetPlaceHolder(locale.T("form.model.drive_placeholder"))

	if model != nil {
		// Edit mode - populate fields that are not select widgets
		if scaleEntry != nil {
			scaleEntry.SetText(model.Scale)
		}
		if modelTypeEntry != nil {
			modelTypeEntry.SetText(model.ModelType)
		}
		if motorTypeEntry != nil && model.MotorType != "" {
			motorTypeEntry.SetText(model.MotorType)
		}
		if driveTypeEntry != nil && model.DriveType != "" {
			driveTypeEntry.SetText(model.DriveType)
		}
	}

	// Create form with fields
	formItems := []*widget.FormItem{
		widget.NewFormItem(locale.T("form.model.brand"), brandWidget),
		widget.NewFormItem(locale.T("form.model.name"), modelNameWidget),
		widget.NewFormItem(locale.T("form.model.scale"), scaleEntry),
		widget.NewFormItem(locale.T("form.model.type"), modelTypeEntry),
		widget.NewFormItem(locale.T("form.model.motor"), motorTypeEntry),
		widget.NewFormItem(locale.T("form.model.drive"), driveTypeEntry),
	}

	form := widget.NewForm(formItems...)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Assign dialog to mainDialog for use in brand selection callback
	mainDialog = d

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Get brand value from Select
		var brand string
		if brandSelect != nil {
			brand = brandSelect.Selected
		}

		// Get model name value
		var modelName string
		if modelNameEntry != nil {
			modelName = modelNameEntry.Text
		}

		// Validate required fields
		if brand == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.brand")), p.window)
			return
		}
		if modelName == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.name")), p.window)
			return
		}
		if scaleEntry.Text == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.scale")), p.window)
			return
		}
		if modelTypeEntry.Text == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.type")), p.window)
			return
		}

		var m *models.RCModel
		if model != nil {
			// Update existing
			m = model
			m.Brand = brand
			m.ModelName = modelName
			m.Scale = scaleEntry.Text
			m.ModelType = modelTypeEntry.Text
			m.MotorType = motorTypeEntry.Text
			m.DriveType = driveTypeEntry.Text
			if err := p.modelService.UpdateModel(m); err != nil {
				fmt.Println("ERROR updating model:", err)
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
			m = &models.RCModel{
				Brand:     brand,
				ModelName: modelName,
				Scale:     scaleEntry.Text,
				ModelType: modelTypeEntry.Text,
				MotorType: motorTypeEntry.Text,
				DriveType: driveTypeEntry.Text,
			}
			if err := p.modelService.CreateModel(m); err != nil {
				fmt.Println("ERROR creating model:", err)
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
	// Get parent window size
	parentSize := p.window.Canvas().Size()

	// Calculate 50% of parent width for dialog
	dialogWidth := parentSize.Width * 0.5
	if dialogWidth < 600 {
		dialogWidth = 600 // Minimum width
	}

	// Calculate dialog height (reasonable portion of parent)
	dialogHeight := parentSize.Height * 0.7
	if dialogHeight < 500 {
		dialogHeight = 500 // Minimum height
	}

	// Resize the dialog window
	d.Resize(fyne.NewSize(dialogWidth, dialogHeight))
}
