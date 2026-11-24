package ui_test

import (
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/database"
	"trackpulse/internal/locale"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
	"trackpulse/internal/ui"
)

func monitoringTestPanel(t *testing.T) (*ui.MonitoringPanel, fyne.Window, *service.LapService, *service.CompetitionService, *database.DB) {
	t.Helper()
	app := newQueuedTestApp()
	t.Cleanup(app.Quit)
	db, err := database.NewDB(filepath.Join(t.TempDir(), "race.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT INTO competition_seasons VALUES ('season','winter','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_model_brands VALUES ('b','Brand','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('c',7,'Test Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_models VALUES ('m','Brand','Model','1:10','Touring','','','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm','c','m','TAG','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitions (id,competition_title,competition_type,model_type,model_scale,track_name,status,season,created_at,updated_at)
		VALUES ('race','Race','heat','','','','scheduled','winter','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p','race','cm',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	compRepo := repository.NewCompetitionRepository(db.DB)
	compService := service.NewCompetitionService(compRepo, repository.NewRCModelTypeRepository(db.DB), repository.NewRCModelScaleRepository(db.DB), repository.NewCompetitionTrackRepository(db.DB), repository.NewCompetitionYearRepository(db.DB), repository.NewCompetitionSeasonRepository(db.DB))
	laps := service.NewLapService(repository.NewRawScanRepository(db.DB), repository.NewCompetitorModelRepository(db.DB), compRepo, repository.NewCompetitionParticipantRepository(db.DB), repository.NewCompetitionLapsRepository(db.DB))
	laps.Start()
	t.Cleanup(laps.Stop)
	window := app.NewWindow("test")
	panel := ui.NewMonitoringPanel(compService, laps, window)
	t.Cleanup(panel.Close)
	window.SetContent(panel.Content())
	window.Resize(fyne.NewSize(1200, 800))
	test.Tap(buttonByText(t, panel.Content(), locale.T("form.competition.select")))
	test.Tap(buttonByText(t, window.Canvas().Overlays().Top(), "Race"))
	return panel, window, laps, compService, db
}

func TestMonitoringTableShowsCompetitorAndLapResults(t *testing.T) {
	locale.SetLocale("ru")
	panel, _, _, _, db := monitoringTestPanel(t)
	_, err := db.Exec(`
		UPDATE competitors SET full_name='Иван Иванов' WHERE id='c';
		INSERT INTO competition_laps (id,competition_participant_id,time_start,number_of_laps,best_lap_time_ms,created_at,updated_at) VALUES ('lap','p','2026-01-01T00:00:00Z',3,5432,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}
	panel.UpdateData()
	table := tableIn(t, panel.Content())
	rows, columns := table.Length()
	if rows != 1 || columns != 5 {
		t.Fatalf("table dimensions = %d x %d", rows, columns)
	}
	for col, expected := range []string{"1", "7", "Иван Иванов", "3", "00:05:432"} {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: 0, Col: col}, label)
		if label.Text != expected {
			t.Errorf("column %d = %q, want %q", col, label.Text, expected)
		}
	}
	for _, tc := range []struct {
		milliseconds int
		want         string
	}{{65007, "01:05:007"}, {0, "–"}} {
		if _, err := db.Exec(`UPDATE competition_laps SET best_lap_time_ms=?`, tc.milliseconds); err != nil {
			t.Fatal(err)
		}
		panel.UpdateData()
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: 0, Col: 4}, label)
		if label.Text != tc.want {
			t.Errorf("lap time = %q, want %q", label.Text, tc.want)
		}
	}
}

func TestMinimumLapTimeInput(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  int
		valid bool
	}{
		{"5", 5000, true}, {"2.5", 2500, true}, {"2,5", 2500, true}, {"0", 0, true},
		{"", 0, false}, {"-1", 0, false}, {"NaN", 0, false}, {"Inf", 0, false},
	} {
		got, err := ui.ParseMinimumLapTime(tc.input)
		if (err == nil) != tc.valid || (tc.valid && got != tc.want) {
			t.Errorf("ParseMinimumLapTime(%q) = %d, %v; want %d, valid=%t", tc.input, got, err, tc.want, tc.valid)
		}
	}
}

func TestMonitoringStartsRaceAndShowsCompletedLap(t *testing.T) {
	locale.SetLocale("en")
	panel, window, laps, competitions, _ := monitoringTestPanel(t)
	test.Tap(buttonByText(t, panel.Content(), locale.T("button.start")))
	if active := laps.GetActiveCompetition(); active == nil || active.ID != "race" {
		t.Fatalf("active competition = %+v", active)
	}
	test.Tap(buttonByText(t, window.Canvas().Overlays().Top(), lang.L("OK")))
	start := time.Now().Add(time.Second)
	for _, at := range []time.Time{start, start.Add(6 * time.Second)} {
		if !laps.QueueScan(service.LapScan{TagValue: "TAG", Timestamp: at, ReaderType: "RFID"}) {
			t.Fatal("failed to queue scan")
		}
	}
	laps.Flush()
	fyne.CurrentApp().(*queuedApp).driver.drain()
	panel.UpdateData()
	table := tableIn(t, panel.Content())
	for col, want := range map[int]string{2: "Test Driver", 3: "1", 4: "00:06:000"} {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: 0, Col: col}, label)
		if label.Text != want {
			t.Errorf("column %d = %q, want %q", col, label.Text, want)
		}
	}
	test.Tap(buttonByText(t, panel.Content(), locale.T("button.stop")))
	competition, err := competitions.GetCompetitionByID("race")
	if err != nil || competition.Status != "finished" || laps.GetActiveCompetition() != nil {
		t.Fatalf("race after stop = %+v, err = %v", competition, err)
	}
}
