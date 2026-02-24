package localization

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Translator struct {
	currentLang string
	translations map[string]string
}

func NewTranslator(lang string) (*Translator, error) {
	t := &Translator{
		currentLang: lang,
		translations: make(map[string]string),
	}

	err := t.loadTranslations(lang)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Translator) loadTranslations(lang string) error {
	var filename string
	switch lang {
	case "en":
		filename = "en.json"
	case "ru":
		fallthrough
	default:
		filename = "ru.json"
	}

	path := filepath.Join("internal", "localization", "locales", filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &t.translations)
	if err != nil {
		return err
	}

	return nil
}

func (t *Translator) T(key string) string {
	if translation, exists := t.translations[key]; exists {
		return translation
	}
	// Return the key itself if no translation is found
	return key
}

func (t *Translator) SetLanguage(lang string) error {
	err := t.loadTranslations(lang)
	if err != nil {
		return err
	}
	t.currentLang = lang
	return nil
}

func (t *Translator) GetCurrentLanguage() string {
	return t.currentLang
}