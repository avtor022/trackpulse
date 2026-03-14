package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// AutoCompleteEntry представляет поле ввода с автодополнением
type AutoCompleteEntry struct {
	entry       *widget.Entry
	popup       *widget.PopUp
	options     []string
	filtered    []string
	onSelected  func(string)
	window      fyne.Window
}

// NewAutoCompleteEntry создает новое поле ввода с автодополнением
func NewAutoCompleteEntry(options []string, window fyne.Window, onSelected func(string)) *AutoCompleteEntry {
	ace := &AutoCompleteEntry{
		entry:      widget.NewEntry(),
		options:    options,
		filtered:   []string{},
		onSelected: onSelected,
		window:     window,
	}

	ace.entry.OnChanged = func(text string) {
		ace.filterOptions(text)
	}

	ace.entry.OnSubmitted = func(text string) {
		ace.hidePopup()
		if ace.onSelected != nil {
			ace.onSelected(text)
		}
	}

	return ace
}

// filterOptions фильтрует опции на основе введенного текста
func (ace *AutoCompleteEntry) filterOptions(text string) {
	if text == "" {
		ace.filtered = []string{}
		ace.hidePopup()
		return
	}

	ace.filtered = []string{}
	textLower := strings.ToLower(text)
	
	for _, opt := range ace.options {
		if strings.Contains(strings.ToLower(opt), textLower) {
			ace.filtered = append(ace.filtered, opt)
		}
	}

	if len(ace.filtered) > 0 {
		ace.showPopup()
	} else {
		ace.hidePopup()
	}
}

