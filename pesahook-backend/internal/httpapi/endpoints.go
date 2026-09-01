package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
	"github.com/Its-Delimas/pesahook/internal/store"
)

type EndpointHandler struct {
	Endpoints store.EndpointStore
}

func NewEndpointHandler(endpoints store.EndpointStore) *EndpointHandler {
	return &EndpointHandler{Endpoints: endpoints}
}

type createEndpointRequest struct {
	Provider       string   `json:"provider"`
	Shortcode      string   `json:"shortcode"`
	EventTypes     []string `json:"event_types"`
	DestinationURL string   `json:"destination_url"`
}

type createEndpointResponse struct {
	ID        string `json:"id"`
	IngestURL string `json:"ingest_url"`
	Secret    string `json:"secret"`
}

func (h *EndpointHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountID, ok := r.Context().Value(accountIDKey).(string)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	var req createEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Provider == "" || req.Shortcode == "" || req.DestinationURL == "" {
		http.Error(w, "provider, shortcode, and destination_url are required", http.StatusBadRequest)
		return
	}

	ep := endpoint.NewEndpoint(accountID, req.Provider, req.Shortcode, req.EventTypes, req.DestinationURL)

	if err := h.Endpoints.Save(ep); err != nil {
		http.Error(w, "failed to save endpoint", http.StatusInternalServerError)
		return
	}

	resp := createEndpointResponse{
		ID:        ep.ID,
		IngestURL: ep.IngestPath,
		Secret:    ep.Secret,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *EndpointHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, ok := r.Context().Value(accountIDKey).(string)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	endpoints, err := h.Endpoints.ListByAccount(accountID)
	if err != nil {
		http.Error(w, "failed to list endpoints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(endpoints)
}
