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

func (s *PostgresAPIKeyStore) GetByHash(hash string) (apikey.APIKey, error) {
	var k apikey.APIKey
	row := s.pool.QueryRow(context.Background(), `
		SELECT id, account_id, key_hash, created_at FROM api_keys WHERE key_hash = $1`, hash)

	err := row.Scan(&k.ID, &k.AccountID, &k.KeyHash, &k.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return apikey.APIKey{}, ErrNotFound
		}
		return apikey.APIKey{}, err
	}
	return k, nil
}
