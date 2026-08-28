package verify

import "testing"

func TestSignature_IsConsistent(t *testing.T) {
	payload := []byte(`{"event":"stk_push"}`)
	secret := "testsecret"

	sig1 := Signature(payload, secret)
	sig2 := Signature(payload, secret)

	if sig1 != sig2 {
		t.Errorf("expected identical payload + scret to be similar signatures, got %s and %s", sig1, sig2)
	}
}
