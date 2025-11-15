package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database DatabaseConfig `json:"database" yaml:"database"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := parseConfig(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseConfig(file *os.File, config *Config) error {
	if err := json.NewDecoder(file).Decode(config); err == nil {
		return nil
	}

	file.Seek(0, 0) // Reset file pointer
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return err
	}

	return nil
}