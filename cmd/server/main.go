package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/httpapi"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

func main() {
	ctx := context.Background()
	pool, err := store.NewPostgresPool(ctx, "postgres://pesahook:devpass@localhost:5432/pesahook")
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	endpointStore := store.NewPostgresEndpointStore(pool)
	eventStore := store.NewPostgresEventStore(pool)

	d := delivery.NewDelivery()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	ingestHandler := httpapi.NewIngestHandler(endpointStore, eventStore, d, deadLetterStore)
	endpointHandler := httpapi.NewEndpointHandler(endpointStore)

	eventHandler := httpapi.NewEventHandler(eventStore, endpointStore, d)

	accountStore := store.NewMemoryAccountStore()
	apiKeyStore := store.NewMemoryAPIKeyStore()
	accountHandler := httpapi.NewAccountHandler(accountStore, apiKeyStore)

	rateLimiter := httpapi.NewAccountRateLimiter(20, 1)

	mux := http.NewServeMux()
	mux.Handle("POST /endpoints", protect(apiKeyStore, rateLimiter, endpointHandler.Create))
	mux.Handle("GET /events", protect(apiKeyStore, rateLimiter, eventHandler.List))
	mux.Handle("GET /events/{id}", protect(apiKeyStore, rateLimiter, eventHandler.Get))
	mux.Handle("POST /events/{id}/replay", protect(apiKeyStore, rateLimiter, eventHandler.Replay))

	// bootsrap route - must say open nio someone getas their first API key
	mux.HandleFunc("POST /accounts", accountHandler.Create)

	// stays unprotected, Daraja calls this directly
	mux.HandleFunc("POST /ingest/{provider}/{id}", func(w http.ResponseWriter, r *http.Request) {
		endpointID := r.PathValue("id")
		ingestHandler.ServeHTTP(w, r, endpointID)
	})

	log.Println("PesaHook listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// helper
func protect(apiKeys store.APIKeyStore, rateLimiter func(http.Handler) http.Handler, h http.HandlerFunc) http.Handler {
	return httpapi.RequireAPIKey(apiKeys)(rateLimiter(h))
}

// *todo: ownership check on endpoints before returning data
