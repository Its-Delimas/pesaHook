package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Its-Delimas/pesahook/internal/store"
)

type DeadLetterHandler struct {
	Endpoints   store.EndpointStore
	DeadLetters store.DeadLetterStore
}

func NewDeadLetterHandler(endpoints store.EndpointStore, deadLetters store.DeadLetterStore) *DeadLetterHandler {
	return &DeadLetterHandler{Endpoints: endpoints, DeadLetters: deadLetters}
}

// GET /endpoints/{id}/dead-letters
func (h *DeadLetterHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, ok := r.Context().Value(accountIDKey).(string)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	endpointID := r.PathValue("id")

	ep, err := h.Endpoints.GetByID(endpointID)
	if err != nil || ep.AccountID != accountID {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	deadLetters, err := h.DeadLetters.ListByEndpoint(endpointID)
	if err != nil {
		http.Error(w, "failed to list dead letters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deadLetters)
}
