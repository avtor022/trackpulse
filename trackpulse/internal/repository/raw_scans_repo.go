package repository

import (
	"database/sql"
	"errors"
	"trackpulse/internal/models"
)

type RawScansRepository interface {
	GetAll() ([]models.RawScan, error)
	GetByID(id string) (*models.RawScan, error)
	Create(scan *models.RawScan) error
	Update(scan *models.RawScan) error
	Delete(id string) error
	GetUnprocessedScans() ([]models.RawScan, error)
	MarkAsProcessed(scanID string) error
	GetByTagValue(tagValue string) ([]models.RawScan, error)
	GetByTimeRange(startTime string, endTime string) ([]models.RawScan, error)
	CleanupOldLogs(retentionDays int) error
}

type rawScansRepo struct {
	db *sql.DB
}

func NewRawScansRepository(db *sql.DB) RawScansRepository {
	return &rawScansRepo{db: db}
}

func (r *rawScansRepo) GetAll() ([]models.RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, 
						is_processed, linked_racer_model_id, created_at 
					FROM raw_scans ORDER BY timestamp DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var scan models.RawScan
		var comPort, signalStrength, linkedRacerModelID *string

		err := rows.Scan(
			&scan.ID,
			&scan.Timestamp,
			&scan.TagValue,
			&scan.ReaderType,
			&comPort,
			&signalStrength,
			&scan.IsProcessed,
			&linkedRacerModelID,
			&scan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if comPort != nil {
			scan.COMPort = comPort
		}
		if signalStrength != nil {
			// Convert string to int pointer
		}
		if linkedRacerModelID != nil {
			scan.LinkedRacerModelID = linkedRacerModelID
		}

		scans = append(scans, scan)
	}

	return scans, nil
}

func (r *rawScansRepo) GetByID(id string) (*models.RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, 
						is_processed, linked_racer_model_id, created_at 
					FROM raw_scans WHERE id = ?`

	var scan models.RawScan
	var comPort, signalStrength, linkedRacerModelID *string

	err := r.db.QueryRow(query, id).Scan(
		&scan.ID,
		&scan.Timestamp,
		&scan.TagValue,
		&scan.ReaderType,
		&comPort,
		&signalStrength,
		&scan.IsProcessed,
		&linkedRacerModelID,
		&scan.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("raw scan not found")
		}
		return nil, err
	}

	// Handle nullable fields
	if comPort != nil {
		scan.COMPort = comPort
	}
	if signalStrength != nil {
		// Convert string to int pointer
	}
	if linkedRacerModelID != nil {
		scan.LinkedRacerModelID = linkedRacerModelID
	}

	return &scan, nil
}

func (r *rawScansRepo) Create(scan *models.RawScan) error {
	query := `INSERT INTO raw_scans (id, timestamp, tag_value, reader_type, com_port, 
																	signal_strength, is_processed, linked_racer_model_id, created_at)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Extract values for nullable fields
	var comPort, signalStrength, linkedRacerModelID *string
	
	if scan.COMPort != nil {
		comPort = scan.COMPort
	}
	if scan.SignalStrength != nil {
		// Convert int pointer to string
	}
	if scan.LinkedRacerModelID != nil {
		linkedRacerModelID = scan.LinkedRacerModelID
	}

	_, err := r.db.Exec(query,
		scan.ID,
		scan.Timestamp,
		scan.TagValue,
		scan.ReaderType,
		comPort,
		signalStrength,
		scan.IsProcessed,
		linkedRacerModelID,
		scan.CreatedAt,
	)
	return err
}

func (r *rawScansRepo) Update(scan *models.RawScan) error {
	query := `UPDATE raw_scans SET timestamp = ?, tag_value = ?, reader_type = ?, 
						com_port = ?, signal_strength = ?, is_processed = ?, 
						linked_racer_model_id = ?, created_at = ?
					WHERE id = ?`

	// Extract values for nullable fields
	var comPort, signalStrength, linkedRacerModelID *string
	
	if scan.COMPort != nil {
		comPort = scan.COMPort
	}
	if scan.SignalStrength != nil {
		// Convert int pointer to string
	}
	if scan.LinkedRacerModelID != nil {
		linkedRacerModelID = scan.LinkedRacerModelID
	}

	_, err := r.db.Exec(query,
		scan.Timestamp,
		scan.TagValue,
		scan.ReaderType,
		comPort,
		signalStrength,
		scan.IsProcessed,
		linkedRacerModelID,
		scan.CreatedAt,
		scan.ID,
	)
	return err
}

