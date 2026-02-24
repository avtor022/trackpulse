package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type AuditLogRepository interface {
	GetAll() ([]models.AuditLog, error)
	GetByID(id string) (*models.AuditLog, error)
	Create(log *models.AuditLog) error
	Delete(id string) error
	GetByActionType(actionType string) ([]models.AuditLog, error)
	GetByEntityType(entityType string) ([]models.AuditLog, error)
	GetByTimeRange(startTime string, endTime string) ([]models.AuditLog, error)
	CleanupOldLogs(retentionDays int) error
}

type auditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (a *auditLogRepo) GetAll() ([]models.AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, 
						ip_address, details, created_at 
					FROM audit_log ORDER BY timestamp DESC`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var entityID, userName, ipAddress, details *string

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.ActionType,
			&log.EntityType,
			&entityID,
			&userName,
			&ipAddress,
			&details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		log.EntityID = entityID
		log.UserName = userName
		log.IPAddress = ipAddress
		log.Details = details

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *auditLogRepo) GetByID(id string) (*models.AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, 
						ip_address, details, created_at 
					FROM audit_log WHERE id = ?`

	var log models.AuditLog
	var entityID, userName, ipAddress, details *string

	err := a.db.QueryRow(query, id).Scan(
		&log.ID,
		&log.Timestamp,
		&log.ActionType,
		&log.EntityType,
		&entityID,
		&userName,
		&ipAddress,
		&details,
		&log.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}

	// Handle nullable fields
	log.EntityID = entityID
	log.UserName = userName
	log.IPAddress = ipAddress
	log.Details = details

	return &log, nil
}

func (a *auditLogRepo) Create(log *models.AuditLog) error {
	query := `INSERT INTO audit_log (id, timestamp, action_type, entity_type, entity_id, 
																	user_name, ip_address, details, created_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var entityID, userName, ipAddress, details *string
	
	if log.EntityID != nil {
		entityID = log.EntityID
	}
	if log.UserName != nil {
		userName = log.UserName
	}
	if log.IPAddress != nil {
		ipAddress = log.IPAddress
	}
	if log.Details != nil {
		details = log.Details
	}

	_, err := a.db.Exec(query,
		log.ID,
		log.Timestamp,
		log.ActionType,
		log.EntityType,
		entityID,
		userName,
		ipAddress,
		details,
		log.CreatedAt,
	)
	return err
}

func (a *auditLogRepo) Delete(id string) error {
	query := `DELETE FROM audit_log WHERE id = ?`

	result, err := a.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("audit log not found")
	}

	return nil
}

func (a *auditLogRepo) GetByActionType(actionType string) ([]models.AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, 
						ip_address, details, created_at 
					FROM audit_log WHERE action_type = ? ORDER BY timestamp DESC`

	rows, err := a.db.Query(query, actionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var entityID, userName, ipAddress, details *string

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.ActionType,
			&log.EntityType,
			&entityID,
			&userName,
			&ipAddress,
			&details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		log.EntityID = entityID
		log.UserName = userName
		log.IPAddress = ipAddress
		log.Details = details

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *auditLogRepo) GetByEntityType(entityType string) ([]models.AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, 
						ip_address, details, created_at 
					FROM audit_log WHERE entity_type = ? ORDER BY timestamp DESC`

	rows, err := a.db.Query(query, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var entityID, userName, ipAddress, details *string

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.ActionType,
			&log.EntityType,
			&entityID,
			&userName,
			&ipAddress,
			&details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		log.EntityID = entityID
		log.UserName = userName
		log.IPAddress = ipAddress
		log.Details = details

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *auditLogRepo) GetByTimeRange(startTime string, endTime string) ([]models.AuditLog, error) {
	query := `SELECT id, timestamp, action_type, entity_type, entity_id, user_name, 
						ip_address, details, created_at 
					FROM audit_log WHERE timestamp BETWEEN ? AND ? ORDER BY timestamp DESC`

	rows, err := a.db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var entityID, userName, ipAddress, details *string

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&log.ActionType,
			&log.EntityType,
			&entityID,
			&userName,
			&ipAddress,
			&details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		log.EntityID = entityID
		log.UserName = userName
		log.IPAddress = ipAddress
		log.Details = details

		logs = append(logs, log)
	}

	return logs, nil
}

func (a *auditLogRepo) CleanupOldLogs(retentionDays int) error {
	query := `DELETE FROM audit_log WHERE timestamp < datetime('now', '-' || ? || ' days')`

	_, err := a.db.Exec(query, retentionDays)
	return err
}