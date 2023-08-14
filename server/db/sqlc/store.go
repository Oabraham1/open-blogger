package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

/* Store defines all functions to execute db queries and transactions */
type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	CreatePostTx(ctx context.Context, arg CreatePostTxParams) (CreatePostTxResult, error)
}

/* SQLStore provides all functions to execute SQL queries and transactions */
type SQLStore struct {
	ConnectionPool *pgxpool.Pool
	*Queries
}

/* NewStore creates a new store */
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		ConnectionPool: connPool,
		Queries:        New(connPool),
	}
}
