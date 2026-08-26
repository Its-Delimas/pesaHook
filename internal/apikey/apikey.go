package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type APIKey struct {
	ID        string
	AccountID string
	KeyHash   string
	CreatedAt time.Time
}

// generate raw key - show once to te user and store
func NewAPIKey(accountID string) (rawKey string, ak APIKey) {
	raw := generateRawKey()
	return raw, APIKey{
		ID:        generateID(),
		AccountID: accountID,
		KeyHash:   hashKey(raw),
		CreatedAt: time.Now(),
	}
}

func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func generateRawKey() string {
	b := make([]byte, 24)
	rand.Read(b)
	return "ph_" + hex.EncodeToString(b)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
