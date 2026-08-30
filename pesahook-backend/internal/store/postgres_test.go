// file: internal/store/postgres_test.go
//go:build integration

package store

import (
	"context"
	"testing"
	"time"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
	"github.com/Its-Delimas/pesahook/internal/event"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testPool(t *testing.T) *pgxpool.Pool {
	pool, err := NewPostgresPool(context.Background(), "postgres://pesahook:devpass@localhost:5432/pesahook")
	if err != nil {
		t.Fatalf("failed to connect to test postgres: %v", err)
	}
	return pool
}

func TestPostgresEndpointStore_SaveAndGet(t *testing.T) {
	pool := testPool(t)
	defer pool.Close()

	store := NewPostgresEndpointStore(pool)

	ep := endpoint.Endpoint{
		ID:             "test-endpoint-1",
		AccountID:      "test-account-1",
		Provider:       "daraja",
		Shortcode:      "600000",
		EventTypes:     []string{"stk_push", "c2b_confirmation"},
		DestinationURL: "https://example.com/webhook",
		Secret:         "testsecret",
		IngestPath:     "/ingest/daraja/test-endpoint-1",
		CreatedAt:      time.Now(),
	}

	if err := store.Save(ep); err != nil {
		t.Fatalf("failed to save endpoint: %v", err)
	}
	defer pool.Exec(context.Background(), "DELETE FROM endpoints WHERE id = $1", ep.ID)

	fetched, err := store.GetByID(ep.ID)
	if err != nil {
		t.Fatalf("failed to fetch endpoint: %v", err)
	}

	if fetched.AccountID != ep.AccountID {
		t.Errorf("expected account ID: %s, got %s", ep.AccountID, fetched.AccountID)
	}

	if fetched.DestinationURL != ep.DestinationURL {
		t.Errorf("expected destination URL %s, got %s", ep.DestinationURL, fetched.DestinationURL)
	}
	if len(fetched.EventTypes) != 2 {
		t.Errorf("expected 2 event types, got %d", len(fetched.EventTypes))
	}
}

func TestPostgresEndpointStore_GetByID_NotFound(t *testing.T) {
	pool := testPool(t)
	defer pool.Close()

	store := NewPostgresEndpointStore(pool)

	_, err := store.GetByID("nonexistent-id")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresEventStore_SaveAndGet(t *testing.T) {
	pool := testPool(t)
	defer pool.Close()

	// events references endpoints via foreign key, so an endpoint must exist first
	endpointStore := NewPostgresEndpointStore(pool)
	ep := endpoint.Endpoint{
		ID: "test-endpoint-2", AccountID: "test-account-1", Provider: "daraja", Shortcode: "600001",
		EventTypes: []string{"stk_push"}, DestinationURL: "https://example.com",
		Secret: "secret", IngestPath: "/ingest/daraja/test-endpoint-2", CreatedAt: time.Now(),
	}
	endpointStore.Save(ep)
	defer pool.Exec(context.Background(), "DELETE FROM endpoints WHERE id = $1", ep.ID)

	eventStore := NewPostgresEventStore(pool)
	ev := event.NormalizedEvent{
		EventID:      "test-event-1",
		EndpointID:   ep.ID,
		EventType:    "stk_push",
		Provider:     "daraja",
		Amount:       100.0,
		Status:       "success",
		ReceivedAt:   time.Now(),
		ProviderMeta: map[string]string{"checkout_request_id": "abc123"},
		Raw:          []byte(`{"foo":"bar"}`),
	}

	if err := eventStore.Save(ev); err != nil {
		t.Fatalf("failed to save event: %v", err)
	}
	defer pool.Exec(context.Background(), "DELETE FROM events WHERE event_id = $1", ev.EventID)

	fetched, err := eventStore.GetByID(ev.EventID)
	if err != nil {
		t.Fatalf("failed to fetch event: %v", err)
	}

	if fetched.Amount != 100.0 {
		t.Errorf("expected amount 100.0, got %f", fetched.Amount)
	}

	if fetched.ProviderMeta["checkout_request_id"] != "abc123" {
		t.Errorf("expected provider_meta checkout_request_id abc123, got %s", fetched.ProviderMeta["checkout_request_id"])
	}
}

func TestPostgresEventStore_ListByEndpoint(t *testing.T) {
	pool := testPool(t)
	defer pool.Close()

	endpointStore := NewPostgresEndpointStore(pool)
	ep := endpoint.Endpoint{
		ID: "test-endpoint-3", Provider: "daraja", Shortcode: "600002",
		EventTypes: []string{"c2b_confirmation"}, DestinationURL: "https://example.com",
		Secret: "secret", IngestPath: "/ingest/daraja/test-endpoint-3", CreatedAt: time.Now(),
	}
	endpointStore.Save(ep)
	defer pool.Exec(context.Background(), "DELETE FROM endpoints WHERE id = $1", ep.ID)

	eventStore := NewPostgresEventStore(pool)
	for i := 0; i < 3; i++ {
		ev := event.NormalizedEvent{
			EventID: "list-test-event-" + string(rune('1'+i)), EndpointID: ep.ID,
			EventType: "c2b_confirmation", Provider: "daraja", Status: "success",
			ReceivedAt: time.Now(), ProviderMeta: map[string]string{}, Raw: []byte(`{}`),
		}
		eventStore.Save(ev)
		defer pool.Exec(context.Background(), "DELETE FROM events WHERE event_id = $1", ev.EventID)
	}

	results, err := eventStore.ListByEndpoint(ep.ID)
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 events, got %d", len(results))
	}
}
