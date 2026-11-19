package service_test

import (
	"path/filepath"
	"testing"
	"time"

	"trackpulse/internal/database"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
)

func lapServiceFixture(t *testing.T) (*service.LapService, *repository.CompetitionLapsRepository, *database.DB) {
	t.Helper()
	db, err := database.NewDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT INTO rc_model_brands VALUES ('b','Brand','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('c',7,'Test Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_models VALUES ('m','Brand','Model','1:10','Touring','','','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm','c','m','TAG','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitions (id, competition_title, status, created_at, updated_at)
		VALUES ('race','Race','in_progress','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p','race','cm',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	laps := repository.NewCompetitionLapsRepository(db.DB)
	svc := service.NewLapService(repository.NewRawScanRepository(db.DB), repository.NewCompetitorModelRepository(db.DB), repository.NewCompetitionRepository(db.DB), repository.NewCompetitionParticipantRepository(db.DB), laps)
	if err := svc.SetActiveCompetition(&models.Competition{ID: "race", Status: "in_progress", MinLapTimeMs: 5000}); err != nil {
		t.Fatal(err)
	}
	return svc, laps, db
}

func TestLapServiceCountsCompletedLapsAndRejectsEarlyReadings(t *testing.T) {
	svc, _, _ := lapServiceFixture(t)
	start := time.Date(2026, 1, 1, 12, 0, 0, 123000000, time.UTC)
	scan := func(tag string, offset time.Duration) {
		svc.ProcessScan(service.LapScan{TagValue: tag, Timestamp: start.Add(offset)})
	}
	scan("TAG", 0)
	if got := svc.GetParticipantResults()["p"]; got.LapCount != 0 || !got.StartTime.Equal(start) {
		t.Fatalf("first passage should only start the first lap: %+v", got)
	}
	scan("TAG", 2*time.Second)
	scan("UNKNOWN", 3*time.Second)
	if got := svc.GetParticipantResults()["p"]; got.LapCount != 0 || !got.LastPassTime.Equal(start) {
		t.Fatalf("early or unknown readings changed the lap: %+v", got)
	}
	scan("TAG", 5*time.Second)
	scan("TAG", 12*time.Second)
	got := svc.GetParticipantResults()["p"]
	if got.LapCount != 2 || got.BestLapTimeMs != 5000 || got.LastLapTimeMs != 7000 {
		t.Fatalf("unexpected lap results: %+v", got)
	}
}

func TestScansFromBeforeRaceStartAreIgnored(t *testing.T) {
	svc, _, _ := lapServiceFixture(t)
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	if err := svc.SetActiveCompetition(&models.Competition{ID: "race", Status: "in_progress", TimeStart: &start, MinLapTimeMs: 5000}); err != nil {
		t.Fatal(err)
	}
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: start.Add(-time.Second)})
	if got := svc.GetParticipantResults()["p"]; !got.LastPassTime.IsZero() {
		t.Fatalf("pre-start scan started a lap: %+v", got)
	}
}

func TestLapResultsCanBePersistedMoreThanOnce(t *testing.T) {
	svc, laps, _ := lapServiceFixture(t)
	start := time.Date(2026, 1, 1, 12, 0, 0, 123000000, time.UTC)
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: start})
	if err := svc.PersistResults(); err != nil {
		t.Fatal(err)
	}
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: start.Add(5500 * time.Millisecond)})
	if err := svc.PersistResults(); err != nil {
		t.Fatal(err)
	}
	got, err := laps.GetByParticipantID("p")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.NumberOfLaps != 1 || got.BestLapTimeMs != 5500 || !got.TimeStart.Equal(start) {
		t.Fatalf("persisted results = %+v", got)
	}
}

