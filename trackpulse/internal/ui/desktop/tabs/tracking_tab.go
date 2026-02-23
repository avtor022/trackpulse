package tabs

import (
	"database/sql"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"trackpulse/internal/models"
)

func CreateRaceTrackingTab(db *sql.DB) *fyne.Container {
	// Create controls for race management
	startButton := widget.NewButton("Старт", func() {
		// Start the active race
		performStartRace(db)
	})
	stopButton := widget.NewButton("Стоп", func() {
		// Stop the active race
		performStopRace(db)
	})
	restartButton := widget.NewButton("Перезапуск", func() {
		// Restart the race
		performRestartRace(db)
	})

	// Create status label
	statusLabel := widget.NewLabel("Статус: Нет активного заезда")

	// Create timer display
	timerLabel := widget.NewLabel("Таймер: 00:00:00")

	// Create a basic layout for the tracking controls
	controlsContainer := container.NewVBox(
		widget.NewCard("Управление заездом", "", container.NewGridWithColumns(4, startButton, stopButton, restartButton)),
		widget.NewCard("Статус", "", statusLabel),
		widget.NewCard("Таймер", "", timerLabel),
	)

	// Create a table to show live lap times
	liveResultsTable := widget.NewTable(
		func() (int, int) { 
			// Get the number of active participants
			activeRace := getActiveRace(db)
			if activeRace.ID == "" {
				return 0, 5 // No active race
			}
			participants := getRaceParticipants(db, activeRace.ID)
			return len(participants), 5 // Participant count x 5 columns (Position, Number, Name, Laps, Last Lap Time)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template"))
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			activeRace := getActiveRace(db)
			if activeRace.ID == "" {
				return // No active race
			}
			
			participants := getRaceParticipants(db, activeRace.ID)
			if id.Row >= len(participants) {
				return
			}
			
			participant := participants[id.Row]
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText("")
			
			switch id.Col {
			case 0:
				// Position - calculated based on laps and last lap time
				position := calculatePosition(db, participant.ID)
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(position))
			case 1:
				// Racer number
				racerNumber := getRacerNumberFromParticipant(db, participant.RacerModelID)
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(racerNumber))
			case 2:
				// Racer name
				racerName := getRacerNameFromParticipant(db, participant.RacerModelID)
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(racerName)
			case 3:
				// Number of laps
				lapCount := getLapCountForParticipant(db, participant.ID)
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(string(lapCount))
			case 4:
				// Last lap time
				lastLapTime := getLastLapTimeForParticipant(db, participant.ID)
				obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(formatMilliseconds(lastLapTime))
			}
		},
	)
	liveResultsTable.SetColumnWidth(0, 80)  // Position
	liveResultsTable.SetColumnWidth(1, 80)  // Number
	liveResultsTable.SetColumnWidth(2, 150) // Name
	liveResultsTable.SetColumnWidth(3, 80)  // Laps
	liveResultsTable.SetColumnWidth(4, 120) // Last Lap Time

	resultsContainer := container.NewBorder(nil, nil, nil, nil, liveResultsTable)

	// Main layout
	mainContainer := container.NewBorder(controlsContainer, nil, nil, nil, resultsContainer)

	return mainContainer
}

// Helper function to get active race
func getActiveRace(db *sql.DB) models.Race {
	// Query for race with status 'active'
	query := `SELECT id, race_title, race_type, model_type, model_scale, track_name, lap_count_target, 
	                 time_limit_minutes, time_start, time_finish, status, created_at, updated_at 
	          FROM races WHERE status = 'active'`
	
	var race models.Race
	err := db.QueryRow(query).Scan(
		&race.ID, &race.RaceTitle, &race.RaceType, &race.ModelType, &race.ModelScale, &race.TrackName, 
		&race.LapCountTarget, &race.TimeLimitMinutes, &race.TimeStart, &race.TimeFinish, &race.Status, 
		&race.CreatedAt, &race.UpdatedAt,
	)
	
	if err != nil {
		// Return empty race if no active race found
		return models.Race{}
	}
	
	return race
}