// showPopup отображает выпадающий список с вариантами
func (ace *AutoCompleteEntry) showPopup() {
	if ace.popup != nil {
		ace.popup.Hide()
	}

	list := widget.NewList(
		func() int { return len(ace.filtered) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(ace.filtered[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selected := ace.filtered[id]
		ace.entry.SetText(selected)
		ace.hidePopup()
		if ace.onSelected != nil {
			ace.onSelected(selected)
		}
	}

	ace.popup = widget.NewPopUp(list, ace.window.Canvas())
	
	// Позиционируем popup под полем ввода
	entryPos := fyne.NewPos(0, ace.entry.Size().Height)
	ace.popup.ShowAtPosition(ace.entry.Position().Add(entryPos))
}

// hidePopup скрывает выпадающий список
func (ace *AutoCompleteEntry) hidePopup() {
	if ace.popup != nil {
		ace.popup.Hide()
	}
}

// GetWidget возвращает основной виджет Entry
func (ace *AutoCompleteEntry) GetWidget() *widget.Entry {
	return ace.entry
}

// GetText возвращает текущий текст
func (ace *AutoCompleteEntry) GetText() string {
	return ace.entry.Text
}

// SetText устанавливает текст
func (ace *AutoCompleteEntry) SetText(text string) {
	ace.entry.SetText(text)
}

// SetPlaceHolder устанавливает текст подсказки
func (ace *AutoCompleteEntry) SetPlaceHolder(text string) {
	ace.entry.SetPlaceHolder(text)
}

// Disable disables the entry
func (ace *AutoCompleteEntry) Disable() {
	ace.entry.Disable()
}

// Enable enables the entry
func (ace *AutoCompleteEntry) Enable() {
	ace.entry.Enable()
}

// FocusGained is called when this item gained the focus
func (ace *AutoCompleteEntry) FocusGained() {
	ace.entry.FocusGained()
	if ace.entry.Text != "" {
		ace.filterOptions(ace.entry.Text)
	}
}

// FocusLost is called when this item lost the focus
func (ace *AutoCompleteEntry) FocusLost() {
	ace.entry.FocusLost()
	// Небольшая задержка перед скрытием popup, чтобы успеть выбрать элемент
	time.AfterFunc(200*time.Millisecond, func() {
		ace.hidePopup()
	})
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (ace *AutoCompleteEntry) Tapped(_ *fyne.PointEvent) {
	ace.entry.Tapped(nil)
}

// MinSize returns the size that this widget should not shrink below
func (ace *AutoCompleteEntry) MinSize() fyne.Size {
	return ace.entry.MinSize()
}

// Resize sets a new size for a widget.
func (ace *AutoCompleteEntry) Resize(size fyne.Size) {
	ace.entry.Resize(size)
}

// Position returns the position of this widget.
func (ace *AutoCompleteEntry) Position() fyne.Position {
	return ace.entry.Position()
}

// Size returns the current size of this widget.
func (ace *AutoCompleteEntry) Size() fyne.Size {
	return ace.entry.Size()
}

// Refresh triggers a redraw of this widget.
func (ace *AutoCompleteEntry) Refresh() {
	ace.entry.Refresh()
}

// Visible returns whether this widget is visible or not.
func (ace *AutoCompleteEntry) Visible() bool {
	return ace.entry.Visible()
}

// Hide hides this widget.
func (ace *AutoCompleteEntry) Hide() {
	ace.entry.Hide()
	ace.hidePopup()
}

// Show shows this widget.
func (ace *AutoCompleteEntry) Show() {
	ace.entry.Show()
}

// CreateComboBoxWithAutoComplete создает комбинированный виджет для выбора бренда с автодополнением
func CreateComboBoxWithAutoComplete(options []string, window fyne.Window, placeHolder string, initialValue string, onSelected func(string)) fyne.CanvasObject {
	// Создаем контейнер с горизонтальной компоновкой
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeHolder)
	
	if initialValue != "" {
		entry.SetText(initialValue)
	}

	// Кнопка для открытия полного списка
	dropdownBtn := widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
		showFullDropdown(options, window, entry, onSelected)
	})

	// Контейнер для поля ввода и кнопки
	content := container.NewHBox(
		container.NewStack(entry), // Stack позволяет размещать элементы друг над другом
		dropdownBtn,
	)

	// Обработчик изменения текста для фильтрации
	var popup *widget.PopUp
	
	entry.OnChanged = func(text string) {
		if popup != nil {
			popup.Hide()
		}
		
		if text == "" {
			return
		}

		// Фильтруем опции
		var filtered []string
		textLower := strings.ToLower(text)
		
		for _, opt := range options {
			if strings.Contains(strings.ToLower(opt), textLower) {
				filtered = append(filtered, opt)
			}
		}

		if len(filtered) > 0 {
			showFilterPopup(entry, filtered, window, func(selected string) {
				entry.SetText(selected)
				if popup != nil {
					popup.Hide()
				}
				if onSelected != nil {
					onSelected(selected)
				}
			})
		}
	}

	entry.OnSubmitted = func(text string) {
		if popup != nil {
			popup.Hide()
		}
		if onSelected != nil {
			onSelected(text)
		}
	}

	// Сохраняем ссылку на popup в замыкании через переменную окружения
	_ = popup

	return content
}

// showFilterPopup показывает всплывающее окно с отфильтрованными вариантами
func showFilterPopup(entry *widget.Entry, options []string, window fyne.Window, onSelect func(string)) {
	list := widget.NewList(
		func() int { return len(options) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(options[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selected := options[id]
		onSelect(selected)
	}

	popup := widget.NewPopUp(list, window.Canvas())
	
	// Вычисляем позицию под полем ввода
	entryPos := entry.Position()
	entrySize := entry.Size()
	
	popupPos := fyne.NewPos(entryPos.X, entryPos.Y+entrySize.Height+5)
	popup.ShowAtPosition(popupPos)
}

// showFullDropdown показывает полный список всех опций
func showFullDropdown(options []string, window fyne.Window, entry *widget.Entry, onSelected func(string)) {
	list := widget.NewList(
		func() int { return len(options) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Template")
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(options[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selected := options[id]
		entry.SetText(selected)
		if onSelected != nil {
			onSelected(selected)
		}
	}

	popup := widget.NewPopUp(list, window.Canvas())
	
	entryPos := entry.Position()
	entrySize := entry.Size()
	
	popupPos := fyne.NewPos(entryPos.X, entryPos.Y+entrySize.Height+5)
	popup.ShowAtPosition(popupPos)
}
