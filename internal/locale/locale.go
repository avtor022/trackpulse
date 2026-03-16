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
	"ru": "Русский",
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
	case "ru":
		CurrentLocale = &Locale{
			code: code,
			data: getRussianStrings(),
		}
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
		
		// Common Info Dialog Title
		"common.info":                  "Information",
		
		// Errors
		"error.database":               "Database error: %s",
		"error.network":                "Network error: %s",
		"error.generic":                "An error occurred: %s",
		"error.permission":             "Permission denied",
		"error.not_found":              "Not found",
		
		// Settings
		"settings.language":            "Language",
		"settings.language.en":         "English",
		"settings.language.ru":         "Russian",
	}
}

// getRussianStrings returns all Russian UI strings
func getRussianStrings() map[string]string {
	return map[string]string{
		// Application
		"app.title":                    "TrackPulse",
		"app.welcome":                  "Добро пожаловать в TrackPulse",
		"app.version":                  "Версия 1.0.0",
		
		// Main Menu
		"menu.file":                    "Файл",
		"menu.file.new_brand":          "Новый бренд",
		"menu.file.new_model":          "Новая модель",
		"menu.file.exit":               "Выход",
		"menu.tools":                   "Инструменты",
		"menu.tools.import":            "Импорт данных",
		"menu.tools.export":            "Экспорт данных",
		"menu.help":                    "Помощь",
		"menu.help.about":              "О программе",
		
		// Tabs
		"tab.monitoring":               "Мониторинг",
		"tab.racers":                   "Гонщики",
		"tab.models":                   "Модели",
		"tab.transponders":             "Транспондеры",
		"tab.races":                    "Гонки",
		"tab.logs":                     "Логи",
		"tab.settings":                 "Настройки",
		
		// Brand Panel
		"brand.panel.title":            "Бренды",
		"brand.panel.search":           "Поиск брендов...",
		"brand.panel.add":              "Добавить бренд",
		"brand.panel.edit":             "Редактировать",
		"brand.panel.delete":           "Удалить",
		"brand.panel.no_brands":        "Бренды не найдены",
		
		// Model Panel
		"model.panel.title":            "Модели",
		"model.panel.search":           "Поиск моделей...",
		"model.panel.add":              "Добавить модель",
		"model.panel.edit":             "Редактировать",
		"model.panel.delete":           "Удалить",
		"model.panel.no_models":        "Модели не найдены",
		"model.panel.select_brand":     "Пожалуйста, выберите бренд",
		
		// Model Table Headers
		"common.id":                    "ID",
		"model.header.brand":           "Бренд",
		"model.header.name":            "Название модели",
		"model.header.scale":           "Масштаб",
		"model.header.type":            "Тип",
		"model.header.motor":           "Мотор",
		"model.header.drive":           "Привод",
		"model.header.created":         "Создано",
		"model.header.updated":         "Обновлено",
		
		// Racer Table Headers
		"racer.header.number":          "Номер",
		"racer.header.name":            "Имя",
		"racer.header.country":         "Страна",
		"racer.header.city":            "Город",
		"racer.header.birthday":        "Дата рождения",
		"racer.header.rating":          "Рейтинг",
		
		// New Brand Dialog
		"dialog.new_brand.title":       "Новый бренд",
		"dialog.new_brand.label":       "Название бренда:",
		"dialog.new_brand.placeholder": "Введите название бренда",
		"dialog.new_brand.create":      "Создать",
		"dialog.new_brand.cancel":      "Отмена",
		"dialog.new_brand.success":     "Бренд успешно создан",
		"dialog.new_brand.error_empty": "Название бренда не может быть пустым",
		"dialog.new_brand.error_exists": "Бренд уже существует",
		
		// New Model Dialog
		"dialog.new_model.title":       "Новая модель",
		"dialog.new_model.brand_label": "Бренд:",
		"dialog.new_model.name_label":  "Название модели:",
		"dialog.new_model.placeholder": "Введите название модели",
		"dialog.new_model.create":      "Создать",
		"dialog.new_model.cancel":      "Отмена",
		"dialog.new_model.success":     "Модель успешно создана",
		"dialog.new_model.error_empty": "Название модели не может быть пустым",
		"dialog.new_model.error_exists": "Модель уже существует",
		"dialog.new_model.error_no_brand": "Пожалуйста, выберите бренд",
		
		// Add Brand Dialog (nested)
		"dialog.add_brand.title":       "Добавить новый бренд",
		"dialog.add_brand.label":       "Название бренда:",
		"dialog.add_brand.placeholder": "Введите название нового бренда",
		"dialog.add_brand.save":        "Сохранить",
		
		// Edit Dialog
		"dialog.edit.title":            "Редактировать",
		"dialog.edit.save":             "Сохранить",
		"dialog.edit.cancel":           "Отмена",
		"dialog.edit.success":          "Успешно обновлено",
		
		// Delete Confirmation
		"dialog.delete.title":          "Подтверждение удаления",
		"dialog.delete.message":        "Вы уверены, что хотите удалить \"%s\"?",
		"dialog.delete.confirm":        "Удалить",
		"dialog.delete.cancel":         "Отмена",
		"dialog.delete.success":        "Успешно удалено",
		
		// Common
		"common.ok":                    "OK",
		"common.cancel":                "Отмена",
		"common.yes":                   "Да",
		"common.no":                    "Нет",
		"common.error":                 "Ошибка",
		"common.success":               "Успешно",
		"common.loading":               "Загрузка...",
		"common.search":                "Поиск",
		"common.filter":                "Фильтр",
		"common.refresh":               "Обновить",
		"common.save":                  "Сохранить",
		"common.create":                "Создать",
		"common.delete":                "Удалить",
		"common.edit":                  "Редактировать",
		"common.add":                   "Добавить",
		"common.close":                 "Закрыть",
		"common.select":                "Выбрать",
		"common.required":              "обязательно",
		
		// Form Labels - Models
		"form.model.brand":             "Бренд",
		"form.model.name":              "Название модели",
		"form.model.scale":             "Масштаб",
		"form.model.type":              "Тип модели",
		"form.model.motor":             "Тип мотора",
		"form.model.drive":             "Тип привода",
		"form.model.scale_placeholder": "например, 1:8",
		"form.model.type_placeholder":  "например, Monster Truck",
		"form.model.motor_placeholder": "например, Brushless",
		"form.model.drive_placeholder": "например, 4WD",
		
		// Form Labels - Racers
		"form.racer.number":            "Номер гонщика",
		"form.racer.name":              "Полное имя",
		"form.racer.country":           "Страна",
		"form.racer.city":              "Город",
		"form.racer.birthday":          "Дата рождения (ММ.ДД.ГГГГ)",
		"form.racer.rating":            "Рейтинг",
		"form.racer.number_placeholder": "например, 7",
		"form.racer.name_placeholder":  "Иван Иванов",
		"form.racer.country_placeholder": "Россия",
		"form.racer.city_placeholder":  "Москва",
		"form.racer.birthday_placeholder": "ММ.ДД.ГГГГ",
		"form.racer.rating_placeholder": "0",
		
		// Validation Errors
		"error.required.brand":         "Бренд обязателен",
		"error.required.name":          "Название модели обязательно",
		"error.required.scale":         "Масштаб обязателен",
		"error.required.type":          "Тип модели обязателен",
		"error.required.number":        "Номер гонщика обязателен",
		"error.required.racer_name":    "Полное имя обязательно",
		"error.invalid.number":         "Неверный номер гонщика",
		"error.invalid.date":           "Неверный формат даты (используйте ММ.ДД.ГГГГ)",
		"error.invalid.rating":         "Неверный рейтинг",
		
		// Status Messages
		"status.ready":                 "Готов",
		"status.processing":            "Обработка...",
		"status.completed":             "Завершено",
		"status.failed":                "Ошибка",
		"status.loaded_models":         "Загружено моделей: %d",
		"status.loaded_racers":         "Загружено гонщиков: %d",
		"status.no_models":             "Модели не найдены",
		"status.no_racers":             "Гонщики не найдены",
		"status.selected_model":        "Выбрано: %s %s",
		"status.selected_racer":        "Выбран: %s",
		"status.operation_cancelled":   "Операция отменена",
		"status.created_success":       "Успешно создано",
		"status.updated_success":       "Успешно обновлено",
		"status.deleted_success":       "Успешно удалено",
		"status.create_failed":         "Ошибка создания: %s",
		"status.update_failed":         "Ошибка обновления: %s",
		"status.delete_failed":         "Ошибка удаления: %s",
		"status.refresh_error":         "Ошибка обновления данных",
		
		// Info Messages
		"info.select_first":            "Пожалуйста, сначала выберите элемент в таблице",
		"info.not_found":               "Выбранный элемент не найден",
		"info.enter_brand_name":        "Введите название бренда",
		
		// Common Info Dialog Title
		"common.info":                  "Информация",
		
		// Errors
		"error.database":               "Ошибка базы данных: %s",
		"error.network":                "Ошибка сети: %s",
		"error.generic":                "Произошла ошибка: %s",
		"error.permission":             "Доступ запрещен",
		"error.not_found":              "Не найдено",
		
		// Settings
		"settings.language":            "Язык",
		"settings.language.en":         "Английский",
		"settings.language.ru":         "Русский",
	}
}

// BindText creates a binding for dynamic text updates (optional helper)
func BindText(key string) string {
	return T(key)
}

// Ensure fyne import is used
var _ fyne.App
