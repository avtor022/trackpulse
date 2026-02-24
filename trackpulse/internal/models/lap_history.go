package models

import (
	"errors"
	"time"
)

type LapHistory struct {
	ID                  string    `db:"id" json:"id"`
	RaceParticipantID   string    `db:"race_participant_id" json:"race_participant_id"`
	LapNumber           int       `db:"lap_number" json:"lap_number"`
	LapTimeMs           int       `db:"lap_time_ms" json:"lap_time_ms"`
	StartTime           time.Time `db:"start_time" json:"start_time"`
	EndTime             time.Time `db:"end_time" json:"end_time"`
	IsValid             bool      `db:"is_valid" json:"is_valid"`
	InvalidationReason  *string   `db:"invalidation_reason" json:"invalidation_reason,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

// TableName возвращает имя таблицы
func (LapHistory) TableName() string {
	return "lap_history"
}

// Validate валидация данных
func (lh *LapHistory) Validate() error {
	if lh.RaceParticipantID == "" || lh.LapNumber <= 0 || lh.LapTimeMs <= 0 {
		return ErrLapHistoryRequiredFields
	}
	return nil
}

var ErrLapHistoryRequiredFields = errors.New("race_participant_id, lap_number and lap_time_ms are required")