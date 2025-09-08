package config

import (
	"strings"
	"testing"
)

func TestParseConfig_ValidConfig(t *testing.T) {
	configJSON := `{
		"database": "postgresql://user:pass@localhost:5432/db",
		"devices": [
			{"name": "Test Device", "file": "/dev/ttyUSB0", "id": 1}
		],
		"registers": [
			{"name": "Test Register", "address": 100}
		]
	}`

	config, err := ParseConfig(strings.NewReader(configJSON))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Database != "postgresql://user:pass@localhost:5432/db" {
		t.Errorf("Expected database to be set, got %s", config.Database)
	}

	if len(config.Devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(config.Devices))
	}

	if config.Devices[0].Name != "Test Device" {
		t.Errorf("Expected device name 'Test Device', got %s", config.Devices[0].Name)
	}

	if len(config.Registers) != 1 {
		t.Errorf("Expected 1 register, got %d", len(config.Registers))
	}

	if config.Registers[0].Address != 100 {
		t.Errorf("Expected register address 100, got %d", config.Registers[0].Address)
	}
}

func TestParseConfig_InvalidJSON(t *testing.T) {
	invalidJSON := `{"database": "test", "devices": [`

	_, err := ParseConfig(strings.NewReader(invalidJSON))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseConfig_EmptyConfig(t *testing.T) {
	emptyJSON := `{}`

	config, err := ParseConfig(strings.NewReader(emptyJSON))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Database != "" {
		t.Errorf("Expected empty database, got %s", config.Database)
	}

	if len(config.Devices) != 0 {
		t.Errorf("Expected 0 devices, got %d", len(config.Devices))
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	config := Config{
		Database: "postgresql://localhost:5432/test",
		Devices: []Device{
			{Name: "Test", File: "/dev/test", ID: 1},
		},
		Registers: []Register{
			{Name: "Test Reg", Address: 100},
		},
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestValidateConfig_MissingDatabase(t *testing.T) {
	config := Config{
		Devices: []Device{
			{Name: "Test", File: "/dev/test", ID: 1},
		},
		Registers: []Register{
			{Name: "Test Reg", Address: 100},
		},
	}

	err := config.Validate()
	if err != ErrMissingDatabase {
		t.Errorf("Expected ErrMissingDatabase, got %v", err)
	}
}

func TestValidateConfig_NoDevices(t *testing.T) {
	config := Config{
		Database: "postgresql://localhost:5432/test",
		Registers: []Register{
			{Name: "Test Reg", Address: 100},
		},
	}

	err := config.Validate()
	if err != ErrNoDevices {
		t.Errorf("Expected ErrNoDevices, got %v", err)
	}
}

func TestValidateConfig_NoRegisters(t *testing.T) {
	config := Config{
		Database: "postgresql://localhost:5432/test",
		Devices: []Device{
			{Name: "Test", File: "/dev/test", ID: 1},
		},
	}

	err := config.Validate()
	if err != ErrNoRegisters {
		t.Errorf("Expected ErrNoRegisters, got %v", err)
	}
}
