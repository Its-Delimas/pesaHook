package endpoint

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func NewEndpoint(provider, shortcode string, eventTypes []string, destinationURL string) Endpoint {
	id := generateID()
	return Endpoint{
		ID:             id,
		Provider:       provider,
		Shortcode:      shortcode,
		EventTypes:     eventTypes,
		DestinationURL: destinationURL,
		Secret:         generateSecret(),
		IngestPath:     "/ingest/" + provider + "/" + id,
		CreatedAt:      time.Now(),
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSecret() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
