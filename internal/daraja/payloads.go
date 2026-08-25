package daraja

type CallbackItem struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}

type CallbackMetadata struct {
	Item []CallbackItem `json:"Item"`
}

type STKCallbackPayload struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string            `json:"MerchantRequestID"`
			CheckoutRequestID string            `json:"CheckoutRequestID"`
			ResultCode        int               `json:"ResultCode"`
			ResultDesc        string            `json:"ResultDesc"`
			CallbackMetadata  *CallbackMetadata `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

type C2BPayload struct {
	TransactionType   string `json:"TransactionType"`
	TransID           string `json:"TransID"`
	TransAmount       string `json:"TransAmount"`
	BusinessShortCode string `json:"BusinessShortCode"`
	BillRefNumber     string `json:"BillRefNumber"`
	MSISDN            string `json:"MSISDN"`
	TransTime         string `json:"TransTime"`
}

type B2CResultParameter struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}

type B2CResultParameters struct {
	ResultParameter []B2CResultParameter `json:"ResultParameter"`
}

type B2CPayload struct {
	Result struct {
		ResultType               int                  `json:"ResultType"`
		ResultCode               int                  `json:"ResultCode"`
		ResultDesc               string               `json:"ResultDesc"`
		OriginatorConversationID string               `json:"OriginatorConversationID"`
		ConversationID           string               `json:"ConversationID"`
		TransactionID            string               `json:"TransactionID"`
		ResultParameters         *B2CResultParameters `json:"ResultParameters"`
	} `json:"Result"`
}
