package service

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainPayment "github.com/omkar273/codegeeky/internal/domain/payment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// PaymentService defines the interface for payment operations
type PaymentService interface {
	// Core payment CRUD operations
	Create(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*dto.PaymentResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdatePaymentRequest) (*dto.PaymentResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.PaymentFilter) (*dto.ListPaymentResponse, error)

	// Payment attempt operations
	CreateAttempt(ctx context.Context, req *dto.PaymentAttemptRequest) (*dto.PaymentAttemptResponse, error)
	GetAttempt(ctx context.Context, id string) (*dto.PaymentAttemptResponse, error)
	ListAttempts(ctx context.Context, paymentID string) ([]*dto.PaymentAttemptResponse, error)
	GetLatestAttempt(ctx context.Context, paymentID string) (*dto.PaymentAttemptResponse, error)

	// Status management
	UpdateStatus(ctx context.Context, paymentID string, status types.PaymentStatus, metadata map[string]string) (*dto.PaymentResponse, error)
	MarkAsSuccess(ctx context.Context, paymentID string, gatewayPaymentID *string, metadata map[string]string) (*dto.PaymentResponse, error)
	MarkAsFailed(ctx context.Context, paymentID string, errorMessage string, metadata map[string]string) (*dto.PaymentResponse, error)
	MarkAsRefunded(ctx context.Context, paymentID string, metadata map[string]string) (*dto.PaymentResponse, error)
}

type paymentService struct {
	ServiceParams ServiceParams
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(params ServiceParams) PaymentService {
	return &paymentService{
		ServiceParams: params,
	}
}

// Create creates a new payment record
func (s *paymentService) Create(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check for existing payment with same idempotency key
	if req.IdempotencyKey != "" {
		existingPayment, err := s.ServiceParams.PaymentRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
		if err != nil && !ierr.IsNotFound(err) {
			return nil, err
		}
		if existingPayment != nil {
			return &dto.PaymentResponse{
				Payment: *existingPayment,
			}, nil
		}
	}

	var createdPayment *domainPayment.Payment

	// Create payment in transaction
	err := s.ServiceParams.DB.WithTx(ctx, func(ctx context.Context) error {
		// Convert request to domain model
		payment := req.ToPayment(ctx)

		// Create payment in database
		if err := s.ServiceParams.PaymentRepo.Create(ctx, payment); err != nil {
			return err
		}

		// If tracking attempts, create initial attempt
		if payment.TrackAttempts {
			attemptReq := &dto.PaymentAttemptRequest{
				PaymentID:     payment.ID,
				PaymentStatus: types.PaymentStatusPending,
				Metadata:      types.MetadataFromEnt(payment.Metadata),
			}

			attempt := attemptReq.ToPaymentAttempt(ctx, 1)
			if err := s.ServiceParams.PaymentRepo.CreateAttempt(ctx, attempt); err != nil {
				s.ServiceParams.Logger.Errorw("failed to create initial payment attempt",
					"payment_id", payment.ID, "error", err)
				// Don't fail the payment creation for attempt creation failure
			}
		}

		createdPayment = payment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *createdPayment,
	}, nil
}

// GetByID retrieves a payment by its ID
func (s *paymentService) GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error) {
	payment, err := s.ServiceParams.PaymentRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *payment,
	}, nil
}

// GetByIdempotencyKey retrieves a payment by its idempotency key
func (s *paymentService) GetByIdempotencyKey(ctx context.Context, key string) (*dto.PaymentResponse, error) {
	payment, err := s.ServiceParams.PaymentRepo.GetByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *payment,
	}, nil
}

