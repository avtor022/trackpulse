package locale

import (
	"reflect"
	"regexp"
	"testing"
)

func TestTranslationsHaveMatchingKeysAndFormats(t *testing.T) {
	en, err := loadTranslations("en")
	if err != nil {
		t.Fatal(err)
	}
	ru, err := loadTranslations("ru")
	if err != nil {
		t.Fatal(err)
	}
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
