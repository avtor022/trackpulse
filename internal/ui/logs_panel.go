package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	header     *widget.Label
	clearBtn   *widget.Button
}

// NewLogsPanel creates a new logs panel
func NewLogsPanel(mainWindow fyne.Window) *LogsPanel {
	p := &LogsPanel{
		mainWindow: mainWindow,
	}

	p.content = p.createContent()
	return p
}

// Content returns the log panel's UI.
func (p *LogsPanel) Content() fyne.CanvasObject {
	return p.content
}

// createContent builds the logs panel content
func (p *LogsPanel) createContent() *fyne.Container {
	// Создание виджета лога
	p.logText = widget.NewRichText()
	p.logScroll = container.NewVScroll(p.logText)

	// Кнопка очистки лога
	p.clearBtn = widget.NewButtonWithIcon(locale.T("logs.clear"), theme.DeleteIcon(), func() {
		p.clearLog()
	})

	// Заголовок лога
	p.header = widget.NewLabel(locale.T("logs.device_log"))
	logHeader := container.NewHBox(
		p.header,
		layout.NewSpacer(),
		p.clearBtn,
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

// AddLog adds a new log entry with timestamp
func (p *LogsPanel) AddLog(message string) {
	if p.logText == nil {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	p.logText.Segments = append(p.logText.Segments, &widget.TextSegment{
		Text: logEntry,
		Style: widget.RichTextStyle{
			Inline: true,
		},
	})

	// Auto-scroll to bottom
	p.logScroll.Offset = fyne.NewPos(0, p.logText.MinSize().Height)
	p.logText.Refresh()
}

// clearLog clears the log display
func (p *LogsPanel) clearLog() {
	p.logText.Segments = []widget.RichTextSegment{}
	p.logText.Refresh()
}

// Refresh updates the panel with new locale strings
func (p *LogsPanel) Refresh() {
	p.header.SetText(locale.T("logs.device_log"))
	p.clearBtn.SetText(locale.T("logs.clear"))
}
