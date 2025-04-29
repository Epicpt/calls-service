package app

import (
	"calls-service/auth-service/config"
	"calls-service/auth-service/internal/controller"
	"calls-service/auth-service/internal/repository"
	"calls-service/auth-service/internal/usecase"
	"calls-service/pkg/grpcserver"
	"calls-service/pkg/logger"
	"calls-service/pkg/postgres"

	authpb "calls-service/auth-service/proto"

	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Level)

	l.Info().Msg("Logger initialized")

	pg, err := postgres.New(cfg.URL, cfg.PoolMax)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pg.Close()

	l.Info().Msg("PostgreSQL initialized")

	authUseCase := usecase.New(repository.New(pg))

	server := grpcserver.New(cfg.Port)

	authService := controller.New(authUseCase, l)
	authpb.RegisterAuthServiceServer(server, authService)

	reflection.Register(server.GrpcServer) // local testing

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
