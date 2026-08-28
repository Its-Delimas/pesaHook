package endpoint

import "time"

type Endpoint struct {
	ID             string
	AccountID      string
	Provider       string //daraja, but maybe airtel or later
	Shortcode      string
	EventTypes     []string //stkpush,c2bconfirmation
	DestinationURL string
	Secret         string
	IngestPath     string //unique per endpoint,example: "/ingest/daraja/{id}"
	CreatedAt      time.Time
}