func (r *rawScansRepo) Delete(id string) error {
	query := `DELETE FROM raw_scans WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("raw scan not found")
	}

	return nil
}

func (r *rawScansRepo) GetUnprocessedScans() ([]models.RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, 
						is_processed, linked_racer_model_id, created_at 
					FROM raw_scans WHERE is_processed = 0 ORDER BY timestamp`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var scan models.RawScan
		var comPort, signalStrength, linkedRacerModelID *string

		err := rows.Scan(
			&scan.ID,
			&scan.Timestamp,
			&scan.TagValue,
			&scan.ReaderType,
			&comPort,
			&signalStrength,
			&scan.IsProcessed,
			&linkedRacerModelID,
			&scan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if comPort != nil {
			scan.COMPort = comPort
		}
		if signalStrength != nil {
			// Convert string to int pointer
		}
		if linkedRacerModelID != nil {
			scan.LinkedRacerModelID = linkedRacerModelID
		}

		scans = append(scans, scan)
	}

	return scans, nil
}

func (r *rawScansRepo) MarkAsProcessed(scanID string) error {
	query := `UPDATE raw_scans SET is_processed = 1, 
						linked_racer_model_id = (SELECT id FROM racer_models WHERE transponder_number = 
																		(SELECT tag_value FROM raw_scans WHERE id = ?))
					WHERE id = ?`

	_, err := r.db.Exec(query, scanID, scanID)
	return err
}

func (r *rawScansRepo) GetByTagValue(tagValue string) ([]models.RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, 
						is_processed, linked_racer_model_id, created_at 
					FROM raw_scans WHERE tag_value = ? ORDER BY timestamp DESC`

	rows, err := r.db.Query(query, tagValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var scan models.RawScan
		var comPort, signalStrength, linkedRacerModelID *string

		err := rows.Scan(
			&scan.ID,
			&scan.Timestamp,
			&scan.TagValue,
			&scan.ReaderType,
			&comPort,
			&signalStrength,
			&scan.IsProcessed,
			&linkedRacerModelID,
			&scan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if comPort != nil {
			scan.COMPort = comPort
		}
		if signalStrength != nil {
			// Convert string to int pointer
		}
		if linkedRacerModelID != nil {
			scan.LinkedRacerModelID = linkedRacerModelID
		}

		scans = append(scans, scan)
	}

	return scans, nil
}

func (r *rawScansRepo) GetByTimeRange(startTime string, endTime string) ([]models.RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, 
						is_processed, linked_racer_model_id, created_at 
					FROM raw_scans WHERE timestamp BETWEEN ? AND ? ORDER BY timestamp`

	rows, err := r.db.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var scan models.RawScan
		var comPort, signalStrength, linkedRacerModelID *string

		err := rows.Scan(
			&scan.ID,
			&scan.Timestamp,
			&scan.TagValue,
			&scan.ReaderType,
			&comPort,
			&signalStrength,
			&scan.IsProcessed,
			&linkedRacerModelID,
			&scan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if comPort != nil {
			scan.COMPort = comPort
		}
		if signalStrength != nil {
			// Convert string to int pointer
		}
		if linkedRacerModelID != nil {
			scan.LinkedRacerModelID = linkedRacerModelID
		}

		scans = append(scans, scan)
	}

	return scans, nil
}

func (r *rawScansRepo) CleanupOldLogs(retentionDays int) error {
	query := `DELETE FROM raw_scans WHERE timestamp < datetime('now', '-' || ? || ' days')`

	_, err := r.db.Exec(query, retentionDays)
	return err
}