// file: internal/httpapi/ingest.go
package httpapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Its-Delimas/pesaHook/internal/daraja"
)

func IngestSTKPush(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var raw daraja.STKCallbackPayload
	if err := json.Unmarshal(body, &raw); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	ev := daraja.NormalizeSTKPush(raw, body)

	// TODO: persist ev via store, then hand off to delivery worker
	_ = ev

	w.WriteHeader(http.StatusOK)
}

func IngestC2B(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var raw daraja.C2BPayload
	if err := json.Unmarshal(body, &raw); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	ev := daraja.NormalizeC2B(raw, body)

	// TODO: persist ev via store, then hand off to delivery worker
	_ = ev

	w.WriteHeader(http.StatusOK)
}
