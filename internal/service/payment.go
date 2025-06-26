package service

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainPayment "github.com/omkar273/codegeeky/internal/domain/payment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	gateway "github.com/omkar273/codegeeky/internal/payment"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type PaymentService interface {
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

	// Payment processing
	ProcessPayment(ctx context.Context, paymentID string) (*dto.PaymentResponse, error)
}

type paymentService struct {
	ServiceParams   ServiceParams
	GatewayRegistry gateway.GatewayRegistryService
}

func NewPaymentService(params ServiceParams, gatewayRegistry gateway.GatewayRegistryService) PaymentService {
	return &paymentService{
		ServiceParams:   params,
		GatewayRegistry: gatewayRegistry,
	}
}

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
	var gatewayResponse *dto.PaymentGatewayResponse

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

		// Process payment with gateway if available
		if s.GatewayRegistry != nil {
			provider, err := s.GatewayRegistry.GetProviderByName(ctx, payment.PaymentGatewayProvider)

			if err != nil {
				s.ServiceParams.Logger.Errorw("failed to get payment gateway provider",
					"provider", payment.PaymentGatewayProvider, "error", err)
				return err
			}

			gwReq := &dto.PaymentRequest{
				ReferenceID:            payment.ID,
				ReferenceType:          payment.DestinationType,
				DestinationID:          payment.DestinationID,
				DestinationType:        payment.DestinationType,
				Amount:                 payment.Amount.IntPart(),
				Currency:               string(payment.Currency),
				PaymentGatewayProvider: payment.PaymentGatewayProvider,
				PaymentMethodType:      payment.PaymentMethodType,
				PaymentMethodID:        &payment.PaymentMethodID,
				SuccessURL:             req.SuccessURL,
				CancelURL:              req.CancelURL,
				IdempotencyKey:         payment.IdempotencyKey,
				Metadata:               payment.Metadata,
				TrackAttempts:          payment.TrackAttempts,
			}

			gwResp, err := provider.CreatePaymentOrder(ctx, gwReq)
			if err != nil {
				s.ServiceParams.Logger.Errorw("failed to create payment with gateway",
					"payment_id", payment.ID, "error", err)

				// Update payment status to failed
				payment.PaymentStatus = types.PaymentStatusFailed
				payment.ErrorMessage = lo.ToPtr(err.Error())
				payment.FailedAt = lo.ToPtr(time.Now())

				if updateErr := s.ServiceParams.PaymentRepo.Update(ctx, payment); updateErr != nil {
					s.ServiceParams.Logger.Errorw("failed to update payment status to failed",
						"payment_id", payment.ID, "error", updateErr)
				}

				return err
			}

			// Update payment with gateway response
			if gwResp != nil {
				payment.GatewayPaymentID = lo.ToPtr(gwResp.GatewayResponse.ProviderPaymentID)
				if err := s.ServiceParams.PaymentRepo.Update(ctx, payment); err != nil {
					s.ServiceParams.Logger.Errorw("failed to update payment with gateway ID",
						"payment_id", payment.ID, "error", err)
				}

				gatewayResponse = &dto.PaymentGatewayResponse{
					ProviderPaymentID: gwResp.GatewayResponse.ProviderPaymentID,
					RedirectURL:       gwResp.GatewayResponse.RedirectURL,
					Status:            gwResp.GatewayResponse.Status,
					Raw:               gwResp.GatewayResponse.Raw,
				}
			}
		}

		createdPayment = payment
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish webhook event
	// if s.ServiceParams.WebhookPublisher != nil {
	// 	go func() {
	// 		ctx := context.Background()
	// 		webhookEvent := &types.WebhookEvent{
	// 			ID:        types.GenerateUUID(),
	// 			EventName: "payment.created",
	// 			UserID:    lo.ToPtr(createdPayment.CreatedBy),
	// 			Payload:   lo.ToPtr(createdPayment),
	// 		}

	// 		if err := s.ServiceParams.WebhookPublisher.PublishWebhook(ctx, webhookEvent); err != nil {
	// 			s.ServiceParams.Logger.Errorw("failed to publish payment created webhook",
	// 				"payment_id", createdPayment.ID, "error", err)
	// 		}
	// 	}()
	// }

	return &dto.PaymentResponse{
		Payment:         *createdPayment,
		GatewayResponse: gatewayResponse,
	}, nil
}

