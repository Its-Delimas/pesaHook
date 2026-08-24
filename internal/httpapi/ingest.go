package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/daraja"
	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/event"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

type IngestHandler struct {
	Endpoints  store.EndpointStore
	Events     store.EventStore
	Delivery   *delivery.Delivery
	DeadLetter store.DeadLetterStore
}

func NewIngestHandler(endpoints store.EndpointStore, events store.EventStore, d *delivery.Delivery, dl store.DeadLetterStore) *IngestHandler {
	return &IngestHandler{Endpoints: endpoints, Events: events, Delivery: d, DeadLetter: dl}
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

	go func() {
		attempts, err := h.Delivery.Deliver(ep, ev)
		if err != nil {
			h.DeadLetter.Save(store.DeadLetter{
				Event:      ev,
				EndpointID: ep.ID,
				LastError:  err.Error(),
				FailedAt:   time.Now(),
				Attempts:   attempts,
			})
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func generateEventID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
