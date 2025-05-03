package repository

import (
	"context"
	"errors"
	"fmt"

	"calls-service/auth-service/internal/entity"
	"calls-service/pkg/postgres"
)

const (
	querySaveUser = `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	queryGetUser  = `SELECT id, username, password_hash FROM users WHERE username = $1 LIMIT 1`
)

var ErrUserAlreadyExists = errors.New("user already exists")

func (r *AuthRepo) SaveUser(user entity.User) error {
	ctx := context.Background()

	_, err := r.Pool.Exec(ctx, querySaveUser, user.Username, user.Password)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}

func (r *AuthRepo) GetUser(login string) (*entity.User, error) {
	ctx := context.Background()

	var user entity.User
	err := r.Pool.QueryRow(ctx, queryGetUser, login).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if postgres.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
