package account

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Account struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func NewAccount(email string) Account {
	return Account{
		ID:        generateID(),
		Email:     email,
		CreatedAt: time.Now(),
	}
}
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
