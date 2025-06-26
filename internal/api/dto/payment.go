package dto

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/payment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/idempotency"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// PaymentRequest represents a request to create a payment
type PaymentRequest struct {
	// reference id
	// ID of the resource (internship ID, subscription ID, etc.)
	ReferenceID   string                       `json:"reference_id"`
	ReferenceType types.PaymentDestinationType `json:"reference_type"`

	// destination id
	// ID of the destination (internship ID, subscription ID, etc.)
	DestinationID   string                       `json:"destination_id"`
	DestinationType types.PaymentDestinationType `json:"destination_type"`

	// amount
	// Amount in smallest unit (e.g. paisa)
	Amount int64 `json:"amount"`

	// currency
	// "INR", "USD", etc.
	Currency string `json:"currency"`

	// payment gateway provider
	// razorpay, stripe, etc.
	PaymentGatewayProvider types.PaymentGatewayProvider `json:"payment_gateway_provider"`

	// payment method type
	// upi, card, wallet, etc.
	PaymentMethodType *types.PaymentMethodType `json:"payment_method_type,omitempty"`
	PaymentMethodID   *string                  `json:"payment_method_id,omitempty"`

	// success url
	// Frontend success redirect
	SuccessURL string `json:"success_url,omitempty"`

	// cancel url
	// Frontend cancel redirect
	CancelURL string `json:"cancel_url,omitempty"`

	// return url
	// For legacy gateway naming compatibility
	ReturnURL string `json:"return_url,omitempty"`

	// idempotency key
	// Optional, for safe retries
	IdempotencyKey string `json:"idempotency_key,omitempty"`

	// metadata
	// Additional tracking info (origin, cohort, etc.)
	Metadata map[string]string `json:"metadata,omitempty"`

	// track attempts
	// Whether to track payment attempts
	TrackAttempts bool `json:"track_attempts,omitempty"`
}

// Validate validates the payment request
func (r *PaymentRequest) Validate() error {
	if r.Amount <= 0 {
		return ierr.NewError("invalid amount").
			WithHint("Amount must be greater than 0").
			Mark(ierr.ErrValidation)
	}

	if r.DestinationID == "" {
		return ierr.NewError("invalid destination id").
			WithHint("Destination id is required").
			Mark(ierr.ErrValidation)
	}

	if err := r.DestinationType.Validate(); err != nil {
		return ierr.NewError("invalid destination type").
			WithHint("Destination type is invalid").
			Mark(ierr.ErrValidation)
	}

	if r.Currency == "" {
		return ierr.NewError("invalid currency").
			WithHint("Currency is required").
			Mark(ierr.ErrValidation)
	}

	if err := r.PaymentGatewayProvider.Validate(); err != nil {
		return ierr.NewError("invalid payment gateway provider").
			WithHint("Payment gateway provider is invalid").
			Mark(ierr.ErrValidation)
	}

	if r.PaymentMethodType != nil && r.PaymentMethodType.Validate() != nil {
		return ierr.NewError("invalid payment method type").
			WithHint("Payment method type is invalid").
			Mark(ierr.ErrValidation)
	}

	return nil
}

