package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}
type TransferTxParms struct {
	FromAccountID int64 `json:"form_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"form_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return fmt.Errorf("tx error: %v, rbError: %v", err, rberr)
		}
		return err
	}
	return tx.Commit()
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParms) (TransferTxResult, error) {
	var result TransferTxResult
	var err error
	err = store.execTx(ctx,func(q *Queries) error {
		result.Transfer ,err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		result.FromEntry ,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		// accounts after update //TODO
		return nil
	})
	return result,err
}
