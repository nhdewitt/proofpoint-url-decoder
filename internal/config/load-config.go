package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port string `json:"port"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	configFile, err := os.Open(file)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %w", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return config, fmt.Errorf("error decoding config file: %w", err)
	}

	return config, nil
}
