package db

import (
	"context"

	"github.com/google/uuid"
)

type RotateSessionTxParams struct {
	FromSessionID uuid.UUID `json:"from_session_id"`
	ToSession     CreateSessionParams   `json:"to_session"`
}

type RotateSessionTxResult Session

func (store *SqlStore) RotateSessionTx(ctx context.Context, arg RotateSessionTxParams) (RotateSessionTxResult, error) {
	var result RotateSessionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.BlockSession(ctx, arg.FromSessionID)
		if err != nil {
			return err
		}

		sessionArgs := CreateSessionParams{
			ID:           arg.ToSession.ID,
			Username:     arg.ToSession.Username,
			RefreshToken: arg.ToSession.RefreshToken,
			UserAgent:    arg.ToSession.UserAgent,
			ClientIp:     arg.ToSession.ClientIp,
			IsBlocked:    arg.ToSession.IsBlocked,
			ExpiresAt:    arg.ToSession.ExpiresAt,
		}
		newSession, err := q.CreateSession(ctx, sessionArgs)
		result = RotateSessionTxResult(newSession)

		return err
	})
	return result, err
}
