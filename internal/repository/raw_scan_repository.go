package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"trackpulse/internal/models"
)

// RawScanRepository handles data access for raw RFID scans
type RawScanRepository struct {
	db *sql.DB
}

// NewRawScanRepository creates a new raw scan repository
func NewRawScanRepository(db *sql.DB) *RawScanRepository {
	return &RawScanRepository{db: db}
}

// Create inserts a single raw scan record
func (r *RawScanRepository) Create(scan *models.RawScan) error {
	now := time.Now().Format(time.RFC3339)
	timestamp := scan.Timestamp.Format(time.RFC3339)

	var signalStrength sql.NullInt64
	if scan.SignalStrength != nil {
		signalStrength = sql.NullInt64{Int64: int64(*scan.SignalStrength), Valid: true}
	}

	var linkedCompetitorModelID sql.NullString
	if scan.LinkedCompetitorModelID != nil {
		linkedCompetitorModelID = sql.NullString{String: *scan.LinkedCompetitorModelID, Valid: true}
	}

	result, err := r.db.Exec(`
		INSERT INTO raw_scans (id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_competitor_model_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		scan.ID,
		timestamp,
		scan.TagValue,
		scan.ReaderType,
		scan.COMPort,
		signalStrength,
		scan.IsProcessed,
		linkedCompetitorModelID,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create raw scan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating raw scan")
	}

	return nil
}

// CreateBulk inserts multiple raw scan records in a single transaction
// This is optimized for high-frequency RFID readings (10-15 scans/sec)
func (r *RawScanRepository) CreateBulk(scans []*models.RawScan) error {
	if len(scans) == 0 {
		return nil
	}

	// Start transaction for atomic bulk insert
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Prepare statement for efficient repeated execution
	stmt, err := tx.Prepare(`
		INSERT INTO raw_scans (id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_competitor_model_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now().Format(time.RFC3339)

	for _, scan := range scans {
		if scan.ID == "" {
			scan.ID = uuid.New().String()
		}
		if scan.Timestamp.IsZero() {
			scan.Timestamp = time.Now()
		}

		timestamp := scan.Timestamp.Format(time.RFC3339)

		var signalStrength sql.NullInt64
		if scan.SignalStrength != nil {
			signalStrength = sql.NullInt64{Int64: int64(*scan.SignalStrength), Valid: true}
		}

		var linkedCompetitorModelID sql.NullString
		if scan.LinkedCompetitorModelID != nil {
			linkedCompetitorModelID = sql.NullString{String: *scan.LinkedCompetitorModelID, Valid: true}
		}

		_, err := stmt.Exec(
			scan.ID,
			timestamp,
			scan.TagValue,
			scan.ReaderType,
			scan.COMPort,
			signalStrength,
			scan.IsProcessed,
			linkedCompetitorModelID,
			now,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute insert: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUnprocessed returns raw scans that haven't been processed yet
func (r *RawScanRepository) GetUnprocessed(limit int) ([]models.RawScan, error) {
	rows, err := r.db.Query(`
		SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_competitor_model_id, created_at
		FROM raw_scans
		WHERE is_processed = 0
		ORDER BY timestamp ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query unprocessed raw scans: %w", err)
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var s models.RawScan
		var timestampStr, createdAtStr string
		var signalStrength sql.NullInt64
		var linkedCompetitorModelID sql.NullString

		err := rows.Scan(
			&s.ID,
			&timestampStr,
			&s.TagValue,
			&s.ReaderType,
			&s.COMPort,
			&signalStrength,
			&s.IsProcessed,
			&linkedCompetitorModelID,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan raw scan: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			s.Timestamp = t
		}
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			s.CreatedAt = t
		}
		if signalStrength.Valid {
			ss := int(signalStrength.Int64)
			s.SignalStrength = &ss
		}
		if linkedCompetitorModelID.Valid {
			s.LinkedCompetitorModelID = &linkedCompetitorModelID.String
		}

		scans = append(scans, s)
	}

	return scans, rows.Err()
}

// MarkAsProcessed marks a raw scan as processed and optionally links it to a competitor model
func (r *RawScanRepository) MarkAsProcessed(id string, competitorModelID *string) error {
	var err error
	if competitorModelID != nil {
		_, err = r.db.Exec(`
			UPDATE raw_scans
			SET is_processed = 1, linked_competitor_model_id = ?
			WHERE id = ?
		`, *competitorModelID, id)
	} else {
		_, err = r.db.Exec(`
			UPDATE raw_scans
			SET is_processed = 1
			WHERE id = ?
		`, id)
	}

	if err != nil {
		return fmt.Errorf("failed to mark raw scan as processed: %w", err)
	}

	return nil
}

// GetByTagValue returns raw scans filtered by tag value
func (r *RawScanRepository) GetByTagValue(tagValue string, limit int) ([]models.RawScan, error) {
	rows, err := r.db.Query(`
		SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_competitor_model_id, created_at
		FROM raw_scans
		WHERE tag_value = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`, tagValue, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query raw scans by tag value: %w", err)
	}
	defer rows.Close()

	var scans []models.RawScan
	for rows.Next() {
		var s models.RawScan
		var timestampStr, createdAtStr string
		var signalStrength sql.NullInt64
		var linkedCompetitorModelID sql.NullString

		err := rows.Scan(
			&s.ID,
			&timestampStr,
			&s.TagValue,
			&s.ReaderType,
			&s.COMPort,
			&signalStrength,
			&s.IsProcessed,
			&linkedCompetitorModelID,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan raw scan: %w", err)
		}

		if t, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			s.Timestamp = t
		}
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			s.CreatedAt = t
		}
		if signalStrength.Valid {
			ss := int(signalStrength.Int64)
			s.SignalStrength = &ss
		}
		if linkedCompetitorModelID.Valid {
			s.LinkedCompetitorModelID = &linkedCompetitorModelID.String
		}

		scans = append(scans, s)
	}

	return scans, rows.Err()
}

// Count returns total number of raw scans
func (r *RawScanRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM raw_scans`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count raw scans: %w", err)
	}
	return count, nil
}

// DeleteOld removes raw scans older than the specified duration
func (r *RawScanRepository) DeleteOld(olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	cutoffStr := cutoff.Format(time.RFC3339)

	result, err := r.db.Exec(`
		DELETE FROM raw_scans
		WHERE timestamp < ? AND is_processed = 1
	`, cutoffStr)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old raw scans: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
