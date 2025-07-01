package dto

import (
	"github.com/shopspring/decimal"
)

// PricingRequest represents a request to calculate pricing for an internship
type PricingRequest struct {
	InternshipID string `json:"internship_id" validate:"required"`
	DiscountCode string `json:"discount_code,omitempty"`
	UserID       string `json:"user_id,omitempty"`
}

// PricingResponse provides pricing information for enrollment and coupon validation
type PricingResponse struct {
	InternshipID string `json:"internship_id"`
	UserID       string `json:"user_id,omitempty"`

	// Pricing breakdown
	OriginalPrice  decimal.Decimal `json:"original_price"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	FinalPrice     decimal.Decimal `json:"final_price"`
	Currency       string          `json:"currency"`

	// Discount information (if any discount is applied)
	AppliedDiscount *DiscountInfo `json:"applied_discount,omitempty"`

	// User-friendly information
	PaymentRequired bool    `json:"payment_required"`
	SavingsPercent  float64 `json:"savings_percent,omitempty"`

	// Messages for UI display
	PricingMessage string `json:"pricing_message,omitempty"` // "Final price: ₹2000" or "You're saving ₹500 (20%)"
	CouponMessage  string `json:"coupon_message,omitempty"`  // "SAVE20 applied successfully" or "Invalid coupon code"
}

// DiscountInfo contains information about applied discounts
type DiscountInfo struct {
	Type        string          `json:"type"`           // "internship", "coupon"
	Code        string          `json:"code,omitempty"` // coupon code if applicable
	Amount      decimal.Decimal `json:"amount"`         // discount amount
	Description string          `json:"description"`    // user-friendly description
	IsValid     bool            `json:"is_valid"`       // for coupon validation
}
