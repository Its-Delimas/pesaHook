package daraja

import (
	"fmt"
	"strconv"

	"github.com/Its-Delimas/pesaHook/internal/event"
)

func NormalizeSTKPush(raw STKCallbackPayload, rawBytes []byte) event.NormalizedEvent {
	cb := raw.Body.StkCallback

	ev := event.NormalizedEvent{
		EventType:    "stk_push",
		Provider:     "daraja",
		ResultCode:   cb.ResultCode,
		StatusReason: cb.ResultDesc,
		ProviderMeta: map[string]string{
			"checkout_request_id": cb.CheckoutRequestID,
			"merchant_request_id": cb.MerchantRequestID,
		},
	}

	if cb.ResultCode == 0 && cb.CallbackMetadata != nil {
		ev.Status = "success"

		fields := make(map[string]interface{})
		for _, item := range cb.CallbackMetadata.Item {
			fields[item.Name] = item.Value
		}

		if amount, ok := fields["Amount"].(float64); ok {
			ev.Amount = amount
		}

		if receipt, ok := fields["MpesaReceiptNumber"].(string); ok {
			ev.TransactionID = receipt
		}

		if phone, ok := fields["PhoneNumber"].(float64); ok {
			ev.PhoneNumber = fmt.Sprintf("%.0f", phone)
		}
	} else {
		ev.Status = "failed"
	}
	ev.Raw = rawBytes
	return ev
}

func NormalizeC2B(raw C2BPayload, rawBytes []byte) event.NormalizedEvent {
	amount, _ := strconv.ParseFloat(raw.TransAmount, 64)

	ev := event.NormalizedEvent{
		EventType:     "c2b_confirmation",
		Provider:      "daraja",
		Shortcode:     raw.BusinessShortCode,
		TransactionID: raw.TransID,
		Amount:        amount,
		PhoneNumber:   raw.MSISDN,
		Status:        "success",
		ProviderMeta:  map[string]string{},
	}
	ev.Raw = rawBytes
	return ev
}