// Helper function to get race participants
func getRaceParticipants(db *sql.DB, raceID string) []models.RaceParticipant {
	query := `SELECT id, race_id, racer_model_id, grid_position, is_finished, disqualified, dnf_reason, created_at, updated_at 
	          FROM race_participants WHERE race_id = ?`
	
	rows, err := db.Query(query, raceID)
	if err != nil {
		return []models.RaceParticipant{}
	}
	defer rows.Close()
	
	var participants []models.RaceParticipant
	for rows.Next() {
		var participant models.RaceParticipant
		err := rows.Scan(
			&participant.ID, &participant.RaceID, &participant.RacerModelID, &participant.GridPosition, 
			&participant.IsFinished, &participant.Disqualified, &participant.DNFReason, 
			&participant.CreatedAt, &participant.UpdatedAt,
		)
		if err != nil {
			continue
		}
		participants = append(participants, participant)
	}
	
	return participants
}

// Helper function to get racer number from participant
func getRacerNumberFromParticipant(db *sql.DB, racerModelID string) int {
	// First get the racer ID from the racer_model
	racerModel := &models.RacerModel{}
	err := racerModel.GetByID(db, racerModelID)
	if err != nil {
		return 0
	}
	
	// Then get the racer number
	racer := &models.Racer{}
	err = racer.GetByID(db, racerModel.RacerID)
	if err != nil {
		return 0
	}
	
	return racer.RacerNumber
}

// Helper function to get racer name from participant
func getRacerNameFromParticipant(db *sql.DB, racerModelID string) string {
	// First get the racer ID from the racer_model
	racerModel := &models.RacerModel{}
	err := racerModel.GetByID(db, racerModelID)
	if err != nil {
		return "Unknown"
	}
	
	// Then get the racer name
	racer := &models.Racer{}
	err = racer.GetByID(db, racerModel.RacerID)
	if err != nil {
		return "Unknown"
	}
	
	return racer.FullName
}

// Helper function to get lap count for participant
func getLapCountForParticipant(db *sql.DB, participantID string) int {
	// Get the lap count from race_laps table
	query := `SELECT number_of_laps FROM race_laps WHERE race_participant_id = ?`
	
	var lapCount int
	err := db.QueryRow(query, participantID).Scan(&lapCount)
	if err != nil {
		return 0
	}
	
	return lapCount
}

// Helper function to get last lap time for participant
func getLastLapTimeForParticipant(db *sql.DB, participantID string) int {
	// Get the last lap time from race_laps table
	query := `SELECT last_lap_time_ms FROM race_laps WHERE race_participant_id = ?`
	
	var lastLapTime int
	err := db.QueryRow(query, participantID).Scan(&lastLapTime)
	if err != nil {
		return 0
	}
	
	return lastLapTime
}

// Helper function to calculate position
func calculatePosition(db *sql.DB, participantID string) int {
	// This is a simplified position calculation
	// In a real implementation, we would calculate based on laps and lap times
	currentRace := getActiveRace(db)
	if currentRace.ID == "" {
		return 0
	}
	
	participants := getRaceParticipants(db, currentRace.ID)
	
	// Get the lap count for the current participant
	currentLapCount := getLapCountForParticipant(db, participantID)
	
	// Count how many participants have more laps
	position := 1
	for _, p := range participants {
		if p.ID != participantID {
			participantLapCount := getLapCountForParticipant(db, p.ID)
			if participantLapCount > currentLapCount {
				position++
			}
		}
	}
	
	return position
}

// Helper function to format milliseconds to MM:SS format
func formatMilliseconds(ms int) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return string(minutes) + ":" + string(seconds)
}

