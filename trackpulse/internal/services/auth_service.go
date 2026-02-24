package services

import (
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
	"trackpulse/internal/security"
)

type AuthService interface {
	Login(username, password string) bool
	Logout()
	IsAuthenticated() bool
	ChangePassword(oldPass, newPass string) error
	ChangeUsername(newUsername string) error
}

type authService struct {
	settingsRepo repository.SystemSettingsRepository
	session      *security.Session
}

func NewAuthService(settingsRepo repository.SystemSettingsRepository) AuthService {
	return &authService{
		settingsRepo: settingsRepo,
		session:      nil,
	}
}

func (s *authService) Login(username, password string) bool {
	// Получение хеша пароля из БД
	storedHash, err := s.settingsRepo.Get(models.SettingKeyAuthPassword)
	if err != nil {
		return false
	}

	storedUser, err := s.settingsRepo.Get(models.SettingKeyAuthUser)
	if err != nil {
		return false
	}

	// Проверка логина
	if username != storedUser.Value {
		return false
	}

	// Проверка пароля
	if err := security.ComparePassword(password, storedHash.Value); err != nil {
		return false
	}

	// Создание сессии
	s.session = security.NewSession(storedUser.Value)
	return true
}

func (s *authService) Logout() {
	if s.session != nil {
		s.session.Logout()
		s.session = nil
	}
}

func (s *authService) IsAuthenticated() bool {
	if s.session == nil {
		return false
	}
	return s.session.IsValid()
}

func (s *authService) ChangePassword(oldPass, newPass string) error {
	// Проверяем текущий пароль
	currentHash, err := s.settingsRepo.Get(models.SettingKeyAuthPassword)
	if err != nil {
		return err
	}

	// Проверяем текущий логин и пароль
	if err := security.ComparePassword(oldPass, currentHash.Value); err != nil {
		return err
	}

	// Хешируем новый пароль
	newHash, err := security.HashPassword(newPass)
	if err != nil {
		return err
	}

	// Обновляем пароль в базе
	setting := &models.SystemSetting{
		Key:       models.SettingKeyAuthPassword,
		Value:     newHash,
		ValueType: "string",
	}

	return s.settingsRepo.Update(setting)
}

func (s *authService) ChangeUsername(newUsername string) error {
	if !s.IsAuthenticated() {
		return ErrNotAuthorized
	}

	setting := &models.SystemSetting{
		Key:       models.SettingKeyAuthUser,
		Value:     newUsername,
		ValueType: "string",
	}

	return s.settingsRepo.Update(setting)
}

import "errors"

var ErrNotAuthorized = errors.New("not authorized")