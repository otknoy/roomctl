package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port      int
	SwitchBot SwitchBot
}

type SwitchBot struct {
	Token      string
	Secret     string
	DeviceId   []string
	DeviceName []string
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("app", &c); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &c, nil
}
