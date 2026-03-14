package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// ModelPanel represents the RC Models management panel
type ModelPanel struct {
	modelService    *service.RCModelService
	content         *fyne.Container
	table           *widget.Table
	statusLabel     *widget.Label
	window          fyne.Window      // Ссылка на окно для диалогов
	selectedModelID string           // ID выбранной модели
	allModels       []models.RCModel // Кэш всех моделей
}

// NewModelPanel creates a new RC model management panel
func NewModelPanel(modelService *service.RCModelService, window fyne.Window) fyne.CanvasObject {
	panel := &ModelPanel{
		modelService: modelService,
		window:       window,
	}
	return panel.buildUI()
}

// buildUI constructs the model panel UI
func (p *ModelPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel("Ready")

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
	// Сначала загружаем данные
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

	// Создаем заголовки
	headers := []string{"ID", "Brand", "Model Name", "Scale", "Type", "Motor", "Drive", "Created At", "Updated At"}
	table.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("Header")
		label.Truncation = fyne.TextTruncateEllipsis
		return label
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(headers) {
			o.(*widget.Label).SetText(headers[id.Col])
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		}
	}

	// Включаем отображение строки заголовков
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
		// Обновляем кэш данных
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
		dialog.ShowInformation("Info", "Please select a model in the table first", p.window)
		return
	}

	// Ищем выбранную модель в кэше
	for _, model := range p.allModels {
		if model.ID == p.selectedModelID {
			p.showModelDialog("Edit RC Model", &model)
			return
		}
	}

	dialog.ShowInformation("Info", "Selected model not found", p.window)
}

// deleteSelected deletes the selected model
func (p *ModelPanel) deleteSelected() {
	if p.selectedModelID == "" {
		dialog.ShowInformation("Info", "Please select a model in the table first", p.window)
		return
	}

	// Ищем выбранную модель в кэше
	var modelToDelete *models.RCModel
	for i, model := range p.allModels {
		if model.ID == p.selectedModelID {
			modelToDelete = &p.allModels[i]
			break
		}
	}

	if modelToDelete == nil {
		dialog.ShowInformation("Info", "Selected model not found", p.window)
		return
	}

	// Показываем диалог подтверждения
	dialog.ShowConfirm(
		"Confirm Delete",
		"Are you sure you want to delete model "+modelToDelete.Brand+" "+modelToDelete.ModelName+"?",
		func(confirmed bool) {
			if confirmed {
				if err := p.modelService.DeleteModel(modelToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText("Delete failed: " + err.Error())
				} else {
					p.refreshData()
					p.selectedModelID = ""
					p.statusLabel.SetText("Model deleted successfully")
				}
			}
		},
		p.window,
	)
}

