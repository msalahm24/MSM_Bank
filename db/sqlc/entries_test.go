package db

import (
	"context"
	"testing"

	"github.com/mahmoud24598salah/MSM_Bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntriesByAccountId(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 3; i++ {
		entry := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomMoney(),
		}
		_, _ = testQueries.CreateEntry(context.Background(), entry)
	}
	entries, err := testQueries.GetEntriesByAccountId(context.Background(), account.ID)
	require.NoError(t, err)
	require.Len(t, entries, 3)
	for _, ent := range entries {
		require.NotEmpty(t, ent)
		require.NotZero(t, ent.ID)
		require.NotZero(t, ent.Amount)
		require.NotZero(t, ent.AccountID)
		require.Equal(t, account.ID, ent.AccountID)
		require.NotZero(t, ent.CreatedAt)
	}

}

func TestGetEntriesByAmount(t *testing.T) {
	var Tempentries []Entry
	amount := util.RandomMoney()
	for i := 0; i < 3; i++ {
		temp := createRandomAccount(t)
		arg := CreateEntryParams{
			AccountID: temp.ID,
			Amount:    amount,
		}
		ent, _ := testQueries.CreateEntry(context.Background(), arg)
		Tempentries = append(Tempentries, ent)
	}
	entries, err := testQueries.GetEntriesByAmount(context.Background(), amount)
	require.NoError(t, err)
	require.Len(t, entries, 3)
	for _, ent := range entries {
		require.Contains(t, Tempentries, ent)
	}
}

func TestGetEntryById(t *testing.T) {
	entry := CreateRandomEntry(t)
	ent, err := testQueries.GetEntryById(context.Background(),entry.ID)
	require.NoError(t,err)
	require.Equal(t,entry,ent)
}
