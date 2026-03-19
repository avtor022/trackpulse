package models

import "time"

// Athlete represents a pilot/participant in competitions
type Athlete struct {
	ID           string     `json:"id" db:"id"`
	RacerNumber  int        `json:"racer_number" db:"racer_number"`
	FullName     string     `json:"full_name" db:"full_name"`
	Birthday     *time.Time `json:"birthday,omitempty" db:"birthday"`
	Country      string     `json:"country" db:"country"`
	City         string     `json:"city" db:"city"`
	Rating       int        `json:"rating" db:"rating"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// RCModel represents an RC model in the catalog
type RCModel struct {
	ID         string    `json:"id" db:"id"`
	Brand      string    `json:"brand" db:"brand"`
	ModelName  string    `json:"model_name" db:"model_name"`
	Scale      string    `json:"scale" db:"scale"`
	ModelType  string    `json:"model_type" db:"model_type"`
	MotorType  string    `json:"motor_type" db:"motor_type"`
	DriveType  string    `json:"drive_type" db:"drive_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// AthleteModel links athlete with model and transponder
type AthleteModel struct {
	ID               string    `json:"id" db:"id"`
	AthleteID        string    `json:"athlete_id" db:"athlete_id"`
	RCModelID        string    `json:"rc_model_id" db:"rc_model_id"`
	AthleteModelNumber string   `json:"athlete_model_number" db:"athlete_model_number"`
	AthleteModelType  string    `json:"athlete_model_type" db:"athlete_model_type"`
	TransponderNumber string   `json:"transponder_number" db:"transponder_number"`
	TransponderType  string    `json:"transponder_type" db:"transponder_type"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Race represents a race/heat event
type Race struct {
	ID                string     `json:"id" db:"id"`
	RaceTitle         string     `json:"race_title" db:"race_title"`
	RaceType          string     `json:"race_type" db:"race_type"`
	ModelType         string     `json:"model_type" db:"model_type"`
	ModelScale        string     `json:"model_scale" db:"model_scale"`
	TrackName         string     `json:"track_name" db:"track_name"`
	LapCountTarget    *int       `json:"lap_count_target,omitempty" db:"lap_count_target"`
	TimeLimitMinutes  *int       `json:"time_limit_minutes,omitempty" db:"time_limit_minutes"`
	TimeStart         *time.Time `json:"time_start,omitempty" db:"time_start"`
	TimeFinish        *time.Time `json:"time_finish,omitempty" db:"time_finish"`
	Status            string     `json:"status" db:"status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// RaceParticipant links race with athlete-model
type RaceParticipant struct {
	ID              string    `json:"id" db:"id"`
	RaceID          string    `json:"race_id" db:"race_id"`
	AthleteModelID  string    `json:"athlete_model_id" db:"athlete_model_id"`
	GridPosition    *int      `json:"grid_position,omitempty" db:"grid_position"`
	IsFinished      bool      `json:"is_finished" db:"is_finished"`
	Disqualified    bool      `json:"disqualified" db:"disqualified"`
	DNFReason       string    `json:"dnf_reason" db:"dnf_reason"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// RaceLaps aggregated results for a participant in a race
type RaceLaps struct {
	ID                 string    `json:"id" db:"id"`
	RaceParticipantID  string    `json:"race_participant_id" db:"race_participant_id"`
	TimeStart          time.Time `json:"time_start" db:"time_start"`
	TimeFinish         *time.Time `json:"time_finish,omitempty" db:"time_finish"`
	NumberOfLaps       int       `json:"number_of_laps" db:"number_of_laps"`
	BestLapTimeMs      int       `json:"best_lap_time_ms" db:"best_lap_time_ms"`
	BestLapNumber      int       `json:"best_lap_number" db:"best_lap_number"`
	LastLapTimeMs      int       `json:"last_lap_time_ms" db:"last_lap_time_ms"`
	LastPassTime       *time.Time `json:"last_pass_time,omitempty" db:"last_pass_time"`
	TotalRaceTimeMs    int       `json:"total_race_time_ms" db:"total_race_time_ms"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// LapHistory detailed history of each lap
type LapHistory struct {
	ID                  string    `json:"id" db:"id"`
	RaceParticipantID   string    `json:"race_participant_id" db:"race_participant_id"`
	LapNumber           int       `json:"lap_number" db:"lap_number"`
	LapTimeMs           int       `json:"lap_time_ms" db:"lap_time_ms"`
	StartTime           time.Time `json:"start_time" db:"start_time"`
	EndTime             time.Time `json:"end_time" db:"end_time"`
	IsValid             bool      `json:"is_valid" db:"is_valid"`
	InvalidationReason  string    `json:"invalidation_reason" db:"invalidation_reason"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// RawScan raw RFID scan log
type RawScan struct {
	ID                 string    `json:"id" db:"id"`
	Timestamp          time.Time `json:"timestamp" db:"timestamp"`
	TagValue           string    `json:"tag_value" db:"tag_value"`
	ReaderType         string    `json:"reader_type" db:"reader_type"`
	COMPort            string    `json:"com_port" db:"com_port"`
	SignalStrength     *int      `json:"signal_strength,omitempty" db:"signal_strength"`
	IsProcessed        bool      `json:"is_processed" db:"is_processed"`
	LinkedAthleteModelID *string   `json:"linked_athlete_model_id,omitempty" db:"linked_athlete_model_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// SystemSetting system configuration parameter
type SystemSetting struct {
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	ValueType   string    `json:"value_type" db:"value_type"`
	Description string    `json:"description" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AuditLog audit trail entry
type AuditLog struct {
	ID          string    `json:"id" db:"id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	ActionType  string    `json:"action_type" db:"action_type"`
	EntityType  string    `json:"entity_type" db:"entity_type"`
	EntityID    string    `json:"entity_id" db:"entity_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	Details     string    `json:"details" db:"details"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
