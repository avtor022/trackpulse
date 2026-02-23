package models

import (
	"time"
)

// Racer represents a single racer
type Racer struct {
	ID          string    `json:"id" db:"id"`
	RacerNumber int       `json:"racer_number" db:"racer_number"`
	FullName    string    `json:"full_name" db:"full_name"`
	Birthday    *string   `json:"birthday,omitempty" db:"birthday"`
	Country     *string   `json:"country,omitempty" db:"country"`
	City        *string   `json:"city,omitempty" db:"city"`
	Rating      int       `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RCModel represents an RC model
type RCModel struct {
	ID         string    `json:"id" db:"id"`
	Brand      string    `json:"brand" db:"brand"`
	ModelName  string    `json:"model_name" db:"model_name"`
	Scale      string    `json:"scale" db:"scale"`
	ModelType  string    `json:"model_type" db:"model_type"`
	MotorType  *string   `json:"motor_type,omitempty" db:"motor_type"`
	DriveType  *string   `json:"drive_type,omitempty" db:"drive_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// RacerModel represents the link between a racer and their model
type RacerModel struct {
	ID                string    `json:"id" db:"id"`
	RacerID           string    `json:"racer_id" db:"racer_id"`
	RCModelID         string    `json:"rc_model_id" db:"rc_model_id"`
	TransponderNumber string    `json:"transponder_number" db:"transponder_number"`
	TransponderType   string    `json:"transponder_type" db:"transponder_type"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Race represents a single race event
type Race struct {
	ID               string     `json:"id" db:"id"`
	RaceTitle        string     `json:"race_title" db:"race_title"`
	RaceType         string     `json:"race_type" db:"race_type"`
	ModelType        *string    `json:"model_type,omitempty" db:"model_type"`
	ModelScale       *string    `json:"model_scale,omitempty" db:"model_scale"`
	TrackName        *string    `json:"track_name,omitempty" db:"track_name"`
	LapCountTarget   *int       `json:"lap_count_target,omitempty" db:"lap_count_target"`
	TimeLimitMinutes *int       `json:"time_limit_minutes,omitempty" db:"time_limit_minutes"`
	TimeStart        *time.Time `json:"time_start,omitempty" db:"time_start"`
	TimeFinish       *time.Time `json:"time_finish,omitempty" db:"time_finish"`
	Status           string     `json:"status" db:"status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// RaceParticipant represents a participant in a race
type RaceParticipant struct {
	ID             string    `json:"id" db:"id"`
	RaceID         string    `json:"race_id" db:"race_id"`
	RacerModelID   string    `json:"racer_model_id" db:"racer_model_id"`
	GridPosition   *int      `json:"grid_position,omitempty" db:"grid_position"`
	IsFinished     bool      `json:"is_finished" db:"is_finished"`
	Disqualified   bool      `json:"disqualified" db:"disqualified"`
	DNFReason      *string   `json:"dnf_reason,omitempty" db:"dnf_reason"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// RaceLap represents aggregated lap data for a participant
type RaceLap struct {
	ID                 string     `json:"id" db:"id"`
	RaceParticipantID  string     `json:"race_participant_id" db:"race_participant_id"`
	TimeStart          time.Time  `json:"time_start" db:"time_start"`
	TimeFinish         *time.Time `json:"time_finish,omitempty" db:"time_finish"`
	NumberOfLaps       int        `json:"number_of_laps" db:"number_of_laps"`
	BestLapTimeMS      int        `json:"best_lap_time_ms" db:"best_lap_time_ms"`
	BestLapNumber      int        `json:"best_lap_number" db:"best_lap_number"`
	LastLapTimeMS      int        `json:"last_lap_time_ms" db:"last_lap_time_ms"`
	LastPassTime       *time.Time `json:"last_pass_time,omitempty" db:"last_pass_time"`
	TotalRaceTimeMS    int        `json:"total_race_time_ms" db:"total_race_time_ms"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// LapHistory represents individual lap records
type LapHistory struct {
	ID                  string    `json:"id" db:"id"`
	RaceParticipantID   string    `json:"race_participant_id" db:"race_participant_id"`
	LapNumber           int       `json:"lap_number" db:"lap_number"`
	LapTimeMS           int       `json:"lap_time_ms" db:"lap_time_ms"`
	StartTime           time.Time `json:"start_time" db:"start_time"`
	EndTime             time.Time `json:"end_time" db:"end_time"`
	IsValid             bool      `json:"is_valid" db:"is_valid"`
	InvalidationReason  *string   `json:"invalidation_reason,omitempty" db:"invalidation_reason"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// RawScan represents a raw RFID scan
type RawScan struct {
	ID                 string    `json:"id" db:"id"`
	Timestamp          time.Time `json:"timestamp" db:"timestamp"`
	TagValue           string    `json:"tag_value" db:"tag_value"`
	ReaderType         string    `json:"reader_type" db:"reader_type"`
	COMPort            *string   `json:"com_port,omitempty" db:"com_port"`
	SignalStrength     *int      `json:"signal_strength,omitempty" db:"signal_strength"`
	IsProcessed        bool      `json:"is_processed" db:"is_processed"`
	LinkedRacerModelID *string   `json:"linked_racer_model_id,omitempty" db:"linked_racer_model_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         string    `json:"id" db:"id"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	ActionType string    `json:"action_type" db:"action_type"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   *string   `json:"entity_id,omitempty" db:"entity_id"`
	UserName   *string   `json:"user_name,omitempty" db:"user_name"`
	IPAddress  *string   `json:"ip_address,omitempty" db:"ip_address"`
	Details    *string   `json:"details,omitempty" db:"details"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// SystemSetting represents a system setting
type SystemSetting struct {
	Key        string    `json:"key" db:"key"`
	Value      string    `json:"value" db:"value"`
	ValueType  string    `json:"value_type" db:"value_type"`
	Description *string  `json:"description,omitempty" db:"description"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ParticipantWithDetails extends RaceParticipant with related data
type ParticipantWithDetails struct {
	ID                string  `json:"id"`
	RaceID            string  `json:"race_id"`
	RacerNumber       int     `json:"racer_number"`
	FullName          string  `json:"full_name"`
	ModelName         string  `json:"model_name"`
	Brand             string  `json:"brand"`
	TransponderNumber string  `json:"transponder_number"`
	GridPosition      *int    `json:"grid_position,omitempty"`
	IsFinished        bool    `json:"is_finished"`
	Disqualified      bool    `json:"disqualified"`
	DNFReason         *string `json:"dnf_reason,omitempty"`
}

// LeaderboardEntry represents an entry in the leaderboard
type LeaderboardEntry struct {
	Position        int     `json:"position"`
	RaceParticipantID string `json:"race_participant_id"`
	RacerNumber     int     `json:"racer_number"`
	FullName        string  `json:"full_name"`
	ModelName       string  `json:"model_name"`
	Brand           string  `json:"brand"`
	NumberOfLaps    int     `json:"number_of_laps"`
	BestLapTimeMS   int     `json:"best_lap_time_ms"`
	LastLapTimeMS   int     `json:"last_lap_time_ms"`
	TotalRaceTimeMS int     `json:"total_race_time_ms"`
	IsActive        bool    `json:"is_active"`
}

// LapUpdateMessage represents a WebSocket message for lap updates
type LapUpdateMessage struct {
	Event            string `json:"event"`
	RaceID           string `json:"race_id"`
	ParticipantID    string `json:"participant_id"`
	LapNumber        int    `json:"lap_number"`
	LapTimeMS        int    `json:"lap_time_ms"`
	TotalLaps        int    `json:"total_laps"`
	BestLapMS        int    `json:"best_lap_ms"`
	Position         int    `json:"position"`
	Timestamp        string `json:"timestamp"`
}

// RaceUpdateMessage represents a WebSocket message for race updates
type RaceUpdateMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}