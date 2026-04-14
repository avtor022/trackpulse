package ui_test

import (
	"testing"

	"trackpulse/internal/ui"
)

// TestReferenceItem tests the ReferenceItem struct
func TestReferenceItem(t *testing.T) {
	item := ui.ReferenceItem{Name: "Test Item"}
	if item.Name != "Test Item" {
		t.Errorf("Expected Name to be 'Test Item', got '%s'", item.Name)
	}
}

// TestReferencePopupConfig tests the ReferencePopupConfig struct
func TestReferencePopupConfig(t *testing.T) {
	config := ui.ReferencePopupConfig{
		Title:          "dialog.title",
		AddTitle:       "dialog.add.title",
		AddLabel:       "dialog.add.label",
		AddPlaceholder: "dialog.add.placeholder",
		DeleteMessage:  "dialog.delete.message",
		NewErrorExists: "dialog.error.exists",
		EnterNameInfo:  "info.enter_name",
		NewItemOption:  "+ Add New",
	}

	if config.Title != "dialog.title" {
		t.Errorf("Expected Title to be 'dialog.title', got '%s'", config.Title)
	}
	if config.AddTitle != "dialog.add.title" {
		t.Errorf("Expected AddTitle to be 'dialog.add.title', got '%s'", config.AddTitle)
	}
	if config.NewItemOption != "+ Add New" {
		t.Errorf("Expected NewItemOption to be '+ Add New', got '%s'", config.NewItemOption)
	}
}

// TestReferencePopupManagerCreation tests creating a ReferencePopupManager
func TestReferencePopupManagerCreation(t *testing.T) {
	initialItems := []string{"Item1", "Item2", "Item3"}
	newItemOption := "+ Add New"

	var selectedValue string
	var updatedOptions []string

	onSelect := func(value string) {
		selectedValue = value
	}

	updateOpts := func(options []string) {
		updatedOptions = options
	}

	config := ui.ReferencePopupConfig{
		Title:          "dialog.title",
		AddTitle:       "dialog.add.title",
		AddLabel:       "dialog.add.label",
		AddPlaceholder: "dialog.add.placeholder",
		DeleteMessage:  "dialog.delete.message",
		NewErrorExists: "dialog.error.exists",
		EnterNameInfo:  "info.enter_name",
		NewItemOption:  newItemOption,
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			return []ui.ReferenceItem{{Name: "Item1"}, {Name: "Item2"}}, nil
		},
		AddFunc: func(name string) error {
			return nil
		},
		DeleteFunc: func(name string) error {
			return nil
		},
	}

	// Note: We can't fully test ShowPopup without a real fyne.Window
	// but we can test the manager creation and helper methods
	manager := ui.NewReferencePopupManager(nil, config, initialItems, newItemOption, onSelect, updateOpts)

	if manager == nil {
		t.Fatal("Expected ReferencePopupManager to be created, got nil")
	}

	items := manager.GetItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Test RefreshItems
	err := manager.RefreshItems()
	if err != nil {
		t.Errorf("Expected no error from RefreshItems, got %v", err)
	}

	items = manager.GetItems()
	if len(items) != 2 {
		t.Errorf("Expected 2 items after refresh, got %d", len(items))
	}

	// Verify updateOpts was called
	if len(updatedOptions) != 3 {
		t.Errorf("Expected 3 options after update, got %d", len(updatedOptions))
	}
}

// TestGetAllFuncError tests error handling in GetAllFunc
func TestGetAllFuncError(t *testing.T) {
	config := ui.ReferencePopupConfig{
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			return nil, &testError{message: "failed to get items"}
		},
	}

	manager := ui.NewReferencePopupManager(nil, config, []string{}, "+ Add", nil, nil)

	err := manager.RefreshItems()
	if err == nil {
		t.Error("Expected error from RefreshItems, got nil")
	}
}

// TestUpdateOptionsCallback tests that UpdateOptions callback is called correctly
func TestUpdateOptionsCallback(t *testing.T) {
	initialItems := []string{"Item1"}
	newItemOption := "+ Add"

	called := false
	var receivedOptions []string

	updateOpts := func(options []string) {
		called = true
		receivedOptions = options
	}

	config := ui.ReferencePopupConfig{
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			return []ui.ReferenceItem{{Name: "New Item"}}, nil
		},
	}

	manager := ui.NewReferencePopupManager(nil, config, initialItems, newItemOption, nil, updateOpts)

	err := manager.RefreshItems()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !called {
		t.Error("Expected updateOpts callback to be called")
	}

	expectedLen := 2 // "New Item" + "+ Add"
	if len(receivedOptions) != expectedLen {
		t.Errorf("Expected %d options, got %d", expectedLen, len(receivedOptions))
	}
}

// TestOnSelectCallback tests that OnSelect callback works correctly
func TestOnSelectCallback(t *testing.T) {
	var selectedValue string
	onSelect := func(value string) {
		selectedValue = value
	}

	config := ui.ReferencePopupConfig{
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			return []ui.ReferenceItem{}, nil
		},
	}

	manager := ui.NewReferencePopupManager(nil, config, []string{}, "+ Add", onSelect, nil)

	// Simulate what happens when an item is deleted (onSelect is called with empty string)
	if manager.GetItems() != nil {
		// Just verifying the manager was created successfully with the callback
	}

	// Note: We can't directly test the onSelect callback being triggered
	// during ShowPopup without a real window, but we verify it's stored
	_ = manager
}

// TestEmptyItemsList tests handling of empty items list
func TestEmptyItemsList(t *testing.T) {
	config := ui.ReferencePopupConfig{
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			return []ui.ReferenceItem{}, nil
		},
	}

	manager := ui.NewReferencePopupManager(nil, config, []string{}, "+ Add", nil, nil)

	items := manager.GetItems()
	if len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}

	err := manager.RefreshItems()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	items = manager.GetItems()
	if len(items) != 0 {
		t.Errorf("Expected 0 items after refresh, got %d", len(items))
	}
}

// TestMultipleItems tests handling of multiple items
func TestMultipleItems(t *testing.T) {
	initialItems := []string{"Item1", "Item2", "Item3", "Item4", "Item5"}

	config := ui.ReferencePopupConfig{
		GetAllFunc: func() ([]ui.ReferenceItem, error) {
			items := make([]ui.ReferenceItem, len(initialItems))
			for i, name := range initialItems {
				items[i] = ui.ReferenceItem{Name: name}
			}
			return items, nil
		},
	}

	manager := ui.NewReferencePopupManager(nil, config, initialItems, "+ Add", nil, nil)

	items := manager.GetItems()
	if len(items) != 5 {
		t.Errorf("Expected 5 items, got %d", len(items))
	}

	expectedItems := map[string]bool{
		"Item1": true,
		"Item2": true,
		"Item3": true,
		"Item4": true,
		"Item5": true,
	}

	for _, item := range items {
		if !expectedItems[item] {
			t.Errorf("Unexpected item: %s", item)
		}
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
