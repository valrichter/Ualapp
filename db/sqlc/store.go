package db

import (
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
