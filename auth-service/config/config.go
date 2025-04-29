package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPC
	Log
	PG
	JWT
}

type GRPC struct {
	Port string `env-required:"true" env:"GRPC_PORT"`
}

type Log struct {
	Level string `env-required:"true" env:"LOG_LEVEL"`
}

type PG struct {
	URL     string `env-required:"true" env:"POSTGRES_URL"`
	PoolMax int    `env-required:"true" env:"POSTGRES_POOL_MAX"`
}

type JWT struct {
	Secret string `env-required:"true" env:"JWT_SECRET"`
}

func Load() (*Config, error) {
	// Для локального запуска
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Файл .env не найден")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
