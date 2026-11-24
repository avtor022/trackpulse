package ui_test

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/database"
	"trackpulse/internal/locale"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
	"trackpulse/internal/ui"
)

type uiFixture struct {
	app                *ui.App
	window             fyne.Window
	tabs               *container.AppTabs
	db                 *database.DB
	competitionService *service.CompetitionService
}

// Fyne's test driver executes fyne.Do on the calling goroutine. Queue those
// callbacks so tests can dispatch them on the same goroutine as UI events.
type queuedDriver struct {
	fyne.Driver
	mu    sync.Mutex
	queue []func()
}

func (d *queuedDriver) DoFromGoroutine(fn func(), wait bool) {
	var done chan struct{}
	if wait {
		done = make(chan struct{})
		original := fn
		fn = func() { defer close(done); original() }
	}
	d.mu.Lock()
	d.queue = append(d.queue, fn)
	d.mu.Unlock()
	if done != nil {
		<-done
	}
}

func (d *queuedDriver) drain() {
	d.mu.Lock()
	queue := d.queue
	d.queue = nil
	d.mu.Unlock()
	for _, fn := range queue {
		fn()
	}
}

type queuedApp struct {
	fyne.App
	driver *queuedDriver
}

func (a *queuedApp) Driver() fyne.Driver { return a.driver }

func newQueuedTestApp() *queuedApp {
	base := test.NewApp()
	app := &queuedApp{App: base, driver: &queuedDriver{Driver: base.Driver()}}
	fyne.SetCurrentApp(app)
	return app
}

