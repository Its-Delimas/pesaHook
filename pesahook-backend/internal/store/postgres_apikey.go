package store

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresAPIKeyStore struct {
	pool *pgxpool.Pool
}

func NewPostgresAPIKeyStore (pool *pgxpool.Pool) *PostgresAPIKeyStore {
	return &PostgresAPIKeyStore{pool: pool}
}
