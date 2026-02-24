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

