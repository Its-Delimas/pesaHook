package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

type EventHandler struct {
	Events    store.EventStore
	Endpoints store.EndpointStore
	Delivery  *delivery.Delivery
}

func NewEventHandler(events store.EventStore, endpoints store.EndpointStore, d *delivery.Delivery) *EventHandler {
	return &EventHandler{Events: events, Endpoints: endpoints, Delivery: d}
}

// GET /events?endpoint_id={id}
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, ok := r.Context().Value(accountIDKey).(string)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	endpointID := r.URL.Query().Get("endpoint_id")
	if endpointID == "" {
		http.Error(w, "endpoint_id query param is required", http.StatusBadRequest)
		return
	}

	ep, err := h.Endpoints.GetByID(endpointID)
	if err != nil {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}
	
	if ep.AccountID != accountID {
		http.Error(w, "endpoint not found", http.StatusNotFound) //404, not 403 so not to expose others endpoints
		return
	}

	events, err := h.Events.ListByEndpoint(endpointID)
	if err != nil {
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ev, err := h.Events.GetByID(id)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ev)
}

// POST /events/{id}/replay
func (h *EventHandler) Replay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ev, err := h.Events.GetByID(id)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get event", http.StatusInternalServerError)
		return
	}

	ep, err := h.Endpoints.GetByID(ev.EndpointID)
	if err != nil {
		http.Error(w, "endpoint for this event no longer exists", http.StatusNotFound)
		return
	}

	attempts, err := h.Delivery.Deliver(ep, ev)
	if err != nil {
		http.Error(w, "replay delivery failed after retries", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"replayed": true,
		"attempts": attempts,
	})

}
