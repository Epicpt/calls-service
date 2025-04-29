package repository

import (
	"calls-service/auth-service/internal/entity"
	"calls-service/auth-service/pkg/postgres"
)

type Repository interface {
	SaveUser(entity.User) error
	GetUser(string) (*entity.User, error)
}

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}
