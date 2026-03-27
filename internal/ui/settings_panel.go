package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/service"
)

// SettingsPanel represents the Settings panel
type SettingsPanel struct {
	settingsService *service.SettingsService
	content         *fyne.Container
	config          *Config
	window          fyne.Window
	// UI components that need to be updated on language change
	languageSelect *widget.Select
	portScanner    *PortScanner
	languageForm   *widget.Form
}

// NewSettingsPanel creates a new settings panel
func NewSettingsPanel(settingsService *service.SettingsService, config *Config, window fyne.Window) *SettingsPanel {
	panel := &SettingsPanel{
		settingsService: settingsService,
		config:          config,
		window:          window,
	}
	panel.buildUI()
	return panel
}

// updateLocale updates all localized text in the panel
func (p *SettingsPanel) updateLocale() {
	// Language label is now part of the form, updated via refreshUI
}

// Refresh refreshes the panel UI with current locale
func (p *SettingsPanel) Refresh() {
	p.updateLocale()
	p.refreshUI()
}

// buildUI constructs the settings panel UI
func (p *SettingsPanel) buildUI() *fyne.Container {
	// Create language selector
	options := make([]string, 0, len(locale.SupportedLocales))
	for _, name := range locale.SupportedLocales {
		options = append(options, name)
	}

	// Find current language name
	currentName := "English"
	for code, name := range locale.SupportedLocales {
		if code == p.config.Language {
			currentName = name
			break
		}
	}

	// Create select without callback first
	p.languageSelect = widget.NewSelect(options, nil)

	// Set initial value without triggering callback
	p.languageSelect.SetSelected(currentName)

	// Now set the callback for future changes
	p.languageSelect.OnChanged = func(selected string) {
		// Find the code for the selected language
		var selectedCode string
		for code, name := range locale.SupportedLocales {
			if name == selected {
				selectedCode = code
				break
			}
		}

		if selectedCode != "" {
			locale.SetLocale(selectedCode)
			p.config.Language = selectedCode

			// Save to database
			if p.settingsService != nil {
				err := p.settingsService.SetLocale(selectedCode)
				if err != nil {
					// Log error but continue with UI update
					fmt.Printf("Failed to save locale: %v\n", err)
				}
			}

			p.refreshUI()
		}
	}

	// Create port scanner
	p.portScanner = NewPortScanner()
	portScannerUI := p.portScanner.BuildUI()

	// Create language form with same style as port settings
	// Wrap language select in HBox to match the width behavior of port selection
	languageSelectContainer := container.NewHBox(p.languageSelect)
	p.languageForm = widget.NewForm(
		widget.NewFormItem(locale.T("settings.language"), languageSelectContainer),
	)

	// Create language section label
	languageSectionLabel := widget.NewLabel(locale.T("settings.language_section"))
	languageSectionLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create connection settings label
	connectionLabel := widget.NewLabel(locale.T("settings.connection"))
	connectionLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Combine all elements vertically with separators
	p.content = container.NewVBox(
		languageSectionLabel,
		p.languageForm,
		widget.NewSeparator(),
		connectionLabel,
		portScannerUI,
	)

	return p.content
}

// refreshUI updates the UI after locale change
func (p *SettingsPanel) refreshUI() {
	p.updateLocale()

	// Update language form label text on locale change
	if p.languageForm != nil && len(p.languageForm.Items) > 0 {
		p.languageForm.Items[0].Text = locale.T("settings.language")
	}

	// Refresh port scanner UI if needed
	if p.portScanner != nil {
		p.portScanner.RefreshPorts()
	}

	// Re-apply bold style to section labels
	if p.content != nil && len(p.content.Objects) > 0 {
		if label, ok := p.content.Objects[0].(*widget.Label); ok {
			label.TextStyle = fyne.TextStyle{Bold: true}
		}
		if label, ok := p.content.Objects[3].(*widget.Label); ok {
			label.TextStyle = fyne.TextStyle{Bold: true}
		}
	}

	if p.content != nil {
		p.content.Refresh()
	}
}
