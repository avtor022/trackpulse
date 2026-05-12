package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/service"
)

// ParticipantPanel represents the panel for binding transponders to competitions
type ParticipantPanel struct {
	competitionService     *service.CompetitionService
	participantService     *service.CompetitionParticipantService
	competitorModelService *service.CompetitorModelService
	rcModelService         *service.RCModelService
	content                *fyne.Container
	window                 fyne.Window
	selectedCompetitionID  string
	allCompetitions        []models.Competition
	availableTransponders  []models.CompetitorModel
	boundParticipants      []models.CompetitionParticipant
	competitionSelect      *widget.Button
	competitionDisplay     *widget.Label
	transponderList        *fyne.Container
	transponderCheckboxes  map[string]*widget.Check
	boundTable             *widget.Table
	statusLabel            *widget.Label
	currentDialog          dialog.Dialog
}

// NewParticipantPanel creates a new participant binding panel
func NewParticipantPanel(
	competitionService *service.CompetitionService,
	participantService *service.CompetitionParticipantService,
	competitorModelService *service.CompetitorModelService,
	rcModelService *service.RCModelService,
	window fyne.Window,
) *ParticipantPanel {
	panel := &ParticipantPanel{
		competitionService:     competitionService,
		participantService:     participantService,
		competitorModelService: competitorModelService,
		rcModelService:         rcModelService,
		window:                 window,
		transponderCheckboxes:  make(map[string]*widget.Check),
	}
	panel.buildUI()
	return panel
}

// Refresh refreshes the panel UI with current locale
func (p *ParticipantPanel) Refresh() {
	p.updateLocale()
}

// updateLocale updates all localized text in the panel
func (p *ParticipantPanel) updateLocale() {
	if p.statusLabel != nil {
		p.statusLabel.SetText(locale.T("status.ready"))
	}
	if p.competitionDisplay != nil {
		if p.selectedCompetitionID == "" {
			p.competitionDisplay.SetText(locale.T("participants.select.competition.placeholder"))
		} else {
			// Refresh the display with current competition title
			for _, comp := range p.allCompetitions {
				if comp.ID == p.selectedCompetitionID {
					p.competitionDisplay.SetText(comp.CompetitionTitle)
					break
				}
			}
		}
	}
	if p.boundTable != nil {
		p.boundTable.Refresh()
	}
}

