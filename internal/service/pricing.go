package service

import (
	"context"
	"fmt"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/domain/discount"
	"github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// PricingService is the service for calculating pricing for an internship
type PricingService interface {
	CalculateEnrollmentPricing(ctx context.Context, internshipID string, discountCodes []string) (*dto.PricingResponse, error)
}

// pricingService is the implementation of the PricingService interface
type pricingService struct {
	ServiceParams
}

// NewPricingService creates a new PricingService
func NewPricingService(params ServiceParams) PricingService {
	return &pricingService{
		ServiceParams: params,
	}
}

// CalculateEnrollmentPricing calculates the pricing for an internship enrollment
func (s *pricingService) CalculateEnrollmentPricing(ctx context.Context, internshipID string, discountCodes []string) (*dto.PricingResponse, error) {
	// Validate input parameters
	if internshipID == "" {
		return nil, ierr.NewError("internship ID is required").
			WithHint("Please provide a valid internship ID").
			Mark(ierr.ErrValidation)
	}

	// Get the internship
	internship, err := s.ServiceParams.InternshipRepo.Get(ctx, internshipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get internship: %w", err)
	}

	// Validate and get discounts
	var discounts []*discount.Discount
	if len(discountCodes) > 0 {
		discounts, err = s.getValidDiscounts(ctx, discountCodes, internship)
		if err != nil {
			return nil, err
		}
	}

	// Calculate pricing with discounts
	pricing := s.calculatePricingWithDiscounts(internship, discounts)

	return pricing, nil
}

// getValidDiscounts validates discount codes and retrieves discount details
func (s *pricingService) getValidDiscounts(ctx context.Context, discountCodes []string, internship *internship.Internship) ([]*discount.Discount, error) {
	// Validate each discount code
	discountService := NewDiscountService(s.ServiceParams)
	for _, code := range discountCodes {
		if err := discountService.ValidateDiscountCode(ctx, code, internship); err != nil {
			return nil, fmt.Errorf("invalid discount code '%s': %w", code, err)
		}
	}

	// Get all valid discounts in a single query
	discountsResponse, err := discountService.List(ctx, &types.DiscountFilter{
		Codes: discountCodes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve discounts: %w", err)
	}

	// Convert DTO responses to domain discount objects
	discounts := make([]*discount.Discount, len(discountsResponse.Items))
	for i, discountResponse := range discountsResponse.Items {
		discounts[i] = &discountResponse.Discount
	}

	return discounts, nil
}

// calculatePricingWithDiscounts computes the final pricing with applied discounts
func (s *pricingService) calculatePricingWithDiscounts(internship *internship.Internship, discounts []*discount.Discount) *dto.PricingResponse {
	// Initialize pricing calculation
	total := internship.Total
	discountAmount := decimal.Zero
	discountsApplied := []*dto.DiscountInfo{}

	// Apply each discount
	for _, discount := range discounts {
		var discountValue decimal.Decimal

		switch discount.DiscountType {
		case types.DiscountTypeFlat:
			// Ensure flat discount doesn't exceed the current total
			discountValue = decimal.Min(discount.DiscountValue, total)

		case types.DiscountTypePercentage:
			// Calculate percentage discount on current total
			percentageDiscount := total.Mul(discount.DiscountValue).Div(decimal.NewFromInt(100))
			discountValue = decimal.Min(percentageDiscount, total)
		}
		discountAmount = discountAmount.Add(discountValue)
		total = total.Sub(discountValue)

		// Ensure total never goes below zero
		if total.LessThan(decimal.Zero) {
			total = decimal.Zero
		}

		// Add discount info to applied discounts list
		discountsApplied = append(discountsApplied, &dto.DiscountInfo{
			Code:        discount.Code,
			Amount:      discountValue,
			Description: discount.Description,
		})

	}

	// Calculate savings percentage
	savingsPercent := decimal.Zero
	if internship.Subtotal.GreaterThan(decimal.Zero) {
		savingsPercent = discountAmount.Div(internship.Subtotal).Mul(decimal.NewFromInt(100))
	}

	return &dto.PricingResponse{
		InternshipID:     internship.ID,
		Subtotal:         internship.Subtotal,
		Total:            total,
		DiscountAmount:   discountAmount,
		AppliedDiscounts: discountsApplied,
		PaymentRequired:  total.GreaterThan(decimal.Zero),
		SavingsPercent:   savingsPercent.InexactFloat64(),
	}
}
