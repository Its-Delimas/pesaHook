// file: internal/httpapi/deadletters_test.go
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Its-Delimas/pesahook/internal/endpoint"

	"github.com/Its-Delimas/pesahook/internal/event"
	"github.com/Its-Delimas/pesahook/internal/store"
)

func TestDeadLetterHandler_List_ReturnsEndpointsDeadLetters(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	handler := NewDeadLetterHandler(endpointStore, deadLetterStore)

	ep := endpoint.NewEndpoint("account-A", "daraja", "600000", []string{"stk_push"}, "https://example.com")
	endpointStore.Save(ep)

	dl := store.DeadLetter{
		Event: event.NormalizedEvent{
			EventID:    "evt-dl-1",
			EndpointID: ep.ID,
			Status:     "success",
		},
		EndpointID: ep.ID,
		LastError:  "destination returned non-2xx status: Internal Server Error",
		FailedAt:   time.Now(),
		Attempts:   4,
	}
	deadLetterStore.Save(dl)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")
	req := httptest.NewRequest("GET", "/endpoints/"+ep.ID+"/dead-letters", nil).WithContext(ctx)
	req.SetPathValue("id", ep.ID)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var results []store.DeadLetter
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 dead letter, got %d", len(results))
	}
	if results[0].Attempts != 4 {
		t.Errorf("expected 4 attempts, got %d", results[0].Attempts)
	}
}

func TestDeadLetterHandler_List_WrongAccountReturns404(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	handler := NewDeadLetterHandler(endpointStore, deadLetterStore)

	// endpoint belongs to account-A
	ep := endpoint.NewEndpoint("account-A", "daraja", "600000", []string{"stk_push"}, "https://example.com")
	endpointStore.Save(ep)

	// account-B tries to read account-A's dead letters
	ctx := context.WithValue(context.Background(), accountIDKey, "account-B")
	req := httptest.NewRequest("GET", "/endpoints/"+ep.ID+"/dead-letters", nil).WithContext(ctx)
	req.SetPathValue("id", ep.ID)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for mismatched account, got %d", w.Code)
	}
}

func TestDeadLetterHandler_List_NonexistentEndpointReturns404(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	handler := NewDeadLetterHandler(endpointStore, deadLetterStore)

	ctx := context.WithValue(context.Background(), accountIDKey, "account-A")
	req := httptest.NewRequest("GET", "/endpoints/nonexistent/dead-letters", nil).WithContext(ctx)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for nonexistent endpoint, got %d", w.Code)
	}
}

func TestDeadLetterHandler_List_Unauthenticated(t *testing.T) {
	endpointStore := store.NewMemoryEndpointStore()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	handler := NewDeadLetterHandler(endpointStore, deadLetterStore)

	req := httptest.NewRequest("GET", "/endpoints/some-id/dead-letters", nil) // no accountIDKey
	req.SetPathValue("id", "some-id")
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
