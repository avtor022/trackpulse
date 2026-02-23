package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

// RawScan represents a raw RFID scan
type RawScan struct {
	ID                 string    `db:"id" json:"id"`
	Timestamp          time.Time `db:"timestamp" json:"timestamp"`
	TagValue           string    `db:"tag_value" json:"tag_value"`
	ReaderType         string    `db:"reader_type" json:"reader_type"`
	COMPort            *string   `db:"com_port" json:"com_port,omitempty"`
	SignalStrength     *int      `db:"signal_strength" json:"signal_strength,omitempty"`
	IsProcessed        bool      `db:"is_processed" json:"is_processed"`
	LinkedRacerModelID *string   `db:"linked_racer_model_id" json:"linked_racer_model_id,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name for raw_scans
func (RawScan) TableName() string {
	return "raw_scans"
}

// Create inserts a new raw scan record into the database
func (rs *RawScan) Create(db *sql.DB) error {
	if rs.ID == "" {
		rs.ID = uuid.New().String()
	}
	
	now := time.Now().UTC()
	rs.CreatedAt = now
	if rs.Timestamp.IsZero() {
		rs.Timestamp = now
	}

	query := `INSERT INTO raw_scans (id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rs.ID, rs.Timestamp, rs.TagValue, rs.ReaderType, rs.COMPort, rs.SignalStrength, rs.IsProcessed, rs.LinkedRacerModelID, rs.CreatedAt)
	return err
}

// GetByID retrieves a raw scan record by ID
func (rs *RawScan) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at 
	          FROM raw_scans WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&rs.ID, &rs.Timestamp, &rs.TagValue, &rs.ReaderType, &rs.COMPort, &rs.SignalStrength, &rs.IsProcessed, &rs.LinkedRacerModelID, &rs.CreatedAt,
	)
	return err
}

// Update updates an existing raw scan record in the database
func (rs *RawScan) Update(db *sql.DB) error {
	query := `UPDATE raw_scans SET timestamp = ?, tag_value = ?, reader_type = ?, com_port = ?, 
	                 signal_strength = ?, is_processed = ?, linked_racer_model_id = ? WHERE id = ?`
	
	_, err := db.Exec(query, rs.Timestamp, rs.TagValue, rs.ReaderType, rs.COMPort, rs.SignalStrength, rs.IsProcessed, rs.LinkedRacerModelID, rs.ID)
	return err
}

// Delete removes a raw scan record from the database
func (rs *RawScan) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM raw_scans WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll retrieves all raw scan records from the database
func (rs *RawScan) GetAll(db *sql.DB) ([]RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at 
	          FROM raw_scans ORDER BY timestamp DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rawScans []RawScan
	for rows.Next() {
		var rawScan RawScan
		err := rows.Scan(
			&rawScan.ID, &rawScan.Timestamp, &rawScan.TagValue, &rawScan.ReaderType, &rawScan.COMPort, 
			&rawScan.SignalStrength, &rawScan.IsProcessed, &rawScan.LinkedRacerModelID, &rawScan.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rawScans = append(rawScans, rawScan)
	}

	return rawScans, nil
}

// Validate checks if the raw scan data is valid
func (rs *RawScan) Validate() error {
	if rs.TagValue == "" {
		return fmt.Errorf("tag value is required")
	}
	if rs.ReaderType == "" {
		return fmt.Errorf("reader type is required")
	}
	return nil
}