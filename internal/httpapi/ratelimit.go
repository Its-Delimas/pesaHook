package httpapi

import (
	"net/http"

	limiter "github.com/Its-Delimas/rate_limiter"
	"github.com/Its-Delimas/rate_limiter/middleware"
)

// NewAccountRateLimiter builds rate-limiting middleware keyed by the
// authenticated account ID (set by RequireAPIKey), not by IP —
// so shared IPs (NAT, proxies) don't get unfairly throttled together,
// and each PesaHook account gets its own independent limit.

func NewAccountRateLimiter(capacity int, refillPerSec float64) func(http.Handler) http.Handler {
	l := limiter.NewTokenBucketLimiter(float64(capacity), refillPerSec)

	return middleware.RateLimit(l, func(r *http.Request) string {
		accountID, ok := r.Context().Value(accountIDKey).(string)
		if !ok {
			return "unauthenticated"
		}
		return accountID
	})
}
