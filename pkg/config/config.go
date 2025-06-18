package config

import (
	"encoding/json"
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

func GetConfig() Config {
	file, err := os.ReadFile("wattsup.json")
	if err != nil {
		log.WithError(err).Error("Error reading config file")
		os.Exit(1)
	}
	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		log.WithError(err).Error("Error reading config")
	}
	return config
}
