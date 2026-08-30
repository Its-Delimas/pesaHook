package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/endpoint"
	"github.com/Its-Delimas/pesaHook/internal/event"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

func TestEventHandler_Replay_Success(t *testing.T) {
	received := make(chan struct{}, 1)
	mockDestination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockDestination.Close()

	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEventStore()
	d := delivery.NewDelivery()

	ep := endpoint.NewEndpoint("test-account-1", "daraja", "600000", []string{"stk_push"}, mockDestination.URL)
	endpointStore.Save(ep)

	ev := event.NormalizedEvent{
		EventID:    "evt123",
		EndpointID: ep.ID,
		EventType:  "stk_push",
		Status:     "success",
	}
	eventStore.Save(ev)

	handler := NewEventHandler(eventStore, endpointStore, d)

	ctx := context.WithValue(context.Background(), accountIDKey, "test-account-1")
	req := httptest.NewRequest("POST", "/events/evt123/replay", nil).WithContext(ctx)
	req.SetPathValue("id", "evt123")
	w := httptest.NewRecorder()

	handler.Replay(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	select {
	case <-received:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected replay to hit destination, timed out")
	}
}
func TestEventHandler_Replay_EventNotFound(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEventStore()
	d := delivery.NewDelivery()

	handler := NewEventHandler(eventStore, endpointStore, d)

	ctx := context.WithValue(context.Background(), accountIDKey, "test-account-1")
	req := httptest.NewRequest("POST", "/events/nonexistent/replay", nil).WithContext(ctx)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	handler.Replay(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
