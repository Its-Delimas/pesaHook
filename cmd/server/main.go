package server

import (
	"log"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/httpapi"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ingest/stk-push", httpapi.IngestSTKPush)
	mux.HandleFunc("/ingest/c2b", httpapi.IngestC2B)

	log.Println("pesahook listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
