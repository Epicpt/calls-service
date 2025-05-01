package usecase

import (
	authpb "calls-service/auth-service/proto"
	"calls-service/rest-service/internal/repository"
)

type UseCase struct {
	repo       repository.Repository
	authClient authpb.AuthServiceClient
}

func New(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}
