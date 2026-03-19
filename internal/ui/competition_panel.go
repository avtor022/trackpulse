package ui

import (
	"fmt"
	"strconv"
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

// CompetitionPanel represents the Competitions management panel
type CompetitionPanel struct {
	competitionService *service.CompetitionService
	content            *fyne.Container
	table              *widget.Table
	statusLabel        *widget.Label
	window             fyne.Window          // Reference to window for dialogs
	selectedID         string               // ID of selected competition
	allCompetitions    []models.Competition // Cache of all competitions
	headers            []string             // Localized table headers
}

// updateLocale updates all localized text in the panel
func (p *CompetitionPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}

	// Update headers
	headers := []string{
		locale.T("common.id"),
		locale.T("competition.header.title"),
		locale.T("competition.header.type"),
		locale.T("model.header.type"),
		locale.T("model.header.scale"),
		locale.T("competition.header.track"),
		locale.T("competition.header.laps"),
		locale.T("competition.header.time_limit"),
		locale.T("competition.header.start"),
		locale.T("competition.header.finish"),
		locale.T("competition.header.status"),
		locale.T("model.header.created"),
		locale.T("model.header.updated"),
	}
	p.headers = headers

	if p.table != nil {
		p.table.Refresh()
	}
}

// Refresh refreshes the panel UI with current locale
func (p *CompetitionPanel) Refresh() {
	p.updateLocale()
}

// NewCompetitionPanel creates a new competition management panel
func NewCompetitionPanel(competitionService *service.CompetitionService, window fyne.Window) *CompetitionPanel {
	panel := &CompetitionPanel{
		competitionService: competitionService,
		window:             window,
	}
	panel.buildUI()
	return panel
}