// Update updates an existing payment
func (s *paymentService) Update(ctx context.Context, id string, req *dto.UpdatePaymentRequest) (*dto.PaymentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var updatedPayment *domainPayment.Payment

	err := s.ServiceParams.DB.WithTx(ctx, func(ctx context.Context) error {
		// Get existing payment
		existingPayment, err := s.ServiceParams.PaymentRepo.Get(ctx, id)
		if err != nil {
			return err
		}

		// Update fields
		if req.PaymentStatus != nil {
			existingPayment.PaymentStatus = *req.PaymentStatus

			// Set status timestamps
			now := time.Now()
			switch *req.PaymentStatus {
			case types.PaymentStatusSuccess:
				existingPayment.SucceededAt = &now
			case types.PaymentStatusFailed:
				existingPayment.FailedAt = &now
			case types.PaymentStatusRefunded:
				existingPayment.RefundedAt = &now
			}
		}

		if req.GatewayPaymentID != nil {
			existingPayment.GatewayPaymentID = req.GatewayPaymentID
		}

		if req.ErrorMessage != nil {
			existingPayment.ErrorMessage = req.ErrorMessage
		}

		if req.Metadata != nil {
			existingPayment.Metadata = req.Metadata
		}

		// Update in database
		if err := s.ServiceParams.PaymentRepo.Update(ctx, existingPayment); err != nil {
			return err
		}

		updatedPayment = existingPayment
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *updatedPayment,
	}, nil
}

// Delete deletes a payment by its ID
func (s *paymentService) Delete(ctx context.Context, id string) error {
	// Verify payment exists
	_, err := s.ServiceParams.PaymentRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.ServiceParams.PaymentRepo.Delete(ctx, id)
}

// List retrieves a paginated list of payments
func (s *paymentService) List(ctx context.Context, filter *types.PaymentFilter) (*dto.ListPaymentResponse, error) {
	if filter == nil {
		filter = types.NewNoLimitPaymentFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	count, err := s.ServiceParams.PaymentRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	payments, err := s.ServiceParams.PaymentRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &dto.ListPaymentResponse{
		Items:      make([]*dto.PaymentResponse, len(payments)),
		Pagination: lo.ToPtr(types.NewPaginationResponse(count, filter.GetLimit(), filter.GetOffset())),
	}

	for i, payment := range payments {
		response.Items[i] = &dto.PaymentResponse{Payment: *payment}
	}

	return response, nil
}

// CreateAttempt creates a new payment attempt
func (s *paymentService) CreateAttempt(ctx context.Context, req *dto.PaymentAttemptRequest) (*dto.PaymentAttemptResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing attempts to determine attempt number
	attempts, err := s.ServiceParams.PaymentRepo.ListAttempts(ctx, req.PaymentID)
	if err != nil && !ierr.IsNotFound(err) {
		return nil, err
	}

	attemptNumber := len(attempts) + 1
	attempt := req.ToPaymentAttempt(ctx, attemptNumber)

	if err := s.ServiceParams.PaymentRepo.CreateAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	return &dto.PaymentAttemptResponse{
		Attempt: *attempt,
	}, nil
}

// GetAttempt retrieves a payment attempt by its ID
func (s *paymentService) GetAttempt(ctx context.Context, id string) (*dto.PaymentAttemptResponse, error) {
	attempt, err := s.ServiceParams.PaymentRepo.GetAttempt(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentAttemptResponse{
		Attempt: *attempt,
	}, nil
}

// ListAttempts retrieves all attempts for a payment
func (s *paymentService) ListAttempts(ctx context.Context, paymentID string) ([]*dto.PaymentAttemptResponse, error) {
	attempts, err := s.ServiceParams.PaymentRepo.ListAttempts(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.PaymentAttemptResponse, len(attempts))
	for i, attempt := range attempts {
		responses[i] = &dto.PaymentAttemptResponse{Attempt: *attempt}
	}

	return responses, nil
}

// GetLatestAttempt retrieves the latest attempt for a payment
func (s *paymentService) GetLatestAttempt(ctx context.Context, paymentID string) (*dto.PaymentAttemptResponse, error) {
	attempt, err := s.ServiceParams.PaymentRepo.GetLatestAttempt(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentAttemptResponse{
		Attempt: *attempt,
	}, nil
}

// UpdateStatus updates the payment status with optional metadata
func (s *paymentService) UpdateStatus(ctx context.Context, paymentID string, status types.PaymentStatus, metadata map[string]string) (*dto.PaymentResponse, error) {
	req := &dto.UpdatePaymentRequest{
		PaymentStatus: &status,
		Metadata:      metadata,
	}
	return s.Update(ctx, paymentID, req)
}

// MarkAsSuccess marks a payment as successful
func (s *paymentService) MarkAsSuccess(ctx context.Context, paymentID string, gatewayPaymentID *string, metadata map[string]string) (*dto.PaymentResponse, error) {
	status := types.PaymentStatusSuccess
	req := &dto.UpdatePaymentRequest{
		PaymentStatus:    &status,
		GatewayPaymentID: gatewayPaymentID,
		Metadata:         metadata,
	}
	return s.Update(ctx, paymentID, req)
}

// MarkAsFailed marks a payment as failed
func (s *paymentService) MarkAsFailed(ctx context.Context, paymentID string, errorMessage string, metadata map[string]string) (*dto.PaymentResponse, error) {
	status := types.PaymentStatusFailed
	req := &dto.UpdatePaymentRequest{
		PaymentStatus: &status,
		ErrorMessage:  &errorMessage,
		Metadata:      metadata,
	}
	return s.Update(ctx, paymentID, req)
}

// MarkAsRefunded marks a payment as refunded
func (s *paymentService) MarkAsRefunded(ctx context.Context, paymentID string, metadata map[string]string) (*dto.PaymentResponse, error) {
	status := types.PaymentStatusRefunded
	req := &dto.UpdatePaymentRequest{
		PaymentStatus: &status,
		Metadata:      metadata,
	}
	return s.Update(ctx, paymentID, req)
}
