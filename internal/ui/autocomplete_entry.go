package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// ComboBoxWithEntry represents a combined widget: Select + Entry
type ComboBoxWithEntry struct {
	selectWidget *widget.Select
	entry        *widget.Entry
	container    *fyne.Container
	onSelected   func(string)
}

// NewComboBoxWithEntry creates a new combined Select widget with input field
func NewComboBoxWithEntry(options []string, placeHolder string, initialValue string, onSelected func(string)) *ComboBoxWithEntry {
	cbe := &ComboBoxWithEntry{
		onSelected: onSelected,
	}

	// Create input field
	cbe.entry = widget.NewEntry()
	cbe.entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		cbe.entry.SetText(initialValue)
	}

	// Create dropdown list
	cbe.selectWidget = widget.NewSelect(options, func(selected string) {
		cbe.entry.SetText(selected)
		if cbe.onSelected != nil {
			cbe.onSelected(selected)
		}
	})
	cbe.selectWidget.SetSelected("") // Set empty value as placeholder

	// Container with horizontal layout
	cbe.container = container.NewHBox(cbe.entry, cbe.selectWidget)

	return cbe
}

// GetContainer returns the widget container
func (cbe *ComboBoxWithEntry) GetContainer() *fyne.Container {
	return cbe.container
}

// GetText returns the current text
func (cbe *ComboBoxWithEntry) GetText() string {
	return cbe.entry.Text
}

// SetText sets the text
func (cbe *ComboBoxWithEntry) SetText(text string) {
	cbe.entry.SetText(text)
}

// GetSelected returns the selected value from Select
func (cbe *ComboBoxWithEntry) GetSelected() string {
	return cbe.selectWidget.Selected
}

// SetSelected sets the selected value in Select
func (cbe *ComboBoxWithEntry) SetSelected(value string) {
	cbe.selectWidget.SetSelected(value)
}

// Disable disables the widgets
func (cbe *ComboBoxWithEntry) Disable() {
	cbe.entry.Disable()
	cbe.selectWidget.Disable()
}

// Enable enables the widgets
func (cbe *ComboBoxWithEntry) Enable() {
	cbe.entry.Enable()
	cbe.selectWidget.Enable()
}

// CreateComboBox creates a simple and reliable combined widget: input field + Select
func CreateComboBox(options []string, placeHolder string, initialValue string, onSelected func(string)) fyne.CanvasObject {
	// Create input field
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		entry.SetText(initialValue)
	}

	// Create dropdown list
	selectWidget := widget.NewSelect(options, func(selected string) {
		entry.SetText(selected)
		if onSelected != nil {
			onSelected(selected)
		}
	})
	selectWidget.SetSelected("") // Set empty value as placeholder

	// Container with horizontal layout
	return container.NewHBox(entry, selectWidget)
}

// CreateComboBoxWithFilter creates a combined widget with filtered options
func CreateComboBoxWithFilter(options []string, placeHolder string, initialValue string, onSelected func(string)) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		entry.SetText(initialValue)
	}

	// Function to update Select options based on Entry text
	var selectWidget *widget.Select
	updateOptions := func(text string) {
		if text == "" {
			selectWidget.Options = options
		} else {
			var filtered []string
			textLower := strings.ToLower(text)
			for _, opt := range options {
				if strings.Contains(strings.ToLower(opt), textLower) {
					filtered = append(filtered, opt)
				}
			}
			selectWidget.Options = filtered
		}
		selectWidget.Refresh()
	}

	selectWidget = widget.NewSelect(options, func(selected string) {
		entry.SetText(selected)
		if onSelected != nil {
			onSelected(selected)
		}
	})
	selectWidget.SetSelected("") // Set empty value as placeholder

	// Text change handler for filtering
	entry.OnChanged = func(text string) {
		updateOptions(text)
	}

	entry.OnSubmitted = func(text string) {
		if onSelected != nil {
			onSelected(text)
		}
	}

	return container.NewHBox(entry, selectWidget)
}