// ToPayment converts the request to a domain payment model
func (r *PaymentRequest) ToPayment(ctx context.Context) *payment.Payment {
	now := time.Now()
	userID := types.GetUserID(ctx)

	// Generate idempotency key if not provided
	idempotencyKey := r.IdempotencyKey
	if idempotencyKey == "" {
		generator := idempotency.NewGenerator()
		params := map[string]interface{}{
			"destination_id":   r.DestinationID,
			"destination_type": r.DestinationType,
			"amount":           r.Amount,
			"currency":         r.Currency,
			"user_id":          userID,
			"timestamp":        now.Unix(),
		}
		idempotencyKey = generator.GenerateKey(idempotency.ScopePayment, params)
	}

	// Convert amount from int64 to decimal
	// TODO: This is a temporary solution to convert the amount to decimal.
	amount := decimal.NewFromInt(r.Amount).Div(decimal.NewFromInt(100))

	paymentMethodID := ""
	if r.PaymentMethodID != nil {
		paymentMethodID = *r.PaymentMethodID
	}

	return &payment.Payment{
		ID:                     types.GenerateUUIDWithPrefix(types.UUID_PREFIX_PAYMENT),
		IdempotencyKey:         idempotencyKey,
		DestinationType:        r.DestinationType,
		DestinationID:          r.DestinationID,
		PaymentMethodType:      r.PaymentMethodType,
		PaymentMethodID:        paymentMethodID,
		PaymentGatewayProvider: r.PaymentGatewayProvider,
		Amount:                 amount,
		Currency:               types.Currency(r.Currency),
		PaymentStatus:          types.PaymentStatusPending,
		TrackAttempts:          r.TrackAttempts,
		Metadata:               r.Metadata,
		BaseModel: types.BaseModel{
			Status:    types.StatusPublished,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	PaymentRequest
}

// UpdatePaymentRequest represents a request to update a payment
type UpdatePaymentRequest struct {
	PaymentStatus    *types.PaymentStatus `json:"payment_status,omitempty"`
	GatewayPaymentID *string              `json:"gateway_payment_id,omitempty"`
	ErrorMessage     *string              `json:"error_message,omitempty"`
	Metadata         map[string]string    `json:"metadata,omitempty"`
}

// Validate validates the update request
func (r *UpdatePaymentRequest) Validate() error {
	if r.PaymentStatus != nil && r.PaymentStatus.Validate() != nil {
		return ierr.NewError("invalid payment status").
			WithHint("Payment status is invalid").
			Mark(ierr.ErrValidation)
	}
	return nil
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	Payment payment.Payment `json:"payment"`

	// Gateway specific response data
	GatewayResponse *PaymentGatewayResponse `json:"gateway_response,omitempty"`
}

// PaymentGatewayResponse represents the response from payment gateway
type PaymentGatewayResponse struct {
	ProviderPaymentID string                 `json:"provider_payment_id"`
	RedirectURL       string                 `json:"redirect_url,omitempty"`
	Status            string                 `json:"status"`
	Raw               map[string]interface{} `json:"raw,omitempty"` // Raw provider response
}

// ListPaymentResponse represents a paginated list of payments
type ListPaymentResponse struct {
	Items      []*PaymentResponse        `json:"items"`
	Pagination *types.PaginationResponse `json:"pagination"`
}

// PaymentAttemptRequest represents a request to create a payment attempt
type PaymentAttemptRequest struct {
	PaymentID        string              `json:"payment_id"`
	PaymentStatus    types.PaymentStatus `json:"payment_status"`
	GatewayAttemptID *string             `json:"gateway_attempt_id,omitempty"`
	ErrorMessage     *string             `json:"error_message,omitempty"`
	Metadata         types.Metadata      `json:"metadata,omitempty"`
}

// Validate validates the payment attempt request
func (r *PaymentAttemptRequest) Validate() error {
	if r.PaymentID == "" {
		return ierr.NewError("invalid payment id").
			WithHint("Payment id is required").
			Mark(ierr.ErrValidation)
	}

	if err := r.PaymentStatus.Validate(); err != nil {
		return ierr.NewError("invalid payment status").
			WithHint("Payment status is invalid").
			Mark(ierr.ErrValidation)
	}

	return nil
}

// ToPaymentAttempt converts the request to a domain payment attempt model
func (r *PaymentAttemptRequest) ToPaymentAttempt(ctx context.Context, attemptNumber int) *payment.PaymentAttempt {
	now := time.Now()
	userID := types.GetUserID(ctx)

	return &payment.PaymentAttempt{
		ID:               types.GenerateUUIDWithPrefix(types.UUID_PREFIX_PAYMENT_ATTEMPT),
		PaymentID:        r.PaymentID,
		AttemptNumber:    attemptNumber,
		PaymentStatus:    r.PaymentStatus,
		GatewayAttemptID: r.GatewayAttemptID,
		ErrorMessage:     r.ErrorMessage,
		Metadata:         r.Metadata,
		BaseModel: types.BaseModel{
			Status:    types.StatusPublished,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
}

// PaymentAttemptResponse represents a payment attempt response
type PaymentAttemptResponse struct {
	Attempt payment.PaymentAttempt `json:"attempt"`
}

// PaymentStatus represents the current status of a payment
type PaymentStatus struct {
	Status            string
	Reason            string
	ProviderPaymentID string
	Raw               map[string]interface{}
}

// WebhookResult represents the result of processing a webhook
type WebhookResult struct {
	EventName string
	EventID   string
	Payload   map[string]interface{}
	Headers   map[string]string
	Raw       map[string]interface{}
}
