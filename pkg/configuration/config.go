package configuration

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

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

func Load[T any](target *T) error {
	// Для локального запуска
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Файл .env не найден")
	}

	if err := env.Parse(target); err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	return nil
}
