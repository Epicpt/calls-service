package repository

import (
	"context"

	"calls-service/pkg/postgres"
	"calls-service/rest-service/internal/entity"
)

type Repository interface {
	SaveCall(context.Context, entity.Call) error
	GetUserCalls(context.Context, int64) ([]entity.CallResponse, error)
	GetUserCallByID(context.Context, int64, int64) (*entity.CallResponse, error)
	UpdateCallStatus(context.Context, int64, int64, string) error
	DeleteCall(context.Context, int64, int64) error
}

type CallsRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *CallsRepo {
	return &CallsRepo{pg}
}
