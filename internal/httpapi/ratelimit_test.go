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
