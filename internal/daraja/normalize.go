package daraja

import (
	"encoding/json"
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
		EventType:        "c2b_confirmation",
		Provider:         "daraja",
		Shortcode:        raw.BusinessShortCode,
		TransactionID:    raw.TransID,
		Amount:           amount,
		PhoneNumber:      raw.MSISDN,
		AccountReference: raw.BillRefNumber, // fixed mapping on test billrefNumber to Account reference

		Status:       "success",
		ProviderMeta: map[string]string{},
	}

	ev.Raw = rawBytes
	return ev
}

func NormalizeAny(rawBytes []byte) event.NormalizedEvent {

	var probe struct {
		Body struct {
			StkCallback json.RawMessage `json:"stkCallback"`
		} `json:"Body"`
		Result  json.RawMessage `json:"Result"`  //b2c
		TransID string          `json:"TransID"` //c2b
	}

	//enebla consumer to know if its invalid json or unknown daraja event
	if err := json.Unmarshal(rawBytes, &probe); err != nil {
		return event.NormalizedEvent{
			Status: "invalid",
			Raw:    rawBytes,
		}
	}

	if probe.Body.StkCallback != nil {
		var raw STKCallbackPayload
		json.Unmarshal(rawBytes, &raw)
		return NormalizeSTKPush(raw, rawBytes)
	}

	if probe.Result != nil {
		var raw B2CPayload
		json.Unmarshal(rawBytes, &raw)
		return NormalizeB2C(raw, rawBytes)
	}

	if probe.TransID != "" {
		var raw C2BPayload
		json.Unmarshal(rawBytes, &raw)
		return NormalizeC2B(raw, rawBytes)
	}

	return event.NormalizedEvent{
		Status: "unrecognized", Raw: rawBytes}
}

func NormalizeB2C(raw B2CPayload, rawBytes []byte) event.NormalizedEvent {
	res := raw.Result

	ev := event.NormalizedEvent{
		EventType:     "b2c_result",
		Provider:      "daraja",
		TransactionID: res.TransactionID,
		ResultCode:    res.ResultCode,
		StatusReason:  res.ResultDesc,
		ProviderMeta: map[string]string{
			"originator_conversation_id": res.OriginatorConversationID,
			"conversation_id":            res.ConversationID,
		},
		Raw: rawBytes,
	}
	if res.ResultCode == 0 && res.ResultParameters != nil {
		ev.Status = "success"

		fields := make(map[string]interface{})
		for _, p := range res.ResultParameters.ResultParameter {
			fields[p.Key] = p.Value
		}

		if amount, ok := fields["TransactionAmount"].(float64); ok {
			ev.Amount = amount
		}

		if receipt, ok := fields["TransactionReceipt"].(string); ok {
			ev.TransactionID = receipt
		}
	} else {
		ev.Status = "failed"
	}
	return ev
}
