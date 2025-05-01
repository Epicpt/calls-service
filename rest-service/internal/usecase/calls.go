package usecase

import (
	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/repository"
	"context"
	"errors"
	"fmt"
)

var ErrCallNotFound = errors.New("call not found")

func (u *UseCase) SaveCall(ctx context.Context, call entity.Call) error {
	if err := u.repo.SaveCall(ctx, call); err != nil {
		return fmt.Errorf("failed to save call: %w", err)
	}
	return nil
}

func (u *UseCase) GetUserCalls(ctx context.Context, userID int64) ([]entity.Call, error) {
	calls, err := u.repo.GetUserCalls(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user calls: %w", err)
	}
	return calls, nil
}

func (u *UseCase) GetUserCallByID(ctx context.Context, callID, userID int64) (*entity.Call, error) {
	call, err := u.repo.GetUserCallByID(ctx, callID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCallNotFound) {
			return nil, ErrCallNotFound
		}
		return nil, fmt.Errorf("failed to get user call: %w", err)
	}
	return call, nil
}

func (u *UseCase) UpdateCallStatus(ctx context.Context, callID, userID int64, newStatus string) error {
	return u.repo.UpdateCallStatus(ctx, callID, userID, newStatus)
}

func (u *UseCase) DeleteCall(ctx context.Context, callID int64, userID int64) error {
	return u.repo.DeleteCall(ctx, callID, userID)
}
