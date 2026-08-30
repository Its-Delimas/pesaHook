package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/Its-Delimas/pesaHook/internal/apikey"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

type contextkey string

const accountIDKey contextkey = "account_id"

func RequireAPIKey(keys store.APIKeyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			rawKey := strings.TrimPrefix(authHeader, "Bearer ")
			hash := apikey.HashKey(rawKey)

			key, err := keys.GetByHash(hash)
			if err != nil {
				http.Error(w, "invalid API key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), accountIDKey, key.AccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
