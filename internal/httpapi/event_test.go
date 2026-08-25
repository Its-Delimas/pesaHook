package httpapi

import (
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

	ep := endpoint.NewEndpoint("daraja", "600000", []string{"stk_push"}, mockDestination.URL)
	endpointStore.Save(ep)

	ev := event.NormalizedEvent{
		EventID:    "evt123",
		EndpointID: ep.ID,
		EventType:  "stk_push",
		Status:     "success",
	}
	eventStore.Save(ev)

	handler := NewEventHandler(eventStore, endpointStore, d)

	req := httptest.NewRequest("POST", "/events/evt123/replay", nil)
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
