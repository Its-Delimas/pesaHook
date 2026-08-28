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

func TestVerify_ValidSignature(t *testing.T) {
	payload := []byte(`{"event":"stk_push"}`)
	secret := "testsecret"

	sig := Signature(payload, secret)

	if !Verify(payload, secret, sig) {
		t.Error("expected valid signature to verify successfully")
	}
}

func TestVerify_TamperedPayloadFails(t *testing.T) {
	payload := []byte(`{"event":"stk_push"}`)
	secret := "testsecret"

	sig := Signature(payload, secret)
	tamperedPayload := []byte(`{"event":"stk_push","amount":999999}`)

	if Verify(tamperedPayload, secret, sig) {
		t.Error("expected verification to fail when payload is tampered with after signing")
	}
}

func TestVerify_WrongSecretFails(t *testing.T) {
	payload := []byte(`{"event":"stk_push"}`)

	sig := Signature(payload, "correctsecret")

	if Verify(payload, "wrongsecret", sig) {
		t.Error("expected verification to fail with wrong secret")
	}
}
