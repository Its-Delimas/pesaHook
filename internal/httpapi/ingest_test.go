// file: internal/httpapi/ingest_test.go
package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/endpoint"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

func TestIngestHandler_STKPush_DeliversToDestination(t *testing.T) {
	received := make(chan []byte, 1)
	mockDestination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		received <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer mockDestination.Close()

	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEventStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	d := delivery.NewDelivery()

	ep := endpoint.NewEndpoint("daraja", "600000", []string{"stk_push"}, mockDestination.URL)
	endpointStore.Save(ep)

	handler := NewIngestHandler(endpointStore, eventStore, d, deadLetterStore)

	payload := []byte(`{"Body":{"stkCallback":{"ResultCode":0,"ResultDesc":"ok","CallbackMetadata":{"Item":[{"Name":"Amount","Value":1.0},{"Name":"MpesaReceiptNumber","Value":"NLJ7RT61SV"},{"Name":"PhoneNumber","Value":254708374149}]}}}}`)

	req := httptest.NewRequest("POST", "/ingest/daraja/"+ep.ID, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req, ep.ID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	select {
	case <-received:
		// delivery happened, test passes
	case <-time.After(2 * time.Second):
		t.Fatal("expected event to be delivered to destination, timed out waiting")
	}

	events, err := eventStore.ListByEndpoint(ep.ID)
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 saved event, got %d", len(events))
	}
	if events[0].TransactionID != "NLJ7RT61SV" {
		t.Errorf("expected transaction ID NLJ7RT61SV, got %s", events[0].TransactionID)
	}
}

func TestIngestHandler_UnknownEndpoint(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEventStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	d := delivery.NewDelivery()

	handler := NewIngestHandler(endpointStore, eventStore, d, deadLetterStore)

	req := httptest.NewRequest("POST", "/ingest/daraja/nonexistent", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestIngestHandler_FailedDelivery_GoesToDeadLetter(t *testing.T) {
	failingDestination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failingDestination.Close()

	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEventStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	d := delivery.NewDelivery()
	d.Backoffs = []time.Duration{0, 0, 0, 0}

	ep := endpoint.NewEndpoint("daraja", "600000", []string{"stk_push"}, failingDestination.URL)
	endpointStore.Save(ep)

	handler := NewIngestHandler(endpointStore, eventStore, d, deadLetterStore)

	payLoad := []byte(`{"Body":{"stkCallback":{"ResultCode":0,"ResultDesc":"ok"}}}`)
	req := httptest.NewRequest("POST", "/ingest/daraja/"+ep.ID, bytes.NewReader(payLoad))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req, ep.ID)

	// backoff schedule is 0,2,5 and 15s...wait past all 4 attempts bedore proceeding
	time.Sleep(100 * time.Millisecond)

	deadLetters, err := deadLetterStore.ListByEndpoint(ep.ID)
	if err != nil {
		t.Fatalf("failed to list dead letters: %v", err)
	}
	if len(deadLetters) != 1 {
		t.Fatalf("expected 1 dead letter, got %d", len(deadLetters))
	}
	if deadLetters[0].Attempts != 4 {
		t.Errorf("expected 4 attempts, got %d", deadLetters[0].Attempts)
	}
}
