package locale

import (
	"fyne.io/fyne/v2"
)

// Locale represents a localization configuration
type Locale struct {
	code string
	data map[string]string
}

// CurrentLocale holds the active locale
var CurrentLocale *Locale

// SupportedLocales defines all supported languages
var SupportedLocales = map[string]string{
	"en": "English",
}

// Init initializes the localization system with the default language (English)
func Init() {
	CurrentLocale = &Locale{
		code: "en",
		data: getEnglishStrings(),
	}
}

// SetLocale changes the current locale (future implementation for multiple languages)
func SetLocale(code string) {
	switch code {
	case "en":
		CurrentLocale = &Locale{
			code: code,
			data: getEnglishStrings(),
		}
	// Add more cases here when adding new languages
	default:
		CurrentLocale = &Locale{
			code: "en",
			data: getEnglishStrings(),
		}
	}
}

// Get retrieves a localized string by key
func Get(key string) string {
	if CurrentLocale == nil {
		Init()
	}
	if val, ok := CurrentLocale.data[key]; ok {
		return val
	}
	return key // Return key if translation not found
}

// T is a shorthand for Get
func T(key string) string {
	return Get(key)
}

// getEnglishStrings returns all English UI strings
func getEnglishStrings() map[string]string {
	return map[string]string{
		// Application
		"app.title":                    "TrackPulse",
		"app.welcome":                  "Welcome to TrackPulse",
		"app.version":                  "Version 1.0.0",
		
		// Main Menu
		"menu.file":                    "File",
		"menu.file.new_brand":          "New Brand",
		"menu.file.new_model":          "New Model",
		"menu.file.exit":               "Exit",
		"menu.tools":                   "Tools",
		"menu.tools.import":            "Import Data",
		"menu.tools.export":            "Export Data",
		"menu.help":                    "Help",
		"menu.help.about":              "About",
		
		// Tabs
		"tab.monitoring":               "Monitoring",
		"tab.racers":                   "Racers",
		"tab.models":                   "Models",
		"tab.transponders":             "Transponders",
		"tab.races":                    "Races",
		"tab.logs":                     "Logs",
		"tab.settings":                 "Settings",
		
		// Brand Panel
		"brand.panel.title":            "Brands",
		"brand.panel.search":           "Search brands...",
		"brand.panel.add":              "Add Brand",
		"brand.panel.edit":             "Edit",
		"brand.panel.delete":           "Delete",
		"brand.panel.no_brands":        "No brands found",
		
		// Model Panel
		"model.panel.title":            "Models",
		"model.panel.search":           "Search models...",
		"model.panel.add":              "Add Model",
		"model.panel.edit":             "Edit",
		"model.panel.delete":           "Delete",
		"model.panel.no_models":        "No models found",
		"model.panel.select_brand":     "Please select a brand first",
		
		// Model Table Headers
		"common.id":                    "ID",
		"model.header.brand":           "Brand",
		"model.header.name":            "Model Name",
		"model.header.scale":           "Scale",
		"model.header.type":            "Type",
		"model.header.motor":           "Motor",
		"model.header.drive":           "Drive",
		"model.header.created":         "Created At",
		"model.header.updated":         "Updated At",
		
		// Racer Table Headers
		"racer.header.number":          "Number",
		"racer.header.name":            "Name",
		"racer.header.country":         "Country",
		"racer.header.city":            "City",
		"racer.header.birthday":        "Birthday",
		"racer.header.rating":          "Rating",
		
		// New Brand Dialog
		"dialog.new_brand.title":       "New Brand",
		"dialog.new_brand.label":       "Brand Name:",
		"dialog.new_brand.placeholder": "Enter brand name",
		"dialog.new_brand.create":      "Create",
		"dialog.new_brand.cancel":      "Cancel",
		"dialog.new_brand.success":     "Brand created successfully",
		"dialog.new_brand.error_empty": "Brand name cannot be empty",
		"dialog.new_brand.error_exists": "Brand already exists",
		
		// New Model Dialog
		"dialog.new_model.title":       "New Model",
		"dialog.new_model.brand_label": "Brand:",
		"dialog.new_model.name_label":  "Model Name:",
		"dialog.new_model.placeholder": "Enter model name",
		"dialog.new_model.create":      "Create",
		"dialog.new_model.cancel":      "Cancel",
		"dialog.new_model.success":     "Model created successfully",
		"dialog.new_model.error_empty": "Model name cannot be empty",
		"dialog.new_model.error_exists": "Model already exists",
		"dialog.new_model.error_no_brand": "Please select a brand",
		
		// Add Brand Dialog (nested)
		"dialog.add_brand.title":       "Add New Brand",
		"dialog.add_brand.label":       "Brand Name:",
		"dialog.add_brand.placeholder": "Enter new brand name",
		"dialog.add_brand.save":        "Save",
		
		// Edit Dialog
		"dialog.edit.title":            "Edit",
		"dialog.edit.save":             "Save",
		"dialog.edit.cancel":           "Cancel",
		"dialog.edit.success":          "Updated successfully",
		
		// Delete Confirmation
		"dialog.delete.title":          "Confirm Delete",
		"dialog.delete.message":        "Are you sure you want to delete \"%s\"?",
		"dialog.delete.confirm":        "Delete",
		"dialog.delete.cancel":         "Cancel",
		"dialog.delete.success":        "Deleted successfully",
		
		// Common
		"common.ok":                    "OK",
		"common.cancel":                "Cancel",
		"common.yes":                   "Yes",
		"common.no":                    "No",
		"common.error":                 "Error",
		"common.success":               "Success",
		"common.loading":               "Loading...",
		"common.search":                "Search",
		"common.filter":                "Filter",
		"common.refresh":               "Refresh",
		"common.save":                  "Save",
		"common.create":                "Create",
		"common.delete":                "Delete",
		"common.edit":                  "Edit",
		"common.add":                   "Add",
		"common.close":                 "Close",
		"common.select":                "Select",
		"common.required":              "is required",
		
		// Form Labels - Models
		"form.model.brand":             "Brand",
		"form.model.name":              "Model Name",
		"form.model.scale":             "Scale",
		"form.model.type":              "Model Type",
		"form.model.motor":             "Motor Type",
		"form.model.drive":             "Drive Type",
		"form.model.scale_placeholder": "e.g., 1:8",
		"form.model.type_placeholder":  "e.g., Monster Truck",
		"form.model.motor_placeholder": "e.g., Brushless",
		"form.model.drive_placeholder": "e.g., 4WD",
		
		// Form Labels - Racers
		"form.racer.number":            "Racer Number",
		"form.racer.name":              "Full Name",
		"form.racer.country":           "Country",
		"form.racer.city":              "City",
		"form.racer.birthday":          "Birthday (MM.DD.YYYY)",
		"form.racer.rating":            "Rating",
		"form.racer.number_placeholder": "e.g., 7",
		"form.racer.name_placeholder":  "John Doe",
		"form.racer.country_placeholder": "USA",
		"form.racer.city_placeholder":  "New York",
		"form.racer.birthday_placeholder": "MM.DD.YYYY",
		"form.racer.rating_placeholder": "0",
		
		// Validation Errors
		"error.required.brand":         "Brand is required",
		"error.required.name":          "Model name is required",
		"error.required.scale":         "Scale is required",
		"error.required.type":          "Model type is required",
		"error.required.number":        "Racer number is required",
		"error.required.racer_name":    "Full name is required",
		"error.invalid.number":         "Invalid racer number",
		"error.invalid.date":           "Invalid date format (use MM.DD.YYYY)",
		"error.invalid.rating":         "Invalid rating",
		
		// Status Messages
		"status.ready":                 "Ready",
		"status.processing":            "Processing...",
		"status.completed":             "Completed",
		"status.failed":                "Failed",
		"status.loaded_models":         "Loaded %d models",
		"status.loaded_racers":         "Loaded %d racers",
		"status.no_models":             "No models found",
		"status.no_racers":             "No racers found",
		"status.selected_model":        "Selected: %s %s",
		"status.selected_racer":        "Selected: %s",
		"status.operation_cancelled":   "Operation cancelled",
		"status.created_success":       "Created successfully",
		"status.updated_success":       "Updated successfully",
		"status.deleted_success":       "Deleted successfully",
		"status.create_failed":         "Create failed: %s",
		"status.update_failed":         "Update failed: %s",
		"status.delete_failed":         "Delete failed: %s",
		"status.refresh_error":         "Error refreshing data",
		
		// Info Messages
		"info.select_first":            "Please select an item in the table first",
		"info.not_found":               "Selected item not found",
		"info.enter_brand_name":        "Enter brand name",
		
		// Errors
		"error.database":               "Database error: %s",
		"error.network":                "Network error: %s",
		"error.generic":                "An error occurred: %s",
		"error.permission":             "Permission denied",
		"error.not_found":              "Not found",
		
		// Settings
		"settings.language":            "Language",
		"settings.language.en":         "English",
	}
}

// BindText creates a binding for dynamic text updates (optional helper)
func BindText(key string) string {
	return T(key)
}

// Ensure fyne import is used
var _ fyne.App
