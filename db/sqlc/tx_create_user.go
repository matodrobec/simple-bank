package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		user, err := store.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = user
		return arg.AfterCreate(user)

	})
	return result, err
}
