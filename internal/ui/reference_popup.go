package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
)

// ReferenceItem represents an item that can be managed in a reference popup
type ReferenceItem struct {
	Name string
}

// ReferencePopupConfig holds configuration for a reference popup
type ReferencePopupConfig struct {
	Title          string                          // Dialog title key for locale
	AddTitle       string                          // Add dialog title key for locale
	AddLabel       string                          // Add dialog label key for locale
	AddPlaceholder string                          // Add dialog placeholder key for locale
	DeleteMessage  string                          // Delete confirmation message key for locale
	NewErrorExists string                          // Error message key when item already exists
	EnterNameInfo  string                          // Info key for entering name error
	NewItemOption  string                          // Text for "add new" option in list
	GetAllFunc     func() ([]ReferenceItem, error) // Function to get all items
	AddFunc        func(string) error              // Function to add new item
	DeleteFunc     func(string) error              // Function to delete item
	OnItemSelected func(string)                    // Callback when item is selected
	UpdateOptions  func([]string)                  // Callback to update select options
}

// ReferencePopupManager manages reference popups for brands, scales, and model types
type ReferencePopupManager struct {
	window        fyne.Window
	config        ReferencePopupConfig
	items         []string
	newItemOption string
	onSelect      func(string)
	updateOpts    func([]string)
}

// NewReferencePopupManager creates a new reference popup manager
func NewReferencePopupManager(window fyne.Window, config ReferencePopupConfig, initialItems []string, newItemOption string, onSelect func(string), updateOpts func([]string)) *ReferencePopupManager {
	return &ReferencePopupManager{
		window:        window,
		config:        config,
		items:         initialItems,
		newItemOption: newItemOption,
		onSelect:      onSelect,
		updateOpts:    updateOpts,
	}
}

// ShowPopup displays the reference selection popup with add/delete functionality
func (m *ReferencePopupManager) ShowPopup(mainDialog dialog.Dialog, currentDialog *dialog.Dialog, setCurrentDialog func(dialog.Dialog)) {
	// Hide current dialog if any
	if *currentDialog != nil {
		(*currentDialog).Hide()
	}

	// Create container for items
	itemContainer := container.NewVBox()

	// Add existing items with delete buttons
	for _, item := range m.items {
		itemName := item

		// Create delete button
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			// Show confirmation dialog
			dialog.ShowConfirm(
				locale.T("dialog.delete.title"),
				fmt.Sprintf(locale.T(m.config.DeleteMessage), itemName),
				func(confirmed bool) {
					if confirmed {
						// Delete item from database
						if err := m.config.DeleteFunc(itemName); err != nil {
							dialog.ShowError(err, m.window)
							return
						}

						// Remove from items slice
						newItems := []string{}
						for _, i := range m.items {
							if i != itemName {
								newItems = append(newItems, i)
							}
						}
						m.items = newItems

						// Refresh the popup
						m.ShowPopup(mainDialog, currentDialog, setCurrentDialog)

						// Update main dialog select if needed
						if m.onSelect != nil {
							m.onSelect("")
						}
					}
				},
				m.window,
			)
		})
		deleteBtn.Importance = widget.DangerImportance

		// Create button with label
		selectBtn := widget.NewButton(item, func() {
			// Select this item
			if m.onSelect != nil {
				m.onSelect(itemName)
			}
			// Hide popup and return to main dialog
			if *currentDialog != nil {
				(*currentDialog).Hide()
			}
			if mainDialog != nil {
				mainDialog.Show()
			}
		})
		selectBtn.Alignment = widget.ButtonAlignLeading
		selectBtn.Importance = widget.LowImportance

		// Create horizontal container with select button and delete button
		selectBtn.Importance = widget.MediumImportance
		deleteBtn.Importance = widget.DangerImportance

		// Place both buttons on the same row: select button expands, delete button stays compact
		itemRow := container.NewBorder(nil, nil, nil, deleteBtn, selectBtn)
		itemContainer.Add(itemRow)
	}

	// Add "Add new" option
	addNewBtn := widget.NewButton(m.newItemOption, func() {
		// Show dialog to add new item
		if mainDialog != nil {
			mainDialog.Hide()
		}
		if *currentDialog != nil {
			(*currentDialog).Hide()
		}

		newEntry := widget.NewEntry()
		newEntry.SetPlaceHolder(locale.T(m.config.AddPlaceholder))

		// Create label and input field vertically for better width
		label := widget.NewLabel(locale.T(m.config.AddLabel))
		entryContainer := container.NewVBox(label, newEntry)

		newItemDialog := dialog.NewCustomWithoutButtons(locale.T(m.config.AddTitle), entryContainer, m.window)

		cancelBtn := widget.NewButton(locale.T("common.cancel"), func() {
			newItemDialog.Hide()
			// Return to main dialog
			if mainDialog != nil {
				mainDialog.Show()
			}
		})

		saveBtn := widget.NewButton(locale.T("common.save"), func() {
			newItemName := strings.TrimSpace(newEntry.Text)
			if newItemName == "" {
				dialog.ShowError(fmt.Errorf(locale.T(m.config.EnterNameInfo)), m.window)
				return
			}

			// Check if item already exists
			for _, i := range m.items {
				if strings.EqualFold(i, newItemName) {
					dialog.ShowError(fmt.Errorf(locale.T(m.config.NewErrorExists), newItemName), m.window)
					return
				}
			}

			// Add new item to reference table
			if err := m.config.AddFunc(newItemName); err != nil {
				dialog.ShowError(err, m.window)
				return
			}

			// Update item list
			m.items = append(m.items, newItemName)
			allOptions := append(m.items, m.newItemOption)
			if m.updateOpts != nil {
				m.updateOpts(allOptions)
			}

			newItemDialog.Hide()
			// Return to main dialog with new item selected
			if mainDialog != nil {
				mainDialog.Show()
			}
			if m.onSelect != nil {
				m.onSelect(newItemName)
			}
		})

		newItemDialog.SetButtons([]fyne.CanvasObject{cancelBtn, saveBtn})

		// Increase dialog width for new item
		parentSize := m.window.Canvas().Size()
		dialogWidth := parentSize.Width * 0.6
		if dialogWidth < 700 {
			dialogWidth = 700
		}
		newItemDialog.Resize(fyne.NewSize(dialogWidth, newItemDialog.MinSize().Height))

		setCurrentDialog(newItemDialog)
		newItemDialog.Show()
	})
	itemContainer.Add(addNewBtn)

	// Create popup dialog
	popup := dialog.NewCustomWithoutButtons(locale.T("common.select_one"), itemContainer, m.window)

	// Add close button
	closeBtn := widget.NewButton(locale.T("common.close"), func() {
		popup.Hide()
		setCurrentDialog(nil)
		// Return to main dialog
		if mainDialog != nil {
			mainDialog.Show()
		}
	})

	popup.SetButtons([]fyne.CanvasObject{closeBtn})

	// Resize popup
	parentSize := m.window.Canvas().Size()
	popupWidth := parentSize.Width * 0.4
	if popupWidth < 400 {
		popupWidth = 400
	}
	popup.Resize(fyne.NewSize(popupWidth, popup.MinSize().Height))

	setCurrentDialog(popup)
	popup.Show()
}

