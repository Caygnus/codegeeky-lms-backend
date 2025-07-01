package dto

import (
	"time"

	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/shopspring/decimal"
)

type InitializeEnrollmentRequest struct {
	InternshipID string `json:"internship_id" validate:"required"`

	// This is optional, it its empty it will be set from the context
	UserID string `json:"user_id" validate:"required"`

	CouponCodes []string          `json:"coupon_codes,omitempty"`
	SuccessURL  string            `json:"success_url" validate:"required,url"`
	CancelURL   string            `json:"cancel_url" validate:"required,url"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

func (r *InitializeEnrollmentRequest) Validate() error {
	if err := validator.ValidateRequest(r); err != nil {
		return err
	}

	return nil
}

type InitializeEnrollmentResponse struct {
	EnrollmentID    string              `json:"enrollment_id"`
	Status          string              `json:"status"`
	PaymentRequired bool                `json:"payment_required"`
	Pricing         *PricingInfo        `json:"pricing"`
	PaymentSession  *PaymentSessionInfo `json:"payment_session,omitempty"`
}

type PricingInfo struct {
	OriginalAmount decimal.Decimal `json:"original_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	FinalAmount    decimal.Decimal `json:"final_amount"`
	Currency       types.Currency  `json:"currency"`
	TaxAmount      decimal.Decimal `json:"tax_amount"`
	NetPayable     decimal.Decimal `json:"net_payable"`
}

type PaymentSessionInfo struct {
	PaymentID       string    `json:"payment_id"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	RazorpayKey     string    `json:"razorpay_key"`
	PaymentURL      string    `json:"payment_url,omitempty"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type EnrollmentStatusResponse struct {
	EnrollmentID      string     `json:"enrollment_id"`
	EnrollmentStatus  string     `json:"enrollment_status"`
	PaymentStatus     string     `json:"payment_status"`
	PaymentID         string     `json:"payment_id,omitempty"`
	RazorpayPaymentID string     `json:"razorpay_payment_id,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	CourseAccessURL   string     `json:"course_access_url,omitempty"`
}
