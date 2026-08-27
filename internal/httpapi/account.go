package httpapi

import "github.com/Its-Delimas/pesaHook/internal/store"

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
