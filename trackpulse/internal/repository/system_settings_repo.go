package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type SystemSettingsRepository interface {
	Get(key string) (*models.SystemSetting, error)
	Set(setting *models.SystemSetting) error
	Update(setting *models.SystemSetting) error
	Delete(key string) error
	GetAll() ([]models.SystemSetting, error)
	GetPasswordHash() (string, error)
	GetUsername() (string, error)
	GetLanguage() (string, error)
	GetHardwareSettings() (map[string]string, error)
}

type systemSettingsRepo struct {
	db *sql.DB
}

func NewSystemSettingsRepository(db *sql.DB) SystemSettingsRepository {
	return &systemSettingsRepo{db: db}
}

func (s *systemSettingsRepo) Get(key string) (*models.SystemSetting, error) {
	query := `SELECT key, value, value_type, description, updated_at 
						FROM system_settings WHERE key = ?`

	var setting models.SystemSetting

	err := s.db.QueryRow(query, key).Scan(
		&setting.Key,
		&setting.Value,
		&setting.ValueType,
		&setting.Description,
		&setting.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("setting not found")
		}
		return nil, err
	}

	return &setting, nil
}

func (s *systemSettingsRepo) Set(setting *models.SystemSetting) error {
	query := `INSERT INTO system_settings (key, value, value_type, description, updated_at)
						VALUES (?, ?, ?, ?, datetime('now'))
						ON CONFLICT(key) DO UPDATE SET 
						value = excluded.value,
						value_type = excluded.value_type,
						description = excluded.description,
						updated_at = excluded.updated_at`

	_, err := s.db.Exec(query,
		setting.Key,
		setting.Value,
		setting.ValueType,
		setting.Description,
	)
	return err
}

func (s *systemSettingsRepo) Update(setting *models.SystemSetting) error {
	// Same as Set since we're using ON CONFLICT clause
	return s.Set(setting)
}

func (s *systemSettingsRepo) Delete(key string) error {
	query := `DELETE FROM system_settings WHERE key = ?`

	result, err := s.db.Exec(query, key)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("setting not found")
	}

	return nil
}

func (s *systemSettingsRepo) GetAll() ([]models.SystemSetting, error) {
	query := `SELECT key, value, value_type, description, updated_at 
						FROM system_settings ORDER BY key`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SystemSetting
	for rows.Next() {
		var setting models.SystemSetting

		err := rows.Scan(
			&setting.Key,
			&setting.Value,
			&setting.ValueType,
			&setting.Description,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		settings = append(settings, setting)
	}

	return settings, nil
}

func (s *systemSettingsRepo) GetPasswordHash() (string, error) {
	setting, err := s.Get(models.SettingKeyAuthPassword)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *systemSettingsRepo) GetUsername() (string, error) {
	setting, err := s.Get(models.SettingKeyAuthUser)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *systemSettingsRepo) GetLanguage() (string, error) {
	setting, err := s.Get(models.SettingKeyUILanguage)
	if err != nil {
		// If language setting is not found, return default
		return "en", nil
	}
	return setting.Value, nil
}

func (s *systemSettingsRepo) GetHardwareSettings() (map[string]string, error) {
	settingsMap := make(map[string]string)
	
	hardwareKeys := []string{
		models.SettingKeyHardwareCOMPort,
		models.SettingKeyHardwareReader,
		models.SettingKeyDebounceMs,
	}
	
	for _, key := range hardwareKeys {
		setting, err := s.Get(key)
		if err != nil {
			// If a hardware setting is missing, we'll return an empty string
			settingsMap[key] = ""
		} else {
			settingsMap[key] = setting.Value
		}
	}
	
	return settingsMap, nil
}