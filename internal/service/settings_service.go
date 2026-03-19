package service

import (
	"trackpulse/internal/repository"
)

// SettingsService handles business logic for settings
type SettingsService struct {
	repo *repository.SettingsRepository
}

// NewSettingsService creates a new settings service
func NewSettingsService(repo *repository.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

// GetLocale retrieves the locale setting from database
func (s *SettingsService) GetLocale() (string, error) {
	value, err := s.repo.GetValue("locale")
	if err != nil {
		return "en", err
	}
	if value == "" {
		return "en", nil
	}
	return value, nil
}

// SetLocale saves the locale setting to database
func (s *SettingsService) SetLocale(locale string) error {
	return s.repo.UpdateValue("locale", locale)
}
