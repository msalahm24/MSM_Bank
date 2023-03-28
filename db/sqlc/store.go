package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParms) (TransferTxResult, error)
	Querier
}
type SqlStore struct {
	*Queries
	db *sql.DB
}
type TransferTxParms struct {
	FromAccountID int64 `json:"form_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"form_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParms) (TransferTxResult, error) {
	var result TransferTxResult
	var err error

	txName := ctx.Value(txKey)
	fmt.Println(txName, "create transfer ")
	err = store.execTx(ctx, func(q *Queries) error {
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1 ")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2 ")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}


		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount,err= addMoney(ctx,q,arg.FromAccountID,-arg.Amount,arg.ToAccountID,arg.Amount)
			if err != nil{
				return err
			}
		} else {
			result.ToAccount, result.FromAccount,err= addMoney(ctx,q,arg.ToAccountID,arg.Amount,arg.FromAccountID,-arg.Amount)
			if err != nil{
				return err
			}
		}

		// accounts after update //TODO
		return nil
	})
	return result, err
}


func addMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,  
)(account1 Account,account2 Account,err error){
	account1 , err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		Amount: amount1,
		ID: accountId1,
	})
	if err != nil{
		return
	}
	account2 ,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		Amount: amount2,
		ID: accountId2,
	})
	if err != nil {
		return
	}
	
	return
}