func refreshTestApp(t *testing.T) *uiFixture {
	t.Helper()
	fyneApp := newQueuedTestApp()
	t.Cleanup(fyneApp.Quit)
	locale.SetLocale("en")
	db, err := database.NewDB(filepath.Join(t.TempDir(), "trackpulse.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if err := db.Initialize(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT INTO competition_seasons VALUES ('season','Summer','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_model_brands VALUES ('brand','Brand','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitors VALUES ('competitor',7,'Driver',NULL,'','',0,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO rc_models VALUES ('model','Brand','Model','1:10','Touring','','','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitor_models VALUES ('tag-one','competitor','model','TAG1','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');
		INSERT INTO competitions (id,competition_title,competition_type,model_type,model_scale,track_name,status,season,created_at,updated_at) VALUES
		('race-one','Race one','heat','*','*','','scheduled','Summer','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z'),
		('race-two','Race two','heat','*','*','','scheduled','Summer','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	competitorRepo := repository.NewCompetitorRepository(db.DB)
	modelRepo := repository.NewRCModelRepository(db.DB)
	typeRepo := repository.NewRCModelTypeRepository(db.DB)
	scaleRepo := repository.NewRCModelScaleRepository(db.DB)
	competitorModelService := service.NewCompetitorModelService(repository.NewCompetitorModelRepository(db.DB), competitorRepo, modelRepo)
	competitionService := service.NewCompetitionService(repository.NewCompetitionRepository(db.DB), typeRepo, scaleRepo, repository.NewCompetitionTrackRepository(db.DB), repository.NewCompetitionYearRepository(db.DB), repository.NewCompetitionSeasonRepository(db.DB))
	app := ui.NewAppWithFyne(fyneApp, service.NewCompetitorService(competitorRepo), service.NewRCModelService(modelRepo, repository.NewRCModelBrandRepository(db.DB), scaleRepo, typeRepo), service.NewSettingsService(repository.NewSettingsRepository(db.DB)), competitorModelService, competitionService, service.NewCompetitionParticipantService(repository.NewCompetitionParticipantRepository(db.DB), competitorModelService, competitionService), "en", db.DB)
	t.Cleanup(app.Close)
	windows := fyneApp.Driver().AllWindows()
	window := windows[len(windows)-1]
	tabs := app.CreateMainContent()
	window.SetContent(tabs)
	window.Resize(fyne.NewSize(1200, 800))
	return &uiFixture{app: app, window: window, tabs: tabs, db: db, competitionService: competitionService}
}

func buttonByText(t *testing.T, object fyne.CanvasObject, text string) *widget.Button {
	t.Helper()
	var buttons []string
	for _, object := range test.LaidOutObjects(object) {
		if button, ok := object.(*widget.Button); ok {
			buttons = append(buttons, button.Text)
			if button.Text == text {
				return button
			}
		}
	}
	t.Fatalf("button %q not found among %q", text, buttons)
	return nil
}

func findCheck(object fyne.CanvasObject, text string) *widget.Check {
	for _, object := range test.LaidOutObjects(object) {
		if check, ok := object.(*widget.Check); ok && strings.Contains(check.Text, text) {
			return check
		}
	}
	return nil
}

func tableIn(t *testing.T, object fyne.CanvasObject) *widget.Table {
	t.Helper()
	for _, object := range test.LaidOutObjects(object) {
		if table, ok := object.(*widget.Table); ok {
			return table
		}
	}
	t.Fatal("table not found")
	return nil
}

func selectRace(t *testing.T, f *uiFixture, tab int, title string) {
	t.Helper()
	button := locale.T("participants.select.competition.button")
	if tab == 0 {
		button = locale.T("form.competition.select")
	}
	test.Tap(buttonByText(t, f.tabs.Items[tab].Content, button))
	test.Tap(buttonByText(t, f.window.Canvas().Overlays().Top(), title))
}

func TestTransponderTabReloadsRelatedNames(t *testing.T) {
	f := refreshTestApp(t)
	f.tabs.SelectIndex(1)
	if _, err := f.db.Exec(`UPDATE competitors SET full_name='Updated driver' WHERE id='competitor'; UPDATE rc_models SET model_name='Updated model' WHERE id='model'`); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(3)
	table := tableIn(t, f.tabs.Items[3].Content)
	for column, want := range map[int]string{1: "Updated driver", 2: "Brand Updated model"} {
		label := widget.NewLabel("")
		table.UpdateCell(widget.TableCellID{Row: 0, Col: column}, label)
		if label.Text != want {
			t.Errorf("column %d = %q, want %q", column, label.Text, want)
		}
	}
}

func TestParticipantTabReloadsTranspondersAndKeepsSelections(t *testing.T) {
	f := refreshTestApp(t)
	f.tabs.SelectIndex(5)
	selectRace(t, f, 5, "Race one")
	check := findCheck(f.tabs.Items[5].Content, "TAG1")
	if check == nil {
		t.Fatal("initial transponder not available")
	}
	check.SetChecked(true)
	f.tabs.SelectIndex(3)
	if _, err := f.db.Exec(`INSERT INTO competitor_models VALUES ('tag-two','competitor','model','TAG2','RFID',1,'2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(5)
	old, new := findCheck(f.tabs.Items[5].Content, "TAG1"), findCheck(f.tabs.Items[5].Content, "TAG2")
	if old == nil || !old.Checked || new == nil {
		t.Fatal("transponder list did not refresh or discarded the selection")
	}
}

func TestBindingAndSwitchingCompetitionsRefreshAvailableTransponders(t *testing.T) {
	f := refreshTestApp(t)
	f.tabs.SelectIndex(5)
	selectRace(t, f, 5, "Race one")
	check := findCheck(f.tabs.Items[5].Content, "TAG1")
	if check == nil {
		t.Fatal("transponder not available")
	}
	check.SetChecked(true)
	for _, object := range test.LaidOutObjects(f.tabs.Items[5].Content) {
		if toolbar, ok := object.(*widget.Toolbar); ok {
			toolbar.Items[1].(*widget.ToolbarAction).OnActivated()
			break
		}
	}
	var count int
	if err := f.db.QueryRow(`SELECT COUNT(*) FROM competition_participants WHERE competition_id='race-one'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("bound participants = %d, err = %v", count, err)
	}
	if findCheck(f.tabs.Items[5].Content, "TAG1") != nil {
		t.Fatal("bound transponder still appears as available")
	}
	test.Tap(buttonByText(t, f.window.Canvas().Overlays().Top(), "OK"))
	selectRace(t, f, 5, "Race two")
	if findCheck(f.tabs.Items[5].Content, "TAG1") == nil {
		t.Fatal("transponder bound to the first race is not available for the second race")
	}
}

func TestMonitoringClearsDeletedCompetition(t *testing.T) {
	f := refreshTestApp(t)
	selectRace(t, f, 0, "Race one")
	if buttonByText(t, f.tabs.Items[0].Content, locale.T("button.start")).Disabled() {
		t.Fatal("scheduled race should be startable")
	}
	f.tabs.SelectIndex(4)
	if err := f.competitionService.DeleteCompetition("race-one"); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(0)
	buttonByText(t, f.tabs.Items[0].Content, locale.T("form.competition.select"))
	if !buttonByText(t, f.tabs.Items[0].Content, locale.T("button.start")).Disabled() || !buttonByText(t, f.tabs.Items[0].Content, locale.T("button.stop")).Disabled() {
		t.Fatal("deleted race remains startable or stoppable")
	}
}

func TestParticipantTabClearsDeletedCompetition(t *testing.T) {
	f := refreshTestApp(t)
	f.tabs.SelectIndex(5)
	selectRace(t, f, 5, "Race one")
	f.tabs.SelectIndex(4)
	if err := f.competitionService.DeleteCompetition("race-one"); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(5)
	if findCheck(f.tabs.Items[5].Content, "TAG1") != nil {
		t.Fatal("participant panel still shows the deleted race's transponder")
	}
	assertLabel(t, f.tabs.Items[5].Content, locale.T("participants.select.competition.placeholder"))
}

func assertLabel(t *testing.T, object fyne.CanvasObject, text string) {
	t.Helper()
	for _, object := range test.LaidOutObjects(object) {
		if label, ok := object.(*widget.Label); ok && label.Text == text {
			return
		}
	}
	t.Fatalf("label %q not found", text)
}

func TestMonitoringKeepsSelectionOnDatabaseError(t *testing.T) {
	f := refreshTestApp(t)
	selectRace(t, f, 0, "Race one")
	f.tabs.SelectIndex(4)
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(0)
	buttonByText(t, f.tabs.Items[0].Content, "Race one")
	assertLabel(t, f.tabs.Items[0].Content, locale.T("status.error_loading"))
}

func TestMonitoringStopsTrackingDeletedActiveCompetition(t *testing.T) {
	f := refreshTestApp(t)
	selectRace(t, f, 0, "Race one")
	test.Tap(buttonByText(t, f.tabs.Items[0].Content, locale.T("button.start")))
	test.Tap(buttonByText(t, f.window.Canvas().Overlays().Top(), "OK"))
	f.tabs.SelectIndex(4)
	if err := f.competitionService.DeleteCompetition("race-one"); err != nil {
		t.Fatal(err)
	}
	f.tabs.SelectIndex(0)
	buttonByText(t, f.tabs.Items[0].Content, locale.T("form.competition.select"))
	selectRace(t, f, 0, "Race two")
	if buttonByText(t, f.tabs.Items[0].Content, locale.T("button.start")).Disabled() {
		t.Fatal("deleted active race still prevents starting another race")
	}
}
