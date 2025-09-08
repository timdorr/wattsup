package config

import "errors"

var (
	ErrMissingDatabase = errors.New("database connection string is required")
	ErrNoDevices       = errors.New("at least one device must be configured")
	ErrNoRegisters     = errors.New("at least one register must be configured")
)