func TestStandingsOrderByLapsThenFastestLap(t *testing.T) {
	_, laps, db := lapServiceFixture(t)
	_, err := db.Exec(`
		INSERT INTO competitors VALUES ('c2',8,'Second Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('c3',9,'Third Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('c4',6,'Fourth Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm2','c2','m','TAG2','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm3','c3','m','TAG3','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('cm4','c4','m','TAG4','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p2','race','cm2',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p3','race','cm3',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_participants VALUES ('p4','race','cm4',NULL,0,0,'','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competition_laps (id, competition_participant_id, time_start, number_of_laps, best_lap_time_ms, created_at, updated_at)
		VALUES ('l1','p','2026-01-01T00:00:00Z',2,6000,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z'),
		       ('l2','p2','2026-01-01T00:00:00Z',2,5000,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z'),
		       ('l4','p4','2026-01-01T00:00:00Z',0,0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	got, err := laps.GetStandings("race")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 || got[0].FullName != "Second Driver" || got[1].FullName != "Test Driver" || got[2].FullName != "Fourth Driver" || got[3].FullName != "Third Driver" {
		t.Fatalf("unexpected standings: %+v", got)
	}
	if got[3].NumberOfLaps != 0 || got[3].BestLapTimeMs != 0 {
		t.Fatalf("unstarted participant must have no result: %+v", got[3])
	}
}

func TestQueuedScansPersistAndRestoreAnInProgressRace(t *testing.T) {
	svc, _, db := lapServiceFixture(t)
	svc.Start()
	start := time.Date(2026, 1, 1, 12, 0, 0, 123000000, time.UTC)
	for _, at := range []time.Time{start, start.Add(5500 * time.Millisecond)} {
		if !svc.QueueScan(service.LapScan{TagValue: "TAG", Timestamp: at, ReaderType: "RFID"}) {
			t.Fatal("failed to queue scan")
		}
	}
	svc.Flush()
	standings, err := svc.GetStandings("race")
	if err != nil || len(standings) != 1 || standings[0].NumberOfLaps != 1 || standings[0].BestLapTimeMs != 5500 {
		t.Fatalf("standings = %+v, err = %v", standings, err)
	}
	var rawCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM raw_scans`).Scan(&rawCount); err != nil || rawCount != 2 {
		t.Fatalf("raw scans = %d, err = %v", rawCount, err)
	}
	svc.Stop()
	if err := svc.SetActiveCompetition(&models.Competition{ID: "race", Status: "in_progress", MinLapTimeMs: 5000}); err != nil {
		t.Fatal(err)
	}
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: start.Add(11 * time.Second)})
	got := svc.GetParticipantResults()["p"]
	if got.LapCount != 2 || got.LastLapTimeMs != 5500 {
		t.Fatalf("restored result = %+v", got)
	}
}

func TestLapPersistenceFailureIsReported(t *testing.T) {
	svc, _, db := lapServiceFixture(t)
	if _, err := db.Exec(`DROP TABLE competition_laps`); err != nil {
		t.Fatal(err)
	}
	var reported error
	svc.SetErrorHandler(func(err error) { reported = err })
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: time.Now()})
	if reported == nil {
		t.Fatal("persistence failure was not reported")
	}
}

func TestRestoringRaceRejectsInvalidLastPassTime(t *testing.T) {
	svc, _, db := lapServiceFixture(t)
	svc.ProcessScan(service.LapScan{TagValue: "TAG", Timestamp: time.Now()})
	if _, err := db.Exec(`UPDATE competition_laps SET last_pass_time = 'invalid'`); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetActiveCompetition(&models.Competition{ID: "race", Status: "in_progress", MinLapTimeMs: 5000}); err == nil {
		t.Fatal("restoring malformed timing data must return an error")
	}
}

func TestStopDrainsQueuedScansAndRejectsNewOnes(t *testing.T) {
	svc, _, _ := lapServiceFixture(t)
	svc.Start()
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, at := range []time.Time{start, start.Add(5 * time.Second)} {
		if !svc.QueueScan(service.LapScan{TagValue: "TAG", Timestamp: at, ReaderType: "RFID"}) {
			t.Fatal("failed to queue scan")
		}
	}
	svc.Stop()
	if got := svc.GetParticipantResults()["p"]; got.LapCount != 1 {
		t.Fatalf("queued scans were not drained: %+v", got)
	}
	if svc.QueueScan(service.LapScan{TagValue: "TAG", Timestamp: start.Add(10 * time.Second)}) {
		t.Fatal("stopped service accepted a scan")
	}
}
