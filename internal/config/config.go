package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig holds application configuration
type AppConfig struct {
	DBPath              string `json:"db_path"`
	HardwareCOMPort     string `json:"hardware_com_port"`
	HardwareReaderType  string `json:"hardware_reader_type"`
	HardwareDebounceMS  int    `json:"hardware_debounce_ms"`
	UILanguage          string `json:"ui_language"`
	TimeLimitMinutes    int    `json:"time_limit_minutes"`
	DefaultTrackName    string `json:"default_track_name"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *AppConfig {
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)

	return &AppConfig{
		DBPath:             filepath.Join(exeDir, "trackpulse.db"),
		HardwareCOMPort:    "",
		HardwareReaderType: "EM4095",
		HardwareDebounceMS: 2000,
		UILanguage:         "en",
		TimeLimitMinutes:   5,
		DefaultTrackName:   "",
	}
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file
func (c *AppConfig) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