// Helper function to perform start race action
func performStartRace(db *sql.DB) {
	// Find a race with status 'scheduled' to start
	query := `SELECT id FROM races WHERE status = 'scheduled' LIMIT 1`
	
	var raceID string
	err := db.QueryRow(query).Scan(&raceID)
	if err != nil {
		// No scheduled race to start
		return
	}
	
	// Update race status to 'active' and set start time
	updateQuery := `UPDATE races SET status = 'active', time_start = ? WHERE id = ?`
	now := "2026-02-23T12:56:00Z" // In a real implementation, use time.Now().Format(time.RFC3339)
	_, err = db.Exec(updateQuery, now, raceID)
	if err != nil {
		// Handle error
		return
	}
	
	// Initialize race_laps entries for all participants
	initializeRaceLaps(db, raceID)
}

// Helper function to initialize race_laps for participants
func initializeRaceLaps(db *sql.DB, raceID string) {
	participants := getRaceParticipants(db, raceID)
	now := "2026-02-23T12:56:00Z" // In a real implementation, use time.Now().Format(time.RFC3339)
	
	for _, participant := range participants {
		// Check if entry already exists
		var count int
		countQuery := `SELECT COUNT(*) FROM race_laps WHERE race_participant_id = ?`
		db.QueryRow(countQuery, participant.ID).Scan(&count)
		
		if count == 0 {
			// Create new race_lap entry
			insertQuery := `INSERT INTO race_laps (id, race_participant_id, time_start, number_of_laps, best_lap_time_ms, best_lap_number, last_lap_time_ms, total_race_time_ms, created_at, updated_at) 
			                VALUES (?, ?, ?, 0, 0, 0, 0, 0, ?, ?)`
			// Generate a new ID for the race_lap entry
			// In a real implementation, use uuid.New().String()
			raceLapID := "temp_id_for_demo"
			_, err := db.Exec(insertQuery, raceLapID, participant.ID, now, now, now)
			if err != nil {
				// Handle error
			}
		}
	}
}

// Helper function to perform stop race action
func performStopRace(db *sql.DB) {
	// Find the active race
	activeRace := getActiveRace(db)
	if activeRace.ID == "" {
		// No active race to stop
		return
	}
	
	// Update race status to 'finished' and set finish time
	updateQuery := `UPDATE races SET status = 'finished', time_finish = ? WHERE id = ?`
	now := "2026-02-23T12:56:00Z" // In a real implementation, use time.Now().Format(time.RFC3339)
	_, err := db.Exec(updateQuery, now, activeRace.ID)
	if err != nil {
		// Handle error
		return
	}
}

// Helper function to perform restart race action
func performRestartRace(db *sql.DB) {
	// Find a race with status 'finished' to restart
	query := `SELECT id FROM races WHERE status = 'finished' LIMIT 1`
	
	var raceID string
	err := db.QueryRow(query).Scan(&raceID)
	if err != nil {
		// No finished race to restart
		return
	}
	
	// Reset race status to 'scheduled' and clear times
	updateQuery := `UPDATE races SET status = 'scheduled', time_start = NULL, time_finish = NULL WHERE id = ?`
	_, err = db.Exec(updateQuery, raceID)
	if err != nil {
		// Handle error
		return
	}
	
	// Reset race_laps for all participants in this race
	resetRaceLaps(db, raceID)
}

// Helper function to reset race_laps for a race
func resetRaceLaps(db *sql.DB, raceID string) {
	// Get all participants for this race
	participants := getRaceParticipants(db, raceID)
	
	for _, participant := range participants {
		// Reset their race_laps data
		updateQuery := `UPDATE race_laps SET number_of_laps = 0, best_lap_time_ms = 0, best_lap_number = 0, 
		                last_lap_time_ms = 0, last_pass_time = NULL, total_race_time_ms = 0 
		                WHERE race_participant_id = ?`
		_, err := db.Exec(updateQuery, participant.ID)
		if err != nil {
			// Handle error
		}
	}
}