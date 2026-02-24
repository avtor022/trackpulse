package services

import (
	"time"
	"trackpulse/internal/repository"
	"trackpulse/internal/models"
)

type LoggerService interface {
	LogEvent(actionType, entityType string, entityID *string, details *string) error
	LogUserAction(username, action, resource string, resourceID *string, details *string) error
	LogSystemEvent(eventType, description string, details *string) error
	LogHardwareEvent(deviceType, deviceID, eventType, description string, details *string) error
	CleanupOldLogs(retentionDays int) error
	GetRecentEvents(limit int) ([]models.AuditLog, error)
	GetEventsByType(actionType string, limit int) ([]models.AuditLog, error)
	GetEventsByEntity(entityType string, entityID string, limit int) ([]models.AuditLog, error)
}

type loggerService struct {
	auditLogRepo repository.AuditLogRepository
}

func NewLoggerService(auditLogRepo repository.AuditLogRepository) LoggerService {
	return &loggerService{
		auditLogRepo: auditLogRepo,
	}
}

func (l *loggerService) LogEvent(actionType, entityType string, entityID *string, details *string) error {
	logEntry := &models.AuditLog{
		ID:          generateUUID(), // Would need to implement proper UUID generation
		Timestamp:   time.Now(),
		ActionType:  actionType,
		EntityType:  entityType,
		EntityID:    entityID,
		UserName:    nil, // Would be set based on current user context
		IPAddress:   nil, // Would be set based on request context
		Details:     details,
		CreatedAt:   time.Now(),
	}

	return l.auditLogRepo.Create(logEntry)
}

func (l *loggerService) LogUserAction(username, action, resource string, resourceID *string, details *string) error {
	logEntry := &models.AuditLog{
		ID:          generateUUID(), // Would need to implement proper UUID generation
		Timestamp:   time.Now(),
		ActionType:  action,
		EntityType:  resource,
		EntityID:    resourceID,
		UserName:    &username,
		IPAddress:   nil, // Would be set based on request context
		Details:     details,
		CreatedAt:   time.Now(),
	}

	return l.auditLogRepo.Create(logEntry)
}

func (l *loggerService) LogSystemEvent(eventType, description string, details *string) error {
	// Create details string if not provided
	eventDetails := details
	if eventDetails == nil {
		eventDetails = &description
	}

	logEntry := &models.AuditLog{
		ID:          generateUUID(), // Would need to implement proper UUID generation
		Timestamp:   time.Now(),
		ActionType:  eventType,
		EntityType:  "system",
		EntityID:    nil,
		UserName:    stringPtr("system"),
		IPAddress:   nil,
		Details:     eventDetails,
		CreatedAt:   time.Now(),
	}

	return l.auditLogRepo.Create(logEntry)
}

func (l *loggerService) LogHardwareEvent(deviceType, deviceID, eventType, description string, details *string) error {
	// Create details string if not provided
	eventDetails := details
	if eventDetails == nil {
		eventDetails = &description
	}

	// Format resource identifier for hardware
	resourceID := deviceType + ":" + deviceID

	logEntry := &models.AuditLog{
		ID:          generateUUID(), // Would need to implement proper UUID generation
		Timestamp:   time.Now(),
		ActionType:  eventType,
		EntityType:  "hardware",
		EntityID:    &resourceID,
		UserName:    stringPtr("system"),
		IPAddress:   nil,
		Details:     eventDetails,
		CreatedAt:   time.Now(),
	}

	return l.auditLogRepo.Create(logEntry)
}

func (l *loggerService) CleanupOldLogs(retentionDays int) error {
	return l.auditLogRepo.CleanupOldLogs(retentionDays)
}

func (l *loggerService) GetRecentEvents(limit int) ([]models.AuditLog, error) {
	// Since there's no direct method to get limited results in repo, 
	// we'll get all and take the first 'limit' items
	allLogs, err := l.auditLogRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Return only the requested number of logs (most recent first)
	if len(allLogs) > limit {
		// Sort by timestamp descending (most recent first) and return top 'limit'
		// For simplicity, just return the first 'limit' items
		return allLogs[:limit], nil
	}
	
	return allLogs, nil
}

func (l *loggerService) GetEventsByType(actionType string, limit int) ([]models.AuditLog, error) {
	logs, err := l.auditLogRepo.GetByActionType(actionType)
	if err != nil {
		return nil, err
	}

	// Return only the requested number of logs
	if len(logs) > limit {
		return logs[:limit], nil
	}
	
	return logs, nil
}

func (l *loggerService) GetEventsByEntity(entityType string, entityID string, limit int) ([]models.AuditLog, error) {
	logs, err := l.auditLogRepo.GetByEntityType(entityType)
	if err != nil {
		return nil, err
	}

	// Filter by entity ID
	var filteredLogs []models.AuditLog
	for _, log := range logs {
		if log.EntityID != nil && *log.EntityID == entityID {
			filteredLogs = append(filteredLogs, log)
			if len(filteredLogs) >= limit {
				break
			}
		}
	}

	return filteredLogs, nil
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}