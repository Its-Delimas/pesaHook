package store

import (
	"context"

	"github.com/Its-Delimas/pesahook/internal/apikey"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAPIKeyStore struct {
	pool *pgxpool.Pool
}

func NewPostgresAPIKeyStore(pool *pgxpool.Pool) *PostgresAPIKeyStore {
	return &PostgresAPIKeyStore{pool: pool}
}

func (s *PostgresAPIKeyStore) Save(k apikey.APIKey) error {
	_, err := s.pool.Exec(context.Background(), `
		INSERT INTO api_keys (id,account_id,key_hash,created_at)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (id) DO NOTHING`,
		k.ID, k.AccountID, k.KeyHash, k.CreatedAt,
	)
	return err
}
