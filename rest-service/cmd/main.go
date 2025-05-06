// @title Calls service
// @version 1.0
// @description API для обработки заявок
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"log"

	"calls-service/rest-service/internal/app"
	"calls-service/rest-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error creating config: %s", err)
	}

	log.Println("Config initializated")

	app.Run(cfg)
}
