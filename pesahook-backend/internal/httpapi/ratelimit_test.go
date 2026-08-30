package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountRateLimiter_AllowsWithinCapacityThenBlocks(t *testing.T) {
	limiterMW := NewAccountRateLimiter(3, 0) //capacity = 3, no refill
	calls := 0
	handler := limiterMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	// simulating already authenticared request by injectingaccountID key directly, same as RequireAPIKey would after lookup
	ctx := context.WithValue(t.Context(), accountIDKey, "account-1")

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/endpoints", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200 within capacity, got %d", i+1, w.Code)
		}
	}

	// 4th request - shoul block
	req := httptest.NewRequest("POST", "/endpoints", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 after exceeding, ot %d", w.Code)
	}

	if calls != 3 {
		t.Errorf("expected handler to be called exactly 3 times, got %d", calls)
	}
}

func TestAccountRateLimiter_SeparateAccountsHaveSeparateLimits(t *testing.T) {
	limiterMW := NewAccountRateLimiter(1, 0)

	handler := limiterMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ctxA := context.WithValue(context.Background(), accountIDKey, "account-A")
	ctxB := context.WithValue(context.Background(), accountIDKey, "account-B")

	// accountA uses its one allowed request
	reqA := httptest.NewRequest("POST", "/endpoints", nil).WithContext(ctxA)
	wA := httptest.NewRecorder()
	handler.ServeHTTP(wA, reqA)
	if wA.Code != http.StatusOK {
		t.Fatalf("expected account A`s first request to suceed, got %d", wA.Code)
	}

	// account B unaffected by As outage
	reqB := httptest.NewRequest("POST", "/endpoints", nil).WithContext(ctxB)
	wB := httptest.NewRecorder()
	handler.ServeHTTP(wB, reqB)
	if wB.Code != http.StatusOK {
		t.Fatalf("expected account B's first request to succeed independently ")
	}
}
