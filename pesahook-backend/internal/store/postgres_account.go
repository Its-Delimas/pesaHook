package store

import (
	"context"

	"github.com/Its-Delimas/pesahook/internal/account"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAccountStore struct {
	pool *pgxpool.Pool
}

func NewPostgresAccountStore(pool *pgxpool.Pool) *PostgresAccountStore {
	return &PostgresAccountStore{pool: pool}
}

func (s *PostgresAccountStore) Save(a account.Account) error {
	_, err := s.pool.Exec(context.Background(), `
	INSERT INTO accounts (id, email, created_at)
	VALUES ($1,$2,$3)
	ON CONFLICT (id) DO NOTHING`,
		a.ID, a.Email, a.CreatedAt)
	return err
}
