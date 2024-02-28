package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries & transactions
type PostgreSQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// Creates a new Store
func NewPostgreSQLStore(connPool *pgxpool.Pool) Store {
	return &PostgreSQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// Executes a function within a database transaction
func (store *PostgreSQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {

	// Initialize a transaction
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
