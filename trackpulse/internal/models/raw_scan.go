package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

// RawScan represents a record in the raw_scans table
type RawScan struct {
	ID                   string    `db:"id" json:"id"`
	Timestamp            string    `db:"timestamp" json:"timestamp"`
	TagValue             string    `db:"tag_value" json:"tag_value"`
	ReaderType           string    `db:"reader_type" json:"reader_type"`
	COMPort              *string   `db:"com_port" json:"com_port,omitempty"`
	SignalStrength       *int      `db:"signal_strength" json:"signal_strength,omitempty"`
	IsProcessed          bool      `db:"is_processed" json:"is_processed"`
	LinkedRacerModelID   *string   `db:"linked_racer_model_id" json:"linked_racer_model_id,omitempty"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the table name
func (RawScan) TableName() string {
	return "raw_scans"
}

// Create creates a new record
func (rs *RawScan) Create(db *sql.DB) error {
	rs.ID = uuid.New().String()
	rs.IsProcessed = false
	rs.CreatedAt = time.Now()

	query := `INSERT INTO raw_scans (id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := db.Exec(query, rs.ID, rs.Timestamp, rs.TagValue, rs.ReaderType, rs.COMPort, rs.SignalStrength, rs.IsProcessed, rs.LinkedRacerModelID, rs.CreatedAt)
	return err
}

// GetByID gets a record by ID
func (rs *RawScan) GetByID(db *sql.DB, id string) error {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at 
	          FROM raw_scans WHERE id = ?`
	
	return db.QueryRow(query, id).Scan(&rs.ID, &rs.Timestamp, &rs.TagValue, &rs.ReaderType, &rs.COMPort, &rs.SignalStrength, &rs.IsProcessed, &rs.LinkedRacerModelID, &rs.CreatedAt)
}

// Update updates a record
func (rs *RawScan) Update(db *sql.DB) error {
	query := `UPDATE raw_scans SET timestamp = ?, tag_value = ?, reader_type = ?, com_port = ?, signal_strength = ?, is_processed = ?, linked_racer_model_id = ? 
	          WHERE id = ?`
	
	_, err := db.Exec(query, rs.Timestamp, rs.TagValue, rs.ReaderType, rs.COMPort, rs.SignalStrength, rs.IsProcessed, rs.LinkedRacerModelID, rs.ID)
	return err
}

// Delete deletes a record
func (rs *RawScan) Delete(db *sql.DB, id string) error {
	query := `DELETE FROM raw_scans WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// GetAll gets all records
func (rs *RawScan) GetAll(db *sql.DB) ([]RawScan, error) {
	query := `SELECT id, timestamp, tag_value, reader_type, com_port, signal_strength, is_processed, linked_racer_model_id, created_at 
	          FROM raw_scans ORDER BY timestamp DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scans []RawScan
	for rows.Next() {
		var scan RawScan
		err := rows.Scan(&scan.ID, &scan.Timestamp, &scan.TagValue, &scan.ReaderType, &scan.COMPort, &scan.SignalStrength, &scan.IsProcessed, &scan.LinkedRacerModelID, &scan.CreatedAt)
		if err != nil {
			return nil, err
		}
		scans = append(scans, scan)
	}

	return scans, nil
}