package db

import (
	"context"
	"database/sql"
	"fmt"
)


type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB)*Store{
	return &Store{
		db :db,
		Queries: New(db),
	}
}

func(store *Store) execTx(ctx context.Context,fn func(*Queries)error)error{
	tx,err := store.db.BeginTx(ctx,nil)
	if err != nil{
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil{
		rberr := tx.Rollback()
		if rberr != nil{
			return fmt.Errorf("tx error: %v, rbError: %v",err,rberr)
		}
		return err
	}
	return tx.Commit() 
}