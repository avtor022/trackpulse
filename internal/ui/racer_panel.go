package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// RacerPanel represents the Racers management panel
type RacerPanel struct {
	racerService    *service.RacerService
	content         *fyne.Container
	table           *widget.Table
	statusLabel     *widget.Label
	window          fyne.Window    // Reference to window for dialogs
	selectedRacerID string         // ID of selected racer
	allRacers       []models.Racer // Cache of all racers
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
	// First load data
	p.allRacers, _ = p.racerService.GetAllRacers()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allRacers) == 0 {
				return 0, 0
			}
			return len(p.allRacers), 9 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allRacers) {
				o.(*widget.Label).SetText("")
				return
			}
			racer := p.allRacers[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(racer.ID)
			case 1:
				o.(*widget.Label).SetText(strconv.Itoa(racer.RacerNumber))
			case 2:
				o.(*widget.Label).SetText(racer.FullName)
			case 3:
				o.(*widget.Label).SetText(racer.Country)
			case 4:
				o.(*widget.Label).SetText(racer.City)
			case 5:
				if racer.Birthday != nil {
					o.(*widget.Label).SetText(racer.Birthday.Format("02.01.2006"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6:
				o.(*widget.Label).SetText(strconv.Itoa(racer.Rating))
			case 7:
				if !racer.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(racer.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 8:
				if !racer.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(racer.UpdatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
			// Ensure text truncation with ellipsis to prevent overflow
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		},
	)

	// Create headers
	headerLabels := []string{
		locale.T("common.id"),
		locale.T("racer.header.number"),
		locale.T("racer.header.name"),
		locale.T("racer.header.country"),
		locale.T("racer.header.city"),
		locale.T("racer.header.birthday"),
		locale.T("racer.header.rating"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	table.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("Header")
		label.Truncation = fyne.TextTruncateEllipsis
		return label
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if id.Col >= 0 && id.Col < len(headerLabels) {
			o.(*widget.Label).SetText(headerLabels[id.Col])
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		}
	}

	// Enable header row display
	table.ShowHeaderRow = true

	// Set column widths for better visibility
	table.SetColumnWidth(0, 280) // ID
	table.SetColumnWidth(1, 80)  // Racer Number
	table.SetColumnWidth(2, 250) // Full Name
	table.SetColumnWidth(3, 120) // Country
	table.SetColumnWidth(4, 120) // City
	table.SetColumnWidth(5, 140) // Birthday
	table.SetColumnWidth(6, 80)  // Rating
	table.SetColumnWidth(7, 150) // Created At
	table.SetColumnWidth(8, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allRacers) {
			p.selectedRacerID = p.allRacers[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.selected_racer"), p.allRacers[id.Row].FullName))
		}
	}

	return table
}

// refreshData reloads the racer data
func (p *RacerPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allRacers, err = p.racerService.GetAllRacers()
		if err != nil {
			fmt.Println("ERROR refreshing data:", err)
			p.statusLabel.SetText(locale.T("status.refresh_error"))
			return
		}
		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allRacers) == 0 {
			p.statusLabel.SetText(locale.T("status.no_racers"))
		} else {
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_racers"), len(p.allRacers)))
		}
		fmt.Printf("DEBUG: refreshData completed, total racers: %d\n", len(p.allRacers))
		for i, r := range p.allRacers {
			fmt.Printf("DEBUG: Racer[%d]: ID=%s, Number=%d, Name=%s\n", i, r.ID, r.RacerNumber, r.FullName)
		}
	}
}

// showCreateDialog shows the dialog for creating a new racer
func (p *RacerPanel) showCreateDialog() {
	p.showRacerDialog("Create New Racer", nil)
}

// showEditDialog shows the dialog for editing an existing racer
func (p *RacerPanel) showEditDialog() {
	if p.selectedRacerID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected racer in cache
	for _, racer := range p.allRacers {
		if racer.ID == p.selectedRacerID {
			p.showRacerDialog("Edit Racer", &racer)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected racer
func (p *RacerPanel) deleteSelected() {
	if p.selectedRacerID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected racer in cache
	var racerToDelete *models.Racer
	for i, racer := range p.allRacers {
		if racer.ID == p.selectedRacerID {
			racerToDelete = &p.allRacers[i]
			break
		}
	}

	if racerToDelete == nil {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
		return
	}

	// Show confirmation dialog
	dialog.ShowConfirm(
		locale.T("dialog.delete.title"),
		fmt.Sprintf(locale.T("dialog.delete.message"), racerToDelete.FullName),
		func(confirmed bool) {
			if confirmed {
				if err := p.racerService.DeleteRacer(racerToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText(fmt.Sprintf(locale.T("status.delete_failed"), err.Error()))
				} else {
					p.refreshData()
					p.selectedRacerID = ""
					p.statusLabel.SetText(locale.T("status.deleted_success"))
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
	numberEntry.SetPlaceHolder(locale.T("form.racer.number_placeholder"))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(locale.T("form.racer.name_placeholder"))

	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder(locale.T("form.racer.country_placeholder"))

	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder(locale.T("form.racer.city_placeholder"))

	birthdayEntry := widget.NewEntry()
	birthdayEntry.SetPlaceHolder(locale.T("form.racer.birthday_placeholder"))

	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder(locale.T("form.racer.rating_placeholder"))

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

	// Create form
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.racer.number"), numberEntry),
		widget.NewFormItem(locale.T("form.racer.name"), nameEntry),
		widget.NewFormItem(locale.T("form.racer.country"), countryEntry),
		widget.NewFormItem(locale.T("form.racer.city"), cityEntry),
		widget.NewFormItem(locale.T("form.racer.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.racer.rating"), ratingEntry),
	)

	// Set minimum width for input fields via wrapper
	minWidth := float32(400)
	numberEntry.SetPlaceHolder("")
	nameEntry.SetPlaceHolder("")
	countryEntry.SetPlaceHolder("")
	cityEntry.SetPlaceHolder("")
	birthdayEntry.SetPlaceHolder("")
	ratingEntry.SetPlaceHolder("")

	// Set minimum width for each field
	numberEntry.Resize(fyne.NewSize(minWidth, numberEntry.MinSize().Height))
	nameEntry.Resize(fyne.NewSize(minWidth, nameEntry.MinSize().Height))
	countryEntry.Resize(fyne.NewSize(minWidth, countryEntry.MinSize().Height))
	cityEntry.Resize(fyne.NewSize(minWidth, cityEntry.MinSize().Height))
	birthdayEntry.Resize(fyne.NewSize(minWidth, birthdayEntry.MinSize().Height))
	ratingEntry.Resize(fyne.NewSize(minWidth, ratingEntry.MinSize().Height))

	// Create form with fields
	form = widget.NewForm(
		widget.NewFormItem(locale.T("form.racer.number"), numberEntry),
		widget.NewFormItem(locale.T("form.racer.name"), nameEntry),
		widget.NewFormItem(locale.T("form.racer.country"), countryEntry),
		widget.NewFormItem(locale.T("form.racer.city"), cityEntry),
		widget.NewFormItem(locale.T("form.racer.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.racer.rating"), ratingEntry),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("dialog.edit.save"), func() {
		// Debug: print values
		fmt.Printf("DEBUG: Number=%s, Name=%s, Country=%s, City=%s, Birthday=%s, Rating=%s\n",
			numberEntry.Text, nameEntry.Text, countryEntry.Text, cityEntry.Text, birthdayEntry.Text, ratingEntry.Text)

		// Parse values
		number, err := strconv.Atoi(strings.TrimSpace(numberEntry.Text))
		if err != nil {
			errMsg := fmt.Sprintf("invalid racer number: %v (got: '%s')", err, numberEntry.Text)
			fmt.Println("ERROR:", errMsg)
			dialog.ShowError(fmt.Errorf(errMsg), p.window)
			return
		}

		ratingStr := strings.TrimSpace(ratingEntry.Text)
		rating := 0
		if ratingStr != "" {
			rating, err = strconv.Atoi(ratingStr)
			if err != nil {
				rating = 0
			}
		}

		var r *models.Racer
		if racer != nil {
			// Update existing
			r = racer
			r.RacerNumber = number
			r.FullName = strings.TrimSpace(nameEntry.Text)
			r.Country = strings.TrimSpace(countryEntry.Text)
			r.City = strings.TrimSpace(cityEntry.Text)
			if birthdayEntry.Text != "" {
				birthdayStr := strings.TrimSpace(birthdayEntry.Text)
				birthday, err := time.Parse("02.01.2006", birthdayStr)
				if err == nil {
					r.Birthday = &birthday
				} else {
					errMsg := fmt.Sprintf("invalid date format (use DD.MM.YYYY): %v (got: '%s')", err, birthdayStr)
					fmt.Println("ERROR:", errMsg)
					dialog.ShowError(fmt.Errorf(errMsg), p.window)
					return
				}
			} else {
				r.Birthday = nil
			}
			r.Rating = rating
			if err := p.racerService.UpdateRacer(r); err != nil {
				fmt.Println("ERROR updating racer:", err)
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
			r = &models.Racer{
				RacerNumber: number,
				FullName:    strings.TrimSpace(nameEntry.Text),
				Country:     strings.TrimSpace(countryEntry.Text),
				City:        strings.TrimSpace(cityEntry.Text),
				Rating:      rating,
			}
			if birthdayEntry.Text != "" {
				birthdayStr := strings.TrimSpace(birthdayEntry.Text)
				birthday, err := time.Parse("02.01.2006", birthdayStr)
				if err == nil {
					r.Birthday = &birthday
				} else {
					errMsg := fmt.Sprintf("invalid date format (use DD.MM.YYYY): %v (got: '%s')", err, birthdayStr)
					fmt.Println("ERROR:", errMsg)
					dialog.ShowError(fmt.Errorf(errMsg), p.window)
					return
				}
			}
			fmt.Printf("DEBUG: Creating racer: %+v\n", r)
			if err := p.racerService.CreateRacer(r); err != nil {
				fmt.Println("ERROR creating racer:", err)
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
	cancelBtn := widget.NewButton(locale.T("dialog.edit.cancel"), func() {
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
