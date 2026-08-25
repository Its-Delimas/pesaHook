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
	eventStore := store.NewMemoryEVentStore()
	deadLetterScore := store.NewMemoryDeadLetterStore()
	d := delivery.NewDelivery()

	ep := endpoint.Endpoint("daraja", "600000", []string{"stk_push"}, mockDestination.URL)
	endpointStore.Save(ep)

	handler := NewIngestHandler(endpointStore, eventStore, d, deadLetterScore)

	payload := []byte(`{"Body":{"stkCallback":{"ResultCode":0,"ResultDesc":"ok","CallbackMetadata":{"Item":[{"Name":"Amount","Value":1.0},{"Name":"MpesaReceiptNumber","Value":"NLJ7RT61SV"},{"Name":"PhoneNumber","Value":254118333997}]}}}}`)

	req := httptest.NewRequest("POST", "/ingest/daraja/"+ep.ID, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req, ep.ID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	select {
	case <-received:
		//delivery happened - test passed
	case <-time.After(2 * time.Second):
		t.Fatalf("expected event to be delivered to destination, timed out waiting")
	}

	events, err := eventStore.ListByEndpoint(ep.ID)
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 saved events, got %d", len(events))
	}
	if events[0].TransactionID != "NLJ7RT61SV" {
		t.Errorf("expected transaction ID NLJ7RT61SV, got %s", events[0].TransactionID)
	}

}