func (s *paymentService) GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error) {
	payment, err := s.ServiceParams.PaymentRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *payment,
	}, nil
}

func (s *paymentService) GetByIdempotencyKey(ctx context.Context, key string) (*dto.PaymentResponse, error) {
	payment, err := s.ServiceParams.PaymentRepo.GetByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		Payment: *payment,
	}, nil
}

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

	// Publish webhook event for status changes
	// if req.PaymentStatus != nil && s.ServiceParams.WebhookPublisher != nil {
	// 	go func() {
	// 		ctx := context.Background()
	// 		webhookEvent := &types.WebhookEvent{
	// 			ID:        types.GenerateUUID(),
	// 			EventName: "payment.status_updated",
	// 			UserID:    lo.ToPtr(updatedPayment.UpdatedBy),
	// 			Payload:   lo.ToPtr(updatedPayment),
	// 		}

	// 		if err := s.ServiceParams.WebhookPublisher.PublishWebhook(ctx, webhookEvent); err != nil {
	// 			s.ServiceParams.Logger.Errorw("failed to publish payment status updated webhook",
	// 				"payment_id", updatedPayment.ID, "error", err)
	// 		}
	// 	}()
	// }

	return &dto.PaymentResponse{
		Payment: *updatedPayment,
	}, nil
}

func (s *paymentService) Delete(ctx context.Context, id string) error {
	// Verify payment exists
	_, err := s.ServiceParams.PaymentRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.ServiceParams.PaymentRepo.Delete(ctx, id)
}

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

// Payment attempt operations
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

func (s *paymentService) GetAttempt(ctx context.Context, id string) (*dto.PaymentAttemptResponse, error) {
	attempt, err := s.ServiceParams.PaymentRepo.GetAttempt(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentAttemptResponse{
		Attempt: *attempt,
	}, nil
}

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

func (s *paymentService) GetLatestAttempt(ctx context.Context, paymentID string) (*dto.PaymentAttemptResponse, error) {
	attempt, err := s.ServiceParams.PaymentRepo.GetLatestAttempt(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentAttemptResponse{
		Attempt: *attempt,
	}, nil
}

func (s *paymentService) ProcessPayment(ctx context.Context, paymentID string) (*dto.PaymentResponse, error) {
	payment, err := s.ServiceParams.PaymentRepo.Get(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	// Only process pending payments
	if payment.PaymentStatus != types.PaymentStatusPending {
		return nil, ierr.NewError("payment cannot be processed").
			WithHint("Only pending payments can be processed").
			WithReportableDetails(map[string]any{
				"payment_id":     paymentID,
				"current_status": payment.PaymentStatus,
			}).
			Mark(ierr.ErrValidation)
	}

	// Process with payment gateway if available
	if s.GatewayRegistry != nil && payment.GatewayPaymentID != nil {
		provider, err := s.GatewayRegistry.GetProviderByName(ctx, payment.PaymentGatewayProvider)
		if err != nil {
			return nil, err
		}

		gwResp, err := provider.VerifyPaymentStatus(ctx, *payment.GatewayPaymentID)
		if err != nil {
			// Update payment status to failed
			payment.PaymentStatus = types.PaymentStatusFailed
			payment.ErrorMessage = lo.ToPtr(err.Error())
			payment.FailedAt = lo.ToPtr(time.Now())

			if updateErr := s.ServiceParams.PaymentRepo.Update(ctx, payment); updateErr != nil {
				s.ServiceParams.Logger.Errorw("failed to update payment status to failed",
					"payment_id", payment.ID, "error", updateErr)
			}

			return nil, err
		}

		// Update payment status based on gateway response
		if gwResp != nil {
			switch gwResp.Status {
			case "success":
				payment.PaymentStatus = types.PaymentStatusSuccess
				payment.SucceededAt = lo.ToPtr(time.Now())
			case "failed":
				payment.PaymentStatus = types.PaymentStatusFailed
				payment.FailedAt = lo.ToPtr(time.Now())
				if gwResp.Reason != "" {
					payment.ErrorMessage = lo.ToPtr(gwResp.Reason)
				}
			case "processing":
				payment.PaymentStatus = types.PaymentStatusProcessing
			}

			if err := s.ServiceParams.PaymentRepo.Update(ctx, payment); err != nil {
				return nil, err
			}
		}
	}

	return &dto.PaymentResponse{
		Payment: *payment,
	}, nil
}