// GetItems returns the current list of items
func (m *ReferencePopupManager) GetItems() []string {
	return m.items
}

// RefreshItems reloads the items from the data source
func (m *ReferencePopupManager) RefreshItems() error {
	allItems, err := m.config.GetAllFunc()
	if err != nil {
		return err
	}

	m.items = make([]string, len(allItems))
	for i, item := range allItems {
		m.items[i] = item.Name
	}

	allOptions := append(m.items, m.newItemOption)
	if m.updateOpts != nil {
		m.updateOpts(allOptions)
	}

	return nil
}

// ShowPopupWithoutAddDelete displays the reference selection popup WITHOUT add/delete buttons
// This is useful for selecting existing items like competitors where modification is not allowed
func (m *ReferencePopupManager) ShowPopupWithoutAddDelete(mainDialog dialog.Dialog, currentDialog *dialog.Dialog, setCurrentDialog func(dialog.Dialog)) {
	// Hide current dialog if any
	if *currentDialog != nil {
		(*currentDialog).Hide()
	}

	// Create container for items
	itemContainer := container.NewVBox()

	// Add existing items as simple select buttons (no delete buttons)
	for _, item := range m.items {
		itemName := item

		// Create button with label
		selectBtn := widget.NewButton(item, func() {
			// Select this item
			if m.onSelect != nil {
				m.onSelect(itemName)
			}
			// Hide popup and return to main dialog
			if *currentDialog != nil {
				(*currentDialog).Hide()
			}
			if mainDialog != nil {
				mainDialog.Show()
			}
		})
		selectBtn.Alignment = widget.ButtonAlignLeading
		selectBtn.Importance = widget.MediumImportance

		itemContainer.Add(selectBtn)
	}

	// Create popup dialog
	popup := dialog.NewCustomWithoutButtons(locale.T("common.select_one"), itemContainer, m.window)

	// Add close button
	closeBtn := widget.NewButton(locale.T("common.close"), func() {
		popup.Hide()
		setCurrentDialog(nil)
		// Return to main dialog
		if mainDialog != nil {
			mainDialog.Show()
		}
	})

	popup.SetButtons([]fyne.CanvasObject{closeBtn})

	// Resize popup
	parentSize := m.window.Canvas().Size()
	popupWidth := parentSize.Width * 0.4
	if popupWidth < 400 {
		popupWidth = 400
	}
	popup.Resize(fyne.NewSize(popupWidth, popup.MinSize().Height))

	setCurrentDialog(popup)
	popup.Show()
}
