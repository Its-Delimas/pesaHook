package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Its-Delimas/pesaHook/internal/store"
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

	req := httptest.NewRequest("POST", "/endpoints", bytes.NewReader(body))
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

	req := httptest.NewRequest("POST", "/endpoints", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
