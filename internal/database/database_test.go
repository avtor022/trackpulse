package database_test

import (
	"path/filepath"
	"testing"

	"trackpulse/internal/database"
)

func TestInitializeAddsMinimumLapTimeToExistingCompetitions(t *testing.T) {
	db, err := database.NewDB(filepath.Join(t.TempDir(), "trackpulse.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE competitions (
		id TEXT PRIMARY KEY NOT NULL, competition_title TEXT NOT NULL,
		competition_type TEXT, model_type TEXT, model_scale TEXT, track_name TEXT,
		lap_count_target INTEGER, time_limit_minutes INTEGER, time_start TEXT,
		time_finish TEXT, status TEXT, competition_year INTEGER, season TEXT,
		created_at TEXT NOT NULL, updated_at TEXT NOT NULL
	); INSERT INTO competitions (id, competition_title, created_at, updated_at)
		VALUES ('existing', 'Existing', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	var minimum int
	if err := db.QueryRow(`SELECT min_lap_time_ms FROM competitions WHERE id = 'existing'`).Scan(&minimum); err != nil {
		t.Fatal(err)
	}
	if minimum != 5000 {
		t.Errorf("minimum lap time = %d ms, want 5000 ms", minimum)
	}
	if err := db.Initialize(); err != nil {
		t.Fatalf("second initialization: %v", err)
	}
}
