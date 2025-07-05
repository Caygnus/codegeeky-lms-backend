package dto

import (
	"github.com/shopspring/decimal"
)

// PricingRequest represents a request to calculate pricing for an internship
type PricingRequest struct {
	InternshipID string `json:"internship_id" validate:"required"`
	DiscountCode string `json:"discount_code,omitempty"`
}

// PricingResponse provides pricing information for enrollment and coupon validation
type PricingResponse struct {
	InternshipID string `json:"internship_id"`

	// Pricing breakdown
	Subtotal       decimal.Decimal `json:"subtotal"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	Total          decimal.Decimal `json:"total"`
	Currency       string          `json:"currency"`

	// Discount information (if any discount is applied)
	AppliedDiscounts []*DiscountInfo `json:"applied_discount,omitempty"`

	// User-friendly information
	PaymentRequired bool    `json:"payment_required"`
	SavingsPercent  float64 `json:"savings_percent,omitempty"`
}

// DiscountInfo contains information about applied discounts
type DiscountInfo struct {
	Code        string          `json:"code,omitempty"` // coupon code if applicable
	Amount      decimal.Decimal `json:"amount"`         // discount amount
	Description string          `json:"description"`    // user-friendly description
	IsValid     bool            `json:"is_valid"`       // for coupon validation
}
