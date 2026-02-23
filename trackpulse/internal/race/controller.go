package race

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"trackpulse/internal/models"

	uuid "github.com/google/uuid"
)

type Controller struct {
	db          *sql.DB
	mu          sync.Mutex
	activeTimer *time.Timer
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (rc *Controller) StartRace(raceID string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	tx, err := rc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if race exists and is in scheduled status
	var status string
	err = tx.QueryRow("SELECT status FROM races WHERE id = ?", raceID).Scan(&status)
	if err != nil {
		return fmt.Errorf("race not found: %v", err)
	}

	if status != "scheduled" {
		return fmt.Errorf("race is not in scheduled status, current status: %s", status)
	}

	// Get race participants
	rows, err := tx.Query("SELECT id FROM race_participants WHERE race_id = ?", raceID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %v", err)
	}
	defer rows.Close()

	var participantIDs []string
	for rows.Next() {
		var participantID string
		if err := rows.Scan(&participantID); err != nil {
			return fmt.Errorf("failed to scan participant: %v", err)
		}
		participantIDs = append(participantIDs, participantID)
	}

	if len(participantIDs) == 0 {
		return fmt.Errorf("no participants in race")
	}

	// Update race status and start time
	now := time.Now().UTC()
	_, err = tx.Exec("UPDATE races SET status = ?, time_start = ? WHERE id = ?", "active", now.Format(time.RFC3339), raceID)
	if err != nil {
		return fmt.Errorf("failed to update race status: %v", err)
	}

	// Initialize race_laps for all participants
	for _, participantID := range participantIDs {
		lapID := uuid.New().String()
		_, err := tx.Exec(`
			INSERT INTO race_laps (
				id, race_participant_id, time_start, number_of_laps,
				best_lap_time_ms, best_lap_number, last_lap_time_ms,
				total_race_time_ms, created_at, updated_at
			) VALUES (?, ?, ?, 0, 0, 0, 0, 0, ?, ?)`,
			lapID, participantID, now.Format(time.RFC3339),
			now.Format(time.RFC3339), now.Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("failed to initialize race_laps: %v", err)
		}
	}

	// Record in audit log
	err = rc.recordAuditLog(tx, "RACE_START", "race", raceID, "system", "", map[string]interface{}{
		"race_id":   raceID,
		"timestamp": now.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to record audit log: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Start automatic finish timer if time limit is set
	rc.startAutoFinishTimer(raceID)

	return nil
}

func (rc *Controller) StopRace(raceID string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	tx, err := rc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if race exists and is active
	var status string
	err = tx.QueryRow("SELECT status FROM races WHERE id = ?", raceID).Scan(&status)
	if err != nil {
		return fmt.Errorf("race not found: %v", err)
	}

	if status != "active" {
		return fmt.Errorf("race is not active, current status: %s", status)
	}

	// Update race status and finish time
	now := time.Now().UTC()
	_, err = tx.Exec("UPDATE races SET status = ?, time_finish = ? WHERE id = ?", "finished", now.Format(time.RFC3339), raceID)
	if err != nil {
		return fmt.Errorf("failed to update race status: %v", err)
	}

	// Record in audit log
	err = rc.recordAuditLog(tx, "RACE_FINISH", "race", raceID, "system", "", map[string]interface{}{
		"race_id":   raceID,
		"timestamp": now.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to record audit log: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Stop any active timer
	rc.stopAutoFinishTimer()

	return nil
}

func (rc *Controller) RestartRace(raceID string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	tx, err := rc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if race exists and is finished
	var status string
	err = tx.QueryRow("SELECT status FROM races WHERE id = ?", raceID).Scan(&status)
	if err != nil {
		return fmt.Errorf("race not found: %v", err)
	}

	if status != "finished" && status != "cancelled" {
		return fmt.Errorf("race is not in finished or cancelled status, current status: %s", status)
	}

	// Reset race_laps for all participants in this race
	_, err = tx.Exec(`
		UPDATE race_laps
		SET number_of_laps = 0,
			best_lap_time_ms = 0,
			best_lap_number = 0,
			last_lap_time_ms = 0,
			total_race_time_ms = 0,
			time_finish = NULL,
			last_pass_time = NULL,
			updated_at = ?
		WHERE race_participant_id IN (
			SELECT id FROM race_participants WHERE race_id = ?
		)`, time.Now().UTC().Format(time.RFC3339), raceID)
	if err != nil {
		return fmt.Errorf("failed to reset race_laps: %v", err)
	}

	// Clear lap history for this race
	_, err = tx.Exec("DELETE FROM lap_history WHERE race_participant_id IN (SELECT id FROM race_participants WHERE race_id = ?)", raceID)
	if err != nil {
		return fmt.Errorf("failed to clear lap_history: %v", err)
	}

	// Reset race status and times
	_, err = tx.Exec("UPDATE races SET status = 'scheduled', time_start = NULL, time_finish = NULL WHERE id = ?", raceID)
	if err != nil {
		return fmt.Errorf("failed to reset race: %v", err)
	}

	// Record in audit log
	now := time.Now().UTC()
	err = rc.recordAuditLog(tx, "RACE_RESTART", "race", raceID, "system", "", map[string]interface{}{
		"race_id":   raceID,
		"timestamp": now.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to record audit log: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Stop any active timer
	rc.stopAutoFinishTimer()

	return nil
}

func (rc *Controller) ProcessLap(raceID string, transponderNumber string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	tx, err := rc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get race status
	var raceStatus string
	err = tx.QueryRow("SELECT status FROM races WHERE id = ?", raceID).Scan(&raceStatus)
	if err != nil {
		return fmt.Errorf("race not found: %v", err)
	}

	if raceStatus != "active" {
		return fmt.Errorf("race is not active, current status: %s", raceStatus)
	}

	// Find the racer_model by transponder number
	var racerModelID string
	err = tx.QueryRow("SELECT id FROM racer_models WHERE transponder_number = ?", transponderNumber).Scan(&racerModelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no racer found with transponder: %s", transponderNumber)
		}
		return fmt.Errorf("error finding racer: %v", err)
	}

	// Find the race participant
	var participantID string
	err = tx.QueryRow(`
		SELECT id FROM race_participants
		WHERE race_id = ? AND racer_model_id = ?`, raceID, racerModelID).Scan(&participantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("racer is not participating in this race: %s", transponderNumber)
		}
		return fmt.Errorf("error finding participant: %v", err)
	}

	// Get current lap count for this participant
	var currentLapCount int
	var lastPassTime *time.Time
	err = tx.QueryRow("SELECT number_of_laps, last_pass_time FROM race_laps WHERE race_participant_id = ?", participantID).Scan(&currentLapCount, &lastPassTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// Initialize race_laps if not exists
			now := time.Now().UTC()
			lapID := uuid.New().String()
			_, err := tx.Exec(`
				INSERT INTO race_laps (
					id, race_participant_id, time_start, number_of_laps,
					best_lap_time_ms, best_lap_number, last_lap_time_ms,
					total_race_time_ms, created_at, updated_at
				) VALUES (?, ?, ?, 0, 0, 0, 0, 0, ?, ?)`,
				lapID, participantID, now.Format(time.RFC3339),
				now.Format(time.RFC3339), now.Format(time.RFC3339))
			if err != nil {
				return fmt.Errorf("failed to initialize race_laps: %v", err)
			}
			currentLapCount = 0
			lastPassTime = nil
		} else {
			return fmt.Errorf("error getting current lap count: %v", err)
		}
	}

	// Calculate lap time if we have a previous pass time
	now := time.Now().UTC()
	var lapTimeMS int64 = 0
	if lastPassTime != nil {
		lapTimeMS = now.Sub(*lastPassTime).Milliseconds()
		// Minimum lap time check (3 seconds)
		if lapTimeMS < 3000 {
			// Just update the last pass time but don't count as a new lap
			_, err := tx.Exec("UPDATE race_laps SET last_pass_time = ?, updated_at = ? WHERE race_participant_id = ?",
				now.Format(time.RFC3339), now.Format(time.RFC3339), participantID)
			if err != nil {
				return fmt.Errorf("failed to update last pass time: %v", err)
			}
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit transaction: %v", err)
			}
			return nil
		}
	}

	// Increment lap count
	newLapCount := currentLapCount + 1
	_, err = tx.Exec(`
		UPDATE race_laps
		SET number_of_laps = ?,
			last_lap_time_ms = ?,
			last_pass_time = ?,
			total_race_time_ms = ?,
			updated_at = ?
		WHERE race_participant_id = ?`,
		newLapCount, lapTimeMS, now.Format(time.RFC3339),
		now.Sub(getTimeFromStr(raceStatus)).Milliseconds(), // Total race time since start
		now.Format(time.RFC3339), participantID)
	if err != nil {
		return fmt.Errorf("failed to update race_laps: %v", err)
	}

	// Update best lap time if this is better
	if lapTimeMS > 0 && (lapTimeMS < int64(getCurrentBestLapTime(tx, participantID)) || getCurrentBestLapTime(tx, participantID) == 0) {
		_, err := tx.Exec(`
			UPDATE race_laps
			SET best_lap_time_ms = ?,
				best_lap_number = ?
			WHERE race_participant_id = ?`,
			lapTimeMS, newLapCount, participantID)
		if err != nil {
			return fmt.Errorf("failed to update best lap: %v", err)
		}
	}

	// Record lap in lap_history
	historyID := uuid.New().String()
	var startTime string
	if lastPassTime != nil {
		startTime = lastPassTime.Format(time.RFC3339)
	} else {
		// For the first lap, we need to get the race start time
		var raceStartTime string
		err := tx.QueryRow("SELECT time_start FROM races WHERE id = ?", raceID).Scan(&raceStartTime)
		if err != nil {
			return fmt.Errorf("failed to get race start time: %v", err)
		}
		startTime = raceStartTime
	}

	_, err = tx.Exec(`
		INSERT INTO lap_history (
			id, race_participant_id, lap_number, lap_time_ms,
			start_time, end_time, is_valid, created_at
		) VALUES (?, ?, ?, ?, ?, ?, 1, ?)`,
		historyID, participantID, newLapCount, lapTimeMS,
		startTime, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to record lap history: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (rc *Controller) startAutoFinishTimer(raceID string) {
	// First, stop any existing timer
	rc.stopAutoFinishTimer()

	// Get race time limit
	var timeLimitMinutes sql.NullInt64
	err := rc.db.QueryRow("SELECT time_limit_minutes FROM races WHERE id = ?", raceID).Scan(&timeLimitMinutes)
	if err != nil {
		log.Printf("Error getting race time limit: %v", err)
		return
	}

	if !timeLimitMinutes.Valid || timeLimitMinutes.Int64 <= 0 {
		// No time limit set, don't start timer
		return
	}

	duration := time.Duration(timeLimitMinutes.Int64) * time.Minute
	rc.activeTimer = time.AfterFunc(duration, func() {
		log.Printf("Auto-finishing race %s due to time limit", raceID)
		err := rc.StopRace(raceID)
		if err != nil {
			log.Printf("Error auto-finishing race: %v", err)
		}
	})
}

func (rc *Controller) stopAutoFinishTimer() {
	if rc.activeTimer != nil {
		rc.activeTimer.Stop()
		rc.activeTimer = nil
	}
}

func (rc *Controller) recordAuditLog(tx *sql.Tx, actionType, entityType, entityID, userName, ipAddress string, details map[string]interface{}) error {
	id := uuid.New().String()
	now := time.Now().UTC()

	detailsJSON := ""
	if details != nil {
		// Convert map to JSON string
		// For simplicity, we'll just store a basic representation
		// In a real implementation, you'd want to properly marshal the JSON
		detailsJSON = fmt.Sprintf("%v", details)
	}

	_, err := tx.Exec(`
		INSERT INTO audit_log (id, timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, now.Format(time.RFC3339), actionType, entityType, entityID, userName, ipAddress, detailsJSON, now.Format(time.RFC3339))

	return err
}

func getCurrentBestLapTime(tx *sql.Tx, participantID string) int {
	var bestLapTime int
	err := tx.QueryRow("SELECT best_lap_time_ms FROM race_laps WHERE race_participant_id = ?", participantID).Scan(&bestLapTime)
	if err != nil {
		return 0
	}
	return bestLapTime
}

func getTimeFromStr(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}