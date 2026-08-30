package daraja

import (
	"testing"
)

func TestNormalizesSTKPush_Success(t *testing.T) {
	raw := STKCallbackPayload{}
	raw.Body.StkCallback.MerchantRequestID = "29115-34620561-1"
	raw.Body.StkCallback.CheckoutRequestID = "ws_CO_191220191020363925"
	raw.Body.StkCallback.ResultCode = 0
	raw.Body.StkCallback.ResultDesc = "the service request is processed successfully"
	raw.Body.StkCallback.CallbackMetadata = &CallbackMetadata{
		Item: []CallbackItem{
			{Name: "Amount", Value: 1.00},
			{Name: "MpesaReceiptNumber", Value: "NLJ7RT61SV"},
			{Name: "PhoneNumber", Value: 254118333997.0},
		},
	}

	ev := NormalizeSTKPush(raw, []byte(`{}`))

	if ev.Status != "success" {
		t.Errorf("expected status success, got %s", ev.Status)
	}
	if ev.Amount != 1.00 {
		t.Errorf("expected amount 1.00, got %f", ev.Amount)
	}
	if ev.TransactionID != "NLJ7RT61SV" {
		t.Errorf("expected transaction ID NLJ7RT61SV, got %s", ev.TransactionID)
	}
	if ev.PhoneNumber != "254118333997" {
		t.Errorf("expected phone 254118333997, got %s", ev.PhoneNumber)
	}
}

func TestNormalizeSTKPush_failure(t *testing.T) {
	raw := STKCallbackPayload{}
	raw.Body.StkCallback.ResultCode = 1032
	raw.Body.StkCallback.ResultDesc = "Request cancelled by the user"
	// empty callback metadata - trap case

	ev := NormalizeSTKPush(raw, []byte(`{}`))

	if ev.Status != "failed" {
		t.Errorf("expected status failed, got %s", ev.Status)
	}

	if ev.StatusReason != "Request cancelled by the user" {
		t.Errorf("expected status reason to match ResultDesc, got %s", ev.Status)
	}

	if ev.Amount != 0 {
		t.Errorf("expected amount 0 on failure, got %f", ev.Amount)
	}
}

func TestNormalizeC2B(t *testing.T) {
	raw := C2BPayload{
		TransID:           "NLJ7RT61SV",
		TransAmount:       "1.00",
		BusinessShortCode: "600000",
		BillRefNumber:     "invoice001",
		MSISDN:            "254118333997",
	}

	ev := NormalizeC2B(raw, []byte(`{}`))

	if ev.Status != "success" {
		t.Errorf("expected status success, got %s", ev.Status)
	}

	if ev.Amount != 1.00 {
		t.Errorf("expected amount 1.00, got %f", ev.Amount)
	}
	if ev.AccountReference != "invoice001" {
		t.Errorf("expected account reference invoice001, got %s", ev.AccountReference)
	}
}

func TestNormalizeC2B_MalformedAmount(t *testing.T) {
	raw := C2BPayload{
		TransID:     "NLJ7RT61SV",
		TransAmount: "not-a-number",
	}

	ev := NormalizeC2B(raw, []byte(`{}`))

	if ev.Amount != 0 {
		t.Errorf("expected amount 0 for malformed input, got %f", ev.Amount)
	}
}

func TestNormalizeB2C_Success(t *testing.T) {
	raw := B2CPayload{}
	raw.Result.ResultCode = 0
	raw.Result.ResultDesc = "The service request is processed successfully."
	raw.Result.TransactionID = "NLJ7RT61SV"
	raw.Result.OriginatorConversationID = "10571-7910404-1"
	raw.Result.ResultParameters = &B2CResultParameters{
		ResultParameter: []B2CResultParameter{
			{Key: "TransactionAmount", Value: 10.0},
			{Key: "TransactionReceipt", Value: "NLJ7RT61SV"},
		},
	}

	ev := NormalizeB2C(raw, []byte(`{}`))

	if ev.Status != "success" {
		t.Errorf("expected status success, got %s", ev.Status)
	}
	if ev.Amount != 10.0 {
		t.Errorf("expected amount 10.o, got %f", ev.Amount)
	}
	if ev.TransactionID != "NLJ7RT61SV" {
		t.Errorf("expected transaction ID NLJ7RT61SV, got %s", ev.TransactionID)
	}
	if ev.ProviderMeta["originator_conversation_id"] != "10571-7910404-1" {
		t.Errorf("expected originator conversation ID in provider_meta, got %s", ev.ProviderMeta["originator_conversation_id"])
	}
}

func TestNormalizeB2C_Failure(t *testing.T) {
	raw := B2CPayload{}
	raw.Result.ResultCode = 1
	raw.Result.ResultDesc = "Insufficient funds in the utility account"
	// result parameters nil again - trap case

	ev := NormalizeB2C(raw, []byte(`{}`))

	if ev.Status != "failed" {
		t.Errorf("expected status failed, got %s", ev.Status)
	}

	if ev.Amount != 0 {
		t.Errorf("expected amount 0 in failure, got %f", ev.Amount)
	}
}

func TestNormalizeAny_DispatchesCorrectly(t *testing.T) {
	stkPayload := []byte(`{"Body":{"stkCallback":{"ResultCode":0,"ResultDesc":"ok"}}}`)
	ev := NormalizeAny(stkPayload)
	if ev.EventType != "stk_push" {
		t.Errorf("expected stk_push, got %s", ev.EventType)
	}

	// we were checking the wrong value and and wrong display message
	c2bPayload := []byte(`{"TransID":"ABC123","TransAmount":"5.00"}`)
	ev = NormalizeAny(c2bPayload)
	if ev.EventType != "c2b_confirmation" {
		t.Errorf("expected c2b_confirmation,got %s", ev.EventType)
	}

	b2cpayload := []byte(`{"Result":{"ResultCode":0}}`)
	ev = NormalizeAny(b2cpayload)
	if ev.EventType != "b2c_result" {
		t.Errorf("Expected b2c_result,got %s", ev.EventType)
	}

	somethingelse := []byte(`{"la":"peace"}`)
	ev = NormalizeAny(somethingelse)
	if ev.Status != "unrecognized" {
		t.Errorf("expected unrecognized, got %s", ev.Status)
	}

}
