package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SwitchBot SwitchBot
	MfLight   MfLight
}

type SwitchBot struct {
	Token    string
	DeviceId string
}

type MfLight struct {
	URL      string
	MobileId string
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("app", &c); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &c, nil
}
