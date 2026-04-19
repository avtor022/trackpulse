package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
)

// LogsPanel represents the logs panel UI
type LogsPanel struct {
	content    *fyne.Container
	logText    *widget.RichText
	logScroll  *container.Scroll
	mainWindow fyne.Window
}

// NewLogsPanel creates a new logs panel
func NewLogsPanel(mainWindow fyne.Window) *LogsPanel {
	p := &LogsPanel{
		mainWindow: mainWindow,
	}

	p.content = p.createContent()
	return p
}

// createContent builds the logs panel content
func (p *LogsPanel) createContent() *fyne.Container {
	// Создание виджета лога
	p.logText = widget.NewRichText()
	p.logScroll = container.NewVScroll(p.logText)

	// Кнопка очистки лога
	clearBtn := widget.NewButtonWithIcon(locale.T("logs.clear"), theme.DeleteIcon(), func() {
		p.clearLog()
	})

	// Заголовок лога
	logHeader := container.NewHBox(
		widget.NewLabel(locale.T("logs.device_log")),
		container.NewSpacer(),
		clearBtn,
	)
	logHeader.Objects[0].(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}

	// Основная компоновка
	content := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			logHeader,
		),
		nil,
		nil,
		nil,
		p.logScroll,
	)

	return content
}

// clearLog clears the log display
func (p *LogsPanel) clearLog() {
	p.logText.Segments = []widget.RichTextSegment{}
	p.logText.Refresh()
}

// Refresh updates the panel with new locale strings
func (p *LogsPanel) Refresh() {
	p.content = p.createContent()
}
