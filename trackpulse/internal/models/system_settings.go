package models

import "time"

type SystemSetting struct {
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	ValueType   string    `db:"value_type" json:"value_type"`
	Description *string   `db:"description" json:"description,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TableName возвращает имя таблицы
func (SystemSetting) TableName() string {
	return "system_settings"
}

// Ключи настроек
const (
	SettingKeyUILanguage      = "ui_language"
	SettingKeyDBPath          = "db_path"
	SettingKeyHardwareCOMPort = "hardware_com_port"
	SettingKeyHardwareReader  = "hardware_reader_type"
	SettingKeyDebounceMs      = "hardware_debounce_ms"
	SettingKeyAuthUser        = "auth_user"
	SettingKeyAuthPassword    = "auth_password_hash"
	SettingKeyLogRetention    = "log_retention_years"
)

// Validate валидация данных
func (ss *SystemSetting) Validate() error {
	if ss.Key == "" {
		return ErrSystemSettingKeyRequired
	}
	return nil
}

var ErrSystemSettingKeyRequired = "key is required"