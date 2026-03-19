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

// AthletePanel represents the Athletes management panel
type AthletePanel struct {
	athleteService  *service.AthleteService
	content         *fyne.Container
	table           *widget.Table
	statusLabel     *widget.Label
	window          fyne.Window    // Reference to window for dialogs
	selectedAthleteID string       // ID of selected athlete
	allAthletes     []models.Athlete // Cache of all athletes
	// UI components that need to be updated on language change
	toolbar   *widget.Toolbar
	headers   []string
	formItems []*widget.FormItem
}

// updateLocale updates all localized text in the panel
func (p *AthletePanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("athlete.header.number"),
		locale.T("athlete.header.name"),
		locale.T("athlete.header.country"),
		locale.T("athlete.header.city"),
		locale.T("athlete.header.birthday"),
		locale.T("athlete.header.rating"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *AthletePanel) Refresh() {
	p.updateLocale()
}

// NewAthletePanel creates a new athlete management panel
func NewAthletePanel(athleteService *service.AthleteService, window fyne.Window) *AthletePanel {
	panel := &AthletePanel{
		athleteService: athleteService,
		window:       window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the athlete panel UI
func (p *AthletePanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	p.toolbar = p.createToolbar()

	// Table for displaying athletes
	p.table = p.createAthleteTable()

	// Initialize headers
	p.updateLocale()

	// Layout
	content := container.NewBorder(
		container.NewHBox(p.toolbar, p.statusLabel), // Top
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
func (p *AthletePanel) createToolbar() *widget.Toolbar {
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

// createAthleteTable creates the data table for athletes
func (p *AthletePanel) createAthleteTable() *widget.Table {
	// First load data
	p.allAthletes, _ = p.athleteService.GetAllAthletes()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allAthletes) == 0 {
				return 0, 0
			}
			return len(p.allAthletes), 9 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allAthletes) {
				o.(*widget.Label).SetText("")
				return
			}
			athlete := p.allAthletes[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(athlete.ID)
			case 1:
				o.(*widget.Label).SetText(strconv.Itoa(athlete.RacerNumber))
			case 2:
				o.(*widget.Label).SetText(athlete.FullName)
			case 3:
				o.(*widget.Label).SetText(athlete.Country)
			case 4:
				o.(*widget.Label).SetText(athlete.City)
			case 5:
				if athlete.Birthday != nil {
					o.(*widget.Label).SetText(athlete.Birthday.Format("02.01.2006"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6:
				o.(*widget.Label).SetText(strconv.Itoa(athlete.Rating))
			case 7:
				if !athlete.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(athlete.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 8:
				if !athlete.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(athlete.UpdatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
			// Ensure text truncation with ellipsis to prevent overflow
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		},
	)

	// Create headers using localized strings from p.headers
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
	table.SetColumnWidth(1, 80)  // Athlete Number
	table.SetColumnWidth(2, 250) // Full Name
	table.SetColumnWidth(3, 120) // Country
	table.SetColumnWidth(4, 120) // City
	table.SetColumnWidth(5, 140) // Birthday
	table.SetColumnWidth(6, 80)  // Rating
	table.SetColumnWidth(7, 150) // Created At
	table.SetColumnWidth(8, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allAthletes) {
			p.selectedAthleteID = p.allAthletes[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allAthletes[id.Row].FullName))
		}
	}

	return table
}

// refreshData reloads the athlete data
func (p *AthletePanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allAthletes, err = p.athleteService.GetAllAthletes()
		if err != nil {
			fmt.Println("ERROR refreshing data:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}
		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allAthletes) == 0 {
			p.statusLabel.SetText("No athletes found")
		} else {
			p.statusLabel.SetText(fmt.Sprintf("Loaded %d athletes", len(p.allAthletes)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new athlete
func (p *AthletePanel) showCreateDialog() {
	p.showAthleteDialog(locale.T("dialog.new_athlete.title"), nil)
}

// showEditDialog shows the dialog for editing an existing athlete
func (p *AthletePanel) showEditDialog() {
	if p.selectedAthleteID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected athlete in cache
	for _, athlete := range p.allAthletes {
		if athlete.ID == p.selectedAthleteID {
			p.showAthleteDialog(locale.T("dialog.edit.title"), &athlete)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected athlete
func (p *AthletePanel) deleteSelected() {
	if p.selectedAthleteID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected athlete in cache
	var athleteToDelete *models.Athlete
	for i, athlete := range p.allAthletes {
		if athlete.ID == p.selectedAthleteID {
			athleteToDelete = &p.allAthletes[i]
			break
		}
	}

	if athleteToDelete == nil {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
		return
	}

	// Show confirmation dialog
	dialog.ShowConfirm(
		locale.T("dialog.delete.title"),
		fmt.Sprintf(locale.T("dialog.delete.message"), athleteToDelete.FullName),
		func(confirmed bool) {
			if confirmed {
				if err := p.athleteService.DeleteAthlete(athleteToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText(locale.T("status.delete_failed") + ": " + err.Error())
				} else {
					p.refreshData()
					p.selectedAthleteID = ""
					p.statusLabel.SetText(locale.T("status.deleted_success"))
				}
			}
		},
		p.window,
	)
}

// showAthleteDialog shows a dialog for creating or editing an athlete
func (p *AthletePanel) showAthleteDialog(title string, athlete *models.Athlete) {
	// Create form fields with placeholders and increased width
	numberEntry := widget.NewEntry()
	numberEntry.SetPlaceHolder(locale.T("form.athlete.number_placeholder"))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(locale.T("form.athlete.name_placeholder"))

	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder(locale.T("form.athlete.country_placeholder"))

	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder(locale.T("form.athlete.city_placeholder"))

	birthdayEntry := widget.NewEntry()
	birthdayEntry.SetPlaceHolder(locale.T("form.athlete.birthday_placeholder"))

	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder(locale.T("form.athlete.rating_placeholder"))

	if athlete != nil {
		// Edit mode - populate fields
		numberEntry.SetText(strconv.Itoa(athlete.RacerNumber))
		nameEntry.SetText(athlete.FullName)
		countryEntry.SetText(athlete.Country)
		cityEntry.SetText(athlete.City)
		if athlete.Birthday != nil {
			birthdayEntry.SetText(athlete.Birthday.Format("02.01.2006"))
		}
		ratingEntry.SetText(strconv.Itoa(athlete.Rating))
	}

	// Create form with localized labels
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.athlete.number"), numberEntry),
		widget.NewFormItem(locale.T("form.athlete.name"), nameEntry),
		widget.NewFormItem(locale.T("form.athlete.country"), countryEntry),
		widget.NewFormItem(locale.T("form.athlete.city"), cityEntry),
		widget.NewFormItem(locale.T("form.athlete.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.athlete.rating"), ratingEntry),
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

	// Re-create form with fields (Fyne quirk)
	form = widget.NewForm(
		widget.NewFormItem(locale.T("form.athlete.number"), numberEntry),
		widget.NewFormItem(locale.T("form.athlete.name"), nameEntry),
		widget.NewFormItem(locale.T("form.athlete.country"), countryEntry),
		widget.NewFormItem(locale.T("form.athlete.city"), cityEntry),
		widget.NewFormItem(locale.T("form.athlete.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.athlete.rating"), ratingEntry),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Parse values
		number, err := strconv.Atoi(strings.TrimSpace(numberEntry.Text))
		if err != nil {
			errMsg := fmt.Sprintf("invalid athlete number: %v (got: '%s')", err, numberEntry.Text)
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

		var a *models.Athlete
		if athlete != nil {
			// Update existing
			a = athlete
			a.RacerNumber = number
			a.FullName = strings.TrimSpace(nameEntry.Text)
			a.Country = strings.TrimSpace(countryEntry.Text)
			a.City = strings.TrimSpace(cityEntry.Text)
			if birthdayEntry.Text != "" {
				birthdayStr := strings.TrimSpace(birthdayEntry.Text)
				birthday, err := time.Parse("02.01.2006", birthdayStr)
				if err == nil {
					a.Birthday = &birthday
				} else {
					errMsg := fmt.Sprintf("invalid date format (use DD.MM.YYYY): %v (got: '%s')", err, birthdayStr)
					fmt.Println("ERROR:", errMsg)
					dialog.ShowError(fmt.Errorf(errMsg), p.window)
					return
				}
			} else {
				a.Birthday = nil
			}
			a.Rating = rating
			if err := p.athleteService.UpdateAthlete(a); err != nil {
				fmt.Println("ERROR updating athlete:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Athlete updated successfully")

			// Close dialog and refresh data in main thread
			d.Hide()
			fyne.Do(func() {
				p.refreshData()
			})
		} else {
			// Create new
			a = &models.Athlete{
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
					a.Birthday = &birthday
				} else {
					errMsg := fmt.Sprintf("invalid date format (use DD.MM.YYYY): %v (got: '%s')", err, birthdayStr)
					fmt.Println("ERROR:", errMsg)
					dialog.ShowError(fmt.Errorf(errMsg), p.window)
					return
				}
			}
			if err := p.athleteService.CreateAthlete(a); err != nil {
				fmt.Println("ERROR creating athlete:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Athlete created successfully")

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
