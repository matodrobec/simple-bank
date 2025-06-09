package db

import (
	"context"
	"errors"
)

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// tx, err := store.connPool.BeginTx(ctx, nil)
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			// return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			return errors.Join(err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}