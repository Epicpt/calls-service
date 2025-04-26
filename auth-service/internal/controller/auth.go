package controller

import (
	authpb "calls-service/auth-service/proto"
	"context"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Ping(ctx context.Context, req *authpb.PingRequest) (*authpb.PingResponse, error) {
	return &authpb.PingResponse{Message: "pong"}, nil
}
