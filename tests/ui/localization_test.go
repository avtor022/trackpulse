package ui_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/ui"
)

func TestMonitoringRefreshKeepsVisibleContentAndUpdatesButtons(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	p := ui.NewMonitoringPanel(nil, nil, app.NewWindow("test"))
	defer p.Close()
	content := p.Content()
	locale.SetLocale("ru")
	p.Refresh()
	if p.Content() != content {
		t.Fatal("refresh replaced the container already attached to the tab")
	}
	buttonByText(t, content, "Старт")
}

func TestFyneDialogButtonsFollowApplicationLanguage(t *testing.T) {
	refreshTestApp(t)
	locale.SetLocale("ru")
	if got := lang.L("Yes"); got != "Да" {
		t.Errorf("confirmation button = %q, want Да", got)
	}
	if got := lang.L("Error"); got != "Ошибка" {
		t.Errorf("error title = %q, want Ошибка", got)
	}
	locale.SetLocale("en")
	if got := lang.L("Yes"); got != "Yes" {
		t.Errorf("confirmation button = %q, want Yes", got)
	}
}

func TestCompetitionTableLocalizesTypeAndStatus(t *testing.T) {
	f := refreshTestApp(t)
	if _, err := f.db.Exec(`UPDATE competitions SET competition_type='qualifying', status='in_progress' WHERE id='race-one'`); err != nil {
		t.Fatal(err)
	}
	locale.SetLocale("ru")
	f.tabs.SelectIndex(4)
	table := tableIn(t, f.tabs.Items[4].Content)
	rows, _ := table.Length()
	row := -1
	for i := 0; i < rows; i++ {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: i, Col: 1}, label)
		if label.Text == "Race one" {
			row = i
			break
		}
	}
	if row == -1 {
		t.Fatal("competition not found in table")
	}
	for col, want := range map[int]string{2: "Квалификация", 3: "Все типы (Массовый заезд)", 4: "Все масштабы", 10: "Идёт"} {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: row, Col: col}, label)
		if label.Text != want {
			t.Errorf("column %d = %q, want %q", col, label.Text, want)
		}
	}
}

func TestYearDeleteMessageAcceptsStringInBothLanguages(t *testing.T) {
	for _, language := range []string{"en", "ru"} {
		locale.SetLocale(language)
		message := fmt.Sprintf(locale.T("dialog.delete_year.message"), "2025")
		if strings.Contains(message, "%!") || !strings.Contains(message, "2025") {
			t.Errorf("%s: invalid message %q", language, message)
		}
	}
}

func TestLocalizedServiceErrors(t *testing.T) {
	locale.SetLocale("ru")
	for _, tc := range []struct{ message, want string }{
		{"transponder number '123' already exists", "Номер транспондера «123» уже занят"},
		{"competitor number 7 already exists", "Номер спортсмена 7 уже занят"},
		{"competitor not found", "Спортсмен не найден"},
		{"cannot delete brand 'XRAY': it is used by model 'X4'", "Нельзя удалить бренд «XRAY»: он используется моделью «X4»"},
		{"database is locked", "Не удалось выполнить операцию: database is locked"},
	} {
		if got := ui.LocalizedError(errors.New(tc.message)).Error(); got != tc.want {
			t.Errorf("LocalizedError(%q) = %q, want %q", tc.message, got, tc.want)
		}
	}
}

func TestAllFilterSelectionSurvivesLanguageChange(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	window := app.NewWindow("test")
	filter := ui.NewCompetitionFilter(window, nil)
	content := filter.CreateContent()
	window.SetContent(content)
	filter.SetCompetitions([]models.Competition{{CompetitionTitle: "Race"}})
	test.Tap(buttonByText(t, content, locale.T("filter.all_years")))
	test.Tap(buttonByText(t, window.Canvas().Overlays().Top(), locale.T("filter.all_years")))
	locale.SetLocale("ru")
	filter.RefreshLocale()
	if len(filter.GetFilteredCompetitions()) != 1 {
		t.Error("selecting all years should not filter competitions after language change")
	}
}

func TestLogsRefreshKeepsVisibleContentAndLog(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	p := ui.NewLogsPanel(app.NewWindow("test"))
	p.AddLog("scan")
	content := p.Content()
	locale.SetLocale("ru")
	p.Refresh()
	if p.Content() != content {
		t.Fatal("refresh replaced the visible log and its content")
	}
	for _, object := range test.LaidOutObjects(content) {
		if log, ok := object.(*widget.RichText); ok && len(log.Segments) == 1 {
			return
		}
	}
	t.Fatal("refresh discarded the existing log")
}

func TestParticipantRefreshUpdatesSelectionButton(t *testing.T) {
	f := refreshTestApp(t)
	f.tabs.SelectIndex(7)
	for _, object := range test.LaidOutObjects(f.tabs.Items[7].Content) {
		if selectLanguage, ok := object.(*widget.Select); ok {
			for _, option := range selectLanguage.Options {
				if option == "Русский" {
					selectLanguage.SetSelected(option)
					buttonByText(t, f.tabs.Items[5].Content, "Выбрать...")
					return
				}
			}
		}
	}
	t.Fatal("language selector not found")
}

func TestFilterRefreshShowsTranslatedAllOptions(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	filter := ui.NewCompetitionFilter(app.NewWindow("test"), nil)
	content := filter.CreateContent()
	locale.SetLocale("ru")
	filter.RefreshButtons()
	buttonByText(t, content, "Все годы")
}