// buildUI constructs the competition panel UI
func (p *CompetitionPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Toolbar with actions
	toolbar := p.createToolbar()

	// Table for displaying competitions
	p.table = p.createCompetitionTable()

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
func (p *CompetitionPanel) createToolbar() *widget.Toolbar {
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

// createCompetitionTable creates the data table for competitions
func (p *CompetitionPanel) createCompetitionTable() *widget.Table {
	// First load data
	p.refreshData()

	table := widget.NewTable(
		func() (int, int) {
			if len(p.allCompetitions) == 0 {
				return 0, 0
			}
			return len(p.allCompetitions), 12 // rows, columns
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.allCompetitions) {
				o.(*widget.Label).SetText("")
				return
			}
			c := p.allCompetitions[i.Row]

			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(c.ID)
			case 1:
				o.(*widget.Label).SetText(c.CompetitionTitle)
			case 2:
				o.(*widget.Label).SetText(c.CompetitionType)
			case 3:
				o.(*widget.Label).SetText(c.ModelType)
			case 4:
				o.(*widget.Label).SetText(c.ModelScale)
			case 5:
				o.(*widget.Label).SetText(c.TrackName)
			case 6:
				if c.LapCountTarget != nil {
					o.(*widget.Label).SetText(strconv.Itoa(*c.LapCountTarget))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 7:
				if c.TimeLimitMinutes != nil {
					o.(*widget.Label).SetText(strconv.Itoa(*c.TimeLimitMinutes))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 8:
				if c.TimeStart != nil {
					o.(*widget.Label).SetText(c.TimeStart.Format("2006-01-02 15:04"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 9:
				if c.TimeFinish != nil {
					o.(*widget.Label).SetText(c.TimeFinish.Format("2006-01-02 15:04"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 10:
				o.(*widget.Label).SetText(c.Status)
			case 11:
				if !c.CreatedAt.IsZero() {
					o.(*widget.Label).SetText(c.CreatedAt.Format("2006-01-02 15:04:05"))
				} else {
					o.(*widget.Label).SetText("-")
				}
			case 12:
				if !c.UpdatedAt.IsZero() {
					o.(*widget.Label).SetText(c.UpdatedAt.Format("2006-01-02 15:04:05"))
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
	table.SetColumnWidth(0, 280)  // ID
	table.SetColumnWidth(1, 250)  // Title
	table.SetColumnWidth(2, 120)  // Type
	table.SetColumnWidth(3, 100)  // Model Type
	table.SetColumnWidth(4, 80)   // Scale
	table.SetColumnWidth(5, 150)  // Track Name
	table.SetColumnWidth(6, 60)   // Laps
	table.SetColumnWidth(7, 80)   // Time Limit
	table.SetColumnWidth(8, 150)  // Start
	table.SetColumnWidth(9, 150)  // Finish
	table.SetColumnWidth(10, 100) // Status
	table.SetColumnWidth(11, 150) // Created At
	table.SetColumnWidth(12, 150) // Updated At

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(p.allCompetitions) {
			p.selectedID = p.allCompetitions[id.Row].ID
			p.statusLabel.SetText(fmt.Sprintf("Selected: %s", p.allCompetitions[id.Row].CompetitionTitle))
		}
	}

	return table
}

// refreshData reloads the competition data
func (p *CompetitionPanel) refreshData() {
	if p.table != nil {
		// Update data cache
		var err error
		p.allCompetitions, err = p.competitionService.GetAllCompetitions()
		if err != nil {
			fmt.Println("ERROR refreshing competitions:", err)
			p.statusLabel.SetText("Error refreshing data")
			return
		}

		// Force table to recalculate rows count and update cell contents
		p.table.Refresh()
		if len(p.allCompetitions) == 0 {
			p.statusLabel.SetText(locale.T("status.no_competitions"))
		} else {
			p.statusLabel.SetText(fmt.Sprintf(locale.T("status.loaded_competitions"), len(p.allCompetitions)))
		}
	}
}

// showCreateDialog shows the dialog for creating a new competition
func (p *CompetitionPanel) showCreateDialog() {
	p.showCompetitionDialog(locale.T("dialog.new_competition.title"), nil)
}

// showEditDialog shows the dialog for editing an existing competition
func (p *CompetitionPanel) showEditDialog() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competition in cache
	for _, c := range p.allCompetitions {
		if c.ID == p.selectedID {
			p.showCompetitionDialog(locale.T("dialog.edit.title"), &c)
			return
		}
	}

	dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
}

// deleteSelected deletes the selected competition
func (p *CompetitionPanel) deleteSelected() {
	if p.selectedID == "" {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.select_first"), p.window)
		return
	}

	// Look for selected competition in cache
	var cToDelete *models.Competition
	for i, c := range p.allCompetitions {
		if c.ID == p.selectedID {
			cToDelete = &p.allCompetitions[i]
			break
		}
	}

	if cToDelete == nil {
		dialog.ShowInformation(locale.T("common.info"), locale.T("info.not_found"), p.window)
		return
	}

	// Show confirmation dialog
	dialog.ShowConfirm(
		locale.T("dialog.delete.title"),
		fmt.Sprintf(locale.T("dialog.delete.message"), cToDelete.CompetitionTitle),
		func(confirmed bool) {
			if confirmed {
				if err := p.competitionService.DeleteCompetition(cToDelete.ID); err != nil {
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

// showCompetitionDialog shows a dialog for creating or editing a competition
func (p *CompetitionPanel) showCompetitionDialog(title string, competition *models.Competition) {
	// Create form fields
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder(locale.T("form.competition.title_placeholder"))

	// Competition type select
	competitionTypes := []string{"qualifying", "final", "practice", "heat"}
	typeSelect := widget.NewSelect(competitionTypes, nil)
	typeSelect.PlaceHolder = locale.T("form.competition.type_placeholder")

	// Model type entry
	modelTypeEntry := widget.NewEntry()
	modelTypeEntry.SetPlaceHolder(locale.T("form.competition.model_type_placeholder"))

	// Model scale entry
	modelScaleEntry := widget.NewEntry()
	modelScaleEntry.SetPlaceHolder(locale.T("form.competition.model_scale_placeholder"))

	// Track name entry
	trackEntry := widget.NewEntry()
	trackEntry.SetPlaceHolder(locale.T("form.competition.track_placeholder"))

	// Lap count target entry
	lapCountEntry := widget.NewEntry()
	lapCountEntry.SetPlaceHolder(locale.T("form.competition.lap_count_placeholder"))

	// Time limit entry
	timeLimitEntry := widget.NewEntry()
	timeLimitEntry.SetPlaceHolder(locale.T("form.competition.time_limit_placeholder"))

	// Status select
	statuses := []string{"scheduled", "running", "finished", "cancelled"}
	statusSelect := widget.NewSelect(statuses, nil)
	statusSelect.PlaceHolder = locale.T("form.competition.status_placeholder")

	if competition != nil {
		// Edit mode - populate fields
		titleEntry.SetText(competition.CompetitionTitle)
		typeSelect.SetSelected(competition.CompetitionType)
		modelTypeEntry.SetText(competition.ModelType)
		modelScaleEntry.SetText(competition.ModelScale)
		trackEntry.SetText(competition.TrackName)
		if competition.LapCountTarget != nil {
			lapCountEntry.SetText(strconv.Itoa(*competition.LapCountTarget))
		}
		if competition.TimeLimitMinutes != nil {
			timeLimitEntry.SetText(strconv.Itoa(*competition.TimeLimitMinutes))
		}
		statusSelect.SetSelected(competition.Status)
	}

	// Create form with localized labels
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.competition.title"), titleEntry),
		widget.NewFormItem(locale.T("form.competition.type"), typeSelect),
		widget.NewFormItem(locale.T("form.competition.model_type"), modelTypeEntry),
		widget.NewFormItem(locale.T("form.competition.model_scale"), modelScaleEntry),
		widget.NewFormItem(locale.T("form.competition.track"), trackEntry),
		widget.NewFormItem(locale.T("form.competition.lap_count"), lapCountEntry),
		widget.NewFormItem(locale.T("form.competition.time_limit"), timeLimitEntry),
		widget.NewFormItem(locale.T("form.competition.status"), statusSelect),
	)

	// Create dialog without buttons first so we can reference it in the callback
	d := dialog.NewCustomWithoutButtons(title, form, p.window)

	// Create save button with callback that has access to 'd'
	saveBtn := widget.NewButton(locale.T("common.save"), func() {
		// Validate title
		compTitle := strings.TrimSpace(titleEntry.Text)
		if compTitle == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.title")), p.window)
			return
		}

		// Validate type
		compType := typeSelect.Selected
		if compType == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.type")), p.window)
			return
		}

		// Parse lap count
		var lapCountTarget *int
		if lapCountEntry.Text != "" {
			lct, err := strconv.Atoi(strings.TrimSpace(lapCountEntry.Text))
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid lap count: %v", err), p.window)
				return
			}
			lapCountTarget = &lct
		}

		// Parse time limit
		var timeLimitMinutes *int
		if timeLimitEntry.Text != "" {
			tlm, err := strconv.Atoi(strings.TrimSpace(timeLimitEntry.Text))
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid time limit: %v", err), p.window)
				return
			}
			timeLimitMinutes = &tlm
		}

		var newC *models.Competition
		if competition != nil {
			// Update existing
			newC = competition
			newC.CompetitionTitle = compTitle
			newC.CompetitionType = compType
			newC.ModelType = strings.TrimSpace(modelTypeEntry.Text)
			newC.ModelScale = strings.TrimSpace(modelScaleEntry.Text)
			newC.TrackName = strings.TrimSpace(trackEntry.Text)
			newC.LapCountTarget = lapCountTarget
			newC.TimeLimitMinutes = timeLimitMinutes
			newC.Status = statusSelect.Selected

			if err := p.competitionService.UpdateCompetition(newC); err != nil {
				fmt.Println("ERROR updating competition:", err)
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
			newC = &models.Competition{
				CompetitionTitle: compTitle,
				CompetitionType:  compType,
				ModelType:        strings.TrimSpace(modelTypeEntry.Text),
				ModelScale:       strings.TrimSpace(modelScaleEntry.Text),
				TrackName:        strings.TrimSpace(trackEntry.Text),
				LapCountTarget:   lapCountTarget,
				TimeLimitMinutes: timeLimitMinutes,
				Status:           statusSelect.Selected,
			}

			if err := p.competitionService.CreateCompetition(newC); err != nil {
				fmt.Println("ERROR creating competition:", err)
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
	dialogHeight := parentSize.Height * 0.7
	if dialogHeight < 500 {
		dialogHeight = 500 // Minimum height
	}

	// Resize the dialog window
	d.Resize(fyne.NewSize(dialogWidth, dialogHeight))
}
