package main

import (
	"fmt"
	"os"
	"path/filepath"

	"trackpulse/internal/config"
	"trackpulse/internal/database"
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
	racerRepo := repository.NewRacerRepository(db.DB)
	modelRepo := repository.NewRCModelRepository(db.DB)
	brandRepo := repository.NewRCModelBrandRepository(db.DB)

	// Initialize services
	racerService := service.NewRacerService(racerRepo)
	modelService := service.NewRCModelService(modelRepo, brandRepo)

	log.Info("TrackPulse initialization complete!")

	// Start UI
	fmt.Println("Starting TrackPulse UI...")
	fmt.Printf("Database: %s\n", cfg.DBPath)
	fmt.Printf("Language: %s\n", cfg.UILanguage)

	uiApp := ui.NewApp(racerService, modelService, cfg.UILanguage)
	
	// Add error handling for UI startup
	defer func() {
		if r := recover(); r != nil {
			log.Error("UI panic: %v", r)
			fmt.Printf("UI Error: %v\n", r)
		}
	}()
	
	uiApp.Run()
}
