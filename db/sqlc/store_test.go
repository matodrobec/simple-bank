package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// var testStore SqlStore

func TestTransferTx(t *testing.T) {
	// store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	log.Println(">> before: ", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errCh := make(chan error)
	resultCh := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errCh <- err
			resultCh <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)

		result := <-resultCh
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Test Accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// Test Balance
		log.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)

		diffFrom := account1.Balance - fromAccount.Balance
		diffTo := toAccount.Balance - account2.Balance

		require.Equal(t, diffFrom, diffTo)
		require.True(t, diffFrom > 0)
		require.True(t, diffFrom%amount == 0)

		k := int(diffFrom / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	log.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	log.Println(">> before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errCh := make(chan error)

	for i := 0; i < n; i++ {
		fromId := account1.ID
		toId := account2.ID

		if i%2 == 0 {
			fromId = account2.ID
			toId = account1.ID
		}
		// txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			// ctx := context.WithValue(context.Background(), txKey, txName)
			ctx := context.Background()
			_, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromId,
				ToAccountID:   toId,
				Amount:        amount,
			})
			errCh <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	log.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
