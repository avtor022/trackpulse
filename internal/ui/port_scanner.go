package ui

import (
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jacobsa/go-serial/serial"
	"go.bug.st/serial/enumerator"
)

// PortScanner handles serial port scanning and connection UI
type PortScanner struct {
	port        io.ReadCloser
	isConnected bool
	statusText  *widget.RichText
	connectBtn  *widget.Button
	portSelect  *widget.Select
	baudEntry   *widget.Entry
	refreshBtn  *widget.Button
}

// scanPorts scans for available serial ports
func scanPorts() ([]string, string) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Println("Ошибка сканирования портов:", err)
		return []string{"Ошибка сканирования"}, ""
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
			// Проверяем, является ли порт Arduino
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
		return []string{"Порты не найдены"}, ""
	}

	// Если Arduino не найден, используем первый порт
	if arduinoPort == "" && firstPort != "" {
		arduinoPort = firstPort
	}

	return portNames, arduinoPort
}

// extractPortName extracts the port name from the selected option
func extractPortName(selected string) string {
	// Убираем информацию в скобках, оставляя только имя порта
	if idx := strings.Index(selected, " ("); idx != -1 {
		return selected[:idx]
	}
	// Если это сообщение об ошибке или пустой список
	if selected == "Порты не найдены" || selected == "Ошибка сканирования" {
		return ""
	}
	return selected
}

// NewPortScanner creates a new PortScanner instance
func NewPortScanner() *PortScanner {
	return &PortScanner{
		isConnected: false,
	}
}

// connect handles the connection/disconnection logic
func (p *PortScanner) connect() {
	if p.isConnected {
		p.disconnect()
		return
	}

	portName := extractPortName(p.portSelect.Selected)
	if portName == "" {
		p.statusText.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text: "Выберите порт из списка",
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
				Text: fmt.Sprintf("Ошибка подключения: %v", err),
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
	p.connectBtn.SetText("Отключиться")
	p.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: "Подключено",
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameSuccess,
				Inline:    true,
			},
		},
	}
	p.statusText.Refresh()
}

// disconnect handles disconnection from the serial port
func (p *PortScanner) disconnect() {
	if p.port != nil {
		p.port.Close()
		p.port = nil
	}
	p.isConnected = false
	p.connectBtn.SetText("Подключиться")
	p.statusText.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: "Отключено",
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameError,
				Inline:    true,
			},
		},
	}
	p.statusText.Refresh()
}

// BuildUI creates the port scanner UI components
func (p *PortScanner) BuildUI() fyne.CanvasObject {
	// Сканирование доступных портов при запуске
	portList, defaultPort := scanPorts()

	// Создание виджетов
	p.statusText = widget.NewRichText(
		&widget.TextSegment{
			Text: "Статус: Отключено",
			Style: widget.RichTextStyle{
				Inline: true,
			},
		},
	)

	// Выпадающий список портов
	p.portSelect = widget.NewSelect(portList, nil)
	if defaultPort != "" {
		// Устанавливаем порт Arduino (или первый порт) по умолчанию
		for _, port := range portList {
			portName := extractPortName(port)
			if portName == defaultPort {
				p.portSelect.SetSelected(port)
				break
			}
		}
	}

	p.baudEntry = widget.NewEntry()
	p.baudEntry.SetPlaceHolder("9600")
	p.baudEntry.SetText("9600")

	p.connectBtn = widget.NewButton("Подключиться", p.connect)

	p.refreshBtn = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
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
	})

	// Панель настроек подключения
	settingsForm := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabel("Порт:"),
			p.portSelect,
			p.refreshBtn,
		),
		container.NewHBox(
			widget.NewLabel("Скорость (бод):"),
			p.baudEntry,
		),
		container.NewHBox(
			p.connectBtn,
			p.statusText,
		),
		widget.NewSeparator(),
	)

	return settingsForm
}

// RefreshPorts refreshes the list of available ports
func (p *PortScanner) RefreshPorts() {
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
