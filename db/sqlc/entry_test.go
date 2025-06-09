package db

import (
	"context"
	"testing"
	"time"

	"github.com/matodrobec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:     util.RandomMoney(),
	}
	entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomAccount(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.ID, entry2.ID)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	entry1 := createRandomEntry(t)
	// for i:=0; i<10; i++ {
	//     createRandomEntry(t)
	// }

	arg := ListEntriesParams{
		AccountID: entry1.AccountID,
		Limit:      10,
		Offset:     0,
	}
	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
