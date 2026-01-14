package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MIDI struct {
		PortName      string `yaml:"port_name"`
		Channel       uint8  `yaml:"channel"`
		TempoChangeCC uint8  `yaml:"tempo_change_cc"`
	} `yaml:"midi"`
	Smoothing struct {
		WindowWidth int `yaml:"window_width"`
	} `yaml:"smoothing"`
	DataSource struct {
		Type string `yaml:"type"`
		CSV  struct {
			Path string `yaml:"path"`
		} `yaml:"csv"`
	} `yaml:"data_source"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &config, nil
}
