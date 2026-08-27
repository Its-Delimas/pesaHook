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

	mux := http.NewServeMux()

	accountStore := store.NewMemoryAccountStore()
	apiKeyStore := store.NewMemoryAPIKeyStore()
	accountHandler := httpapi.NewAccountHandler(accountStore, apiKeyStore)

	mux.HandleFunc("POST /accounts", accountHandler.Create)
	mux.Handle("POST /endpoints", httpapi.RequireAPIKey(apiKeyStore)(http.HandlerFunc(endpointHandler.Create)))

	mux.Handle("POST /endpoints", httpapi.RequireAPIKey(apiKeyStore)(http.HandlerFunc(endpointHandler.Create)))

	mux.HandleFunc("POST /ingest/{provider}/{id}", func(w http.ResponseWriter, r *http.Request) {
		endpointID := r.PathValue("id")
		ingestHandler.ServeHTTP(w, r, endpointID)
	})

	mux.HandleFunc("GET /events", eventHandler.List)
	mux.HandleFunc("GET /events/{id}", eventHandler.Get)
	mux.HandleFunc("POST /events/{id}/replay", eventHandler.Replay)

	log.Println("PesaHook listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
