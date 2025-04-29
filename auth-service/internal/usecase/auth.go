package usecase

import (
	"calls-service/auth-service/internal/entity"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")

func (uc *UseCase) Create(user entity.User) error {
	if err := uc.repo.SaveUser(user); err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) GetUser(login string) (*entity.User, error) {
	user, err := uc.repo.GetUser(login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
