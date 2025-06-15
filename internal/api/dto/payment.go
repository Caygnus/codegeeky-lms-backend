package dto

type PaymentRequest struct {
	Amount      int64
	Currency    string
	CustomerID  string
	ReturnURL   string
	Metadata    map[string]string
	PaymentMode string // e.g. "card", "upi", "netbanking"
}

type PaymentResponse struct {
	ProviderPaymentID string
	RedirectURL       string
	Status            string
	Raw               map[string]interface{} // Raw provider response
}

type PaymentStatus struct {
	Status            string
	Reason            string
	ProviderPaymentID string
	Raw               map[string]interface{}
}

type WebhookResult struct {
	EventName string
	EventID   string
	Payload   map[string]interface{}
	Headers   map[string]string
	Raw       map[string]interface{}
}
