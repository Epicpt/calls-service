package controller

import (
	"calls-service/auth-service/internal/entity"
	"calls-service/auth-service/internal/services"
	"calls-service/auth-service/internal/usecase"
	authpb "calls-service/auth-service/proto"
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	u usecase.UseCase
	l zerolog.Logger
}

func New(u *usecase.UseCase, l zerolog.Logger) *AuthService {
	return &AuthService{u: *u, l: l}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	username, password, err := validateAndCleanCredentials(req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hashedPass, err := services.HashPassword(password)
	if err != nil {
		s.l.Err(err).Msg("Failed to hash password")
		return nil, status.Error(codes.Internal, "Failed to hash password")
	}

	user := entity.User{
		Username: username,
		Password: hashedPass,
	}

	if err := s.u.Create(user); err != nil {
		s.l.Err(err).Msg("Failed to create user")
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	s.l.Info().Str("username", user.Username).Msg("User registered successfully")

	return &authpb.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	username, password, err := validateAndCleanCredentials(req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.u.GetUser(username)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, "Invalid username")
		}
		s.l.Err(err).Msg("failed to get user")
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	if !services.CheckPassword(password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid password")
	}

	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		s.l.Err(err).Msg("failed to generate token")
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	s.l.Info().Interface("user", user).Str("token", token).Msg("User logged in successfully")
	return &authpb.LoginResponse{Token: token}, nil
}

func validateAndCleanCredentials(username, password string) (string, string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if len(username) == 0 || len(password) == 0 {
		return "", "", errors.New("username and password must be provided")
	}
	if len(username) > 32 || len(password) > 72 {
		return "", "", errors.New("username or password too long")
	}
	return username, password, nil
}
