package daraja

import "testing"

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
		t.Errorf("expected phone 254708374149, got %s", ev.PhoneNumber)
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
