package event

import (
	"encoding/json"
	"time"
)

type NormalizedEvent struct {
	EventID          string            `json:"event_id"`
	EndpointID       string            `json:"endpoint_id"`
	EventType        string            `json:"event_type"`
	Provider         string            `json:"provider"`
	Shortcode        string            `json:"shortcode"`
	TransactionID    string            `json:"transaction_id"`
	Amount           float64           `json:"amount"`
	PhoneNumber      string            `json:"phone_number"`
	AccountReference string            `json:"account_reference"`
	Status           string            `json:"status"`
	ResultCode       int               `json:"result_code"`
	StatusReason     string            `json:"status_reason"`
	OccurredAt       time.Time         `json:"occurred_at"`
	ReceivedAt       time.Time         `json:"received_at"`
	ProviderMeta     map[string]string `json:"provider_meta"`
	Raw              json.RawMessage   `json:"raw"`
}
