package store

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresAccountStore struct {
	pool *pgxpool.Pool
}

func NewPostgresAccountStore(pool *pgxpool.Pool) *PostgresAccountStore {
	return &PostgresAccountStore{pool: pool}
}
