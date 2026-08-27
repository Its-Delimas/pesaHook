package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/endpoint"
	"github.com/Its-Delimas/pesaHook/internal/event"
	"github.com/Its-Delimas/pesaHook/pkg/verify"
)

type Delivery struct {
	Client   *http.Client
	Backoffs []time.Duration
}

func NewDelivery() *Delivery {
	return &Delivery{
		Client:   &http.Client{Timeout: 10 * time.Second},
		Backoffs: []time.Duration{0, 2 * time.Second, 5 * time.Second, 15 * time.Second},
	}
}

//deliver sends one event to the endpoint's destination, signed with its secret.
// retries with exponential backoff; returns error only if all attempts fail.

func (d *Delivery) Deliver(ep endpoint.Endpoint, ev event.NormalizedEvent) (attempts int, err error) {
	payload, err := json.Marshal(ev)
	if err != nil {
		return 0, err
	}

	signature := verify.Signature(payload, ep.Secret)

	var lastErr error
	for i, wait := range d.Backoffs {
		time.Sleep(wait)
		attempts = i + 1

		req, reqErr := http.NewRequest("POST", ep.DestinationURL, bytes.NewReader(payload))
		if reqErr != nil {
			lastErr = reqErr
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-PesaHook-Signature", signature)

		resp, doErr := d.Client.Do(req)
		if doErr != nil {
			lastErr = doErr
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return attempts, nil
		}
		lastErr = errStatus(resp.StatusCode)
	}
	return attempts, lastErr
}

func errStatus(code int) error {
	return &statusError{code}
}

type statusError struct{ code int }

func (e *statusError) Error() string {
	return "destination returned non-2xx status: " + http.StatusText(e.code)
}
