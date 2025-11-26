package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout_seconds"`
		WriteTimeout int    `yaml:"write_timeout_seconds"`
		IdleTimeout  int    `yaml:"idle_timeout_seconds"`
	} `yaml:"server"`
	DatabaseUrl string `yaml:"database_url"`
}

func Load(pathToFile string) (*Config, error) {
	raw, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err = yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
