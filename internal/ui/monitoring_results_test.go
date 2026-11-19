package ui

import (
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/database"
	"trackpulse/internal/locale"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
)

func TestMonitoringTableShowsCompetitorAndLapResults(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	locale.SetLocale("ru")
	p := NewMonitoringPanel(nil, nil, app.NewWindow("test"))
	defer p.Close()
	p.standings = []models.CompetitionStanding{{CompetitorNumber: 7, FullName: "Иван Иванов", NumberOfLaps: 3, BestLapTimeMs: 5432}}
	rows, columns := p.standingsTable.Length()
	if rows != 1 || columns != 5 {
		t.Fatalf("table dimensions = %d x %d", rows, columns)
	}
	want := []string{"1", "7", "Иван Иванов", "3", "00:05:432"}
	for col, expected := range want {
		label := widget.NewLabel("")
		p.standingsTable.UpdateCell(widget.TableCellID{Row: 0, Col: col}, label)
		if label.Text != expected {
			t.Errorf("column %d = %q, want %q", col, label.Text, expected)
		}
	}
	p.standings[0].BestLapTimeMs = 65007
	label := widget.NewLabel("")
	p.standingsTable.UpdateCell(widget.TableCellID{Row: 0, Col: 4}, label)
	if label.Text != "01:05:007" {
		t.Errorf("lap time = %q, want 01:05:007", label.Text)
	}
	p.standings[0].BestLapTimeMs = 0
	label = widget.NewLabel("")
	p.standingsTable.UpdateCell(widget.TableCellID{Row: 0, Col: 4}, label)
	if label.Text != "–" {
		t.Errorf("no completed lap should display a dash, got %q", label.Text)
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
		got, err := parseMinimumLapTime(tc.input)
		if (err == nil) != tc.valid || (tc.valid && got != tc.want) {
			t.Errorf("parseMinimumLapTime(%q) = %d, %v; want %d, valid=%t", tc.input, got, err, tc.want, tc.valid)
		}
	}
}

func TestMonitoringStartsRaceAndShowsCompletedLap(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	db, err := database.NewDB(filepath.Join(t.TempDir(), "race.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT INTO competition_seasons VALUES ('season', 'winter', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');
		INSERT INTO rc_model_brands VALUES ('b','Brand','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('c',7,'Test Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_models VALUES ('m','Brand','Model','1:10','Touring','','','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm','c','m','TAG','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitions (id, competition_title, competition_type, model_type, model_scale, track_name, status, season, created_at, updated_at)
		VALUES ('race','Race','heat','','','','scheduled','winter','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p','race','cm',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	compRepo := repository.NewCompetitionRepository(db.DB)
	compService := service.NewCompetitionService(compRepo, repository.NewRCModelTypeRepository(db.DB), repository.NewRCModelScaleRepository(db.DB), repository.NewCompetitionTrackRepository(db.DB), repository.NewCompetitionYearRepository(db.DB), repository.NewCompetitionSeasonRepository(db.DB))
	lapService := service.NewLapService(repository.NewRawScanRepository(db.DB), repository.NewCompetitorModelRepository(db.DB), compRepo, repository.NewCompetitionParticipantRepository(db.DB), repository.NewCompetitionLapsRepository(db.DB))
	lapService.Start()
	defer lapService.Stop()
	p := NewMonitoringPanel(compService, lapService, app.NewWindow("test"))
	defer p.Close()
	p.onCompetitionSelected("race")
	p.startMonitoring()
	if active := lapService.GetActiveCompetition(); active == nil || active.ID != "race" {
		t.Fatalf("active competition = %+v", active)
	}
	start := time.Now().Add(time.Second)
	for _, at := range []time.Time{start, start.Add(6 * time.Second)} {
		if !lapService.QueueScan(service.LapScan{TagValue: "TAG", Timestamp: at, ReaderType: "RFID"}) {
			t.Fatal("failed to queue scan")
		}
	}
	lapService.Flush()
	p.refreshStandings()
	if len(p.standings) != 1 || p.standings[0].FullName != "Test Driver" || p.standings[0].NumberOfLaps != 1 || p.standings[0].BestLapTimeMs != 6000 {
		t.Fatalf("standings = %+v", p.standings)
	}
	p.stopMonitoring()
	competition, err := compService.GetCompetitionByID("race")
	if err != nil || competition.Status != "finished" || lapService.GetActiveCompetition() != nil {
		t.Fatalf("race after stop = %+v, err = %v", competition, err)
	}
}
