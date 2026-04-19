package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// RFIDLogger handles reading tag IDs from a serial device and logging them to a file
type RFIDLogger struct {
	port        io.ReadCloser
	isConnected bool
	logFile     *os.File
	mu          sync.Mutex
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// NewRFIDLogger creates a new RFIDLogger instance
func NewRFIDLogger() *RFIDLogger {
	return &RFIDLogger{
		isConnected: false,
		stopChan:    make(chan struct{}),
	}
}

// Connect opens the serial port and starts reading data
func (r *RFIDLogger) Connect(portName string, baudRate uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isConnected {
		return fmt.Errorf("already connected")
	}

	// Open log file for appending
	logFile, err := os.OpenFile("rfid_tags.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	r.logFile = logFile

	// Open serial port
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		logFile.Close()
		return fmt.Errorf("failed to open serial port: %v", err)
	}

	r.port = port
	r.isConnected = true
	r.stopChan = make(chan struct{})

	// Start reading data in a goroutine
	r.wg.Add(1)
	go r.readData()

	return nil
}

// Disconnect closes the serial port and log file
func (r *RFIDLogger) Disconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isConnected {
		return nil
	}

	// Signal to stop reading
	close(r.stopChan)
	r.wg.Wait()

	if r.port != nil {
		r.port.Close()
		r.port = nil
	}

	if r.logFile != nil {
		r.logFile.Close()
		r.logFile = nil
	}

	r.isConnected = false
	return nil
}

// readData reads from the serial port and logs tag IDs
func (r *RFIDLogger) readData() {
	defer r.wg.Done()

	scanner := bufio.NewScanner(r.port)
	for scanner.Scan() {
		select {
		case <-r.stopChan:
			return
		default:
			line := scanner.Text()
			tagID := r.extractTagID(line)
			if tagID != "" {
				r.writeLog(tagID)
			}
		}
	}
}

// extractTagID extracts the tag ID from the received line
// This function should be customized based on the actual device output format
func (r *RFIDLogger) extractTagID(line string) string {
	// Trim whitespace
	line = strings.TrimSpace(line)

	// Common RFID reader formats:
	// 1. Just the ID: "04A3B2C1"
	// 2. With prefix: "UID: 04A3B2C1" or "Tag ID: 04A3B2C1"
	// 3. Hex values: "04 A3 B2 C1"

	// Try to find hex pattern
	fields := strings.Fields(line)
	for _, field := range fields {
		// Remove common prefixes
		field = strings.TrimPrefix(field, "UID:")
		field = strings.TrimPrefix(field, "ID:")
		field = strings.TrimPrefix(field, "Tag:")
		field = strings.TrimSpace(field)

		// Check if it's a valid hex string (tag ID)
		if r.isValidHexTagID(field) {
			return strings.ToUpper(field)
		}
	}

	// If the whole line looks like a tag ID
	if r.isValidHexTagID(line) {
		return strings.ToUpper(line)
	}

	return ""
}

// isValidHexTagID checks if a string is a valid hexadecimal tag ID
func (r *RFIDLogger) isValidHexTagID(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Remove spaces and colons
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ":", "")

	// Tag IDs are typically 4-10 hex characters
	if len(s) < 4 || len(s) > 20 {
		return false
	}

	// Check if all characters are hex
	for _, ch := range s {
		if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F') || (ch >= 'a' && ch <= 'f')) {
			return false
		}
	}

	return true
}

// writeLog writes a tag ID to the log file with the specified format
// Format: "unix_timestamp::YYYY-MM-DD HH:MM:SS::tag_id"
func (r *RFIDLogger) writeLog(tagID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.logFile == nil {
		return
	}

	// Get current unix timestamp
	unixTimestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Get current date and time in system format
	systemTime := time.Now().Format("2006-01-02 15:04:05")

	// Format: "unix_timestamp::YYYY-MM-DD HH:MM:SS::tag_id"
	logLine := fmt.Sprintf("%s::%s::%s\n", unixTimestamp, systemTime, tagID)

	_, err := r.logFile.WriteString(logLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
	}

	// Sync to ensure data is written
	r.logFile.Sync()
}

// IsConnected returns the connection status
func (r *RFIDLogger) IsConnected() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isConnected
}

// SetPort sets the serial port (called after successful connection)
func (r *RFIDLogger) SetPort(port io.ReadCloser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.port = port
}

// GetLogFile returns the current log file path
func (r *RFIDLogger) GetLogFile() string {
	return "rfid_tags.log"
}
