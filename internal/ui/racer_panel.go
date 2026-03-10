package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// RacerPanel represents the Racers management panel
type RacerPanel struct {
	racerService *service.RacerService
	content      *fyne.Container
	table        *widget.Table
	statusLabel  *widget.Label
}

// NewRacerPanel creates a new racer management panel
func NewRacerPanel(racerService *service.RacerService) fyne.CanvasObject {
	panel := &RacerPanel{
		racerService: racerService,
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
	newButton := widget.NewButtonWithIcon("New", widget.ContentAdd, p.showCreateDialog)
	editButton := widget.NewButtonWithIcon("Edit", widget.ContentRedo, p.showEditDialog)
	deleteButton := widget.NewButtonWithIcon("Delete", widget.ContentRemove, p.deleteSelected)
	refreshButton := widget.NewButtonWithIcon("Refresh", widget.ContentRefresh, p.refreshData)

	return widget.NewToolbar(
		newButton,
		editButton,
		deleteButton,
		widget.NewSeparator(),
		refreshButton,
	)
}

// createRacerTable creates the data table for racers
func (p *RacerPanel) createRacerTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			racers, _ := p.racerService.GetAllRacers()
			if len(racers) == 0 {
				return 0, 0
			}
			return len(racers), 7 // rows, columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			racers, err := p.racerService.GetAllRacers()
			if err != nil || i.Row >= len(racers) {
				o.(*widget.Label).SetText("")
				return
			}
			racer := racers[i.Row]
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
				if racer.Birthday != "" {
					o.(*widget.Label).SetText(racer.Birthday)
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 5:
				o.(*widget.Label).SetText(strconv.Itoa(racer.Rating))
			case 6:
				o.(*widget.Label).SetText(racer.UpdatedAt)
			}
		},
	)

	// Set column headers
	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header")
	}
	table.UpdateHeader = func(i widget.TableCellID, o fyne.CanvasObject) {
		headers := []string{"#", "Full Name", "Country", "City", "Birthday", "Rating", "Updated"}
		o.(*widget.Label).SetText(headers[i.Col])
		o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
	}

	return table
}

// refreshData reloads the racer data
func (p *RacerPanel) refreshData() {
	if p.table != nil {
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
	selectedRow := p.table.SelectedRow
	if selectedRow < 0 {
		dialog.ShowInformation("Info", "Please select a racer to edit", p.content)
		return
	}

	racers, err := p.racerService.GetAllRacers()
	if err != nil || selectedRow >= len(racers) {
		dialog.ShowError(err, p.content)
		return
	}

	p.showRacerDialog("Edit Racer", racers[selectedRow])
}

// deleteSelected deletes the selected racer
func (p *RacerPanel) deleteSelected() {
	selectedRow := p.table.SelectedRow
	if selectedRow < 0 {
		dialog.ShowInformation("Info", "Please select a racer to delete", p.content)
		return
	}

	racers, err := p.racerService.GetAllRacers()
	if err != nil || selectedRow >= len(racers) {
		dialog.ShowError(err, p.content)
		return
	}

	racer := racers[selectedRow]
	confirmDialog := dialog.ShowConfirm(
		"Confirm Delete",
		"Are you sure you want to delete racer "+racer.FullName+"?",
		func(confirmed bool) {
			if confirmed {
				if err := p.racerService.DeleteRacer(racer.ID); err != nil {
					dialog.ShowError(err, p.content)
					p.statusLabel.SetText("Delete failed: " + err.Error())
				} else {
					p.refreshData()
					p.statusLabel.SetText("Racer deleted successfully")
				}
			}
		},
		p.content,
	)
	confirmDialog.Resize(fyne.NewSize(300, 150))
}

// showRacerDialog shows a dialog for creating or editing a racer
func (p *RacerPanel) showRacerDialog(title string, racer *models.Racer) {
	// Create form fields
	numberEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	countryEntry := widget.NewEntry()
	cityEntry := widget.NewEntry()
	birthdayEntry := widget.NewEntry()
	ratingEntry := widget.NewEntry()

	if racer != nil {
		// Edit mode - populate fields
		numberEntry.SetText(strconv.Itoa(racer.RacerNumber))
		nameEntry.SetText(racer.FullName)
		countryEntry.SetText(racer.Country)
		cityEntry.SetText(racer.City)
		if racer.Birthday != "" {
			birthdayEntry.SetText(racer.Birthday)
		}
		ratingEntry.SetText(strconv.Itoa(racer.Rating))
	}

	// Create form
	form := widget.NewForm(
		widget.NewFormItem("Racer Number", numberEntry),
		widget.NewFormItem("Full Name", nameEntry),
		widget.NewFormItem("Country", countryEntry),
		widget.NewFormItem("City", cityEntry),
		widget.NewFormItem("Birthday (YYYY-MM-DD)", birthdayEntry),
		widget.NewFormItem("Rating", ratingEntry),
	)

	d := dialog.NewCustom(title, "Save", form, p.content)
	
	// Set submit action
	form.OnSubmit = func() {
		// Parse values
		number, err := strconv.Atoi(numberEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid racer number"), p.content)
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
			r.Birthday = birthdayEntry.Text
			r.Rating = rating
			if err := p.racerService.UpdateRacer(r); err != nil {
				dialog.ShowError(err, p.content)
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
				Birthday:    birthdayEntry.Text,
				Rating:      rating,
			}
			if err := p.racerService.CreateRacer(r); err != nil {
				dialog.ShowError(err, p.content)
				return
			}
			p.statusLabel.SetText("Racer created successfully")
		}

		d.Hide()
		p.refreshData()
	}

	d.Show()
}
