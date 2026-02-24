package services

import (
	"fmt"
	"time"
	"trackpulse/internal/config"
	"github.com/tarm/serial"
)

type HardwareService interface {
	Initialize() error
	ReadTag() (string, error)
	Close() error
	GetDeviceInfo() (map[string]interface{}, error)
	TestConnection() error
	GetSignalStrength() (int, error)
	SendCommand(command string) error
	ResetDevice() error
}

type hardwareService struct {
	config     *config.Config
	port       *serial.Port
	deviceType string
	isOpen     bool
}

func NewHardwareService(cfg *config.Config) HardwareService {
	return &hardwareService{
		config: cfg,
		isOpen: false,
	}
}

func (h *hardwareService) Initialize() error {
	// Get hardware settings from configuration
	settings := h.getHardwareSettings()

	// Configure serial port
	c := &serial.Config{
		Name: settings["com_port"],
		Baud: 9600,
		ReadTimeout: time.Second * 2,
	}

	p, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %v", err)
	}

	h.port = p
	h.isOpen = true
	h.deviceType = settings["reader_type"]

	// Initialize the RFID reader based on type
	switch h.deviceType {
	case "EM4095":
		err = h.initializeEM4095()
	case "RC522":
		err = h.initializeRC522()
	default:
		err = fmt.Errorf("unsupported RFID reader type: %s", h.deviceType)
	}

	if err != nil {
		h.Close()
		return fmt.Errorf("failed to initialize RFID reader: %v", err)
	}

	return nil
}

func (h *hardwareService) ReadTag() (string, error) {
	if !h.isOpen {
		return "", fmt.Errorf("hardware not initialized")
	}

	// Send command to read RFID tag based on device type
	var response []byte
	var err error

	switch h.deviceType {
	case "EM4095":
		response, err = h.readTagEM4095()
	case "RC522":
		response, err = h.readTagRC522()
	default:
		return "", fmt.Errorf("unsupported RFID reader type: %s", h.deviceType)
	}

	if err != nil {
		return "", err
	}

	// Parse the response to extract tag value
	tagValue := h.parseTagResponse(response)

	return tagValue, nil
}

func (h *hardwareService) Close() error {
	if h.port != nil && h.isOpen {
		err := h.port.Close()
		h.isOpen = false
		return err
	}
	return nil
}

func (h *hardwareService) GetDeviceInfo() (map[string]interface{}, error) {
	if !h.isOpen {
		return nil, fmt.Errorf("hardware not initialized")
	}

	info := make(map[string]interface{})

	// Get device information based on type
	switch h.deviceType {
	case "EM4095":
		version, err := h.getEM4095Version()
		if err != nil {
			return nil, err
		}
		info["version"] = version
	case "RC522":
		// Get RC522-specific information
		info["version"] = "RC522"
	}

	info["device_type"] = h.deviceType
	info["is_connected"] = h.isOpen

	return info, nil
}

func (h *hardwareService) TestConnection() error {
	if !h.isOpen {
		return fmt.Errorf("hardware not initialized")
	}

	// Send a simple ping/test command to the device
	switch h.deviceType {
	case "EM4095":
		return h.testEM4095Connection()
	case "RC522":
		return h.testRC522Connection()
	default:
		return fmt.Errorf("unsupported RFID reader type: %s", h.deviceType)
	}
}

func (h *hardwareService) GetSignalStrength() (int, error) {
	if !h.isOpen {
		return 0, fmt.Errorf("hardware not initialized")
	}

	// Get signal strength based on device type
	switch h.deviceType {
	case "EM4095":
		return h.getEM4095SignalStrength()
	case "RC522":
		return h.getRC522SignalStrength()
	default:
		return 0, fmt.Errorf("unsupported RFID reader type: %s", h.deviceType)
	}
}

func (h *hardwareService) SendCommand(command string) error {
	if !h.isOpen {
		return fmt.Errorf("hardware not initialized")
	}

	// Write command to the serial port
	_, err := h.port.Write([]byte(command))
	return err
}

