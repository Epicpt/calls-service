package repository

import (
	"calls-service/pkg/postgres"
	"calls-service/rest-service/internal/entity"
	"context"
)

type Repository interface {
	SaveCall(context.Context, entity.Call) error
	GetUserCalls(context.Context, int64) ([]entity.Call, error)
	GetUserCallByID(context.Context, int64, int64) (*entity.Call, error)
	UpdateCallStatus(context.Context, int64, int64, string) error
	DeleteCall(context.Context, int64, int64) error
}

type CallsRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *CallsRepo {
	return &CallsRepo{pg}
}
