package main

import (
	"log"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/delivery"
	"github.com/Its-Delimas/pesaHook/internal/httpapi"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

func main() {
	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEVentStore()

	d := delivery.NewDelivery()
	deadLetterStore := store.NewMemoryDeadLetterStore()
	ingestHandler := httpapi.NewIngestHandler(endpointStore, eventStore, d, deadLetterStore)
	endpointHandler := httpapi.NewEndpointHandler(endpointStore)

	eventHandler := httpapi.NewEventHandler(eventStore, endpointStore, d)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /endpoints", endpointHandler.Create)

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
