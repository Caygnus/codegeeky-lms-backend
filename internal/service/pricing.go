package service

import (
	"context"
	"fmt"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/domain/discount"
	"github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// PricingConfig holds configuration for pricing calculations
type PricingConfig struct {
	// Allow zero or negative final prices (for full discounts)
	AllowZeroPrice bool
	// Maximum discount percentage allowed (e.g., 100 for 100%)
	MaxDiscountPercent decimal.Decimal
	// Minimum final price (e.g., 0 for free internships)
	MinFinalPrice decimal.Decimal
	// Currency formatting options
	DefaultCurrency string
}

// DefaultPricingConfig returns a sensible default configuration
func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		AllowZeroPrice:     true,
		MaxDiscountPercent: decimal.NewFromInt(100), // Allow up to 100% discount
		MinFinalPrice:      decimal.Zero,
		DefaultCurrency:    "INR",
	}
}

type PricingService interface {
	// Calculate pricing with optional discount code validation
	CalculateInternshipPricing(ctx context.Context, req *dto.PricingRequest) (*dto.PricingResponse, error)

	// Validate discount code without calculating final pricing (for UI feedback)
	ValidateCouponCode(ctx context.Context, internshipID, discountCode string) (*dto.DiscountInfo, error)

	// Calculate pricing breakdown for enrollment service
	CalculateEnrollmentPricing(ctx context.Context, internshipID, discountCode, userID string) (*PricingBreakdown, error)
}

// PricingBreakdown provides detailed pricing information for internal use
type PricingBreakdown struct {
	InternshipID       string
	UserID             string
	OriginalPrice      decimal.Decimal
	InternshipDiscount decimal.Decimal
	CouponDiscount     decimal.Decimal
	TotalDiscount      decimal.Decimal
	FinalPrice         decimal.Decimal
	Currency           string
	AppliedDiscount    *discount.Discount
	IsPaymentRequired  bool
	CalculatedAt       time.Time
}

type pricingService struct {
	ServiceParams
	config *PricingConfig
}

func NewPricingService(params ServiceParams) PricingService {
	return &pricingService{
		ServiceParams: params,
		config:        DefaultPricingConfig(),
	}
}

func NewPricingServiceWithConfig(params ServiceParams, config *PricingConfig) PricingService {
	if config == nil {
		config = DefaultPricingConfig()
	}
	return &pricingService{
		ServiceParams: params,
		config:        config,
	}
}

