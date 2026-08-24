package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/daraja"
	"github.com/Its-Delimas/pesaHook/internal/event"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

type IngestHandler struct {
	Endpoints store.EndpointStore
	Events    store.EventStore
}

func NewIngestHandler(endpoints store.EndpointStore, events store.EventStore) *IngestHandler {
	return &IngestHandler{Endpoints: endpoints, Events: events}
}

// ServeHTTP handles POST /ingest/{provider}/{endpointID}
func (h *IngestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, endpointID string) {
	ep, err := h.Endpoints.GetByID(endpointID)
	if err != nil {
		http.Error(w, "unknown endpoint", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var ev event.NormalizedEvent

	switch ep.Provider {
	case "daraja":
		ev = daraja.NormalizeAny(body)
	default:
		http.Error(w, "unsupported provider", http.StatusBadRequest)
		return
	}

	ev.EventID = generateEventID()
	ev.EndpointID = ep.ID
	ev.ReceivedAt = time.Now()

	if err := h.Events.Save(ev); err != nil {
		http.Error(w, "failed to save event", http.StatusInternalServerError)
		return
	}

	// todo: hand off to delivery worker to forward to ep.DestinationURL

	w.WriteHeader(http.StatusOK)
}

func generateEventID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
