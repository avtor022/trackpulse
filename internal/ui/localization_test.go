package ui

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
)

func TestMonitoringRefreshKeepsVisibleContentAndUpdatesButtons(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	p := NewMonitoringPanel(nil, app.NewWindow("test"))
	content := p.content

	locale.SetLocale("ru")
	p.Refresh()

	if p.content != content {
		t.Fatal("refresh replaced the container already attached to the tab")
	}
	if got := p.startButton.Text; got != "Старт" {
		t.Errorf("start button = %q, want Старт", got)
	}
}

func TestFyneDialogButtonsFollowApplicationLanguage(t *testing.T) {
	initializeFyneLocale()
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
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("ru")
	p := &CompetitionPanel{allCompetitions: []models.Competition{{CompetitionType: "qualifying", Status: "in_progress", ModelType: "*", ModelScale: "*"}}}
	table := p.createCompetitionTable()
	for col, want := range map[int]string{2: "Квалификация", 3: "Все типы (Массовый заезд)", 4: "Все масштабы", 10: "Идёт"} {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: 0, Col: col}, label)
		if label.Text != want {
			t.Errorf("column %d = %q, want %q", col, label.Text, want)
		}
	}
}

func TestYearDeleteMessageAcceptsStringInBothLanguages(t *testing.T) {
	for _, lang := range []string{"en", "ru"} {
		locale.SetLocale(lang)
		message := fmt.Sprintf(locale.T("dialog.delete_year.message"), "2025")
		if strings.Contains(message, "%!") || !strings.Contains(message, "2025") {
			t.Errorf("%s: invalid message %q", lang, message)
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
		if got := localizedError(errors.New(tc.message)).Error(); got != tc.want {
			t.Errorf("localizedError(%q) = %q, want %q", tc.message, got, tc.want)
		}
	}
}

func TestAllFilterSelectionSurvivesLanguageChange(t *testing.T) {
	locale.SetLocale("en")
	filter := &CompetitionFilter{selectedYear: filterValue(locale.T("filter.all_years"), "filter.all_years")}
	locale.SetLocale("ru")
	if !filter.matchesFilters(models.Competition{}) {
		t.Error("selecting all years should not filter competitions after language change")
	}
}

func TestLogsRefreshKeepsVisibleContentAndLog(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	p := NewLogsPanel(app.NewWindow("test"))
	p.AddLog("scan")
	content, log := p.content, p.logText

	locale.SetLocale("ru")
	p.Refresh()

	if p.content != content || p.logText != log {
		t.Fatal("refresh replaced the visible log and its content")
	}
	if len(p.logText.Segments) != 1 {
		t.Error("refresh discarded the existing log")
	}
}

func TestParticipantRefreshUpdatesSelectionButton(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	p := &ParticipantPanel{competitionSelect: widget.NewButton(locale.T("participants.select.competition.button"), nil)}

	locale.SetLocale("ru")
	p.Refresh()

	if got := p.competitionSelect.Text; got != "Выбрать..." {
		t.Errorf("selection button = %q, want Выбрать...", got)
	}
}

func TestFilterRefreshShowsTranslatedAllOptions(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("en")
	filter := NewCompetitionFilter(app.NewWindow("test"), nil)
	filter.CreateContent()

	locale.SetLocale("ru")
	filter.RefreshButtons()

	if got := filter.yearButton.Text; got != "Все годы" {
		t.Errorf("year filter = %q, want Все годы", got)
	}
}
