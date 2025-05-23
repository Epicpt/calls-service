package config

import (
	c "calls-service/pkg/configuration"
)

type Config struct {
	GRPC
	Log c.Log
	PG  c.PG
	JWT c.JWT
}

type GRPC struct {
	Port string `env-required:"true" env:"GRPC_PORT"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := c.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
