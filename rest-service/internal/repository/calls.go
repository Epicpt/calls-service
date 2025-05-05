package repository

import (
	"context"
	"errors"
	"fmt"

	"calls-service/pkg/postgres"
	"calls-service/rest-service/internal/entity"

	"github.com/rs/zerolog/log"
)

var ErrCallNotFound = errors.New("call not found")

const (
	querySaveCall         = `INSERT INTO calls (client_name, phone_number, description, user_id) VALUES ($1, $2, $3, $4)`
	queryGetUserCalls     = `SELECT id, client_name, phone_number, description, status, created_at FROM calls WHERE user_id = $1 ORDER BY created_at DESC`
	queryGetUserCallByID  = `SELECT id, client_name, phone_number, description, status, created_at FROM calls WHERE id = $1 AND user_id = $2`
	queryUpdateCallStatus = `UPDATE calls SET status = $1 WHERE id = $2 AND user_id = $3`
	queryDeleteCall       = `DELETE FROM calls WHERE id = $1 AND user_id = $2`
)

func (r *CallsRepo) SaveCall(ctx context.Context, call entity.Call) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !postgres.IsTxClosed(err) {
			log.Error().Err(err).Msg("failed to rollback transaction")
		}
	}()

	_, err = tx.Exec(ctx, querySaveCall,
		call.ClientName,
		call.PhoneNumber,
		call.Description,
		call.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *CallsRepo) GetUserCalls(ctx context.Context, userID int64) ([]entity.CallResponse, error) {
	rows, err := r.Pool.Query(ctx, queryGetUserCalls, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []entity.CallResponse
	for rows.Next() {
		var call entity.CallResponse
		if err := rows.Scan(
			&call.ID,
			&call.ClientName,
			&call.PhoneNumber,
			&call.Description,
			&call.Status,
			&call.CreatedAt,
		); err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return calls, nil
}

func (r *CallsRepo) GetUserCallByID(ctx context.Context, callID, userID int64) (*entity.CallResponse, error) {
	var call entity.CallResponse

	err := r.Pool.QueryRow(ctx, queryGetUserCallByID, callID, userID).Scan(
		&call.ID,
		&call.ClientName,
		&call.PhoneNumber,
		&call.Description,
		&call.Status,
		&call.CreatedAt,
	)

	if err != nil {
		if postgres.IsNotFoundError(err) {
			return nil, ErrCallNotFound
		}
		return nil, err
	}

	return &call, nil
}

func (r *CallsRepo) UpdateCallStatus(ctx context.Context, callID int64, userID int64, newStatus string) error {
	cmdTag, err := r.Pool.Exec(ctx, queryUpdateCallStatus, newStatus, callID, userID)
	if err != nil {
		return fmt.Errorf("failed to update call status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrCallNotFound
	}

	return nil
}

func (r *CallsRepo) DeleteCall(ctx context.Context, callID int64, userID int64) error {
	cmdTag, err := r.Pool.Exec(ctx, queryDeleteCall, callID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete call: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrCallNotFound
	}

	return nil
}