func (h *hardwareService) ResetDevice() error {
	if h.isOpen {
		h.Close()
	}

	// Reinitialize the device
	return h.Initialize()
}

// Helper methods for EM4095
func (h *hardwareService) initializeEM4095() error {
	// Send initialization sequence for EM4095
	initCmd := []byte{0x02, 0x00, 0x00, 0x03} // Example initialization command
	_, err := h.port.Write(initCmd)
	if err != nil {
		return err
	}

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (h *hardwareService) readTagEM4095() ([]byte, error) {
	// Send read command for EM4095
	readCmd := []byte{0x02, 0x01, 0x00, 0x03} // Example read command
	_, err := h.port.Write(readCmd)
	if err != nil {
		return nil, err
	}

	// Read response (typically takes some time)
	buffer := make([]byte, 128)
	n, err := h.port.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (h *hardwareService) getEM4095Version() (string, error) {
	// Send version command for EM4095
	versionCmd := []byte{0x02, 0x03, 0x00, 0x03} // Example version command
	_, err := h.port.Write(versionCmd)
	if err != nil {
		return "", err
	}

	// Read response
	buffer := make([]byte, 128)
	n, err := h.port.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}

func (h *hardwareService) testEM4095Connection() error {
	// Send test command and verify response
	version, err := h.getEM4095Version()
	if err != nil {
		return err
	}

	if version == "" {
		return fmt.Errorf("empty version response from EM4095")
	}

	return nil
}

func (h *hardwareService) getEM4095SignalStrength() (int, error) {
	// EM4095 may not support signal strength directly
	// This is a placeholder implementation
	return 0, nil
}

// Helper methods for RC522
func (h *hardwareService) initializeRC522() error {
	// Send initialization sequence for RC522
	// RC522 initialization typically involves SPI communication
	// This is a simplified representation
	initCmd := []byte{0x01, 0x00, 0xFF, 0x00} // Example initialization command
	_, err := h.port.Write(initCmd)
	if err != nil {
		return err
	}

	// Wait for response
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (h *hardwareService) readTagRC522() ([]byte, error) {
	// Send read command for RC522
	readCmd := []byte{0x01, 0x01, 0xFF, 0x00} // Example read command
	_, err := h.port.Write(readCmd)
	if err != nil {
		return nil, err
	}

	// Read response
	buffer := make([]byte, 128)
	n, err := h.port.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (h *hardwareService) testRC522Connection() error {
	// Send a simple command and verify response
	cmd := []byte{0x01, 0x03, 0xFF, 0x00} // Example test command
	_, err := h.port.Write(cmd)
	if err != nil {
		return err
	}

	// Read response
	buffer := make([]byte, 128)
	n, err := h.port.Read(buffer)
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("no response from RC522")
	}

	return nil
}

func (h *hardwareService) getRC522SignalStrength() (int, error) {
	// RC522 may not support signal strength directly
	// This is a placeholder implementation
	return 0, nil
}

// Helper method to parse tag response
func (h *hardwareService) parseTagResponse(response []byte) string {
	// The exact format depends on the RFID reader
	// This is a simplified implementation
	if len(response) < 4 {
		return ""
	}

	// Assuming tag ID starts at position 2 and is 4 bytes long
	// This varies depending on the RFID reader protocol
	tagBytes := response[2:6]
	
	// Convert to hex string
	tagStr := ""
	for _, b := range tagBytes {
		tagStr += fmt.Sprintf("%02X", b)
	}

	return tagStr
}

// Helper method to get hardware settings from configuration
func (h *hardwareService) getHardwareSettings() map[string]string {
	settings := make(map[string]string)
	
	// In a real implementation, these would come from the database
	// For now, we'll use placeholders
	settings["com_port"] = h.config.DBPath // This would actually come from system settings
	settings["reader_type"] = "EM4095"
	
	return settings
}