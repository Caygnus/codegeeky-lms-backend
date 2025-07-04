package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainInternshipEnrollment "github.com/omkar273/codegeeky/internal/domain/internshipenrollment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/idempotency"
	"github.com/omkar273/codegeeky/internal/types"
)

// InternshipEnrollmentService is the service for managing internship enrollments
type InternshipEnrollmentService interface {
	InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.InitializeEnrollmentResponse, error)
}

// internshipEnrollmentService is the implementation of the InternshipEnrollmentService interface
type internshipEnrollmentService struct {
	ServiceParams
	PricingService PricingService
}

// NewInternshipEnrollmentService creates a new InternshipEnrollmentService
func NewInternshipEnrollmentService(params ServiceParams, pricingService PricingService) InternshipEnrollmentService {
	return &internshipEnrollmentService{
		ServiceParams:  params,
		PricingService: pricingService,
	}
}

// InitializeEnrollment initializes an internship enrollment
func (s *internshipEnrollmentService) InitializeEnrollment(ctx context.Context, req *dto.InitializeEnrollmentRequest) (*dto.InitializeEnrollmentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	batch, err := s.ServiceParams.InternshipBatchRepo.Get(ctx, req.InternshipBatchID)
	if err != nil {
		return nil, err
	}

	// Generate idempotency key
	generator := idempotency.NewGenerator()
	idempotencyKey := generator.GenerateKey(idempotency.ScopeInternship, map[string]interface{}{
		"user_id":       types.GetUserID(ctx),
		"internship_id": batch.InternshipID,
		"batch_id":      batch.ID,
	})

	// Check for existing enrollment
	existingEnrollment, err := s.ServiceParams.InternshipEnrollmentRepo.GetByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !ierr.IsNotFound(err) {
		return nil, err
	}

	if existingEnrollment != nil {
		if existingEnrollment.EnrollmentStatus == types.InternshipEnrollmentStatusCompleted {
			return nil, ierr.NewError("enrollment already completed").
				WithHint("Enrollment already completed").
				Mark(ierr.ErrAlreadyExists)
		}

		// Return existing enrollment if still pending
		return &dto.InitializeEnrollmentResponse{
			EnrollmentID:     existingEnrollment.ID,
			EnrollmentStatus: existingEnrollment.EnrollmentStatus,
			PaymentRequired:  existingEnrollment.PaymentStatus == types.PaymentStatusPending,
		}, nil
	}

	// Calculate pricing with coupon (use first coupon code if multiple provided)
	var discountCode string
	if len(req.CouponCodes) > 0 {
		discountCode = req.CouponCodes[0]
	}

	pricingResponse, err := s.PricingService.CalculateEnrollmentPricing(ctx, req.InternshipBatchID, []string{discountCode})
	if err != nil {
		return nil, err
	}

	// Create enrollment
	enrollmentData := &domainInternshipEnrollment.InternshipEnrollment{
		UserID:           types.GetUserID(ctx),
		InternshipID:     batch.InternshipID,
		EnrollmentStatus: types.InternshipEnrollmentStatusPending,
		PaymentStatus:    types.PaymentStatusPending,
		IdempotencyKey:   &idempotencyKey,
		Metadata:         req.Metadata,
	}

	// If payment is not required, mark as completed
	if !pricingResponse.PaymentRequired {
		enrollmentData.EnrollmentStatus = types.InternshipEnrollmentStatusCompleted
		enrollmentData.PaymentStatus = types.PaymentStatusSuccess
	}

	err = s.ServiceParams.InternshipEnrollmentRepo.Create(ctx, enrollmentData)
	if err != nil {
		return nil, err
	}

	return &dto.InitializeEnrollmentResponse{
		EnrollmentID:     enrollmentData.ID,
		EnrollmentStatus: enrollmentData.EnrollmentStatus,
		PaymentRequired:  pricingResponse.PaymentRequired,
	}, nil
}
