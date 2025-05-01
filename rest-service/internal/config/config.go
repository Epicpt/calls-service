package config

import (
	c "calls-service/pkg/configuration"
)

type Config struct {
	HTTP
	Log c.Log
	PG  c.PG
	JWT c.JWT
}

type HTTP struct {
	Port string `env-required:"true" env:"HTTP_PORT"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := c.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
