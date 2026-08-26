package store

import (
	"context"
	"encoding/json"

	"github.com/Its-Delimas/pesaHook/internal/event"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresEventStore struct {
	pool *pgxpool.Pool
}

func NewPostgresEventStore(pool *pgxpool.Pool) *PostgresEventStore {
	return &PostgresEventStore{pool: pool}
}

func (s *PostgresEventStore) Save(e event.NormalizedEvent) error {
	providerMeta, err := json.Marshal(e.ProviderMeta)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(context.Background(), `
		INSERT INTO events (
			event_id, endpoint_id, event_type, provider, shortcode,
			transaction_id, amount, phone_number, account_reference,
			status, result_code, status_reason, occurred_at, received_at,
			provider_meta, raw
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, e.EventID, e.EndpointID, e.EventType, e.Provider, e.Shortcode, e.TransactionID,
		e.Amount, e.PhoneNumber, e.AccountReference, e.Status, e.ResultCode, e.StatusReason, e.OccurredAt,
		e.ReceivedAt, providerMeta, []byte(e.Raw),
	)
	return err
}

func (s *PostgresEventStore) GetByID(id string) (event.NormalizedEvent, error) {
	var e event.NormalizedEvent
	var providerMeta []byte
	var raw []byte

	row := s.pool.QueryRow(context.Background(), `
		SELECT event_id, endpoint_id, event_type, provider, shortcode,
		transaction_id, amount, phone_number, account_reference,
		status, result_code, status_reason, occurred_at, received_at,
		provider_meta, raw
		FROM events WHERE event_id = $1	`, id)

	err := row.Scan(
		&e.EventID, &e.EndpointID, &e.EventType, &e.Provider, &e.Shortcode,
		&e.TransactionID, &e.Amount, &e.PhoneNumber, &e.AccountReference,
		&e.Status, &e.ResultCode, &e.StatusReason, &e.OccurredAt, &e.ReceivedAt,
		&providerMeta, &raw,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return event.NormalizedEvent{}, ErrNotFound
		}
		return event.NormalizedEvent{}, err
	}

	json.Unmarshal(providerMeta, &e.ProviderMeta)
	e.Raw = json.RawMessage(raw)

	return e, nil
}

func (s *PostgresEventStore) ListByEndpoint(endpointID string) ([]event.NormalizedEvent, error) {
	rows, err := s.pool.Query(context.Background(), `
		SELECT event_id, endpoint_id, event_type, provider, shortcode,
		transaction_id, amount, phone_number, account_reference,
		status, result_code, status_reason, occurred_at, received_at,
		provider_meta, raw
		FROM events WHERE endpoint_id = $1
		ORDER BY received_at DESC
	`, endpointID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []event.NormalizedEvent
	for rows.Next() {
		var e event.NormalizedEvent
		var providerMeta []byte
		var raw []byte

		if err := rows.Scan(
			&e.EventID, &e.EndpointID, &e.EventType, &e.Provider, &e.Shortcode,
			&e.TransactionID, &e.Amount, &e.PhoneNumber, &e.AccountReference,
			&e.Status, &e.ResultCode, &e.StatusReason, &e.OccurredAt, &e.ReceivedAt,
			&providerMeta, &raw,
		); err != nil {
			return nil, err
		}

		json.Unmarshal(providerMeta, &e.ProviderMeta)
		e.Raw = json.RawMessage(raw)
		results = append(results, e)
	}
	return results, rows.Err()
}
