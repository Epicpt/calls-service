package config

import (
	c "calls-service/pkg/configuration"
	"time"
)

type Config struct {
	HTTP
	GRPC
	Log c.Log
	PG  c.PG
	JWT c.JWT
}

type HTTP struct {
	Port string `env-required:"true" env:"HTTP_PORT"`
}

type GRPC struct {
	Port              string        `env-required:"true" env:"GRPC_PORT"`
	Name              string        `env-required:"true" env:"GRPC_NAME"`
	ConnectionTimeout time.Duration `env-required:"true" env:"GRPC_CLIENT_CONN_TIMEOUT"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := c.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
