package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
)

// EntityConfig holds configuration for different entity types (brand, scale, model type)
type EntityConfig struct {
	EntityType        string // e.g., "brand", "scale", "model_type"
	ExistingItems     *[]string
	NewOptionText     string
	AddPlaceholderKey string
	AddLabelKey       string
	AddTitleKey       string
	DeleteMessageKey  string
	ErrorExistsKey    string
	EnterNameInfoKey  string
	ServiceAddFunc    func(name string) error
	ServiceDeleteFunc func(name string) error
}

// ShowEntityPopup shows a popup dialog for managing entities (brands, scales, model types)
// It displays existing items with delete buttons and an option to add new items
func ShowEntityPopup(window fyne.Window, config EntityConfig, currentDialog *dialog.Dialog, mainDialog dialog.Dialog, selectWidget **widget.Select, updateSelectOptions func([]string)) {
	if *currentDialog != nil {
		(*currentDialog).Hide()
	}

	// Create container for entity options
	entityContainer := container.NewVBox()

	// Add existing items with delete buttons
	for _, item := range *config.ExistingItems {
		itemName := item
		
		// Create label that can be clicked to select the item
		itemLabel := widget.NewLabel(item)
		itemLabel.OnTapped = func() {
			// Select this item and close popup
			if mainDialog != nil {
				mainDialog.Hide()
			}
			if *currentDialog != nil {
				(*currentDialog).Hide()
				*currentDialog = nil
			}
			// Update the select widget with selected item
			(*selectWidget).SetSelected(itemName)
			// Return to main dialog
			if mainDialog != nil {
				mainDialog.Show()
			}
		}
		
		// Create horizontal layout for item name and delete button
		itemRow := container.NewHBox(
			itemLabel,
			layout.NewSpacer(),
		)

		// Create delete button
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			// Show confirmation dialog
			dialog.ShowConfirm(
				locale.T("dialog.delete.title"),
				fmt.Sprintf(locale.T(config.DeleteMessageKey), itemName),
				func(confirmed bool) {
					if confirmed {
						// Delete item from database
						if err := config.ServiceDeleteFunc(itemName); err != nil {
							dialog.ShowError(err, window)
							return
						}

						// Remove from existing items slice
						newExistingItems := []string{}
						for _, i := range *config.ExistingItems {
							if i != itemName {
								newExistingItems = append(newExistingItems, i)
							}
						}
						*config.ExistingItems = newExistingItems

						// Refresh the popup
						ShowEntityPopup(window, config, currentDialog, mainDialog, selectWidget, updateSelectOptions)

						// Update main dialog select if needed
						if (*selectWidget).Selected == itemName {
							(*selectWidget).SetSelected("")
						}
					}
				},
				window,
			)
		})
		deleteBtn.Importance = widget.DangerImportance

		itemRow.Add(deleteBtn)
		entityContainer.Add(itemRow)
	}

	// Add "Add new" option
	addNewBtn := widget.NewButton(config.NewOptionText, func() {
		// Show dialog to add new item
		if mainDialog != nil {
			mainDialog.Hide()
		}
		if *currentDialog != nil {
			(*currentDialog).Hide()
		}

		newEntry := widget.NewEntry()
		newEntry.SetPlaceHolder(locale.T(config.AddPlaceholderKey))

		// Create label and input field vertically for better width
		label := widget.NewLabel(locale.T(config.AddLabelKey))
		entryContainer := container.NewVBox(label, newEntry)

		newItemDialog := dialog.NewCustomWithoutButtons(locale.T(config.AddTitleKey), entryContainer, window)

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
				dialog.ShowError(fmt.Errorf(locale.T(config.EnterNameInfoKey)), window)
				return
			}

			// Check if item already exists
			for _, i := range *config.ExistingItems {
				if strings.EqualFold(i, newItemName) {
					dialog.ShowError(fmt.Errorf(locale.T(config.ErrorExistsKey), newItemName), window)
					return
				}
			}

			// Add new item to reference table
			if err := config.ServiceAddFunc(newItemName); err != nil {
				dialog.ShowError(err, window)
				return
			}

			// Update item list
			*config.ExistingItems = append(*config.ExistingItems, newItemName)
			updateSelectOptions(*config.ExistingItems)

			newItemDialog.Hide()
			// Return to main dialog with new item selected
			if mainDialog != nil {
				mainDialog.Show()
			}
			(*selectWidget).SetSelected(newItemName)
		})

		newItemDialog.SetButtons([]fyne.CanvasObject{cancelBtn, saveBtn})

		// Increase dialog width for new item
		parentSize := window.Canvas().Size()
		dialogWidth := parentSize.Width * 0.6
		if dialogWidth < 700 {
			dialogWidth = 700
		}
		newItemDialog.Resize(fyne.NewSize(dialogWidth, newItemDialog.MinSize().Height))

		*currentDialog = newItemDialog
		newItemDialog.Show()
	})
	entityContainer.Add(addNewBtn)

	// Create popup dialog
	popup := dialog.NewCustomWithoutButtons(locale.T("common.select_one"), entityContainer, window)

	// Add close button
	closeBtn := widget.NewButton(locale.T("common.close"), func() {
		popup.Hide()
		*currentDialog = nil
		// Return to main dialog
		if mainDialog != nil {
			mainDialog.Show()
		}
	})

	popup.SetButtons([]fyne.CanvasObject{closeBtn})

	// Resize popup
	parentSize := window.Canvas().Size()
	popupWidth := parentSize.Width * 0.4
	if popupWidth < 400 {
		popupWidth = 400
	}
	popup.Resize(fyne.NewSize(popupWidth, popup.MinSize().Height))

	*currentDialog = popup
	popup.Show()
}

// NewBrandOption creates the "Add new brand" option text
func NewBrandOption() string {
	return "+ " + locale.T("common.add") + " " + strings.TrimSuffix(locale.T("form.model.brand"), ":")
}

// NewScaleOption creates the "Add new scale" option text
func NewScaleOption() string {
	return "+ " + locale.T("common.add") + " " + strings.TrimSuffix(locale.T("form.model.scale"), ":")
}

// NewModelTypeOption creates the "Add new model type" option text
func NewModelTypeOption() string {
	return "+ " + locale.T("common.add") + " " + strings.TrimSuffix(locale.T("form.model.type"), ":")
}
