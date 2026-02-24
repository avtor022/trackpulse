package models

// Common interfaces and utilities for models

// Model interface that all models should implement
type Model interface {
	TableName() string
	Validate() error
}

// Common validation errors
var (
	ErrValidationFailed = "validation failed"
)