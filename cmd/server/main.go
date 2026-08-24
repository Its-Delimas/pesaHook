package main

import (
	"log"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/httpapi"
	"github.com/Its-Delimas/pesaHook/internal/store"
)

func main() {
	endpointStore := store.NewMemoryEndpointStore()
	eventStore := store.NewMemoryEVentStore()

	ingestHandler := httpapi.NewIngestHandler(endpointStore, eventStore)
	endpointHandler := httpapi.NewEndpointHandler(endpointStore)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /endpoints", endpointHandler.Create)

	mux.HandleFunc("POST /ingest/{provider}/{id}", func(w http.ResponseWriter, r *http.Request) {
		endpointID := r.PathValue("id")
		ingestHandler.ServeHTTP(w, r, endpointID)
	})

	log.Println("PesaHook listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
