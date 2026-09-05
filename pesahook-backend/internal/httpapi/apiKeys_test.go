// file: internal/httpapi/apikeys_test.go
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Its-Delimas/pesahook/internal/store"
)

func TestAPIKeyHandler_Create_ReturnsRawKeyOnce(t *testing.T) {
	apiKeyStore := store.NewMemoryAPIKeyStore()
	handler := NewAPIKeyHandler(apiKeyStore)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")
	req := httptest.NewRequest("POST", "/api-keys", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["api_key"] == "" {
		t.Error("expected a raw api_key in the creation response")
	}
	if !strings.HasPrefix(resp["api_key"], "ph_") {
		t.Errorf("expected key prefixed with ph_, got %s", resp["api_key"])
	}
}

func TestAPIKeyHandler_List_NeverExposesRawKeyOrHash(t *testing.T) {
	apiKeyStore := store.NewMemoryAPIKeyStore()
	createHandler := NewAPIKeyHandler(apiKeyStore)
	listHandler := NewAPIKeyHandler(apiKeyStore)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")

	createReq := httptest.NewRequest("POST", "/api-keys", nil).WithContext(ctx)
	createW := httptest.NewRecorder()
	createHandler.Create(createW, createReq)

	var createResp map[string]string
	json.NewDecoder(createW.Body).Decode(&createResp)
	rawKey := createResp["api_key"]

	listReq := httptest.NewRequest("GET", "/api-keys", nil).WithContext(ctx)
	listW := httptest.NewRecorder()
	listHandler.List(listW, listReq)

	body := listW.Body.String()
	if strings.Contains(body, rawKey) {
		t.Error("List response contains the raw API key — must never be exposed after creation")
	}
	if strings.Contains(body, "key_hash") {
		t.Error("List response exposes key_hash field — should only return id and created_at")
	}
}

func TestAPIKeyHandler_List_ReturnsOnlyAccountsOwnKeys(t *testing.T) {
	apiKeyStore := store.NewMemoryAPIKeyStore()
	handler := NewAPIKeyHandler(apiKeyStore)

	ctxA := context.WithValue(context.Background(), accountIDKey, "account-A")
	ctxB := context.WithValue(context.Background(), accountIDKey, "account-B")

	handler.Create(httptest.NewRecorder(), httptest.NewRequest("POST", "/api-keys", nil).WithContext(ctxA))
	handler.Create(httptest.NewRecorder(), httptest.NewRequest("POST", "/api-keys", nil).WithContext(ctxA))
	handler.Create(httptest.NewRecorder(), httptest.NewRequest("POST", "/api-keys", nil).WithContext(ctxB))

	req := httptest.NewRequest("GET", "/api-keys", nil).WithContext(ctxA)
	w := httptest.NewRecorder()
	handler.List(w, req)

	var results []apikeyResponse
	json.NewDecoder(w.Body).Decode(&results)

	if len(results) != 2 {
		t.Fatalf("expected 2 keys for account-A, got %d", len(results))
	}
}
