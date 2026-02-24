package models

import "errors"
import "time"

type RaceParticipant struct {
	ID              string    `db:"id" json:"id"`
	RaceID          string    `db:"race_id" json:"race_id"`
	RacerModelID    string    `db:"racer_model_id" json:"racer_model_id"`
	GridPosition    *int      `db:"grid_position" json:"grid_position,omitempty"`
	IsFinished      bool      `db:"is_finished" json:"is_finished"`
	Disqualified    bool      `db:"disqualified" json:"disqualified"`
	DNFReason       *string   `db:"dnf_reason" json:"dnf_reason,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// TableName возвращает имя таблицы
func (RaceParticipant) TableName() string {
	return "race_participants"
}

// Validate валидация данных
func (rp *RaceParticipant) Validate() error {
	if rp.RaceID == "" || rp.RacerModelID == "" {
		return ErrRaceParticipantRequiredFields
	}
	return nil
}

var ErrRaceParticipantRequiredFields = errors.New("race_id and racer_model_id are required")