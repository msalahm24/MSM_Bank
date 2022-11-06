package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)
	 errors := make(chan error)
	 results := make(chan TransferTxResult)
	for i:= 0 ;i < n ;i++{
		go func ()  {
			result ,err := store.TransferTx(context.Background(),TransferTxParms{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errors <- err
			results <- result
		}()
	}

	for i := 0 ; i< n ; i++{
		err := <- errors
		require.NoError(t,err)

		result:= <- results
		require.NotEmpty(t,result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,account1.ID,transfer.FromAccountID)
		require.Equal(t,account2.ID,transfer.ToAccountID)
		require.Equal(t,amount,transfer.Amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)

		_,err = store.GetTransferById(context.Background(),transfer.ID)
		require.NoError(t,err)

		fromEntry := result.FromEntry
		require.NotEmpty(t,fromEntry)
		require.Equal(t,account1.ID,fromEntry.AccountID)
		require.Equal(t,-amount,fromEntry.Amount)
		require.NotZero(t,fromEntry.ID)
		require.NotZero(t,fromEntry.CreatedAt)
		_,err = store.GetEntryById(context.Background(),fromEntry.ID)
		require.NoError(t,err)


		toEntry := result.ToEntry
		require.NotEmpty(t,toEntry)
		require.Equal(t,account2.ID,toEntry.AccountID)
		require.Equal(t,amount,toEntry.Amount)
		require.NotZero(t,toEntry.ID)
		require.NotZero(t,toEntry.CreatedAt)
		_,err = store.GetEntryById(context.Background(),toEntry.ID)
		require.NoError(t,err)
	}

}