// showModelDialog shows a dialog for creating or editing a model
func (p *ModelPanel) showModelDialog(title string, model *models.RCModel) {
	// Получаем список всех брендов из отдельной таблицы справочника
	allBrands, err := p.modelService.GetAllBrands()
	if err != nil {
		fmt.Println("ERROR getting brands:", err)
		// Продолжаем работу даже если не удалось получить бренды
	}

	// Извлекаем имена брендов
	var existingBrands []string
	for _, brand := range allBrands {
		existingBrands = append(existingBrands, brand.Name)
	}

	// Получаем список всех названий моделей
	allModelNames, err := p.modelService.GetAllModelNames()
	if err != nil {
		fmt.Println("ERROR getting model names:", err)
		// Продолжаем работу даже если не удалось получить названия моделей
	}

	// Создаем виджет для выбора бренда с автодополнением
	var brandWidget fyne.CanvasObject
	var brandEntry *widget.Entry

	if len(existingBrands) > 0 {
		// Используем Entry с автодополнением
		brandEntry = widget.NewEntry()
		brandEntry.SetPlaceHolder("Например: Traxxas")
		brandEntry.Resize(fyne.NewSize(250, brandEntry.MinSize().Height))
		
		if model != nil && model.Brand != "" {
			brandEntry.SetText(model.Brand)
		}

		// Создаем контейнер с полем ввода и кнопкой dropdown
		dropdownBtn := widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
			showFullDropdown(existingBrands, p.window, brandEntry, func(selected string) {
				brandEntry.SetText(selected)
			})
		})

		brandWidget = container.NewHBox(
			container.NewStack(brandEntry),
			dropdownBtn,
		)

		// Обработчик изменения текста для фильтрации
		var popup *widget.PopUp
		
		brandEntry.OnChanged = func(text string) {
			if popup != nil {
				popup.Hide()
			}
			
			if text == "" {
				return
			}

			// Фильтруем опции
			var filtered []string
			textLower := strings.ToLower(text)
			
			for _, opt := range existingBrands {
				if strings.Contains(strings.ToLower(opt), textLower) {
					filtered = append(filtered, opt)
				}
			}

			if len(filtered) > 0 {
				showFilterPopup(brandEntry, filtered, p.window, func(selected string) {
					brandEntry.SetText(selected)
					if popup != nil {
						popup.Hide()
					}
				})
			}
		}

		brandEntry.OnSubmitted = func(text string) {
			if popup != nil {
				popup.Hide()
			}
		}

		_ = popup
	} else {
		// Если брендов нет, используем обычное поле ввода
		brandEntry = widget.NewEntry()
		brandEntry.SetPlaceHolder("Например: Traxxas")
		brandEntry.Resize(fyne.NewSize(250, brandEntry.MinSize().Height))
		if model != nil && model.Brand != "" {
			brandEntry.SetText(model.Brand)
		}
		brandWidget = brandEntry
	}

	// Создаем виджет для выбора названия модели с автодополнением
	var modelNameWidget fyne.CanvasObject
	var modelNameEntry *widget.Entry

	if len(allModelNames) > 0 {
		// Используем Entry с автодополнением
		modelNameEntry = widget.NewEntry()
		modelNameEntry.SetPlaceHolder("Например: X-Maxx")
		modelNameEntry.Resize(fyne.NewSize(250, modelNameEntry.MinSize().Height))
		
		if model != nil && model.ModelName != "" {
			modelNameEntry.SetText(model.ModelName)
		}

		// Создаем контейнер с полем ввода и кнопкой dropdown
		dropdownBtn := widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
			showFullDropdown(allModelNames, p.window, modelNameEntry, func(selected string) {
				modelNameEntry.SetText(selected)
			})
		})

		modelNameWidget = container.NewHBox(
			container.NewStack(modelNameEntry),
			dropdownBtn,
		)

		// Обработчик изменения текста для фильтрации
		var popup *widget.PopUp
		
		modelNameEntry.OnChanged = func(text string) {
			if popup != nil {
				popup.Hide()
			}
			
			if text == "" {
				return
			}

			// Фильтруем опции
			var filtered []string
			textLower := strings.ToLower(text)
			
			for _, opt := range allModelNames {
				if strings.Contains(strings.ToLower(opt), textLower) {
					filtered = append(filtered, opt)
				}
			}

			if len(filtered) > 0 {
				showFilterPopup(modelNameEntry, filtered, p.window, func(selected string) {
					modelNameEntry.SetText(selected)
					if popup != nil {
						popup.Hide()
					}
				})
			}
		}

		modelNameEntry.OnSubmitted = func(text string) {
			if popup != nil {
				popup.Hide()
			}
		}

		_ = popup
	} else {
		// Если названий моделей нет, используем обычное поле ввода
		modelNameEntry = widget.NewEntry()
		modelNameEntry.SetPlaceHolder("Например: X-Maxx")
		modelNameEntry.Resize(fyne.NewSize(250, modelNameEntry.MinSize().Height))
		if model != nil && model.ModelName != "" {
			modelNameEntry.SetText(model.ModelName)
		}
		modelNameWidget = modelNameEntry
	}

	scaleEntry := widget.NewEntry()
	scaleEntry.SetPlaceHolder("Например: 1:8")

	modelTypeEntry := widget.NewEntry()
	modelTypeEntry.SetPlaceHolder("Например: Monster Truck")

	motorTypeEntry := widget.NewEntry()
	motorTypeEntry.SetPlaceHolder("Например: Brushless")

	driveTypeEntry := widget.NewEntry()
	driveTypeEntry.SetPlaceHolder("Например: 4WD")

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

	// Создаем форму с полями
	formItems := []*widget.FormItem{
		widget.NewFormItem("Brand", brandWidget),
		widget.NewFormItem("Model Name", modelNameWidget),
		widget.NewFormItem("Scale", scaleEntry),
		widget.NewFormItem("Model Type", modelTypeEntry),
		widget.NewFormItem("Motor Type", motorTypeEntry),
		widget.NewFormItem("Drive Type", driveTypeEntry),
	}
	
	form := widget.NewForm(formItems...)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton("Save", func() {
		// Получаем значение бренда
		var brand string
		if brandEntry != nil {
			brand = brandEntry.Text
		}

		// Получаем значение названия модели
		var modelName string
		if modelNameEntry != nil {
			modelName = modelNameEntry.Text
		}

		// Validate required fields
		if brand == "" {
			dialog.ShowError(fmt.Errorf("brand is required"), p.window)
			return
		}
		if modelName == "" {
			dialog.ShowError(fmt.Errorf("model name is required"), p.window)
			return
		}
		if scaleEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("scale is required"), p.window)
			return
		}
		if modelTypeEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("model type is required"), p.window)
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
			p.statusLabel.SetText("Model updated successfully")

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
			p.statusLabel.SetText("Model created successfully")

			// Close dialog and refresh data in main thread
			d.Hide()
			fyne.Do(func() {
				p.refreshData()
			})
		}
	})

	// Create cancel button
	cancelBtn := widget.NewButton("Cancel", func() {
		p.statusLabel.SetText("Operation cancelled")
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
