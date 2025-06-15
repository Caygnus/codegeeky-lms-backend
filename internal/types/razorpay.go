package types

import "encoding/json"

type AccountInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RazorpayWebhookPayload struct {
	Event   string          `json:"event"`
	Entity  string          `json:"entity"`
	Account AccountInfo     `json:"account"`
	Payment json.RawMessage `json:"payment,omitempty"`
	Order   json.RawMessage `json:"order,omitempty"`
	Refund  json.RawMessage `json:"refund,omitempty"`
}

type RazorpayPaymentEvent struct {
	Event   string `json:"event"`
	Entity  string `json:"entity"`
	Payment struct {
		ID          string `json:"id"`
		Amount      int    `json:"amount"`
		Currency    string `json:"currency"`
		Status      string `json:"status"`
		OrderID     string `json:"order_id"`
		Email       string `json:"email"`
		Contact     string `json:"contact"`
		Fee         int    `json:"fee"`
		Tax         int    `json:"tax"`
		Description string `json:"description"`
		CreatedAt   int64  `json:"created_at"`
	} `json:"payment"`
}
