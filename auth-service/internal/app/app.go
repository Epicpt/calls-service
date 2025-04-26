package app

import (
	"calls-service/auth-service/config"
	"calls-service/auth-service/internal/controller"
	"calls-service/auth-service/pkg/grpcserver"
	"calls-service/auth-service/pkg/logger"
	"calls-service/auth-service/pkg/postgres"

	authpb "calls-service/auth-service/proto"

	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/reflection"
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

	// usecase init

	server := grpcserver.New(cfg.Port)

	authService := controller.NewAuthService()
	authpb.RegisterAuthServiceServer(server, authService)

	reflection.Register(server.GrpcServer)

	server.Start()

	l.Info().Msg("Server start")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info().Msgf("app - Run - signal: %s", s.String())
	case err := <-server.Notify():
		l.Error().Err(err).Msg("app - Run - grpcServer.Notify")
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		l.Error().Err(err).Msg("app - Run - grpcServer.Shutdown")
	}
}
