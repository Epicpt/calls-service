package app

import (
	"calls-service/pkg/httpserver"
	"calls-service/pkg/logger"
	"calls-service/pkg/postgres"
	"calls-service/rest-service/internal/config"
	"calls-service/rest-service/internal/controller"
	"calls-service/rest-service/internal/repository"
	"calls-service/rest-service/internal/usecase"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Info().Msg("Logger initialized")

	pg, err := postgres.New(cfg.PG.URL, cfg.PG.PoolMax)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pg.Close()

	l.Info().Msg("PostgreSQL initialized")

	// Use case
	usecase := usecase.New(repository.New(pg))

	// Run server
	httpServer := httpserver.New(cfg.Port)

	handler := controller.New(usecase, l)
	controller.NewCallsRoutes(httpServer.Engine, handler)

	httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info().Msgf("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error().Err(err).Msg("app - Run - httpServer.Shutdown")
	}
}
