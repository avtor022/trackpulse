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
	window             fyne.Window                 // Reference to window for dialogs
	selectedID         string                      // ID of selected competition
	allCompetitions    []models.Competition        // Cache of all competitions
	headers            []string                    // Localized table headers
	allModelTypes      []models.RCModelType        // Cache of all model types
	allModelScales     []models.RCModelScale       // Cache of all model scales
	allCompetitionTracks []models.CompetitionTrack // Cache of all competition tracks
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
			p.statusLabel.SetText(fmt.Sprintf("%s: %s", locale.T("common.selected"), p.allCompetitions[id.Row].CompetitionTitle))
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

		// Load model types for dropdown
		p.allModelTypes, err = p.competitionService.GetAllModelTypes()
		if err != nil {
			fmt.Println("ERROR loading model types:", err)
		}

		// Load model scales for dropdown
		p.allModelScales, err = p.competitionService.GetAllModelScales()
		if err != nil {
			fmt.Println("ERROR loading model scales:", err)
		}

		// Load competition tracks for dropdown
		p.allCompetitionTracks, err = p.competitionService.GetAllCompetitionTracks()
		if err != nil {
			fmt.Println("ERROR loading competition tracks:", err)
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

	// Competition type select with localized options
	competitionTypes := []string{
		locale.T("competition.type.qualifying"),
		locale.T("competition.type.final"),
		locale.T("competition.type.practice"),
		locale.T("competition.type.heat"),
	}
	typeSelect := widget.NewSelect(competitionTypes, nil)
	typeSelect.PlaceHolder = locale.T("common.select_one")

	// Map display names to internal values
	typeMap := map[string]string{
		locale.T("competition.type.qualifying"): "qualifying",
		locale.T("competition.type.final"):      "final",
		locale.T("competition.type.practice"):   "practice",
		locale.T("competition.type.heat"):       "heat",
	}

	// Status select with localized options
	statuses := []string{
		locale.T("competition.status.scheduled"),
		locale.T("competition.status.running"),
		locale.T("competition.status.finished"),
		locale.T("competition.status.cancelled"),
	}
	statusSelect := widget.NewSelect(statuses, nil)
	statusSelect.PlaceHolder = locale.T("common.select_one")

	// Map display names to internal values
	statusMap := map[string]string{
		locale.T("competition.status.scheduled"): "scheduled",
		locale.T("competition.status.running"):   "running",
		locale.T("competition.status.finished"):  "finished",
		locale.T("competition.status.cancelled"): "cancelled",
	}

	// Model type select - populate from database with "All Types" option for mass race
	modelTypeOptions := []string{locale.T("competition.model_type.all")} // "All types" option first
	modelTypeNames := []string{"*"}                                      // Internal value for "all types"
	for _, mt := range p.allModelTypes {
		modelTypeOptions = append(modelTypeOptions, mt.Name)
		modelTypeNames = append(modelTypeNames, mt.Name)
	}
	modelTypeSelect := widget.NewSelect(modelTypeOptions, nil)
	modelTypeSelect.PlaceHolder = locale.T("common.select_one")

	// Model scale select - populate from database with "All Scales" option for mass race
	scaleOptions := []string{locale.T("competition.model_scale.all")} // "All scales" option first
	scaleNames := []string{"*"}                                       // Internal value for "all scales"
	for _, ms := range p.allModelScales {
		scaleOptions = append(scaleOptions, ms.Name)
		scaleNames = append(scaleNames, ms.Name)
	}
	modelScaleSelect := widget.NewSelect(scaleOptions, nil)
	modelScaleSelect.PlaceHolder = locale.T("common.select_one")

	// Track name select with popup manager (similar to brand/scale/type in rc_model_panel)
	var trackSelect *widget.Select
	var currentDialog dialog.Dialog

	// Extract track names
	var existingTracks []string
	for _, track := range p.allCompetitionTracks {
		existingTracks = append(existingTracks, track.Name)
	}

	// Add option to create new track
	newTrackOption := "+ " + locale.T("common.add") + " " + strings.TrimSuffix(locale.T("form.competition.track"), ":")
	trackSelectOptions := append(existingTracks, newTrackOption)

	var mainDialog dialog.Dialog
	var trackPopupManager *ReferencePopupManager

	// Helper function to update track button text
	var trackButton *widget.Button
	updateTrackButton := func(selected string) {
		if trackButton == nil {
			return
		}
		if selected == "" {
			trackButton.SetText(locale.T("common.select_one"))
		} else {
			trackButton.SetText(selected)
		}
	}

	// Create the hidden Select widget to maintain compatibility
	trackSelect = widget.NewSelect(trackSelectOptions, func(selected string) {
		updateTrackButton(selected)
		if selected == newTrackOption {
			// This case is now handled in the popup
			trackSelect.SetSelected("")
		}
	})

	var showTrackPopup func()
	showTrackPopup = func() {
		if trackPopupManager == nil {
			// Convert existingTracks to ReferenceItem slice
			items := make([]ReferenceItem, len(existingTracks))
			for i, t := range existingTracks {
				items[i] = ReferenceItem{Name: t}
			}

			trackPopupManager = NewReferencePopupManager(
				p.window,
				ReferencePopupConfig{
					Title:          "common.select_one",
					AddTitle:       "dialog.add_track.title",
					AddLabel:       "dialog.add_track.label",
					AddPlaceholder: "dialog.add_track.placeholder",
					DeleteMessage:  "dialog.delete_track.message",
					NewErrorExists: "dialog.new_track.error_exists",
					EnterNameInfo:  "info.enter_track_name",
					GetAllFunc: func() ([]ReferenceItem, error) {
						allTracks, err := p.competitionService.GetAllCompetitionTracks()
						if err != nil {
							return nil, err
						}
						result := make([]ReferenceItem, len(allTracks))
						for i, t := range allTracks {
							result[i] = ReferenceItem{Name: t.Name}
						}
						return result, nil
					},
					AddFunc: func(name string) error {
						return p.competitionService.AddCompetitionTrack(name)
					},
					DeleteFunc: func(name string) error {
						return p.competitionService.DeleteCompetitionTrack(name)
					},
					OnItemSelected: func(selected string) {
						trackSelect.SetSelected(selected)
						updateTrackButton(selected)
					},
					UpdateOptions: func(opts []string) {
						trackSelectOptions = opts
						trackSelect.Options = trackSelectOptions
					},
				},
				existingTracks,
				newTrackOption,
				func(selected string) {
					trackSelect.SetSelected(selected)
					updateTrackButton(selected)
				},
				func(opts []string) {
					trackSelectOptions = opts
					trackSelect.Options = trackSelectOptions
				},
			)
		}
		trackPopupManager.ShowPopup(mainDialog, &currentDialog, func(d dialog.Dialog) {
			currentDialog = d
		})
	}

	// Create a button that shows the track popup when clicked
	initialTrackText := locale.T("common.select_one")
	if competition != nil && competition.TrackName != "" {
		initialTrackText = competition.TrackName
	}
	trackButton = widget.NewButton(initialTrackText, func() {
		if mainDialog != nil {
			mainDialog.Hide()
		}
		showTrackPopup()
	})

	// Use the button as the track widget
	var trackWidget = trackButton

	// Lap count target entry
	lapCountEntry := widget.NewEntry()
	lapCountEntry.SetPlaceHolder(locale.T("form.competition.lap_count_placeholder"))

	// Time limit entry
	timeLimitEntry := widget.NewEntry()
	timeLimitEntry.SetPlaceHolder(locale.T("form.competition.time_limit_placeholder"))

	if competition != nil {
		// Edit mode - populate fields
		titleEntry.SetText(competition.CompetitionTitle)
		// Map internal value to localized display
		if localizedType, ok := reverseMap(typeMap, competition.CompetitionType); ok {
			typeSelect.SetSelected(localizedType)
		}
		// Map internal model type to display value
		if competition.ModelType == "*" {
			modelTypeSelect.SetSelected(locale.T("competition.model_type.all"))
		} else {
			for i, name := range modelTypeNames {
				if name == competition.ModelType {
					modelTypeSelect.SetSelected(modelTypeOptions[i])
					break
				}
			}
		}
		// Map internal model scale to display value
		if competition.ModelScale == "*" {
			modelScaleSelect.SetSelected(locale.T("competition.model_scale.all"))
		} else if competition.ModelScale != "" {
			for i, name := range scaleNames {
				if name == competition.ModelScale {
					modelScaleSelect.SetSelected(scaleOptions[i])
					break
				}
			}
		}
		// Set track button text for edit mode
		if competition.TrackName != "" {
			updateTrackButton(competition.TrackName)
		}
		if competition.LapCountTarget != nil {
			lapCountEntry.SetText(strconv.Itoa(*competition.LapCountTarget))
		}
		if competition.TimeLimitMinutes != nil {
			timeLimitEntry.SetText(strconv.Itoa(*competition.TimeLimitMinutes))
		}
		// Map internal status to localized display
		if localizedStatus, ok := reverseMap(statusMap, competition.Status); ok {
			statusSelect.SetSelected(localizedStatus)
		}
	}

	// Create form with localized labels in the requested order:
	// 1. Competition Type, 2. Model Type, 3. Model Scale, 4. Title, 5. Track, 6. Lap Count, 7. Time Limit, 8. Status
	form := widget.NewForm(
		widget.NewFormItem(locale.T("form.competition.type"), typeSelect),
		widget.NewFormItem(locale.T("form.competition.model_type"), modelTypeSelect),
		widget.NewFormItem(locale.T("form.competition.model_scale"), modelScaleSelect),
		widget.NewFormItem(locale.T("form.competition.title"), titleEntry),
		widget.NewFormItem(locale.T("form.competition.track"), trackWidget),
		widget.NewFormItem(locale.T("form.competition.lap_count"), lapCountEntry),
		widget.NewFormItem(locale.T("form.competition.time_limit"), timeLimitEntry),
		widget.NewFormItem(locale.T("form.competition.status"), statusSelect),
	)

	// Store mainDialog for use in popup
	mainDialog = d

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

		// Validate type and map to internal value
		compType := typeSelect.Selected
		if compType == "" {
			dialog.ShowError(fmt.Errorf(locale.T("error.required.type")), p.window)
			return
		}
		// Map localized display to internal value
		if mappedType, ok := typeMap[compType]; ok {
			compType = mappedType
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

		// Map status to internal value
		statusValue := statusSelect.Selected
		if mappedStatus, ok := statusMap[statusValue]; ok {
			statusValue = mappedStatus
		}

		// Map model type selection to internal value
		modelTypeValue := modelTypeSelect.Selected
		var modelTypeInternal string
		for i, opt := range modelTypeOptions {
			if opt == modelTypeValue {
				modelTypeInternal = modelTypeNames[i]
				break
			}
		}

		// Map model scale selection to internal value
		modelScaleValue := modelScaleSelect.Selected
		var modelScaleInternal string
		for i, opt := range scaleOptions {
			if opt == modelScaleValue {
				modelScaleInternal = scaleNames[i]
				break
			}
		}

		// Get track name from button text
		trackName := ""
		if trackButton != nil {
			trackName = trackButton.Text
			if trackName == locale.T("common.select_one") {
				trackName = ""
			}
		}

		var newC *models.Competition
		if competition != nil {
			// Update existing
			newC = competition
			newC.CompetitionTitle = compTitle
			newC.CompetitionType = compType
			newC.ModelType = modelTypeInternal
			newC.ModelScale = modelScaleInternal
			newC.TrackName = strings.TrimSpace(trackName)
			newC.LapCountTarget = lapCountTarget
			newC.TimeLimitMinutes = timeLimitMinutes
			newC.Status = statusValue

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
				ModelType:        modelTypeInternal,
				ModelScale:       modelScaleInternal,
				TrackName:        strings.TrimSpace(trackName),
				LapCountTarget:   lapCountTarget,
				TimeLimitMinutes: timeLimitMinutes,
				Status:           statusValue,
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

// reverseMap finds the key by value in a map
func reverseMap(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}
