package main

import (
	"fmt"
	"os"
	"path/filepath"

	"trackpulse/internal/config"
	"trackpulse/internal/database"
	"trackpulse/internal/models"
	"trackpulse/internal/repository"
	"trackpulse/internal/service"
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

	// Initialize services
	racerService := service.NewRacerService(racerRepo)

	// Test racer CRUD operations
	log.Info("Testing Racer CRUD operations...")

	// Create a test racer
	testRacer := &models.Racer{
		FullName:    "Test Driver",
		RacerNumber: 1,
		Country:     "USA",
		City:        "New York",
		Rating:      100,
	}

	if err := racerService.CreateRacer(testRacer); err != nil {
		log.Error("Failed to create test racer: %v", err)
	} else {
		log.Info("Test racer created successfully: %s", testRacer.FullName)
	}

	// Get all racers
	racers, err := racerService.GetAllRacers()
	if err != nil {
		log.Error("Failed to get racers: %v", err)
	} else {
		log.Info("Total racers in database: %d", len(racers))
		for _, r := range racers {
			log.Info("  - #%d: %s (%s, %s)", r.RacerNumber, r.FullName, r.City, r.Country)
		}
	}

	log.Info("TrackPulse initialization complete!")
	fmt.Println("TrackPulse initialized successfully!")
	fmt.Printf("Database: %s\n", cfg.DBPath)
	fmt.Printf("Language: %s\n", cfg.UILanguage)
	fmt.Printf("Total racers: %d\n", len(racers))
}
