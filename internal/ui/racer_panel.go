package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"strconv"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// RacerPanel represents the Racers management panel
type RacerPanel struct {
	racerService    *service.RacerService
	content         *fyne.Container
	table           *widget.Table
	statusLabel     *widget.Label
	window          fyne.Window // Ссылка на окно для диалогов
	selectedRacerID string      // ID выбранного гонщика
	allRacers       []models.Racer // Кэш всех гонщиков
}

// NewRacerPanel creates a new racer management panel
func NewRacerPanel(racerService *service.RacerService, window fyne.Window) fyne.CanvasObject {
	panel := &RacerPanel{
		racerService: racerService,
		window:       window,
	}
	return panel.buildUI()
}

// buildUI constructs the racer panel UI
func (p *RacerPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel("Ready")

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying racers
	p.table = p.createRacerTable()

	// Layout
	content := container.NewBorder(
		container.NewHBox(toolbar, p.statusLabel), // Top
		nil, // Bottom
		nil, // Left
		nil, // Right
		p.table, // Content
	)

	p.content = content
	p.refreshData()

	return content
}

// createToolbar creates the action toolbar
func (p *RacerPanel) createToolbar() *widget.Toolbar {
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

// createRacerTable creates the data table for racers
func (p *RacerPanel) createRacerTable() *widget.Table {
	// Сначала загружаем данные
	p.allRacers, _ = p.racerService.GetAllRacers()

	table := widget.NewTableWithHeaders(
		func() (int, int) {
			if len(p.allRacers) == 0 {
				return 0, 0
			}
			return len(p.allRacers), 7 // rows, columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allRacers) {
				o.(*widget.Label).SetText("")
				return
			}
			racer := p.allRacers[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(strconv.Itoa(racer.RacerNumber))
			case 1:
				o.(*widget.Label).SetText(racer.FullName)
			case 2:
				o.(*widget.Label).SetText(racer.Country)
			case 3:
				o.(*widget.Label).SetText(racer.City)
			case 4:
				if racer.Birthday != nil {
					o.(*widget.Label).SetText(racer.Birthday.Format("02.01.2006"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 5:
				o.(*widget.Label).SetText(strconv.Itoa(racer.Rating))
			case 6:
				if !racer.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(racer.UpdatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
		},
	)

        // Set column widths for better visibility
        table.SetColumnWidth(0, 80)   // Number
        table.SetColumnWidth(1, 250)  // Full Name
        table.SetColumnWidth(2, 120)  // Country
        table.SetColumnWidth(3, 120)  // City
        table.SetColumnWidth(4, 140)  // Birthday (DD.MM.YYYY)
        table.SetColumnWidth(5, 80)   // Rating
        table.SetColumnWidth(6, 150)  // Updated

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allRacers) {
			p.selectedRacerID = p.allRacers[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allRacers[id.Row].FullName))
		}
	}

	return table
}

// refreshData reloads the racer data
func (p *RacerPanel) refreshData() {
	if p.table != nil {
		// Обновляем кэш данных
		p.allRacers, _ = p.racerService.GetAllRacers()
		p.table.Refresh()
		p.statusLabel.SetText("Data refreshed")
	}
}

// showCreateDialog shows the dialog for creating a new racer
func (p *RacerPanel) showCreateDialog() {
	p.showRacerDialog("Create New Racer", nil)
}

// showEditDialog shows the dialog for editing an existing racer
func (p *RacerPanel) showEditDialog() {
	if p.selectedRacerID == "" {
		dialog.ShowInformation("Info", "Please select a racer in the table first", p.window)
		return
	}

	// Ищем выбранного гонщика в кэше
	for _, racer := range p.allRacers {
		if racer.ID == p.selectedRacerID {
			p.showRacerDialog("Edit Racer", &racer)
			return
		}
	}

	dialog.ShowInformation("Info", "Selected racer not found", p.window)
}

// deleteSelected deletes the selected racer
func (p *RacerPanel) deleteSelected() {
	if p.selectedRacerID == "" {
		dialog.ShowInformation("Info", "Please select a racer in the table first", p.window)
		return
	}

	// Ищем выбранного гонщика в кэше
	var racerToDelete *models.Racer
	for i, racer := range p.allRacers {
		if racer.ID == p.selectedRacerID {
			racerToDelete = &p.allRacers[i]
			break
		}
	}

	if racerToDelete == nil {
		dialog.ShowInformation("Info", "Selected racer not found", p.window)
		return
	}

	// Показываем диалог подтверждения
	dialog.ShowConfirm(
		"Confirm Delete",
		"Are you sure you want to delete racer "+racerToDelete.FullName+"?",
		func(confirmed bool) {
			if confirmed {
				if err := p.racerService.DeleteRacer(racerToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText("Delete failed: " + err.Error())
				} else {
					p.refreshData()
					p.selectedRacerID = ""
					p.statusLabel.SetText("Racer deleted successfully")
				}
			}
		},
		p.window,
	)
}

// showRacerDialog shows a dialog for creating or editing a racer
func (p *RacerPanel) showRacerDialog(title string, racer *models.Racer) {
	// Create form fields with placeholders and increased width
	numberEntry := widget.NewEntry()
	numberEntry.SetPlaceHolder("Например: 7")
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Иванов Иван Иванович")
	
	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder("Россия")
	
	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder("Москва")
	
	birthdayEntry := widget.NewEntry()
	birthdayEntry.SetPlaceHolder("ДД.ММ.ГГГГ")
	
	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder("0")

	if racer != nil {
		// Edit mode - populate fields
		numberEntry.SetText(strconv.Itoa(racer.RacerNumber))
		nameEntry.SetText(racer.FullName)
		countryEntry.SetText(racer.Country)
		cityEntry.SetText(racer.City)
		if racer.Birthday != nil {
			birthdayEntry.SetText(racer.Birthday.Format("02.01.2006"))
		}
		ratingEntry.SetText(strconv.Itoa(racer.Rating))
	}
	
	// Увеличиваем ширину полей ввода для лучшей видимости
	minWidth := float32(400)
	numberEntry.Resize(fyne.NewSize(minWidth, numberEntry.MinSize().Height))
	nameEntry.Resize(fyne.NewSize(minWidth, nameEntry.MinSize().Height))
	countryEntry.Resize(fyne.NewSize(minWidth, countryEntry.MinSize().Height))
	cityEntry.Resize(fyne.NewSize(minWidth, cityEntry.MinSize().Height))
	birthdayEntry.Resize(fyne.NewSize(minWidth, birthdayEntry.MinSize().Height))
	ratingEntry.Resize(fyne.NewSize(minWidth, ratingEntry.MinSize().Height))

	// Создаем форму
	form := widget.NewForm(
		widget.NewFormItem("Racer Number", numberEntry),
		widget.NewFormItem("Full Name", nameEntry),
		widget.NewFormItem("Country", countryEntry),
		widget.NewFormItem("City", cityEntry),
		widget.NewFormItem("Birthday (DD.MM.YYYY)", birthdayEntry),
		widget.NewFormItem("Rating", ratingEntry),
	)

	// Создаем диалог с двумя кнопками: Save и Cancel
	d := dialog.NewCustomConfirm(title, "Save", "Cancel", form, func(confirmed bool) {
		if !confirmed {
			// User clicked Cancel
			p.statusLabel.SetText("Operation cancelled")
			return
		}

		// Parse values
		number, err := strconv.Atoi(numberEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid racer number"), p.window)
			return
		}

		rating, err := strconv.Atoi(ratingEntry.Text)
		if err != nil {
			rating = 0
		}

		var r *models.Racer
		if racer != nil {
			// Update existing
			r = racer
			r.RacerNumber = number
			r.FullName = nameEntry.Text
			r.Country = countryEntry.Text
			r.City = cityEntry.Text
			if birthdayEntry.Text != "" {
				birthday, err := time.Parse("02.01.2006", birthdayEntry.Text)
				if err == nil {
					r.Birthday = &birthday
				} else {
					dialog.ShowError(fmt.Errorf("invalid date format (use DD.MM.YYYY)"), p.window)
					return
				}
			} else {
				r.Birthday = nil
			}
			r.Rating = rating
			if err := p.racerService.UpdateRacer(r); err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Racer updated successfully")
		} else {
			// Create new
			r = &models.Racer{
				RacerNumber: number,
				FullName:    nameEntry.Text,
				Country:     countryEntry.Text,
				City:        cityEntry.Text,
				Rating:      rating,
			}
			if birthdayEntry.Text != "" {
				birthday, err := time.Parse("02.01.2006", birthdayEntry.Text)
				if err == nil {
					r.Birthday = &birthday
				} else {
					dialog.ShowError(fmt.Errorf("invalid date format (use DD.MM.YYYY)"), p.window)
					return
				}
			}
			if err := p.racerService.CreateRacer(r); err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Racer created successfully")
		}

		p.refreshData()
	}, p.window)

	d.Show()
}
