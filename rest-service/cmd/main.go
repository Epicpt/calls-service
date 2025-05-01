package main

import (
	"calls-service/rest-service/internal/app"
	"calls-service/rest-service/internal/config"

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
