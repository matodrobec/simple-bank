package db

import (
	"cmp"
	"context"
	"slices"
)


type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// for debuging
		// txNane := ctx.Value(txKey)

		// log.Println(txNane, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// log.Println(txNane, "create from entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// log.Println(txNane, "create to entry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, result.ToAccount, err = addMoney(
			ctx,
			q,
			arg.FromAccountID,
			arg.ToAccountID,
			arg.Amount,
		)

		return err

		/**
		 * Update DeadLock
		 * 		update deadlock
		 **/
		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	ID:     arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })

		/**
		 * Select from DeadLock
		 * 		foreign key deadlock
		 **/

		// // log.Println(txNane, "get from account")

		// result.FromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// // log.Println(txNane, "update from account")
		// result.FromAccount, err = q.UpdateAccountAndGet(
		// 	ctx,
		// 	UpdateAccountAndGetParams{
		// 		ID:      result.FromAccount.ID,
		// 		Balance: result.FromAccount.Balance - arg.Amount,
		// 	},
		// )
		// if err != nil {
		// 	return err
		// }

		// // log.Println(txNane, "get to account")
		// result.ToAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// // log.Println(txNane, "update to account")
		// result.ToAccount, err = q.UpdateAccountAndGet(
		// 	ctx,
		// 	UpdateAccountAndGetParams{
		// 		ID:      result.ToAccount.ID,
		// 		Balance: result.ToAccount.Balance + arg.Amount,
		// 	},
		// )

		// return err

	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	fromAccountId int64,
	toAccountId int64,
	amount int64,
) (fromAccount Account, toAccount Account, err error) {
	type data struct {
		id     int64
		amount int64
		dest   *Account
	}
	list := []data{
		{id: fromAccountId, amount: -amount, dest: &fromAccount},
		{id: toAccountId, amount: amount, dest: &toAccount},
	}

	// sort because deadlock
	slices.SortFunc(list, func(a, b data) int {
		return cmp.Compare(a.id, b.id)
	})

	for _, value := range list {
		*value.dest, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     value.id,
			Amount: value.amount,
		})
		if err != nil {
			return
		}
	}
	return
}
