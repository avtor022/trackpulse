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

// CompetitorPanel represents the Competitors management panel
type CompetitorPanel struct {
	competitorService    *service.CompetitorService
	content              *fyne.Container
	table                *widget.Table
	statusLabel          *widget.Label
	window               fyne.Window         // Reference to window for dialogs
	selectedCompetitorID string              // ID of selected competitor
	allCompetitors       []models.Competitor // Cache of all competitors
	// UI components that need to be updated on language change
	toolbar   *widget.Toolbar
	headers   []string
	formItems []*widget.FormItem
}

// updateLocale updates all localized text in the panel
func (p *CompetitorPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("competitor.header.number"),
		locale.T("competitor.header.name"),
		locale.T("competitor.header.country"),
		locale.T("competitor.header.city"),
		locale.T("competitor.header.birthday"),
		locale.T("competitor.header.rating"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *CompetitorPanel) Refresh() {
	p.updateLocale()
}

// NewCompetitorPanel creates a new competitor management panel
func NewCompetitorPanel(competitorService *service.CompetitorService, window fyne.Window) *CompetitorPanel {
	panel := &CompetitorPanel{
		competitorService: competitorService,
		window:            window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the competitor panel UI
func (p *CompetitorPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	p.toolbar = p.createToolbar()

	// Table for displaying competitors
	p.table = p.createCompetitorTable()

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
func (p *CompetitorPanel) createToolbar() *widget.Toolbar {
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

// createRacerTable creates the data table for competitors
func (p *CompetitorPanel) createCompetitorTable() *widget.Table {
	// First load data
	p.allCompetitors, _ = p.competitorService.GetAllCompetitors()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allCompetitors) == 0 {
				return 0, 0
			}
			return len(p.allCompetitors), 9 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allCompetitors) {
				o.(*widget.Label).SetText("")
				return
			}
			competitor := p.allCompetitors[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(competitor.ID)
			case 1:
				o.(*widget.Label).SetText(strconv.Itoa(competitor.CompetitorNumber))
			case 2:
				o.(*widget.Label).SetText(competitor.FullName)
			case 3:
				o.(*widget.Label).SetText(competitor.Country)
			case 4:
				o.(*widget.Label).SetText(competitor.City)
			case 5:
				if competitor.Birthday != nil {
					o.(*widget.Label).SetText(competitor.Birthday.Format("02.01.2006"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 6:
				o.(*widget.Label).SetText(strconv.Itoa(competitor.Rating))
			case 7:
				if !competitor.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(competitor.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 8:
				if !competitor.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(competitor.UpdatedAt.Format("2006-01-02 15:04:05"))
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
	table.SetColumnWidth(1, 80)  // Competitor Number
	table.SetColumnWidth(2, 250) // Full Name
	table.SetColumnWidth(3, 120) // Country
	table.SetColumnWidth(4, 120) // City
	table.SetColumnWidth(5, 140) // Birthday
	table.SetColumnWidth(6, 80)  // Rating
	table.SetColumnWidth(7, 150) // Created At
	table.SetColumnWidth(8, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allCompetitors) {
			p.selectedCompetitorID = p.allCompetitors[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("%s: %s", locale.T("common.selected"), p.allCompetitors[id.Row].FullName))
		}
	}

	return table
}

// refreshData reloads the competitor data
func (p *CompetitorPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allCompetitors, err = p.competitorService.GetAllCompetitors()
		if err != nil {
			fmt.Println("ERROR refreshing data:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}
		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allCompetitors) == 0 {
			p.statusLabel.SetText(locale.T("status.no_competitors"))
		} else {
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_competitors"), len(p.allCompetitors)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new competitor
func (p *CompetitorPanel) showCreateDialog() {
	p.showRacerDialog(locale.T("dialog.new_competitor.title"), nil)
}

// showEditDialog shows the dialog for editing an existing competitor
func (p *CompetitorPanel) showEditDialog() {
	if p.selectedCompetitorID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competitor in cache
	for _, competitor := range p.allCompetitors {
		if competitor.ID == p.selectedCompetitorID {
			p.showRacerDialog(locale.T("dialog.edit.title"), &competitor)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected competitor
func (p *CompetitorPanel) deleteSelected() {
	if p.selectedCompetitorID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competitor in cache
	var racerToDelete *models.Competitor
	for i, competitor := range p.allCompetitors {
		if competitor.ID == p.selectedCompetitorID {
			racerToDelete = &p.allCompetitors[i]
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
				if err := p.competitorService.DeleteCompetitor(racerToDelete.ID); err != nil {
					dialog.ShowError(err, p.window)
					p.statusLabel.SetText(locale.T("status.delete_failed") + ": " + err.Error())
				} else {
					p.refreshData()
					p.selectedCompetitorID = ""
					p.statusLabel.SetText(locale.T("status.deleted_success"))
				}
			}
		},
		p.window,
	)
}

// showRacerDialog shows a dialog for creating or editing a competitor
func (p *CompetitorPanel) showRacerDialog(title string, competitor *models.Competitor) {
	// Create form fields with placeholders and increased width
	numberEntry := widget.NewEntry()
	numberEntry.SetPlaceHolder(locale.T("form.competitor.number_placeholder"))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(locale.T("form.competitor.name_placeholder"))

	countryEntry := widget.NewEntry()
	countryEntry.SetPlaceHolder(locale.T("form.competitor.country_placeholder"))

	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder(locale.T("form.competitor.city_placeholder"))

	birthdayEntry := widget.NewEntry()
	birthdayEntry.SetPlaceHolder(locale.T("form.competitor.birthday_placeholder"))

	ratingEntry := widget.NewEntry()
	ratingEntry.SetPlaceHolder(locale.T("form.competitor.rating_placeholder"))

	if competitor != nil {
		// Edit mode - populate fields
		numberEntry.SetText(strconv.Itoa(competitor.CompetitorNumber))
		nameEntry.SetText(competitor.FullName)
		countryEntry.SetText(competitor.Country)
		cityEntry.SetText(competitor.City)
		if competitor.Birthday != nil {
			birthdayEntry.SetText(competitor.Birthday.Format("02.01.2006"))
		}
		ratingEntry.SetText(strconv.Itoa(competitor.Rating))
	}

	// Create form with localized labels
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.competitor.number"), numberEntry),
		widget.NewFormItem(locale.T("form.competitor.name"), nameEntry),
		widget.NewFormItem(locale.T("form.competitor.country"), countryEntry),
		widget.NewFormItem(locale.T("form.competitor.city"), cityEntry),
		widget.NewFormItem(locale.T("form.competitor.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.competitor.rating"), ratingEntry),
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
		widget.NewFormItem(locale.T("form.competitor.number"), numberEntry),
		widget.NewFormItem(locale.T("form.competitor.name"), nameEntry),
		widget.NewFormItem(locale.T("form.competitor.country"), countryEntry),
		widget.NewFormItem(locale.T("form.competitor.city"), cityEntry),
		widget.NewFormItem(locale.T("form.competitor.birthday"), birthdayEntry),
		widget.NewFormItem(locale.T("form.competitor.rating"), ratingEntry),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Parse values
		number, err := strconv.Atoi(strings.TrimSpace(numberEntry.Text))
		if err != nil {
			errMsg := fmt.Sprintf("invalid competitor number: %v (got: '%s')", err, numberEntry.Text)
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

		var r *models.Competitor
		if competitor != nil {
			// Update existing
			r = competitor
			r.CompetitorNumber = number
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
			if err := p.competitorService.UpdateCompetitor(r); err != nil {
				fmt.Println("ERROR updating competitor:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Competitor updated successfully")

			// Close dialog and refresh data in main thread
			d.Hide()
			fyne.Do(func() {
				p.refreshData()
			})
		} else {
			// Create new
			r = &models.Competitor{
				CompetitorNumber: number,
				FullName:         strings.TrimSpace(nameEntry.Text),
				Country:          strings.TrimSpace(countryEntry.Text),
				City:             strings.TrimSpace(cityEntry.Text),
				Rating:           rating,
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
			if err := p.competitorService.CreateCompetitor(r); err != nil {
				fmt.Println("ERROR creating competitor:", err)
				dialog.ShowError(err, p.window)
				return
			}
			p.statusLabel.SetText("Competitor created successfully")

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
