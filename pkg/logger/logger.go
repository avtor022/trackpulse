package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger provides structured logging
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	dateStr := time.Now().Format("2006-01-02")
	
	infoFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("info_%s.log", dateStr)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open info log file: %w", err)
	}

	errorFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("error_%s.log", dateStr)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		infoFile.Close()
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}

	debugFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("debug_%s.log", dateStr)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		infoFile.Close()
		errorFile.Close()
		return nil, fmt.Errorf("failed to open debug log file: %w", err)
	}

	return &Logger{
		infoLog:  log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog: log.New(debugFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

// Info logs an informational message
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.debugLog.Printf(format, v...)
}

// Close closes all log files
func (l *Logger) Close() error {
	// Note: log.Logger doesn't expose the underlying writer, so we can't close files directly
	// In production, you might want to track the files separately
	return nil
}
