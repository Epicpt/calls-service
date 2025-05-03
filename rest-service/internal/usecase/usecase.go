package usecase

import (
	"context"

	authpb "calls-service/auth-service/proto"

	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/repository"
)

type UseCase interface {
	SaveCall(context.Context, entity.Call) error
	GetUserCalls(context.Context, int64) ([]entity.Call, error)
	GetUserCallByID(context.Context, int64, int64) (*entity.CallResponse, error)
	UpdateCallStatus(context.Context, int64, int64, string) error
	DeleteCall(context.Context, int64, int64) error
	RegisterUser(context.Context, entity.AuthRequest) error
	LoginUser(context.Context, entity.AuthRequest) (string, error)
}

type CallsService struct {
	repo       repository.Repository
	authClient authpb.AuthServiceClient
}

func New(repo repository.Repository, authClient authpb.AuthServiceClient) *CallsService {
	return &CallsService{
		repo:       repo,
		authClient: authClient,
	}
}
