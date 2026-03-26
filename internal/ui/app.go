package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jacobsa/go-serial/serial"
	"go.bug.st/serial/enumerator"
	"trackpulse/internal/locale"
	"trackpulse/internal/service"
)

// App represents the main application UI
type App struct {
	fyneApp                fyne.App
	mainWindow             fyne.Window
	competitorService      *service.CompetitorService
	modelService           *service.RCModelService
	settingsService        *service.SettingsService
	competitorModelService *service.CompetitorModelService
	config                 *Config
	tabs                   *container.AppTabs
	competitorPanel        *CompetitorPanel
	modelPanel             *ModelPanel
	competitorModelPanel   *CompetitorModelPanel
	// Serial connection state
	port        io.ReadCloser
	isConnected bool
	statusText  *widget.RichText
	connectBtn  *widget.Button
	portSelect  *widget.Select
	baudEntry   *widget.Entry
}

// Config holds UI configuration
type Config struct {
	Language string
	Title    string
}

// NewApp creates a new TrackPulse application
func NewApp(competitorService *service.CompetitorService, modelService *service.RCModelService, settingsService *service.SettingsService, competitorModelService *service.CompetitorModelService, language string) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("TrackPulse")

	return &App{
		fyneApp:                fyneApp,
		mainWindow:             mainWindow,
		competitorService:      competitorService,
		modelService:           modelService,
		settingsService:        settingsService,
		competitorModelService: competitorModelService,
		config: &Config{
			Language: language,
			Title:    "TrackPulse",
		},
	}
}

// Run starts the application UI
func (a *App) Run() {
	a.mainWindow.SetContent(a.createMainContent())
	a.mainWindow.Resize(fyne.NewSize(1200, 800))
	a.mainWindow.ShowAndRun()
}

// createMainContent builds the main tabbed interface
func (a *App) createMainContent() *container.AppTabs {
	a.tabs = container.NewAppTabs(
		container.NewTabItem(locale.T("tab.monitoring"), a.createMonitoringTab()),
		container.NewTabItem(locale.T("tab.competitors"), a.createCompetitorsTab()),
		container.NewTabItem(locale.T("tab.models"), a.createModelsTab()),
		container.NewTabItem(locale.T("tab.transponders"), a.createTranspondersTab()),
		container.NewTabItem(locale.T("tab.competitions"), a.createCompetitionsTab()),
		container.NewTabItem(locale.T("tab.logs"), a.createLogsTab()),
		container.NewTabItem(locale.T("tab.settings"), a.createSettingsTab()),
	)

	a.tabs.SetTabLocation(container.TabLocationTop)
	return a.tabs
}

// createMonitoringTab creates the Live Monitoring tab
func (a *App) createMonitoringTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("app.welcome"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createCompetitorsTab creates the Competitors management tab
func (a *App) createCompetitorsTab() fyne.CanvasObject {
	a.competitorPanel = NewCompetitorPanel(a.competitorService, a.mainWindow)
	return a.competitorPanel.content
}

// createModelsTab creates the Models management tab
func (a *App) createModelsTab() fyne.CanvasObject {
	a.modelPanel = NewModelPanel(a.modelService, a.mainWindow)
	return a.modelPanel.content
}

// createTranspondersTab creates the Transponders management tab
func (a *App) createTranspondersTab() fyne.CanvasObject {
	a.competitorModelPanel = NewCompetitorModelPanel(a.competitorModelService, a.competitorService, a.modelService, a.mainWindow)
	return a.competitorModelPanel.content
}

// createCompetitionsTab creates the Competitions management tab
func (a *App) createCompetitionsTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("tab.competitions"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createLogsTab creates the Logs viewing tab
func (a *App) createLogsTab() fyne.CanvasObject {
	content := widget.NewLabel(locale.T("tab.logs"))
	content.Alignment = fyne.TextAlignCenter
	return container.NewCenter(content)
}

// createSettingsTab creates the Settings tab
func (a *App) createSettingsTab() fyne.CanvasObject {
	// Create language selector
	languageLabel := widget.NewLabel(locale.T("settings.language"))

	// Build options for language select
	options := make([]string, 0, len(locale.SupportedLocales))
	for _, name := range locale.SupportedLocales {
		options = append(options, name)
	}

	// Find current language name
	currentName := "English"
	for code, name := range locale.SupportedLocales {
		if code == a.config.Language {
			currentName = name
			break
		}
	}

	// Create select without callback first
	languageSelect := widget.NewSelect(options, nil)

	// Set initial value without triggering callback
	languageSelect.SetSelected(currentName)

	// Now set the callback for future changes
	languageSelect.OnChanged = func(selected string) {
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
			a.config.Language = selectedCode

			// Save to database
			if a.settingsService != nil {
				err := a.settingsService.SetLocale(selectedCode)
				if err != nil {
					// Log error but continue with UI update
					fmt.Printf("Failed to save locale: %v\n", err)
				}
			}

			a.refreshUI()
		}
	}

	// Serial port scanner widgets
	portList, defaultPort := a.scanPorts()

	// Create serial connection widgets
	a.portSelect = widget.NewSelect(portList, nil)
	if defaultPort != "" {
		for _, port := range portList {
			portName := a.extractPortName(port)
			if portName == defaultPort {
				a.portSelect.SetSelected(port)
				break
			}
		}
	}

	a.baudEntry = widget.NewEntry()
	a.baudEntry.SetPlaceHolder("9600")
	a.baudEntry.SetText("9600")
	a.baudEntry.SetMinSize(fyne.NewSize(120, 0))

	a.statusText = widget.NewRichText(
		&widget.TextSegment{
			Text: locale.T("settings.serial.disconnected"),
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
	)

	a.connectBtn = widget.NewButton(locale.T("settings.serial.connect"), a.connectSerial)

	// Refresh ports button
	refreshBtn := widget.NewButton(locale.T("settings.serial.refresh_ports"), func() {
		portList, defaultPort := a.scanPorts()
		a.portSelect.Options = portList
		if defaultPort != "" {
			for _, port := range portList {
				portName := a.extractPortName(port)
				if portName == defaultPort {
					a.portSelect.SetSelected(port)
					break
				}
			}
		}
		a.portSelect.Refresh()
	})

	// Serial settings form
	serialForm := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel(locale.T("settings.serial.port")), a.portSelect, refreshBtn),
		container.NewHBox(widget.NewLabel(locale.T("settings.serial.baud")), a.baudEntry),
		container.NewHBox(a.statusText, a.connectBtn),
		widget.NewSeparator(),
	)

	// Language settings form
	languageForm := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(languageLabel, languageSelect),
		widget.NewSeparator(),
	)

	// Main layout
	content := container.NewBorder(
		container.NewVBox(
			languageForm,
			serialForm,
		),
		nil,
		nil,
		nil,
		widget.NewLabel(""),
	)

	return container.NewPadded(content)
}

