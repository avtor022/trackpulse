package serial

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"
	"trackpulse/internal/race"

	"github.com/tarm/serial"
)

type Handler struct {
	raceController *race.Controller
	portName       string
	baudRate       int
	readerType     string
	debounceMS     int
	scanner        *bufio.Scanner
}

func NewHandler(raceController *race.Controller) *Handler {
	return &Handler{
		raceController: raceController,
		portName:       "COM3", // Default, should be configurable
		baudRate:       115200, // Default baud rate
		readerType:     "EM4095", // Default reader type
		debounceMS:     2000, // Default debounce in milliseconds
	}
}

func (h *Handler) StartListening() {
	for {
		// Attempt to open the serial port
		c := &serial.Config{
			Name:        h.portName,
			Baud:        h.baudRate,
			ReadTimeout: time.Second * 5,
		}
		
		s, err := serial.OpenPort(c)
		if err != nil {
			log.Printf("Error opening serial port %s: %v", h.portName, err)
			
			// Wait before attempting to reconnect
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Connected to serial port %s", h.portName)
		
		// Create a scanner to read from the serial port
		h.scanner = bufio.NewScanner(s)
		
		// Read from the serial port
		for h.scanner.Scan() {
			line := h.scanner.Text()
			
			// Process the received line
			h.processLine(line)
		}
		
		// Handle scanner errors
		if err := h.scanner.Err(); err != nil {
			log.Printf("Error reading from serial port: %v", err)
		}
		
		// Close the port
		s.Close()
		
		// Wait before attempting to reconnect
		time.Sleep(5 * time.Second)
	}
}

func (h *Handler) processLine(line string) {
	// Trim whitespace from the line
	line = strings.TrimSpace(line)
	
	// Check if the line starts with "TAG:"
	if strings.HasPrefix(line, "TAG:") {
		// Extract the tag value
		tagValue := strings.TrimPrefix(line, "TAG:")
		
		// Process the RFID tag
		h.handleRFIDTag(tagValue)
	} else {
		// Log unexpected input
		log.Printf("Unexpected input from serial: %s", line)
	}
}

func (h *Handler) handleRFIDTag(tagValue string) {
	// TODO: Implement logic to find the active race and process the lap
	// For now, we'll just log the received tag
	
	// In a real implementation, you would:
	// 1. Query the database to find the active race
	// 2. Find the racer associated with this tagValue
	// 3. Call the race controller to process the lap
	
	log.Printf("RFID Tag detected: %s (Type: %s)", tagValue, h.readerType)
	
	// Example of how you might call the race controller:
	// Assuming we know the active race ID (this would need to be determined)
	// h.raceController.ProcessLap(activeRaceID, tagValue)
}

// SetPort sets the serial port name
func (h *Handler) SetPort(portName string) {
	h.portName = portName
}

// SetBaudRate sets the baud rate for the serial connection
func (h *Handler) SetBaudRate(baudRate int) {
	h.baudRate = baudRate
}

// SetReaderType sets the RFID reader type (EM4095 or RC522)
func (h *Handler) SetReaderType(readerType string) {
	h.readerType = readerType
}

// SetDebounce sets the debounce time in milliseconds
func (h *Handler) SetDebounce(debounceMS int) {
	h.debounceMS = debounceMS
}