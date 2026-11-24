package ui_test

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
)

func TestParticipantTabShowsNewCompetitionsWithoutRestart(t *testing.T) {
	f := refreshTestApp(t)
	if _, err := f.db.Exec(`UPDATE competitions SET track_name='Track A' WHERE id='race-one'`); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(5)
	test.Tap(buttonByText(t, f.tabs.Items[5].Content, locale.T("filter.all_tracks")))
	test.Tap(buttonByText(t, f.window.Canvas().Overlays().Top(), "Track A"))
	selectRace(t, f, 5, "Race one")
	check := findCheck(f.tabs.Items[5].Content, "TAG1")
	if check == nil {
		t.Fatal("transponder not available")
	}
	check.SetChecked(true)
	f.tabs.SelectIndex(4)
	for _, competition := range []*models.Competition{
		{CompetitionTitle: "New race", CompetitionType: "heat", TrackName: "Track A", Season: "Summer"},
		{CompetitionTitle: "Other race", CompetitionType: "heat", TrackName: "Track B", Season: "Summer"},
	} {
		if err := f.competitionService.CreateCompetition(competition); err != nil {
			t.Fatal(err)
		}
	}
	f.tabs.SelectIndex(5)
	buttonByText(t, f.tabs.Items[5].Content, "Track A")
	if check = findCheck(f.tabs.Items[5].Content, "TAG1"); check == nil || !check.Checked {
		t.Fatal("switching tabs lost the unsaved transponder selection")
	}
	test.Tap(buttonByText(t, f.tabs.Items[5].Content, locale.T("participants.select.competition.button")))
	popup := f.window.Canvas().Overlays().Top()
	buttonByText(t, popup, "New race")
	for _, object := range test.LaidOutObjects(popup) {
		if button, ok := object.(*widget.Button); ok && button.Text == "Other race" {
			t.Fatal("competition not matching the filter was displayed")
		}
	}
}
