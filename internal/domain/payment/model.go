package payment

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// Payment represents a payment transaction
type Payment struct {
	ID             string `json:"id"`
	IdempotencyKey string `json:"idempotency_key"`

	// payment destination
	DestinationType   types.PaymentDestinationType `json:"destination_type"`
	DestinationID     string                       `json:"destination_id"`
	PaymentMethodType types.PaymentMethodType      `json:"payment_method_type"`
	PaymentMethodID   string                       `json:"payment_method_id"`

	// payment gateway
	PaymentGatewayProvider types.PaymentGatewayProvider `json:"payment_gateway_provider"`
	GatewayPaymentID       *string                      `json:"gateway_payment_id,omitempty"`

	// payment amount and currency
	Amount   decimal.Decimal `json:"amount"`
	Currency types.Currency  `json:"currency"`

	// payment status
	PaymentStatus types.PaymentStatus `json:"payment_status"`

	// payment tracking
	TrackAttempts bool       `json:"track_attempts"`
	SucceededAt   *time.Time `json:"succeeded_at,omitempty"`
	FailedAt      *time.Time `json:"failed_at,omitempty"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`

	// payment error
	ErrorMessage *string `json:"error_message,omitempty"`

	// payment metadata
	Metadata types.Metadata    `json:"metadata,omitempty"`
	Attempts []*PaymentAttempt `json:"attempts,omitempty"`
	types.BaseModel
}

// PaymentAttempt represents an attempt to process a payment
type PaymentAttempt struct {
	ID               string              `json:"id"`
	PaymentID        string              `json:"payment_id"`
	AttemptNumber    int                 `json:"attempt_number"`
	PaymentStatus    types.PaymentStatus `json:"payment_status"`
	GatewayAttemptID *string             `json:"gateway_attempt_id,omitempty"`
	ErrorMessage     *string             `json:"error_message,omitempty"`
	Metadata         types.Metadata      `json:"metadata,omitempty"`
	types.BaseModel
}

// Validate validates the payment
func (p *Payment) Validate() error {
	if p.Amount.IsZero() || p.Amount.IsNegative() {
		return ierr.NewError("invalid amount").
			WithHint("Amount must be greater than 0").
			Mark(ierr.ErrValidation)
	}

	if err := p.DestinationType.Validate(); err != nil {
		return ierr.NewError("invalid destination type").
			WithHint("Destination type is invalid").
			Mark(ierr.ErrValidation)
	}

	if p.DestinationID == "" {
		return ierr.NewError("invalid destination id").
			WithHint("Destination id is invalid").
			Mark(ierr.ErrValidation)
	}

	if err := p.PaymentMethodType.Validate(); err != nil {
		return ierr.NewError("invalid payment method type").
			WithHint("Payment method type is invalid").
			Mark(ierr.ErrValidation)
	}

	if err := p.Currency.Validate(); err != nil {
		return ierr.NewError("invalid currency").
			WithHint("Currency is invalid").
			Mark(ierr.ErrValidation)
	}

	if p.PaymentMethodType == types.PaymentMethodTypeOffline && p.PaymentMethodID != "" {
		return ierr.NewError("payment method id is not allowed for offline payment method type").
			WithHint("Payment method id is invalid").
			Mark(ierr.ErrValidation)
	}

	if p.PaymentMethodType != types.PaymentMethodTypeOffline && p.PaymentMethodID == "" {
		return ierr.NewError("invalid payment method id").
			WithHint("Payment method id is required").
			Mark(ierr.ErrValidation)
	}

	return nil
}

// Validate validates the payment attempt
func (pa *PaymentAttempt) Validate() error {
	if pa.PaymentID == "" {
		return ierr.NewError("invalid payment id").
			WithHint("Payment id is invalid").
			Mark(ierr.ErrValidation)
	}

	if pa.AttemptNumber <= 0 {
		return ierr.NewError("invalid attempt number").
			WithHint("Attempt number must be greater than 0").
			Mark(ierr.ErrValidation)
	}

	return nil
}

// TableName returns the table name for Payment
func (p *Payment) TableName() string {
	return "payments"
}

// TableName returns the table name for PaymentAttempt
func (pa *PaymentAttempt) TableName() string {
	return "payment_attempts"
}

// FromEnt converts an Ent payment to domain model
func FromEnt(p *ent.Payment) *Payment {
	if p == nil {
		return nil
	}

	payment := &Payment{
		ID:                     p.ID,
		IdempotencyKey:         p.IdempotencyKey,
		DestinationType:        p.DestinationType,
		DestinationID:          p.DestinationID,
		PaymentMethodType:      p.PaymentMethodType,
		PaymentMethodID:        p.PaymentMethodID,
		PaymentGatewayProvider: lo.FromPtr(p.PaymentGatewayProvider),
		GatewayPaymentID:       p.GatewayPaymentID,
		Amount:                 p.Amount,
		Currency:               p.Currency,
		PaymentStatus:          p.PaymentStatus,
		TrackAttempts:          p.TrackAttempts,
		SucceededAt:            p.SucceededAt,
		FailedAt:               p.FailedAt,
		RefundedAt:             p.RefundedAt,
		ErrorMessage:           p.ErrorMessage,
		Metadata:               p.Metadata,
		BaseModel: types.BaseModel{
			Status:    types.Status(p.Status),
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			CreatedBy: p.CreatedBy,
			UpdatedBy: p.UpdatedBy,
		},
	}

	if attempts := p.Edges.Attempts; attempts != nil {
		payment.Attempts = FromEntAttemptList(attempts)
	}

	return payment
}

// FromEntAttempt converts an Ent payment attempt to domain model
func FromEntAttempt(a *ent.PaymentAttempt) *PaymentAttempt {
	if a == nil {
		return nil
	}

	return &PaymentAttempt{
		ID:               a.ID,
		PaymentID:        a.PaymentID,
		AttemptNumber:    a.AttemptNumber,
		PaymentStatus:    types.PaymentStatus(a.PaymentStatus),
		GatewayAttemptID: a.GatewayAttemptID,
		ErrorMessage:     a.ErrorMessage,
		Metadata:         types.MetadataFromEnt(a.Metadata),
		BaseModel: types.BaseModel{
			Status:    types.Status(a.Status),
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			CreatedBy: a.CreatedBy,
			UpdatedBy: a.UpdatedBy,
		},
	}
}

// FromEntList converts a list of Ent payments to domain payments
func FromEntList(payments []*ent.Payment) []*Payment {
	result := make([]*Payment, 0, len(payments))
	for _, p := range payments {
		result = append(result, FromEnt(p))
	}
	return result
}

// FromEntAttemptList converts a list of Ent payment attempts to domain model
func FromEntAttemptList(attempts []*ent.PaymentAttempt) []*PaymentAttempt {
	result := make([]*PaymentAttempt, 0, len(attempts))
	for _, a := range attempts {
		result = append(result, FromEntAttempt(a))
	}
	return result
}
