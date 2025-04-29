package controller

import (
	"calls-service/auth-service/internal/usecase"
	authpb "calls-service/auth-service/proto"
	"context"

	"github.com/rs/zerolog"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	u usecase.UseCase
	l zerolog.Logger
}

func New(u *usecase.UseCase, l zerolog.Logger) *AuthService {
	return &AuthService{u: *u, l: l}
}

func (s *AuthService) Ping(ctx context.Context, req *authpb.PingRequest) (*authpb.PingResponse, error) {
	return &authpb.PingResponse{Message: "pong"}, nil
}
