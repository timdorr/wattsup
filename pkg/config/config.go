package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/apex/log"
)

type Config struct {
	Database  string     `json:"database"`
	Devices   []Device   `json:"devices"`
	Registers []Register `json:"registers"`
}

type Device struct {
	Name string `json:"name"`
	File string `json:"file"`
	ID   int    `json:"id"`
}

type Register struct {
	Name    string `json:"name"`
	Address uint16 `json:"address"`
}

// ConfigReader interface for testing
type ConfigReader interface {
	ReadConfig(filename string) (Config, error)
}

// FileConfigReader reads config from file system
type FileConfigReader struct{}

func (r *FileConfigReader) ReadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	return ParseConfig(file)
}

// ParseConfig parses config from any io.Reader
func ParseConfig(reader io.Reader) (Config, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// GetConfig returns config using the default file reader
func GetConfig() Config {
	reader := &FileConfigReader{}
	config, err := reader.ReadConfig("wattsup.json")
	if err != nil {
		log.WithError(err).Error("Error reading config file")
		os.Exit(1)
	}
	return config
}

// ValidateConfig validates the configuration
func (c *Config) Validate() error {
	if c.Database == "" {
		return ErrMissingDatabase
	}
	if len(c.Devices) == 0 {
		return ErrNoDevices
	}
	if len(c.Registers) == 0 {
		return ErrNoRegisters
	}
	return nil
}
