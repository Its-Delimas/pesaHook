package verify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// signature computes HMAC-sha256 signature for payload,using same alp peashook uses when signing outgoing deliveries

func Signature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// verify checks match of received signature and expected
func Verify(payload []byte, secret string, receivedSignature string) bool {
	expected := Signature(payload, secret)
	return hmac.Equal([]byte(expected), []byte(receivedSignature))
}
