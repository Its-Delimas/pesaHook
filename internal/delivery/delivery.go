package delivery

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Its-Delimas/pesaHook/internal/endpoint"
	"github.com/Its-Delimas/pesaHook/internal/event"
)

type Delivery struct {
	Client *http.Client
}

func NewDelivery() *Delivery {
	return &Delivery{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

//deliver sends one event to the endpoint's destination, signed with its secret.
// retries with exponential backoff; returns error only if all attempts fail.

func (d *Delivery) Deliver(ep endpoint.Endpoint, ev event.NormalizedEvent) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	signature := sign(payload, ep.Secret)

	backoffs := []time.Duration{0, 2 * time.Second, 5 * time.Second, 15 * time.Second}

	var lastErr error
	for _, wait := range backoffs {
		time.Sleep(wait)

		req, err := http.NewRequest("POST", ep.DestinationURL, bytes.NewReader(payload))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-PesaHook-Signature", signature)

		resp, err := d.Client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		lastErr = errStatus(resp.StatusCode)
	}
	//todo: on final failure, write to dead-letter store instead of just returning error
	return lastErr
}

func sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func errStatus(code int) error {
	return &statusError{code}
}

type statusError struct{ code int }

func (e *statusError) Error() string {
	return "destination returned non-2xx status: " + http.StatusText(e.code)
}
