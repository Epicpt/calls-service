package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	authpb "calls-service/auth-service/proto"

	"calls-service/pkg/grpcserver"
	"calls-service/pkg/httpserver"
	"calls-service/pkg/logger"
	"calls-service/pkg/postgres"
	"calls-service/rest-service/internal/config"
	"calls-service/rest-service/internal/controller"
	"calls-service/rest-service/internal/repository"
	"calls-service/rest-service/internal/usecase"
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

	ctx := context.Background()
	conn, err := grpcserver.NewClient(ctx, cfg.Name+":"+cfg.GRPC.Port, cfg.ConnectionTimeout)
	if err != nil {
		l.Fatal().Msgf("failed to connect to auth service: %v", err)
	}
	defer func() { _ = conn.Close() }()

	l.Info().Msg("GRPC server connected")

	authClient := authpb.NewAuthServiceClient(conn)

	// Use case
	callsService := usecase.New(repository.New(pg), authClient)

	// Run server
	httpServer := httpserver.New(cfg.HTTP.Port)

	handler := controller.New(callsService, l)
	controller.NewCallsRoutes(httpServer.Engine, handler)

	httpServer.Start()

	l.Info().Msg("Server start")

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
