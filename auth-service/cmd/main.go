package main

import (
	"calls-service/auth-service/config"
	"calls-service/auth-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error creating config: %s", err)
	}

	log.Println("Config initializated")

	app.Run(cfg)
}
