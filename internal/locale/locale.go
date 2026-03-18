package locale

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// translationsDir returns the path to the translations directory
func translationsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "translations")
}

// loadTranslations loads translations from a JSON file
func loadTranslations(langCode string) (map[string]string, error) {
	filePath := filepath.Join(translationsDir(), langCode+".json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read translation file %s: %w", filePath, err)
	}
	
	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, fmt.Errorf("failed to parse translation file %s: %w", filePath, err)
	}
	
	return translations, nil
}

// Init initializes the localization system with the default language (English)
func Init() {
	translations, err := loadTranslations("en")
	if err != nil {
		// Fallback to empty map if loading fails
		translations = make(map[string]string)
	}
	CurrentLocale = &Locale{
		code: "en",
		data: translations,
	}
}

// SetLocale changes the current locale
func SetLocale(code string) {
	translations, err := loadTranslations(code)
	if err != nil {
		// Fallback to English if loading fails
		translations, _ = loadTranslations("en")
		code = "en"
	}
	
	CurrentLocale = &Locale{
		code: code,
		data: translations,
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

// BindText creates a binding for dynamic text updates (optional helper)
func BindText(key string) string {
	return T(key)
}

// Ensure fyne import is used
var _ fyne.App
