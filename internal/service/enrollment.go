package service

import (
	"context"
	"fmt"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainEnrollment "github.com/omkar273/codegeeky/internal/domain/enrollment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/idempotency"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type EnrollmentService interface {
	InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.InitializeEnrollmentResponse, error)
	FinalizeEnrollment(ctx context.Context, id string) (*dto.EnrollmentStatusResponse, error)
	GetEnrollment(ctx context.Context, id string) (*dto.EnrollmentStatusResponse, error)
	CancelEnrollment(ctx context.Context, id string) error
}

type enrollmentService struct {
	ServiceParams
	pricingService PricingService
	paymentService PaymentService
}

func NewEnrollmentService(params ServiceParams) EnrollmentService {
	return &enrollmentService{
		ServiceParams:  params,
		pricingService: NewPricingService(params),
		paymentService: NewPaymentService(params, nil), // Assuming gateway registry is optional
	}
}

func (s *enrollmentService) InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.InitializeEnrollmentResponse, error) {
	if req.UserID == "" {
		req.UserID = types.GetUserID(ctx)
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate user exists
	user, err := s.ServiceParams.UserRepo.Get(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ierr.NewError("user not found").
			WithHint("User not found").
			Mark(ierr.ErrNotFound)
	}

	// Validate internship exists
	internship, err := s.ServiceParams.InternshipRepo.Get(ctx, req.InternshipID)
	if err != nil {
		return nil, err
	}

	if internship == nil {
		return nil, ierr.NewError("internship not found").
			WithHint("Internship not found").
			Mark(ierr.ErrNotFound)
	}

	// Generate idempotency key
	generator := idempotency.NewGenerator()
	idempotencyKey := generator.GenerateKey(idempotency.ScopeInternship, map[string]interface{}{
		"user_id":       req.UserID,
		"internship_id": req.InternshipID,
	})

	// Check for existing enrollment
	existingEnrollment, err := s.ServiceParams.EnrollmentRepo.GetByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !ierr.IsNotFound(err) {
		return nil, err
	}

	if existingEnrollment != nil {
		if existingEnrollment.EnrollmentStatus == types.EnrollmentStatusCompleted {
			return nil, ierr.NewError("enrollment already completed").
				WithHint("Enrollment already completed").
				Mark(ierr.ErrAlreadyExists)
		}

		// Return existing enrollment if still pending
		return s.buildEnrollmentResponse(ctx, existingEnrollment, req)
	}

	// Calculate pricing with coupon (use first coupon code if multiple provided)
	var discountCode string
	if len(req.CouponCodes) > 0 {
		discountCode = req.CouponCodes[0]
	}

	pricingBreakdown, err := s.pricingService.CalculateEnrollmentPricing(ctx, req.InternshipID, discountCode, req.UserID)
	if err != nil {
		return nil, err
	}

	// Create enrollment
	enrollmentData := &domainEnrollment.Enrollment{
		UserID:           req.UserID,
		InternshipID:     req.InternshipID,
		EnrollmentStatus: types.EnrollmentStatusPending,
		PaymentStatus:    types.PaymentStatusPending,
		IdempotencyKey:   &idempotencyKey,
		Metadata:         req.Metadata,
	}

	// If payment is not required, mark as completed
	if !pricingBreakdown.IsPaymentRequired {
		enrollmentData.EnrollmentStatus = types.EnrollmentStatusCompleted
		enrollmentData.PaymentStatus = types.PaymentStatusSuccess
	}

	err = s.ServiceParams.EnrollmentRepo.Create(ctx, enrollmentData)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.InitializeEnrollmentResponse{
		EnrollmentID:    enrollmentData.ID,
		Status:          string(enrollmentData.EnrollmentStatus),
		PaymentRequired: pricingBreakdown.IsPaymentRequired,
		Pricing: &dto.PricingInfo{
			OriginalAmount: pricingBreakdown.OriginalPrice,
			DiscountAmount: pricingBreakdown.TotalDiscount,
			FinalAmount:    pricingBreakdown.FinalPrice,
			Currency:       types.Currency(pricingBreakdown.Currency),
			TaxAmount:      decimal.Zero, // TODO: Implement tax calculation if needed
			NetPayable:     pricingBreakdown.FinalPrice,
		},
	}

	// Create payment session if payment is required
	if pricingBreakdown.IsPaymentRequired {
		paymentSession, err := s.createPaymentSession(ctx, enrollmentData, pricingBreakdown, req)
		if err != nil {
			s.ServiceParams.Logger.Errorw("failed to create payment session",
				"enrollment_id", enrollmentData.ID, "error", err)
			// Don't fail enrollment creation, payment can be retried
		} else {
			response.PaymentSession = paymentSession
		}
	}

	return response, nil
}

func (s *enrollmentService) buildEnrollmentResponse(ctx context.Context, enrollment *domainEnrollment.Enrollment, req *dto.InitializeEnrollmentRequest) (*dto.InitializeEnrollmentResponse, error) {
	// Recalculate pricing for existing enrollment
	var discountCode string
	if len(req.CouponCodes) > 0 {
		discountCode = req.CouponCodes[0]
	}

	pricingBreakdown, err := s.pricingService.CalculateEnrollmentPricing(ctx, enrollment.InternshipID, discountCode, enrollment.UserID)
	if err != nil {
		return nil, err
	}

	response := &dto.InitializeEnrollmentResponse{
		EnrollmentID:    enrollment.ID,
		Status:          string(enrollment.EnrollmentStatus),
		PaymentRequired: pricingBreakdown.IsPaymentRequired,
		Pricing: &dto.PricingInfo{
			OriginalAmount: pricingBreakdown.OriginalPrice,
			DiscountAmount: pricingBreakdown.TotalDiscount,
			FinalAmount:    pricingBreakdown.FinalPrice,
			Currency:       types.Currency(pricingBreakdown.Currency),
			TaxAmount:      decimal.Zero,
			NetPayable:     pricingBreakdown.FinalPrice,
		},
	}

	// Get existing payment session if available
	if enrollment.PaymentStatus == types.PaymentStatusPending && pricingBreakdown.IsPaymentRequired {
		// TODO: Get existing payment session or create new one
		paymentSession, err := s.createPaymentSession(ctx, enrollment, pricingBreakdown, req)
		if err != nil {
			s.ServiceParams.Logger.Errorw("failed to create payment session for existing enrollment",
				"enrollment_id", enrollment.ID, "error", err)
		} else {
			response.PaymentSession = paymentSession
		}
	}

	return response, nil
}

func (s *enrollmentService) createPaymentSession(ctx context.Context, enrollment *domainEnrollment.Enrollment, pricing *PricingBreakdown, req *dto.InitializeEnrollmentRequest) (*dto.PaymentSessionInfo, error) {
	// Create payment request
	paymentMethodType := types.PaymentMethodTypeCard
	paymentReq := &dto.CreatePaymentRequest{
		PaymentRequest: dto.PaymentRequest{
			ReferenceID:            enrollment.ID,
			ReferenceType:          types.PaymentDestinationTypeEnrollment,
			DestinationID:          enrollment.InternshipID,
			DestinationType:        types.PaymentDestinationTypeInternship,
			Amount:                 pricing.FinalPrice.Mul(decimal.NewFromInt(100)).IntPart(), // Convert to paisa
			Currency:               pricing.Currency,
			PaymentGatewayProvider: types.PaymentGatewayProviderRazorpay,
			PaymentMethodType:      &paymentMethodType,
			SuccessURL:             req.SuccessURL,
			CancelURL:              req.CancelURL,
			IdempotencyKey:         fmt.Sprintf("enrollment-%s", enrollment.ID),
			Metadata:               req.Metadata,
			TrackAttempts:          true,
		},
	}

	paymentResp, err := s.paymentService.Create(ctx, paymentReq)
	if err != nil {
		return nil, err
	}

	// Update enrollment with payment ID
	enrollment.PaymentID = lo.ToPtr(paymentResp.Payment.ID)
	if err := s.ServiceParams.EnrollmentRepo.Update(ctx, enrollment); err != nil {
		s.ServiceParams.Logger.Errorw("failed to update enrollment with payment ID",
			"enrollment_id", enrollment.ID, "payment_id", paymentResp.Payment.ID, "error", err)
	}

	paymentSession := &dto.PaymentSessionInfo{
		PaymentID: paymentResp.Payment.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // TODO: Make this configurable
	}

	// Add gateway-specific information if available
	if paymentResp.GatewayResponse != nil {
		paymentSession.RazorpayOrderID = paymentResp.GatewayResponse.ProviderPaymentID
		paymentSession.PaymentURL = paymentResp.GatewayResponse.RedirectURL
		// TODO: Add RazorpayKey from config
	}

	return paymentSession, nil
}

func (s *enrollmentService) FinalizeEnrollment(ctx context.Context, id string) (*dto.EnrollmentStatusResponse, error) {
	enrollment, err := s.ServiceParams.EnrollmentRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if enrollment == nil {
		return nil, ierr.NewError("enrollment not found").
			WithHint("Enrollment not found").
			Mark(ierr.ErrNotFound)
	}

	// Check payment status if payment was required
	if enrollment.PaymentID != nil {
		payment, err := s.paymentService.GetByID(ctx, *enrollment.PaymentID)
		if err != nil {
			return nil, err
		}

		// Update enrollment status based on payment status
		switch payment.Payment.PaymentStatus {
		case types.PaymentStatusSuccess:
			enrollment.EnrollmentStatus = types.EnrollmentStatusCompleted
			enrollment.PaymentStatus = types.PaymentStatusSuccess
		case types.PaymentStatusFailed:
			enrollment.EnrollmentStatus = types.EnrollmentStatusFailed
			enrollment.PaymentStatus = types.PaymentStatusFailed
		case types.PaymentStatusCancelled:
			enrollment.EnrollmentStatus = types.EnrollmentStatusCancelled
			enrollment.PaymentStatus = types.PaymentStatusCancelled
		}

		if err := s.ServiceParams.EnrollmentRepo.Update(ctx, enrollment); err != nil {
			return nil, err
		}
	}

	return s.buildEnrollmentStatusResponse(enrollment), nil
}

func (s *enrollmentService) GetEnrollment(ctx context.Context, id string) (*dto.EnrollmentStatusResponse, error) {
	enrollment, err := s.ServiceParams.EnrollmentRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if enrollment == nil {
		return nil, ierr.NewError("enrollment not found").
			WithHint("Enrollment not found").
			Mark(ierr.ErrNotFound)
	}

	return s.buildEnrollmentStatusResponse(enrollment), nil
}

func (s *enrollmentService) CancelEnrollment(ctx context.Context, id string) error {
	enrollment, err := s.ServiceParams.EnrollmentRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if enrollment == nil {
		return ierr.NewError("enrollment not found").
			WithHint("Enrollment not found").
			Mark(ierr.ErrNotFound)
	}

	if enrollment.EnrollmentStatus == types.EnrollmentStatusCompleted {
		return ierr.NewError("cannot cancel completed enrollment").
			WithHint("Cannot cancel completed enrollment").
			Mark(ierr.ErrBadRequest)
	}

	// Cancel payment if exists
	if enrollment.PaymentID != nil {
		// TODO: Implement payment cancellation
	}

	enrollment.EnrollmentStatus = types.EnrollmentStatusCancelled
	enrollment.PaymentStatus = types.PaymentStatusCancelled

	return s.ServiceParams.EnrollmentRepo.Update(ctx, enrollment)
}

func (s *enrollmentService) buildEnrollmentStatusResponse(enrollment *domainEnrollment.Enrollment) *dto.EnrollmentStatusResponse {
	response := &dto.EnrollmentStatusResponse{
		EnrollmentID:     enrollment.ID,
		EnrollmentStatus: string(enrollment.EnrollmentStatus),
		PaymentStatus:    string(enrollment.PaymentStatus),
	}

	if enrollment.PaymentID != nil {
		response.PaymentID = *enrollment.PaymentID
	}

	if enrollment.EnrolledAt != nil {
		response.CompletedAt = enrollment.EnrolledAt
	}

	// TODO: Add course access URL generation logic
	if enrollment.EnrollmentStatus == types.EnrollmentStatusCompleted {
		response.CourseAccessURL = fmt.Sprintf("/courses/%s", enrollment.InternshipID)
	}

	return response
}
