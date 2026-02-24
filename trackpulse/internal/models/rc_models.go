package models

import "errors"
import "time"

type RCModel struct {
	ID         string    `db:"id" json:"id"`
	Brand      string    `db:"brand" json:"brand"`
	ModelName  string    `db:"model_name" json:"model_name"`
	Scale      string    `db:"scale" json:"scale"`
	ModelType  string    `db:"model_type" json:"model_type"`
	MotorType  *string   `db:"motor_type" json:"motor_type,omitempty"`
	DriveType  *string   `db:"drive_type" json:"drive_type,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// TableName возвращает имя таблицы
func (RCModel) TableName() string {
	return "rc_models"
}

// Validate валидация данных
func (m *RCModel) Validate() error {
	if m.Brand == "" || m.ModelName == "" {
		return ErrModelRequiredFields
	}
	return nil
}

var ErrModelRequiredFields = errors.New("brand and model name are required")