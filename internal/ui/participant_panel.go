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
	competitionSelect      *widget.Select
	transponderList        *fyne.Container
	transponderCheckboxes  map[string]*widget.Check
	boundTable             *widget.Table
	statusLabel            *widget.Label
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
	if p.competitionSelect != nil {
		p.competitionSelect.SetOptions([]string{locale.T("participants.select.competition.placeholder")})
	}
	if p.boundTable != nil {
		p.boundTable.Refresh()
	}
}

// buildUI constructs the participant panel UI
func (p *ParticipantPanel) buildUI() *fyne.Container {
	// Status label
	p.statusLabel = widget.NewLabel(locale.T("status.ready"))

	// Competition selector
	p.competitionSelect = widget.NewSelect(
		[]string{locale.T("participants.select.competition.placeholder")},
		func(value string) {
			if value == locale.T("participants.select.competition.placeholder") {
				p.selectedCompetitionID = ""
				p.clearTransponderList()
				p.clearBoundTable()
				return
			}
			// Find selected competition
			for _, comp := range p.allCompetitions {
				if comp.CompetitionTitle == value {
					p.selectedCompetitionID = comp.ID
					p.loadAvailableTransponders()
					p.loadBoundParticipants()
					break
				}
			}
		},
	)

	// Toolbar
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			p.refreshData()
		}),
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			p.bindSelectedTransponders()
		}),
	)

	// Transponder list container
	p.transponderList = container.NewVBox()

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

	// Layout
	leftPanel := container.NewBorder(
		container.NewHBox(
			widget.NewLabel(locale.T("participants.select.competition")),
			p.competitionSelect,
		),
		nil,
		nil,
		nil,
		container.NewScroll(p.transponderList),
	)

	rightPanel := container.NewBorder(
		widget.NewLabelWithStyle(locale.T("participants.already_bound"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewScroll(p.boundTable),
	)

	content := container.NewBorder(
		container.NewHBox(toolbar, p.statusLabel),
		nil,
		nil,
		nil,
		container.NewHSplit(leftPanel, rightPanel),
	)

	p.content = content
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

	options := []string{locale.T("participants.select.competition.placeholder")}
	for _, comp := range competitions {
		options = append(options, comp.CompetitionTitle)
	}
	p.competitionSelect.SetOptions(options)
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

	// Get RC models matching the competition's type and scale
	rcModels, err := p.rcModelService.GetModelsByTypeAndScale(selectedComp.ModelType, selectedComp.ModelScale)
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
			"No transponders selected",
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
