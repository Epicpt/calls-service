package usecase

import (
	"errors"
	"fmt"

	"calls-service/auth-service/internal/entity"
	"calls-service/auth-service/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

func (uc *UseCase) Create(user entity.User) error {
	err := uc.repo.SaveUser(user)
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		return ErrUserAlreadyExists
	}
	return err
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
