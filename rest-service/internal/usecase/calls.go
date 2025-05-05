package usecase

import (
	"context"
	"errors"
	"fmt"

	"calls-service/rest-service/internal/entity"
	"calls-service/rest-service/internal/repository"
)

var ErrCallNotFound = errors.New("call not found")

func (u *CallsService) SaveCall(ctx context.Context, call entity.Call) error {
	if err := u.repo.SaveCall(ctx, call); err != nil {
		return fmt.Errorf("failed to save call: %w", err)
	}
	return nil
}

func (u *CallsService) GetUserCalls(ctx context.Context, userID int64) ([]entity.CallResponse, error) {
	calls, err := u.repo.GetUserCalls(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user calls: %w", err)
	}
	return calls, nil
}

func (u *CallsService) GetUserCallByID(ctx context.Context, callID, userID int64) (*entity.CallResponse, error) {
	call, err := u.repo.GetUserCallByID(ctx, callID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCallNotFound) {
			return nil, ErrCallNotFound
		}
		return nil, fmt.Errorf("failed to get user call: %w", err)
	}
	return call, nil
}

func (u *CallsService) UpdateCallStatus(ctx context.Context, callID, userID int64, newStatus string) error {
	if err := u.repo.UpdateCallStatus(ctx, callID, userID, newStatus); err != nil {
		if errors.Is(err, repository.ErrCallNotFound) {
			return ErrCallNotFound
		}
		return fmt.Errorf("failed to update call status: %w", err)
	}
	return nil
}

func (u *CallsService) DeleteCall(ctx context.Context, callID int64, userID int64) error {
	if err := u.repo.DeleteCall(ctx, callID, userID); err != nil {
		if errors.Is(err, repository.ErrCallNotFound) {
			return ErrCallNotFound
		}
		return fmt.Errorf("failed to delete call: %w", err)
	}
	return nil
}
