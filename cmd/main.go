package main

import (
	"fmt"
	"os"
	"path/filepath"

	"trackpulse/internal/config"
	"trackpulse/internal/database"
	"trackpulse/internal/locale"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
	"trackpulse/internal/ui"
	"trackpulse/pkg/logger"
)

func main() {
	// Initialize logger
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")

	log, err := logger.NewLogger(logDir)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Info("TrackPulse starting...")

	// Load configuration
	cfgPath := filepath.Join(exeDir, "config.json")
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Error("Failed to load config: %v", err)
		os.Exit(1)
	}
	log.Info("Configuration loaded from %s", cfgPath)

	// Initialize database
	db, err := database.NewDB(cfg.DBPath)
	if err != nil {
		log.Error("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("Database opened: %s", cfg.DBPath)

	// Initialize schema
	if err := db.Initialize(); err != nil {
		log.Error("Failed to initialize database schema: %v", err)
		os.Exit(1)
	}
	log.Info("Database schema initialized")

	// Initialize repositories
	competitorRepo := repository.NewCompetitorRepository(db.DB)
	modelRepo := repository.NewRCModelRepository(db.DB)
	brandRepo := repository.NewRCModelBrandRepository(db.DB)
	scaleRepo := repository.NewRCModelScaleRepository(db.DB)
	typeRepo := repository.NewRCModelTypeRepository(db.DB)
	trackRepo := repository.NewCompetitionTrackRepository(db.DB)
	settingsRepo := repository.NewSettingsRepository(db.DB)
	competitorModelRepo := repository.NewCompetitorModelRepository(db.DB)
	competitionRepo := repository.NewCompetitionRepository(db.DB)

	// Initialize services
	competitorService := service.NewCompetitorService(competitorRepo)
	modelService := service.NewRCModelService(modelRepo, brandRepo, scaleRepo, typeRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	competitorModelService := service.NewCompetitorModelService(competitorModelRepo, competitorRepo, modelRepo)
	competitionService := service.NewCompetitionService(competitionRepo, typeRepo, scaleRepo, trackRepo)

	// Load locale from settings
	savedLocale, err := settingsService.GetLocale()
	if err != nil {
		log.Error("Failed to load locale from settings: %v", err)
		savedLocale = "en"
	}
	
	// Initialize locale with saved value
	locale.SetLocale(savedLocale)
	log.Info("Locale loaded: %s", savedLocale)

	log.Info("TrackPulse initialization complete!")

	// Start UI
	uiApp := ui.NewApp(competitorService, modelService, settingsService, competitorModelService, competitionService, savedLocale)
	uiApp.Run()
}