// buildUI constructs the participant panel UI
func (p *ParticipantPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Competition display label
	p.competitionDisplay = widget.NewLabel(locale.T("participants.select.competition.placeholder"))

	// Competition select button using reference_popup pattern
	p.competitionSelect = widget.NewButton(locale.T("participants.select.competition.button"), func() {
		p.showCompetitionPopup()
	})

	// Clear selection button
	clearBtn := widget.NewButton(locale.T("participants.clear_selection"), func() {
		p.selectedCompetitionID = ""
		p.competitionDisplay.SetText(locale.T("participants.select.competition.placeholder"))
		p.clearTransponderList()
		p.clearBoundTable()
	})
	clearBtn.Importance = widget.DangerImportance

	// Toolbar with refresh and save buttons
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			p.refreshData()
		}),
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			p.bindSelectedTransponders()
		}),
	)

	// Competition selection area (without toolbar - moved to main content)
	competitionArea := container.NewVBox(
		container.NewHBox(
			widget.NewLabel(locale.T("participants.select.competition")),
			p.competitionDisplay,
			p.competitionSelect,
			clearBtn,
		),
	)

	// Transponder list container
	p.transponderList = container.NewVBox()

	// Left panel with transponder list
	leftPanel := container.NewScroll(p.transponderList)

	// Bound participants table
	p.boundTable = widget.NewTable(
		func() (int, int) {
			if len(p.boundParticipants) == 0 {
				return 0, 0
			}
			return len(p.boundParticipants), 6
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(p.boundParticipants) {
				o.(*widget.Label).SetText("")
				return
			}
			participant := p.boundParticipants[i.Row]

			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(participant.ID[:8])
			case 1:
				o.(*widget.Label).SetText(p.getCompetitionTitle(participant.CompetitionID))
			case 2:
				o.(*widget.Label).SetText(p.getCompetitorName(participant.CompetitorModelID))
			case 3:
				o.(*widget.Label).SetText(p.getRCModelName(participant.CompetitorModelID))
			case 4:
				o.(*widget.Label).SetText(p.getTransponderNumber(participant.CompetitorModelID))
			case 5:
				if participant.GridPosition != nil {
					o.(*widget.Label).SetText(strconv.Itoa(*participant.GridPosition))
				} else {
					o.(*widget.Label).SetText("-")
				}
			}
			o.(*widget.Label).Truncation = fyne.TextTruncateEllipsis
		},
	)
	p.boundTable.ShowHeaderRow = true
	p.boundTable.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("Header", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}
	p.boundTable.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		headers := []string{
			locale.T("participants.table.id"),
			locale.T("participants.table.competition"),
			locale.T("participants.table.competitor"),
			locale.T("participants.table.model"),
			locale.T("participants.table.transponder"),
			locale.T("participants.table.grid"),
		}
		if id.Col >= 0 && id.Col < len(headers) {
			o.(*widget.Label).SetText(headers[id.Col])
			o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
		}
	}

	// Set column widths based on header content
	p.boundTable.SetColumnWidth(0, 100) // ID
	p.boundTable.SetColumnWidth(1, 200) // Competition
	p.boundTable.SetColumnWidth(2, 150) // Competitor
	p.boundTable.SetColumnWidth(3, 150) // Model
	p.boundTable.SetColumnWidth(4, 120) // Transponder
	p.boundTable.SetColumnWidth(5, 80)  // Grid

	// Right panel with bound participants table
	rightPanel := container.NewBorder(
		widget.NewLabelWithStyle(locale.T("participants.already_bound"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewScroll(p.boundTable),
	)

	// Main content: competition selection area at top, split panel below
	content := container.NewBorder(
		container.NewVBox(
			competitionArea,
			container.NewHBox(toolbar, p.statusLabel),
		),
		nil,
		nil,
		nil,
		container.NewHSplit(leftPanel, rightPanel),
	)

	p.content = content
	// Load competitions on panel initialization
	p.loadCompetitions()

	return content
}

// loadCompetitions loads all competitions into the selector
func (p *ParticipantPanel) loadCompetitions() {
	competitions, err := p.competitionService.GetAllCompetitions()
	if err != nil {
		p.showError(locale.T("status.error_loading"))
		return
	}
	p.allCompetitions = competitions
}

// showCompetitionPopup displays a popup to select a competition using reference_popup pattern
func (p *ParticipantPanel) showCompetitionPopup() {
	// Create container for competition items
	itemContainer := container.NewVBox()

	// Add existing competitions as select buttons
	for _, comp := range p.allCompetitions {
		compItem := comp
		selectBtn := widget.NewButton(comp.CompetitionTitle, func() {
			// Select this competition
			p.selectedCompetitionID = compItem.ID
			p.competitionDisplay.SetText(compItem.CompetitionTitle)

			// Hide popup
			if p.currentDialog != nil {
				p.currentDialog.Hide()
				p.currentDialog = nil
			}

			// Load transponders and bound participants
			p.loadAvailableTransponders()
			p.loadBoundParticipants()
		})
		selectBtn.Alignment = widget.ButtonAlignLeading
		selectBtn.Importance = widget.MediumImportance
		itemContainer.Add(selectBtn)
	}

	// Create popup dialog
	popup := dialog.NewCustomWithoutButtons(locale.T("participants.select.competition.popup_title"), itemContainer, p.window)

	// Add close button
	closeBtn := widget.NewButton(locale.T("common.close"), func() {
		popup.Hide()
		p.currentDialog = nil
	})

	popup.SetButtons([]fyne.CanvasObject{closeBtn})

	// Resize popup
	parentSize := p.window.Canvas().Size()
	popupWidth := parentSize.Width * 0.4
	if popupWidth < 400 {
		popupWidth = 400
	}
	popup.Resize(fyne.NewSize(popupWidth, popup.MinSize().Height))

	p.currentDialog = popup
	popup.Show()
}

// loadAvailableTransponders loads transponders that match the selected competition's model type and scale
func (p *ParticipantPanel) loadAvailableTransponders() {
	if p.selectedCompetitionID == "" {
		return
	}

	// Get selected competition
	var selectedComp *models.Competition
	for _, comp := range p.allCompetitions {
		if comp.ID == p.selectedCompetitionID {
			selectedComp = &comp
			break
		}
	}
	if selectedComp == nil {
		return
	}

	// Convert "*" to "all" for the service method
	modelType := selectedComp.ModelType
	if modelType == "*" {
		modelType = "all"
	}
	modelScale := selectedComp.ModelScale
	if modelScale == "*" {
		modelScale = "all"
	}

	// Get RC models matching the competition's type and scale
	rcModels, err := p.rcModelService.GetModelsByTypeAndScale(modelType, modelScale)
	if err != nil {
		p.showError(locale.T("status.error_loading"))
		return
	}

	// Get RC model IDs
	rcModelIDs := make([]string, 0, len(rcModels))
	for _, model := range rcModels {
		rcModelIDs = append(rcModelIDs, model.ID)
	}

	// Get competitor models (transponders) for these RC models
	transponders, err := p.competitorModelService.GetCompetitorModelsByRCModelIDs(rcModelIDs)
	if err != nil {
		p.showError(locale.T("status.error_loading"))
		return
	}

	// Filter out already bound transponders
	boundModelIDs := make(map[string]bool)
	for _, bp := range p.boundParticipants {
		boundModelIDs[bp.CompetitorModelID] = true
	}

	var available []models.CompetitorModel
	for _, t := range transponders {
		if !boundModelIDs[t.ID] {
			available = append(available, t)
		}
	}

	p.availableTransponders = available
	p.renderTransponderList()
}

// renderTransponderList renders the list of available transponders with checkboxes
func (p *ParticipantPanel) renderTransponderList() {
	p.transponderList.Objects = nil
	p.transponderCheckboxes = make(map[string]*widget.Check)

	if len(p.availableTransponders) == 0 {
		p.transponderList.Add(widget.NewLabel(locale.T("participants.no_transponders")))
		p.transponderList.Refresh()
		return
	}

	for _, t := range p.availableTransponders {
		competitorName := p.getCompetitorName(t.ID)
		modelName := p.getRCModelName(t.ID)
		labelText := fmt.Sprintf("%s - %s (%s)", t.TransponderNumber, competitorName, modelName)

		check := widget.NewCheck(labelText, nil)
		p.transponderCheckboxes[t.ID] = check
		p.transponderList.Add(check)
	}
	p.transponderList.Refresh()
}

// clearTransponderList clears the transponder list
func (p *ParticipantPanel) clearTransponderList() {
	p.transponderList.Objects = nil
	p.transponderCheckboxes = make(map[string]*widget.Check)
	p.availableTransponders = nil
	p.transponderList.Refresh()
}

// loadBoundParticipants loads participants already bound to the selected competition
func (p *ParticipantPanel) loadBoundParticipants() {
	if p.selectedCompetitionID == "" {
		return
	}

	participants, err := p.participantService.GetParticipantsByCompetitionID(p.selectedCompetitionID)
	if err != nil {
		p.showError(locale.T("status.error_loading"))
		return
	}
	p.boundParticipants = participants
	p.boundTable.Refresh()
}

// clearBoundTable clears the bound participants table
func (p *ParticipantPanel) clearBoundTable() {
	p.boundParticipants = nil
	p.boundTable.Refresh()
}

// bindSelectedTransponders binds the selected transponders to the competition
func (p *ParticipantPanel) bindSelectedTransponders() {
	if p.selectedCompetitionID == "" {
		dialog.ShowInformation(
			locale.T("dialog.info"),
			locale.T("participants.no_competition_selected"),
			p.window,
		)
		return
	}

	// Collect selected transponder IDs
	var selectedIDs []string
	for id, check := range p.transponderCheckboxes {
		if check.Checked {
			selectedIDs = append(selectedIDs, id)
		}
	}

	if len(selectedIDs) == 0 {
		dialog.ShowInformation(
			locale.T("dialog.info"),
			locale.T("participants.no_transponders_selected"),
			p.window,
		)
		return
	}

	// Bind transponders
	addedIDs, errors := p.participantService.AddParticipantsBulk(p.selectedCompetitionID, selectedIDs)

	// Show results
	successCount := len(addedIDs)
	errorCount := len(errors)

	message := fmt.Sprintf(locale.T("participants.bind.success"), successCount)
	if errorCount > 0 {
		message += fmt.Sprintf("\n"+locale.T("participants.bind.errors"), errorCount)
		for _, err := range errors {
			message += fmt.Sprintf("\n- %v", err)
		}
	}

	dialog.ShowInformation(locale.T("dialog.success"), message, p.window)

	// Refresh data
	p.refreshData()
}

// refreshData reloads all data
func (p *ParticipantPanel) refreshData() {
	p.loadCompetitions()
	if p.selectedCompetitionID != "" {
		p.loadAvailableTransponders()
		p.loadBoundParticipants()
	}
}

// Helper methods to get related data
func (p *ParticipantPanel) getCompetitionTitle(compID string) string {
	for _, comp := range p.allCompetitions {
		if comp.ID == compID {
			return comp.CompetitionTitle
		}
	}
	return compID
}

func (p *ParticipantPanel) getCompetitorName(cmID string) string {
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return cmID
	}
	// In a real implementation, you would fetch the competitor name from the competitor service
	return cm.TransponderNumber
}

func (p *ParticipantPanel) getRCModelName(cmID string) string {
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return ""
	}
	model, err := p.rcModelService.GetModelByID(cm.RCModelID)
	if err != nil || model == nil {
		return ""
	}
	return model.ModelName
}

func (p *ParticipantPanel) getTransponderNumber(cmID string) string {
	cm, err := p.competitorModelService.GetCompetitorModelByID(cmID)
	if err != nil || cm == nil {
		return ""
	}
	return cm.TransponderNumber
}

func (p *ParticipantPanel) showError(message string) {
	p.statusLabel.SetText(message)
	dialog.ShowError(fmt.Errorf(message), p.window)
}
