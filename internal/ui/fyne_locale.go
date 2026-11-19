package ui

import (
	"sync"

	"fyne.io/fyne/v2/lang"
	"trackpulse/internal/locale"
)

var fyneLocaleOnce sync.Once

// initializeFyneLocale makes Fyne's built-in dialog and edit controls follow
// the language selected in TrackPulse, rather than only the system language.
func initializeFyneLocale() {
	fyneLocaleOnce.Do(func() {
		fallback := lang.L
		keys := map[string]string{
			"OK":         "common.ok",
			"Error":      "common.error",
			"Yes":        "common.yes",
			"No":         "common.no",
			"Cancel":     "common.cancel",
			"Confirm":    "common.confirm",
			"Cut":        "common.cut",
			"Copy":       "common.copy",
			"Paste":      "common.paste",
			"Select all": "common.select_all",
			"Undo":       "common.undo",
			"Redo":       "common.redo",
		}
		lang.L = func(text string, data ...interface{}) string {
			if key, ok := keys[text]; ok && len(data) == 0 {
				return locale.T(key)
			}
			return fallback(text, data...)
		}
	})
}
