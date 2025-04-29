package repository

import (
	"calls-service/auth-service/internal/entity"
	"calls-service/auth-service/pkg/postgres"
	"context"
	"fmt"
)

func (r *AuthRepo) SaveUser(user entity.User) error {
	ctx := context.Background()

	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`

	_, err := r.Pool.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}

func (r *AuthRepo) GetUser(login string) (*entity.User, error) {
	ctx := context.Background()
	query := `SELECT id, username, password_hash FROM users WHERE username = $1 LIMIT 1`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if postgres.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
