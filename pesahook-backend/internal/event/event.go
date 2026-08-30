// file: internal/event/event.go
package event

import (
	"encoding/json"
	"time"
)

type NormalizedEvent struct {
	EventID          string
	EndpointID       string
	EventType        string
	Provider         string
	Shortcode        string
	TransactionID    string
	Amount           float64
	PhoneNumber      string
	AccountReference string
	Status           string
	ResultCode       int
	StatusReason     string
	OccurredAt       time.Time
	ReceivedAt       time.Time
	ProviderMeta     map[string]string
	Raw              json.RawMessage
}
