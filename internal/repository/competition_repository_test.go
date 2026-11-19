package repository_test

import (
	"path/filepath"
	"testing"

	"trackpulse/internal/database"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
)

func TestCompetitionMinimumLapTimeIsPersisted(t *testing.T) {
	db, err := database.NewDB(filepath.Join(t.TempDir(), "trackpulse.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO competition_seasons (id, season, created_at, updated_at) VALUES ('s', 'winter', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`)
	if err != nil {
		t.Fatal(err)
	}
	repo := repository.NewCompetitionRepository(db.DB)
	competition := &models.Competition{ID: "race", CompetitionTitle: "Race", CompetitionType: "heat", Season: "winter", MinLapTimeMs: 5500}
	if err := repo.Create(competition); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByID(competition.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.MinLapTimeMs != 5500 {
		t.Fatalf("minimum lap time = %d, want 5500", got.MinLapTimeMs)
	}
	got.MinLapTimeMs = 3000
	if err := repo.Update(got); err != nil {
		t.Fatal(err)
	}
	all, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].MinLapTimeMs != 3000 {
		t.Fatalf("updated competitions = %+v", all)
	}
}
