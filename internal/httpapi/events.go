package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/store"
)

type EventHandler struct {
	Events store.EventStore
}

func NewEventHandler(events store.EventStore) *EventHandler {
	return &EventHandler{Events: events}
}

// GET /events?endpoint_id={id}
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	endpointID := r.URL.Query().Get("endpoint_id")
	if endpointID == "" {
		http.Error(w, "endpoint_id query param is required", http.StatusBadRequest)
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
