package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// ComboBoxWithEntry представляет комбинированный виджет: Select + Entry
type ComboBoxWithEntry struct {
	selectWidget *widget.Select
	entry        *widget.Entry
	container    *fyne.Container
	onSelected   func(string)
}

// NewComboBoxWithEntry создает новый комбинированный виджет Select с полем ввода
func NewComboBoxWithEntry(options []string, placeHolder string, initialValue string, onSelected func(string)) *ComboBoxWithEntry {
	cbe := &ComboBoxWithEntry{
		onSelected: onSelected,
	}

	// Создаем поле ввода
	cbe.entry = widget.NewEntry()
	cbe.entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		cbe.entry.SetText(initialValue)
	}

	// Создаем выпадающий список
	cbe.selectWidget = widget.NewSelect(options, func(selected string) {
		cbe.entry.SetText(selected)
		if cbe.onSelected != nil {
			cbe.onSelected(selected)
		}
	})
	cbe.selectWidget.SetSelected("") // Устанавливаем пустое значение как placeholder

	// Контейнер с горизонтальной компоновкой
	cbe.container = container.NewHBox(cbe.entry, cbe.selectWidget)

	return cbe
}

// GetContainer возвращает контейнер виджета
func (cbe *ComboBoxWithEntry) GetContainer() *fyne.Container {
	return cbe.container
}

// GetText возвращает текущий текст
func (cbe *ComboBoxWithEntry) GetText() string {
	return cbe.entry.Text
}

// SetText устанавливает текст
func (cbe *ComboBoxWithEntry) SetText(text string) {
	cbe.entry.SetText(text)
}

// GetSelected возвращает выбранное значение из Select
func (cbe *ComboBoxWithEntry) GetSelected() string {
	return cbe.selectWidget.Selected
}

// SetSelected устанавливает выбранное значение в Select
func (cbe *ComboBoxWithEntry) SetSelected(value string) {
	cbe.selectWidget.SetSelected(value)
}

// Disable отключает виджеты
func (cbe *ComboBoxWithEntry) Disable() {
	cbe.entry.Disable()
	cbe.selectWidget.Disable()
}

// Enable включает виджеты
func (cbe *ComboBoxWithEntry) Enable() {
	cbe.entry.Enable()
	cbe.selectWidget.Enable()
}

// CreateComboBox создает простой и надежный комбинированный виджет: поле ввода + Select
func CreateComboBox(options []string, placeHolder string, initialValue string, onSelected func(string)) fyne.CanvasObject {
	// Создаем поле ввода
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		entry.SetText(initialValue)
	}

	// Создаем выпадающий список
	selectWidget := widget.NewSelect(options, func(selected string) {
		entry.SetText(selected)
		if onSelected != nil {
			onSelected(selected)
		}
	})
	selectWidget.SetSelected("") // Устанавливаем пустое значение как placeholder

	// Контейнер с горизонтальной компоновкой
	return container.NewHBox(entry, selectWidget)
}

// CreateComboBoxWithFilter создает комбинированный виджет с фильтрацией опций
func CreateComboBoxWithFilter(options []string, placeHolder string, initialValue string, onSelected func(string)) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeHolder)
	if initialValue != "" {
		entry.SetText(initialValue)
	}

	// Функция для обновления опций Select на основе текста в Entry
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
	selectWidget.SetSelected("") // Устанавливаем пустое значение как placeholder

	// Обработчик изменения текста для фильтрации
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
