package db

import (
	"context"
	"testing"

	"github.com/mahmoud24598salah/MSM_Bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	accountTo := createRandomAccount(t)
	accountFrom := createRandomAccount(t)
	arg := CreateTransferParams{
		ToAccountID:   accountTo.ID,
		FromAccountID: accountFrom.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.Equal(t, accountTo.ID, transfer.ToAccountID)
	require.Equal(t, accountFrom.ID, transfer.FromAccountID)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransferById(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	tran, err := testQueries.GetTransferById(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tran)
	require.Equal(t, transfer, tran)
}

func TestGetTransfersByAmount(t *testing.T) {
	amount := util.RandomMoney()
	for i := 0; i < 3; i++ {
		accountTo := createRandomAccount(t)
		accountFrom := createRandomAccount(t)
		arg := CreateTransferParams{
			FromAccountID: accountFrom.ID,
			ToAccountID:   accountTo.ID,
			Amount:        amount,
		}
		_, _ = testQueries.CreateTransfer(context.Background(), arg)
	}
	transfers, err := testQueries.GetTransfersByAmount(context.Background(), amount)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 3)
}

func TestGetTransfersByFromAccountId(t *testing.T) {
	accountFrom := createRandomAccount(t)
	for i := 0; i < 3; i++ {
		accountTo := createRandomAccount(t)
		arg := CreateTransferParams{
			FromAccountID: accountFrom.ID,
			ToAccountID:   accountTo.ID,
			Amount:        util.RandomMoney(),
		}
		_, _ = testQueries.CreateTransfer(context.Background(), arg)
	}

	transfers, err := testQueries.GetTransfersByFromAccountId(context.Background(), accountFrom.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 3)
	for _, transfer := range transfers {
		require.Equal(t, accountFrom.ID, transfer.FromAccountID)
	}
}

func TestGetTransfersByToAccountId(t *testing.T) {
	accountTo := createRandomAccount(t)
	for i := 0; i < 3; i++ {
		accountFrom := createRandomAccount(t)
		arg := CreateTransferParams{
			FromAccountID: accountFrom.ID,
			ToAccountID:   accountTo.ID,
			Amount:        util.RandomMoney(),
		}
		_, _ = testQueries.CreateTransfer(context.Background(), arg)
	}

	transfers, err := testQueries.GetTransfersByToAccountId(context.Background(), accountTo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 3)
	for _, transfer := range transfers {
		require.Equal(t, accountTo.ID, transfer.ToAccountID)
	}
}
