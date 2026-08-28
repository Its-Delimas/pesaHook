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

func TestSignature_DifferebtSecretsProduceDifferentSignatures(t *testing.T) {
	payload := []byte(`{"event":"stk_push"}`)

	sig1 := Signature(payload, "secret1")
	sig2 := Signature(payload, "secret2")

	if sig2 == sig1 {
		t.Errorf("expected different secrets to produce different signatures")
	}
}

