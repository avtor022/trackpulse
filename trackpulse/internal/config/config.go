package config

import (
	"os"
)

type Config struct {
	DBPath       string
	RequireLogin bool
	ListenAddr   string
}

func Load() *Config {
	// Default configuration
	cfg := &Config{
		DBPath:       "./track_pulse.db",
		RequireLogin: true,
		ListenAddr:   ":8080",
	}

	// Override with environment variables if present
	if dbPath := os.Getenv("TRACKPULSE_DB_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	if addr := os.Getenv("TRACKPULSE_LISTEN_ADDR"); addr != "" {
		cfg.ListenAddr = addr
	}

	// Check if login is disabled via environment variable
	if os.Getenv("TRACKPULSE_NO_LOGIN") == "true" {
		cfg.RequireLogin = false
	}

	return cfg
}