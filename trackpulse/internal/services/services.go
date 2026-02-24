package services

import (
	"trackpulse/internal/config"
	"trackpulse/internal/repository"
)

// Services holds all service instances
type Services struct {
	AuthService    AuthService
	RaceService    RaceService
	LapService     LapService
	HardwareService HardwareService
	LoggerService  LoggerService
	SettingsService SettingsService
}

// NewServices creates and returns all services
func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	authSvc := NewAuthService(repos.SystemSettings)
	settingsSvc := NewSettingsService(repos.SystemSettings)

	return &Services{
		AuthService:    authSvc,
		RaceService:    NewRaceService(repos, cfg),
		LapService:     NewLapService(repos, cfg),
		HardwareService: NewHardwareService(cfg),
		LoggerService:  NewLoggerService(repos.AuditLog),
		SettingsService: settingsSvc,
	}
}

// Interfaces for all services
type AuthService interface {
	Login(username, password string) bool
	Logout()
	IsAuthenticated() bool
	ChangePassword(oldPass, newPass string) error
	ChangeUsername(newUsername string) error
}

type RaceService interface {
	StartRace(raceID string) error
	StopRace(raceID string) error
	GetLiveStandings(raceID string) ([]StandingsRow, error)
}

type LapService interface {
	ProcessScan(tagValue string, readerType string, comPort string) error
	StartRace(raceID string) error
	StopRace(raceID string) error
	GetLiveStandings(raceID string) ([]StandingsRow, error)
}

type HardwareService interface {
	Initialize() error
	ReadTag() (string, error)
	Close() error
}

type LoggerService interface {
	LogEvent(actionType, entityType string, entityID *string, details *string) error
}

type SettingsService interface {
	GetLanguage() (string, error)
	SetLanguage(lang string) error
	GetDBPath() (string, error)
	SetDBPath(path string) error
	GetHardwareSettings() (map[string]string, error)
	UpdateHardwareSettings(settings map[string]string) error
}