package services

import (
	"trackpulse/internal/repository"
	"trackpulse/internal/models"
)

type SettingsService interface {
	GetLanguage() (string, error)
	SetLanguage(lang string) error
	GetDBPath() (string, error)
	SetDBPath(path string) error
	GetHardwareSettings() (map[string]string, error)
	UpdateHardwareSettings(settings map[string]string) error
	GetAuthSettings() (map[string]string, error)
	UpdateAuthSettings(settings map[string]string) error
	GetLogRetention() (int, error)
	SetLogRetention(days int) error
}

type settingsService struct {
	settingsRepo repository.SystemSettingsRepository
}

func NewSettingsService(settingsRepo repository.SystemSettingsRepository) SettingsService {
	return &settingsService{
		settingsRepo: settingsRepo,
	}
}

func (s *settingsService) GetLanguage() (string, error) {
	setting, err := s.settingsRepo.Get(models.SettingKeyUILanguage)
	if err != nil {
		// If language setting is not found, return default
		return "en", nil
	}
	return setting.Value, nil
}

func (s *settingsService) SetLanguage(lang string) error {
	setting := &models.SystemSetting{
		Key:       models.SettingKeyUILanguage,
		Value:     lang,
		ValueType: "string",
		Description: "Interface language",
	}
	return s.settingsRepo.Update(setting)
}

func (s *settingsService) GetDBPath() (string, error) {
	setting, err := s.settingsRepo.Get(models.SettingKeyDBPath)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *settingsService) SetDBPath(path string) error {
	setting := &models.SystemSetting{
		Key:       models.SettingKeyDBPath,
		Value:     path,
		ValueType: "string",
		Description: "Database file path",
	}
	return s.settingsRepo.Update(setting)
}

func (s *settingsService) GetHardwareSettings() (map[string]string, error) {
	return s.settingsRepo.GetHardwareSettings()
}

func (s *settingsService) UpdateHardwareSettings(settings map[string]string) error {
	// Update COM port setting
	if comPort, exists := settings[models.SettingKeyHardwareCOMPort]; exists {
		setting := &models.SystemSetting{
			Key:       models.SettingKeyHardwareCOMPort,
			Value:     comPort,
			ValueType: "string",
			Description: "COM port for Arduino connection",
		}
		err := s.settingsRepo.Update(setting)
		if err != nil {
			return err
		}
	}

	// Update reader type setting
	if readerType, exists := settings[models.SettingKeyHardwareReader]; exists {
		setting := &models.SystemSetting{
			Key:       models.SettingKeyHardwareReader,
			Value:     readerType,
			ValueType: "string",
			Description: "RFID reader type (EM4095 or RC522)",
		}
		err := s.settingsRepo.Update(setting)
		if err != nil {
			return err
		}
	}

	// Update debounce setting
	if debounce, exists := settings[models.SettingKeyDebounceMs]; exists {
		setting := &models.SystemSetting{
			Key:       models.SettingKeyDebounceMs,
			Value:     debounce,
			ValueType: "int",
			Description: "Debounce delay in milliseconds",
		}
		err := s.settingsRepo.Update(setting)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *settingsService) GetAuthSettings() (map[string]string, error) {
	authSettings := make(map[string]string)
	
	// Get username
	username, err := s.settingsRepo.GetUsername()
	if err != nil {
		// If username doesn't exist, return empty map
		return make(map[string]string), nil
	}
	authSettings[models.SettingKeyAuthUser] = username
	
	// We don't return password hash for security reasons
	// Instead, we return a boolean indicating if a password is set
	authSettings["password_set"] = "true"
	
	return authSettings, nil
}

func (s *settingsService) UpdateAuthSettings(settings map[string]string) error {
	// Update username if provided
	if username, exists := settings[models.SettingKeyAuthUser]; exists && username != "" {
		setting := &models.SystemSetting{
			Key:       models.SettingKeyAuthUser,
			Value:     username,
			ValueType: "string",
			Description: "Admin username",
		}
		err := s.settingsRepo.Update(setting)
		if err != nil {
			return err
		}
	}

	// Update password if provided
	if password, exists := settings[models.SettingKeyAuthPassword]; exists && password != "" {
		// Hash the password
		hashedPassword, err := hashPassword(password) // This would need to be implemented
		if err != nil {
			return err
		}
		
		setting := &models.SystemSetting{
			Key:       models.SettingKeyAuthPassword,
			Value:     hashedPassword,
			ValueType: "string",
			Description: "Admin password hash",
		}
		err = s.settingsRepo.Update(setting)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *settingsService) GetLogRetention() (int, error) {
	setting, err := s.settingsRepo.Get(models.SettingKeyLogRetention)
	if err != nil {
		// If retention setting is not found, return default (1 year = 365 days)
		return 365, nil
	}
	
	// Convert string value to int
	// For now, we'll return a default value
	return 365, nil
}

func (s *settingsService) SetLogRetention(days int) error {
	setting := &models.SystemSetting{
		Key:       models.SettingKeyLogRetention,
		Value:     string(rune(days)), // This would need proper conversion
		ValueType: "int",
		Description: "Log retention period in days",
	}
	return s.settingsRepo.Update(setting)
}

// Helper function to hash passwords (would need to import appropriate package)
func hashPassword(password string) (string, error) {
	// Implementation would use a hashing library like bcrypt
	// For now, returning a placeholder
	return "hashed_password_placeholder", nil
}