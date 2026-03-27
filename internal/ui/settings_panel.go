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
	languageLabel  *widget.Label
	languageSelect *widget.Select
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
	if p.languageLabel != nil {
		p.languageLabel.SetText(locale.T("settings.language"))
	}
}

// Refresh refreshes the panel UI with current locale
func (p *SettingsPanel) Refresh() {
	p.updateLocale()
}

// buildUI constructs the settings panel UI
func (p *SettingsPanel) buildUI() *fyne.Container {
	// Create language selector
	p.languageLabel = widget.NewLabel(locale.T("settings.language"))

	// Build options for language select
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
	portScanner := NewPortScanner()
	portScannerUI := portScanner.BuildUI()

	// Create settings form using Form layout like in transponder dialog
	form := widget.NewForm(
		widget.NewFormItem(locale.T("settings.language"), p.languageSelect),
		widget.NewFormItem("Порт", portScannerUI),
	)

	p.content = container.NewPadded(form)
	return p.content
}

// refreshUI updates the UI after locale change
func (p *SettingsPanel) refreshUI() {
	p.updateLocale()
	if p.content != nil {
		p.content.Refresh()
	}
}
