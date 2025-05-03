package usecase

import (
	authpb "calls-service/auth-service/proto"
	"calls-service/rest-service/internal/entity"
	"context"
)

func (u *CallsService) RegisterUser(ctx context.Context, req entity.AuthRequest) error {
	_, err := u.authClient.Register(ctx, &authpb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	return err
}

func (u *CallsService) LoginUser(ctx context.Context, req entity.AuthRequest) (string, error) {
	token, err := u.authClient.Login(ctx, &authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	return token.Token, err
}
