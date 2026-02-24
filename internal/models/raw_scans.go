package models

import (
	"time"
)

// RawScan represents a raw RFID scan record
type RawScan struct {
	ID                  string    `db:"id" json:"id"`
	Timestamp           time.Time `db:"timestamp" json:"timestamp"`
	TagValue            string    `db:"tag_value" json:"tag_value"`
	ReaderType          string    `db:"reader_type" json:"reader_type"`
	COMPort             *string   `db:"com_port" json:"com_port,omitempty"`
	SignalStrength      *int      `db:"signal_strength" json:"signal_strength,omitempty"`
	IsProcessed         bool      `db:"is_processed" json:"is_processed"`
	LinkedRacerModelID  *string   `db:"linked_racer_model_id" json:"linked_racer_model_id,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}