package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Its-Delimas/pesahook/internal/apikey"
	"github.com/Its-Delimas/pesahook/internal/store"
)

type APIKeyHandler struct {
	APIKeys store.APIKeyStore
}

func NewAPIKeyHandler(apiKeys store.APIKeyStore) *APIKeyHandler {
	return &APIKeyHandler{APIKeys: apiKeys}
}

type apikeyResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

// POST /api-keys, creates new key under the aythenticated account, returns raw key once
func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountId, ok := r.Context().Value(accountIDKey).(string)
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	rawKey, key := apikey.NewAPIKey(accountId)
	if err := h.APIKeys.Save(key); err != nil {
		http.Error(w, "failed to save api key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":         key.ID,
		"api_key":    rawKey,
		"created_at": key.CreatedAt.Format(time.RFC3339),
	})

}
