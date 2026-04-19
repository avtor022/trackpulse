package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jacobsa/go-serial/serial"
	"go.bug.st/serial/enumerator"
	"trackpulse/internal/locale"
)

// LogsPanel represents the logs panel UI
type LogsPanel struct {
	content     fyne.CanvasObject
	logText     *widget.RichText
	logScroll   *container.Scroll
	statusText  *widget.RichText
	connectBtn  *widget.Button
	portSelect  *widget.Select
	baudEntry   *widget.Entry
	port        io.ReadCloser
	isConnected bool
	mainWindow  fyne.Window
}

// NewLogsPanel creates a new logs panel
func NewLogsPanel(mainWindow fyne.Window) *LogsPanel {
	p := &LogsPanel{
		isConnected: false,
		mainWindow:  mainWindow,
	}

	p.content = p.createContent()
	return p
}

// scanPorts scans for available serial ports
func scanPorts() ([]string, string) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Println("Ошибка сканирования портов:", err)
		return []string{locale.T("logs.error_scanning")}, ""
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
			// Проверяем, является ли порт Arduino или подобным устройством
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
		return []string{locale.T("logs.no_ports")}, ""
	}

	// Если Arduino не найден, используем первый порт
	if arduinoPort == "" && firstPort != "" {
		arduinoPort = firstPort
	}

	return portNames, arduinoPort
}

// extractPortName extracts the port name from the select option
func extractPortName(selected string) string {
	if idx := strings.Index(selected, " ("); idx != -1 {
		return selected[:idx]
	}
	if selected == locale.T("logs.no_ports") || selected == locale.T("logs.error_scanning") {
		return ""
	}
	return selected
}

// createContent builds the logs panel content
func (p *LogsPanel) createContent() fyne.CanvasObject {
	// Сканирование доступных портов при запуске
	portList, defaultPort := scanPorts()

	// Создание виджета лога
	p.logText = widget.NewRichText()
	p.logScroll = container.NewVScroll(p.logText)

	// Статус подключения
	p.statusText = widget.NewRichText(
		&widget.TextSegment{
			Text: locale.T("logs.status_disconnected"),
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
	)

	// Выпадающий список портов
	p.portSelect = widget.NewSelect(portList, nil)
	if defaultPort != "" {
		for _, port := range portList {
			portName := extractPortName(port)
			if portName == defaultPort {
				p.portSelect.SetSelected(port)
				break
			}
		}
	}

	// Поле ввода скорости (боуд)
	p.baudEntry = widget.NewEntry()
	p.baudEntry.SetPlaceHolder("9600")
	p.baudEntry.SetText("9600")

	// Кнопка подключения
	p.connectBtn = widget.NewButton(locale.T("logs.connect"), p.connect)

	// Кнопка обновления портов
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		p.refreshPorts()
	})

	// Кнопка очистки лога
	clearBtn := widget.NewButtonWithIcon(locale.T("logs.clear"), theme.DeleteIcon(), func() {
		p.clearLog()
	})

	// Панель настроек подключения
	settingsForm := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(locale.T("logs.port"), p.portSelect),
			widget.NewFormItem(locale.T("logs.baud_rate"), p.baudEntry),
		),
		container.NewHBox(p.connectBtn, refreshBtn, clearBtn),
		p.statusText,
	)

	// Заголовок лога
	logHeader := widget.NewLabel(locale.T("logs.device_log"))
	logHeader.TextStyle = fyne.TextStyle{Bold: true}

	// Основная компоновка
	content := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			settingsForm,
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

// connect handles the connect/disconnect action
func (p *LogsPanel) connect() {
	if p.isConnected {
		p.disconnect()
		return
	}

	portName := extractPortName(p.portSelect.Selected)
	if portName == "" {
		p.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: locale.T("logs.select_port"),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		p.statusText.Refresh()
		return
	}

	baudRate := uint(9600)
	if p.baudEntry.Text != "" {
		fmt.Sscanf(p.baudEntry.Text, "%d", &baudRate)
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
		p.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: fmt.Sprintf(locale.T("logs.connect_error"), err),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		p.statusText.Refresh()
		return
	}

	p.port = port
	p.isConnected = true
	p.connectBtn.SetText(locale.T("logs.disconnect"))
	p.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: locale.T("logs.connected"),
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameSuccess,
				Inline:    true,
			},
		},
	}
	p.statusText.Refresh()

	// Запуск чтения данных в горутине
	go p.readData()
}

// disconnect handles disconnection from the device
func (p *LogsPanel) disconnect() {
	if p.port != nil {
		p.port.Close()
		p.port = nil
	}
	p.isConnected = false
	p.connectBtn.SetText(locale.T("logs.connect"))
	p.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: locale.T("logs.disconnected"),
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameError,
				Inline:    true,
			},
		},
	}
	p.statusText.Refresh()
}

// readData reads data from the serial port
func (p *LogsPanel) readData() {
	scanner := bufio.NewScanner(p.port)
	for scanner.Scan() {
		line := scanner.Text()
		timestamp := time.Now().Format("15:04:05")

		// Форматирование строки лога (выделение RFID данных)
		formattedLine := p.formatRFIDLine(line)

		// Добавляем новую строку в лог
		p.logText.Segments = append(p.logText.Segments, &widget.TextSegment{
			Text:  fmt.Sprintf("[%s] %s\n", timestamp, formattedLine),
			Style: widget.RichTextStyleInline,
		})
		p.logText.Refresh()

		// Прокрутка вниз
		p.logScroll.ScrollToBottom()

		// Ограничиваем количество строк в логе
		if len(p.logText.Segments) > 200 {
			p.logText.Segments = p.logText.Segments[len(p.logText.Segments)-200:]
		}
	}

	if err := scanner.Err(); err != nil && p.isConnected {
		p.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: fmt.Sprintf(locale.T("logs.read_error"), err),
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					Inline:    true,
				},
			},
		}
		p.statusText.Refresh()
	}
}

// formatRFIDLine formats RFID tag reading lines
func (p *LogsPanel) formatRFIDLine(line string) string {
	// Простая эвристика для выделения UID карт
	// Обычно RFID-RC522 выводит UID в формате шестнадцатеричных чисел
	if strings.Contains(strings.ToLower(line), "uid") ||
		strings.Contains(strings.ToLower(line), "tag") ||
		strings.Contains(strings.ToLower(line), "card") ||
		p.isHexLine(line) {
		return fmt.Sprintf("** RFID: %s **", line)
	}
	return line
}

// isHexLine checks if a line contains hexadecimal values
func (p *LogsPanel) isHexLine(line string) bool {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}
	// Проверяем, состоит ли строка из шестнадцатеричных значений
	for _, field := range fields {
		for _, ch := range field {
			if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F') || (ch >= 'a' && ch <= 'f')) {
				return false
			}
		}
	}
	return true
}

// refreshPorts refreshes the list of available ports
func (p *LogsPanel) refreshPorts() {
	portList, defaultPort := scanPorts()
	p.portSelect.Options = portList
	if defaultPort != "" {
		for _, port := range portList {
			portName := extractPortName(port)
			if portName == defaultPort {
				p.portSelect.SetSelected(port)
				break
			}
		}
	}
	p.portSelect.Refresh()
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
