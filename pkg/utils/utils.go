package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4 string
func GenerateUUID() string {
	return uuid.New().String()
}

// GetCurrentTimestamp returns current time in RFC3339 format
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// ParseTimestamp parses RFC3339 timestamp to time.Time
func ParseTimestamp(ts string) (time.Time, error) {
	return time.Parse(time.RFC3339, ts)
}

// FormatDuration formats milliseconds to human readable string
func FormatDuration(ms int) string {
	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}
	seconds := float64(ms) / 1000.0
	return fmt.Sprintf("%.3fs", seconds)
}
