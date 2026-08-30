package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Its-Delimas/pesahook/internal/account"
	"github.com/Its-Delimas/pesahook/internal/apikey"
	"github.com/Its-Delimas/pesahook/internal/store"
)

type AccountHandler struct {
	Accounts store.AccountStore
	APIKeys  store.APIKeyStore
}

func NewAccountHandler(accounts store.AccountStore, apikeys store.APIKeyStore) *AccountHandler {
	return &AccountHandler{Accounts: accounts, APIKeys: apikeys}
}

type createAccountRequest struct {
	Email string `json:"email"`
}

type createAccountResponse struct {
	AccountID string `json:"account_id"`
	APIKey    string `json:"api_key"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "valid email is required", http.StatusBadRequest)
		return
	}

	acc := account.NewAccount(req.Email)
	if err := h.Accounts.Save(acc); err != nil {
		http.Error(w, "failed to save account", http.StatusInternalServerError)
		return
	}

	rawKey, key := apikey.NewAPIKey(acc.ID)
	if err := h.APIKeys.Save(key); err != nil {
		http.Error(w, "failed to save api key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createAccountResponse{
		AccountID: acc.ID,
		APIKey:    rawKey,
	})
}
