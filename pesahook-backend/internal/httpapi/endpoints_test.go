package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
	"github.com/Its-Delimas/pesahook/internal/store"
)

func TestEndpointHandler_Create(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	handler := NewEndpointHandler(endpointStore)

	reqBody := createEndpointRequest{
		Provider:       "daraja",
		Shortcode:      "600000",
		EventTypes:     []string{"stk_push", "c2b_confirmation"},
		DestinationURL: "https://example.com/webhook",
	}

	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), accountIDKey, "test-account-1")
	req := httptest.NewRequest("POST", "/endpoints", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp createEndpointResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}

	if resp.Secret == "" {
		t.Error("expected non-empty secret")
	}

	saved, err := endpointStore.GetByID(resp.ID)
	if err != nil {
		t.Fatalf("expected endpoint to be saved, got error: %v", err)
	}
	if saved.DestinationURL != "https://example.com/webhook" {
		t.Errorf("expected destination URL to match, got %s", saved.DestinationURL)
	}
}

func TestEndpointHandler_Create_MissingFields(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	handler := NewEndpointHandler(endpointStore)

	reqBody := createEndpointRequest{Provider: "daraja"} //missing shortcode
	body, _ := json.Marshal(reqBody)

	ctx := context.WithValue(context.Background(), accountIDKey, "test-account-1")
	req := httptest.NewRequest("POST", "/endpoints", bytes.NewReader(body)).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// file: internal/httpapi/endpoints_test.go — add these tests
func TestEndpointHandler_List_ReturnsOnlyAccountsOwnEndpoints(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	handler := NewEndpointHandler(endpointStore)

	epA1 := endpoint.NewEndpoint("account-A", "daraja", "600000", []string{"stk_push"}, "https://a1.example.com")
	epA2 := endpoint.NewEndpoint("account-A", "daraja", "600001", []string{"c2b_confirmation"}, "https://a2.example.com")
	epB1 := endpoint.NewEndpoint("account-B", "daraja", "600002", []string{"stk_push"}, "https://b1.example.com")

	endpointStore.Save(epA1)
	endpointStore.Save(epA2)
	endpointStore.Save(epB1)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")
	req := httptest.NewRequest("GET", "/endpoints", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var results []endpointResponse
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 endpoints for account-A, got %d", len(results))
	}

	for _, r := range results {
		if r.Shortcode != "600000" && r.Shortcode != "600001" {
			t.Errorf("unexpected shortcode in results: %s (belongs to account-B, should not be returned)", r.Shortcode)
		}
	}
}

func TestEndpointHandler_List_OmitsSecret(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	handler := NewEndpointHandler(endpointStore)

	ep := endpoint.NewEndpoint("account-A", "daraja", "600000", []string{"stk_push"}, "https://example.com")
	endpointStore.Save(ep)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")
	req := httptest.NewRequest("GET", "/endpoints", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.List(w, req)

	body := w.Body.String()
	if strings.Contains(body, ep.Secret) {
		t.Error("response body contains the raw secret — endpointResponse must omit it")
	}
}

func TestEndpointHandler_List_Unauthenticated(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	handler := NewEndpointHandler(endpointStore)

	req := httptest.NewRequest("GET", "/endpoints", nil) // no accountIDKey in context
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
