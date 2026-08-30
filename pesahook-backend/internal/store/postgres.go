package store

import (
	"context"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

type PostgresEndpointStore struct {
	pool *pgxpool.Pool
}

func NewPostgresEndpointStore(pool *pgxpool.Pool) *PostgresEndpointStore {
	return &PostgresEndpointStore{pool: pool}
}

func (s *PostgresEndpointStore) Save(e endpoint.Endpoint) error {
	_, err := s.pool.Exec(context.Background(), `
		INSERT INTO endpoints (id, account_id,provider, shortcode, event_types,destination_url,secret,ingest_path, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET 
			destination_url = EXCLUDED.destination_url,
			event_types = EXCLUDED.event_types`,
		e.ID, e.AccountID, e.Provider, e.Shortcode, e.EventTypes, e.DestinationURL, e.Secret, e.IngestPath, e.CreatedAt,
	)
	return err
}

func (s *PostgresEndpointStore) GetByID(id string) (endpoint.Endpoint, error) {
	var e endpoint.Endpoint

	row := s.pool.QueryRow(context.Background(), `
		SELECT id,account_id, provider, shortcode, event_types, destination_url, secret, ingest_path, created_at
		FROM endpoints WHERE id = $1`, id)

	err := row.Scan(&e.ID, &e.AccountID, &e.Provider, &e.Shortcode, &e.EventTypes, &e.DestinationURL, &e.Secret, &e.IngestPath, &e.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return endpoint.Endpoint{}, ErrNotFound
		}
		return endpoint.Endpoint{}, err
	}
	return e, nil
}
