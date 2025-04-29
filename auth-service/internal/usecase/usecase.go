package usecase

import "calls-service/auth-service/internal/repository"

type UseCase struct {
	repo repository.Repository
}

func New(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}
