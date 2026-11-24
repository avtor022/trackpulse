package locale_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func loadTranslations(t *testing.T, language string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "internal", "locale", "translations", language+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		t.Fatal(err)
	}
	return translations
}

func TestTranslationsHaveMatchingKeysAndFormats(t *testing.T) {
	en := loadTranslations(t, "en")
	ru := loadTranslations(t, "ru")
	formats := regexp.MustCompile(`%[+#0 -]*[0-9]*(?:\.[0-9]+)?[a-zA-Z]`)
	for key, value := range en {
		translated, ok := ru[key]
		if !ok {
			t.Errorf("ru translation missing %q", key)
			continue
		}
		if !reflect.DeepEqual(formats.FindAllString(value, -1), formats.FindAllString(translated, -1)) {
			t.Errorf("format placeholders differ for %q: %q / %q", key, value, translated)
		}
	}
	for key := range ru {
		if _, ok := en[key]; !ok {
			t.Errorf("en translation missing %q", key)
		}
	}
}