// refreshUI updates all UI elements with new locale strings
func (a *App) refreshUI() {
	// Refresh tab titles
	for i, tab := range a.tabs.Items {
		switch i {
		case 0:
			tab.Text = locale.T("tab.monitoring")
		case 1:
			tab.Text = locale.T("tab.competitors")
		case 2:
			tab.Text = locale.T("tab.models")
		case 3:
			tab.Text = locale.T("tab.transponders")
		case 4:
			tab.Text = locale.T("tab.competitions")
		case 5:
			tab.Text = locale.T("tab.logs")
		case 6:
			tab.Text = locale.T("tab.settings")
		}
	}

	a.tabs.Refresh()

	// Refresh panels only if they have been created
	if a.competitorPanel != nil {
		a.competitorPanel.Refresh()
	}
	if a.modelPanel != nil {
		a.modelPanel.Refresh()
	}
	if a.competitorModelPanel != nil {
		a.competitorModelPanel.Refresh()
	}

	// Also update settings tab content
	a.tabs.Refresh()
}

// scanPorts scans for available serial ports and returns port list and default Arduino port
func (a *App) scanPorts() ([]string, string) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Println("Ошибка сканирования портов:", err)
		return []string{locale.T("settings.serial.scan_error")}, ""
	}

	var portNames []string
	var arduinoPort string
	var firstPort string

	for _, port := range ports {
		name := port.Name
		if firstPort == "" {
			firstPort = name
		}
		if port.IsUSB {
			info := []string{}
			if port.Product != "" {
				info = append(info, port.Product)
			}
			if len(info) > 0 {
				name += fmt.Sprintf(" (%s)", strings.Join(info, " - "))
			}
			// Check if port is Arduino
			productLower := strings.ToLower(port.Product)
			if strings.Contains(productLower, "arduino") ||
				strings.Contains(productLower, "ch340") ||
				strings.Contains(productLower, "ftdi") ||
				strings.Contains(productLower, "cp210") {
				if arduinoPort == "" {
					arduinoPort = port.Name
				}
			}
		}
		portNames = append(portNames, name)
	}

	if len(portNames) == 0 {
		return []string{locale.T("settings.serial.no_ports")}, ""
	}

	// If Arduino not found, use first port
	if arduinoPort == "" && firstPort != "" {
		arduinoPort = firstPort
	}

	return portNames, arduinoPort
}

// extractPortName extracts the port name from the display string
func (a *App) extractPortName(selected string) string {
	// Remove info in parentheses, leaving only the port name
	if idx := strings.Index(selected, " ("); idx != -1 {
		return selected[:idx]
	}
	// If this is an error message or empty list
	if selected == locale.T("settings.serial.no_ports") || selected == locale.T("settings.serial.scan_error") {
		return ""
	}
	return selected
}

// connectSerial handles the connect/disconnect button click
func (a *App) connectSerial() {
	if a.isConnected {
		a.disconnectSerial()
		return
	}

	portName := a.extractPortName(a.portSelect.Selected)
	if portName == "" {
		a.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: locale.T("settings.serial.select_port_first"),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		a.statusText.Refresh()
		return
	}

	baudRate := uint(9600)
	if a.baudEntry.Text != "" {
		fmt.Sscanf(a.baudEntry.Text, "%d", &baudRate)
	}

	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		a.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: fmt.Sprintf(locale.T("settings.serial.connection_error"), err),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		a.statusText.Refresh()
		return
	}

	a.port = port
	a.isConnected = true
	a.connectBtn.SetText(locale.T("settings.serial.disconnect"))
	a.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: locale.T("settings.serial.connected"),
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameSuccess,
				Inline:    true,
			},
		},
	}
	a.statusText.Refresh()

	// Start reading data in goroutine
	go a.readSerialData()
}

// disconnectSerial disconnects from the serial port
func (a *App) disconnectSerial() {
	if a.port != nil {
		a.port.Close()
		a.port = nil
	}
	a.isConnected = false
	a.connectBtn.SetText(locale.T("settings.serial.connect"))
	a.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: locale.T("settings.serial.disconnected"),
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameError,
				Inline:    true,
			},
		},
	}
	a.statusText.Refresh()
}

// readSerialData reads data from the serial port
func (a *App) readSerialData() {
	scanner := bufio.NewScanner(a.port)
	for scanner.Scan() {
		line := scanner.Text()
		// Data is read but not logged since log panel was removed
	}

	if err := scanner.Err(); err != nil && a.isConnected {
		a.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: fmt.Sprintf(locale.T("settings.serial.read_error"), err),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		a.statusText.Refresh()
	}
}
