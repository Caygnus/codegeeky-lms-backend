package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/payment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryPaymentStore implements payment.Repository
type InMemoryPaymentStore struct {
	*InMemoryStore[*payment.Payment]
	attempts *InMemoryStore[*payment.PaymentAttempt]
}

// NewInMemoryPaymentStore creates a new in-memory payment store
func NewInMemoryPaymentStore() *InMemoryPaymentStore {
	return &InMemoryPaymentStore{
		InMemoryStore: NewInMemoryStore[*payment.Payment](),
		attempts:      NewInMemoryStore[*payment.PaymentAttempt](),
	}
}

// paymentFilterFn implements filtering logic for payments
func paymentFilterFn(ctx context.Context, p *payment.Payment, filter interface{}) bool {
	if p == nil {
		return false
	}

	filter_, ok := filter.(*types.PaymentFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by destination type
	if filter_.DestinationType != nil {
		if string(p.DestinationType) != *filter_.DestinationType {
			return false
		}
	}

	// Filter by destination ID
	if filter_.DestinationID != nil {
		if p.DestinationID != *filter_.DestinationID {
			return false
		}
	}

	// Filter by payment method type
	if filter_.PaymentMethodType != nil {
		if p.PaymentMethodType == nil || string(*p.PaymentMethodType) != *filter_.PaymentMethodType {
			return false
		}
	}

	// Filter by payment gateway
	if filter_.PaymentGateway != nil {
		if string(p.PaymentGatewayProvider) != *filter_.PaymentGateway {
			return false
		}
	}

	// Filter by payment status
	if filter_.PaymentStatus != nil {
		if string(p.PaymentStatus) != *filter_.PaymentStatus {
			return false
		}
	}

	// Filter by currency
	if filter_.Currency != nil {
		if string(p.Currency) != *filter_.Currency {
			return false
		}
	}

	// Filter by payment IDs
	if len(filter_.PaymentIDs) > 0 {
		if !lo.Contains(filter_.PaymentIDs, p.ID) {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active payments
	if filter_.GetStatus() != "" {
		if string(p.Status) != filter_.GetStatus() {
			return false
		}
	} else if string(p.Status) == string(types.StatusDeleted) {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && p.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && p.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// paymentSortFn implements sorting logic for payments
func paymentSortFn(i, j *payment.Payment) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryPaymentStore) Create(ctx context.Context, p *payment.Payment) error {
	if p == nil {
		return ierr.NewError("payment cannot be nil").
			WithHint("Payment data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, p.ID, p)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A payment with this ID already exists").
				WithReportableDetails(map[string]any{
					"payment_id": p.ID,
					"amount":     p.Amount,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create payment").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryPaymentStore) Get(ctx context.Context, id string) (*payment.Payment, error) {
	payment, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Payment with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"payment_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get payment with ID %s", id).
			Mark(ierr.ErrDatabase)
	}

	// Load attempts
	attempts, err := s.ListAttempts(ctx, id)
	if err != nil {
		return nil, err
	}
	payment.Attempts = attempts

	return payment, nil
}

func (s *InMemoryPaymentStore) Update(ctx context.Context, p *payment.Payment) error {
	if p == nil {
		return ierr.NewError("payment cannot be nil").
			WithHint("Payment data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	p.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, p.ID, p)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Payment with ID %s was not found", p.ID).
				WithReportableDetails(map[string]any{
					"payment_id": p.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update payment with ID %s", p.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryPaymentStore) Delete(ctx context.Context, id string) error {
	// Get the payment first
	p, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	p.Status = string(types.StatusDeleted)
	p.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, p)
}

func (s *InMemoryPaymentStore) List(ctx context.Context, filter *types.PaymentFilter) ([]*payment.Payment, error) {
	payments, err := s.InMemoryStore.List(ctx, filter, paymentFilterFn, paymentSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list payments").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return payments, nil
}

func (s *InMemoryPaymentStore) Count(ctx context.Context, filter *types.PaymentFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, paymentFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count payments").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryPaymentStore) GetByIdempotencyKey(ctx context.Context, key string) (*payment.Payment, error) {
	payments, err := s.InMemoryStore.List(ctx, nil, paymentFilterFn, paymentSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get payment by idempotency key").
			Mark(ierr.ErrDatabase)
	}

	for _, p := range payments {
		if p.IdempotencyKey == key && string(p.Status) != string(types.StatusDeleted) {
			return p, nil
		}
	}

	return nil, ierr.NewError("payment not found").
		WithHintf("Payment with idempotency key %s was not found", key).
		WithReportableDetails(map[string]any{
			"idempotency_key": key,
		}).
		Mark(ierr.ErrNotFound)
}

// Payment Attempt methods

func (s *InMemoryPaymentStore) CreateAttempt(ctx context.Context, attempt *payment.PaymentAttempt) error {
	if attempt == nil {
		return ierr.NewError("payment attempt cannot be nil").
			WithHint("Payment attempt data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if attempt.CreatedAt.IsZero() {
		attempt.CreatedAt = now
	}
	if attempt.UpdatedAt.IsZero() {
		attempt.UpdatedAt = now
	}

	err := s.attempts.Create(ctx, attempt.ID, attempt)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A payment attempt with this ID already exists").
				WithReportableDetails(map[string]any{
					"attempt_id": attempt.ID,
					"payment_id": attempt.PaymentID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create payment attempt").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryPaymentStore) GetAttempt(ctx context.Context, id string) (*payment.PaymentAttempt, error) {
	attempt, err := s.attempts.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Payment attempt with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"attempt_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get payment attempt with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return attempt, nil
}

func (s *InMemoryPaymentStore) UpdateAttempt(ctx context.Context, attempt *payment.PaymentAttempt) error {
	if attempt == nil {
		return ierr.NewError("payment attempt cannot be nil").
			WithHint("Payment attempt data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	attempt.UpdatedAt = time.Now().UTC()

	err := s.attempts.Update(ctx, attempt.ID, attempt)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Payment attempt with ID %s was not found", attempt.ID).
				WithReportableDetails(map[string]any{
					"attempt_id": attempt.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update payment attempt with ID %s", attempt.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryPaymentStore) ListAttempts(ctx context.Context, paymentID string) ([]*payment.PaymentAttempt, error) {
	allAttempts, err := s.attempts.List(ctx, nil, nil, nil)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list payment attempts").
			Mark(ierr.ErrDatabase)
	}

	// Filter by payment ID
	attempts := lo.Filter(allAttempts, func(attempt *payment.PaymentAttempt, _ int) bool {
		return attempt.PaymentID == paymentID && attempt.Status != types.StatusDeleted
	})

	return attempts, nil
}

func (s *InMemoryPaymentStore) GetLatestAttempt(ctx context.Context, paymentID string) (*payment.PaymentAttempt, error) {
	attempts, err := s.ListAttempts(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if len(attempts) == 0 {
		return nil, ierr.NewError("no payment attempts found").
			WithHintf("No payment attempts found for payment ID %s", paymentID).
			WithReportableDetails(map[string]any{
				"payment_id": paymentID,
			}).
			Mark(ierr.ErrNotFound)
	}

	// Return the attempt with the highest attempt number
	latestAttempt := attempts[0]
	for _, attempt := range attempts {
		if attempt.AttemptNumber > latestAttempt.AttemptNumber {
			latestAttempt = attempt
		}
	}

	return latestAttempt, nil
}

// Clear clears both payment and attempt stores
func (s *InMemoryPaymentStore) Clear() {
	s.InMemoryStore.Clear()
	s.attempts.Clear()
}