func (s *pricingService) CalculateInternshipPricing(ctx context.Context, req *dto.PricingRequest) (*dto.PricingResponse, error) {
	// Calculate detailed pricing breakdown
	breakdown, err := s.CalculateEnrollmentPricing(ctx, req.InternshipID, req.DiscountCode, req.UserID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := &dto.PricingResponse{
		InternshipID:    breakdown.InternshipID,
		UserID:          breakdown.UserID,
		OriginalPrice:   breakdown.OriginalPrice,
		DiscountAmount:  breakdown.TotalDiscount,
		FinalPrice:      breakdown.FinalPrice,
		Currency:        breakdown.Currency,
		PaymentRequired: breakdown.IsPaymentRequired,
	}

	// Add discount information if applicable
	if breakdown.AppliedDiscount != nil {
		response.AppliedDiscount = &dto.DiscountInfo{
			Type:        "coupon",
			Code:        breakdown.AppliedDiscount.Code,
			Amount:      breakdown.CouponDiscount,
			Description: breakdown.AppliedDiscount.Description,
			IsValid:     true,
		}
		response.CouponMessage = fmt.Sprintf("%s applied successfully", breakdown.AppliedDiscount.Code)
	}

	// Calculate savings percentage
	if breakdown.TotalDiscount.GreaterThan(decimal.Zero) && breakdown.OriginalPrice.GreaterThan(decimal.Zero) {
		savings := breakdown.TotalDiscount.Div(breakdown.OriginalPrice).Mul(decimal.NewFromInt(100))
		response.SavingsPercent = savings.InexactFloat64()
	}

	// Generate user-friendly messages
	response.PricingMessage = s.generatePricingMessage(breakdown)

	return response, nil
}

func (s *pricingService) ValidateCouponCode(ctx context.Context, internshipID, discountCode string) (*dto.DiscountInfo, error) {
	if discountCode == "" {
		return nil, ierr.NewError("discount code is required").
			WithHint("Please provide a discount code").
			Mark(ierr.ErrValidation)
	}

	// Get internship details
	internship, err := s.InternshipRepo.Get(ctx, internshipID)
	if err != nil {
		return nil, err
	}

	// Validate discount code
	discountService := NewDiscountService(s.ServiceParams)
	err = discountService.ValidateDiscountCode(ctx, discountCode, internship)
	if err != nil {
		return &dto.DiscountInfo{
			Type:        "coupon",
			Code:        discountCode,
			Amount:      decimal.Zero,
			Description: "Invalid or expired coupon code",
			IsValid:     false,
		}, nil // Return success with invalid flag instead of error for better UX
	}

	// Get discount details
	discountResp, err := discountService.GetByCode(ctx, discountCode)
	if err != nil {
		return nil, err
	}

	// Calculate discount amount
	discountAmount := s.calculateCouponDiscount(internship, &discountResp.Discount)

	return &dto.DiscountInfo{
		Type:        "coupon",
		Code:        discountResp.Discount.Code,
		Amount:      discountAmount,
		Description: discountResp.Discount.Description,
		IsValid:     true,
	}, nil
}

func (s *pricingService) CalculateEnrollmentPricing(ctx context.Context, internshipID, discountCode, userID string) (*PricingBreakdown, error) {
	// 1. Get internship details
	internship, err := s.InternshipRepo.Get(ctx, internshipID)
	if err != nil {
		return nil, err
	}

	// 2. Initialize pricing breakdown
	breakdown := &PricingBreakdown{
		InternshipID:       internshipID,
		UserID:             userID,
		OriginalPrice:      internship.Price,
		InternshipDiscount: decimal.Zero,
		CouponDiscount:     decimal.Zero,
		TotalDiscount:      decimal.Zero,
		Currency:           internship.Currency,
		CalculatedAt:       time.Now().UTC(),
	}

	// 3. Apply internship-level discounts (flat and percentage)
	breakdown.InternshipDiscount = s.calculateInternshipDiscount(internship)

	// 4. Apply coupon discount if provided
	if discountCode != "" {
		couponDiscount, appliedDiscount, err := s.applyCouponDiscount(ctx, internship, discountCode)
		if err != nil {
			return nil, err
		}
		breakdown.CouponDiscount = couponDiscount
		breakdown.AppliedDiscount = appliedDiscount
	}

	// 5. Calculate totals
	breakdown.TotalDiscount = breakdown.InternshipDiscount.Add(breakdown.CouponDiscount)
	breakdown.FinalPrice = breakdown.OriginalPrice.Sub(breakdown.TotalDiscount)

	// 6. Apply pricing constraints
	if breakdown.FinalPrice.LessThan(s.config.MinFinalPrice) {
		breakdown.FinalPrice = s.config.MinFinalPrice
	}

	// 7. Determine if payment is required
	breakdown.IsPaymentRequired = breakdown.FinalPrice.GreaterThan(decimal.Zero)

	return breakdown, nil
}

func (s *pricingService) calculateInternshipDiscount(internship *internship.Internship) decimal.Decimal {
	totalDiscount := decimal.Zero

	// Apply flat discount
	if !internship.FlatDiscount.IsZero() {
		totalDiscount = totalDiscount.Add(internship.FlatDiscount)
	}

	// Apply percentage discount
	if internship.PercentageDiscount.GreaterThan(decimal.Zero) {
		percentageDiscount := internship.Price.Mul(internship.PercentageDiscount).Div(decimal.NewFromInt(100))
		totalDiscount = totalDiscount.Add(percentageDiscount)
	}

	return totalDiscount
}

func (s *pricingService) applyCouponDiscount(ctx context.Context, internship *internship.Internship, discountCode string) (decimal.Decimal, *discount.Discount, error) {
	discountService := NewDiscountService(s.ServiceParams)

	// Validate discount code
	err := discountService.ValidateDiscountCode(ctx, discountCode, internship)
	if err != nil {
		return decimal.Zero, nil, err
	}

	// Get discount details
	discountResp, err := discountService.GetByCode(ctx, discountCode)
	if err != nil {
		return decimal.Zero, nil, err
	}

	// Calculate discount amount
	discountAmount := s.calculateCouponDiscount(internship, &discountResp.Discount)

	return discountAmount, &discountResp.Discount, nil
}

func (s *pricingService) calculateCouponDiscount(internship *internship.Internship, discount *discount.Discount) decimal.Decimal {
	basePrice := internship.Price

	// Apply internship-level discounts first to get the discounted base price
	internshipDiscount := s.calculateInternshipDiscount(internship)
	basePrice = basePrice.Sub(internshipDiscount)

	// Ensure base price is not negative
	if basePrice.LessThan(decimal.Zero) {
		basePrice = decimal.Zero
	}

	switch discount.DiscountType {
	case types.DiscountTypeFlat:
		return discount.DiscountValue
	case types.DiscountTypePercentage:
		return basePrice.Mul(discount.DiscountValue).Div(decimal.NewFromInt(100))
	default:
		return decimal.Zero
	}
}

func (s *pricingService) generatePricingMessage(breakdown *PricingBreakdown) string {
	if breakdown.TotalDiscount.GreaterThan(decimal.Zero) {
		if breakdown.FinalPrice.IsZero() {
			return "This internship is completely free!"
		}
		return fmt.Sprintf("Final price: %s %s (You're saving %s %s)",
			breakdown.Currency,
			breakdown.FinalPrice.String(),
			breakdown.Currency,
			breakdown.TotalDiscount.String())
	}
	return fmt.Sprintf("Final price: %s %s", breakdown.Currency, breakdown.FinalPrice.String())
}
