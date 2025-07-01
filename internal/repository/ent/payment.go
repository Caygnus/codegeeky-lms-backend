package ent

import (
	"context"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/payment"
	"github.com/omkar273/codegeeky/ent/paymentattempt"
	domainPayment "github.com/omkar273/codegeeky/internal/domain/payment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type paymentRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts PaymentQueryOptions
}

func NewPaymentRepository(client postgres.IClient, logger *logger.Logger) domainPayment.Repository {
	return &paymentRepository{
		client:    client,
		log:       *logger,
		queryOpts: PaymentQueryOptions{},
	}
}

func (r *paymentRepository) Count(ctx context.Context, filter *types.PaymentFilter) (int, error) {
	client := r.client.Querier(ctx)
	query := client.Payment.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count payments").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *paymentRepository) Create(ctx context.Context, p *domainPayment.Payment) error {
	if err := p.Validate(); err != nil {
		return err
	}

	client := r.client.Querier(ctx)

	builder := client.Payment.Create().
		SetID(p.ID).
		SetIdempotencyKey(p.IdempotencyKey).
		SetDestinationType(p.DestinationType).
		SetDestinationID(p.DestinationID).
		SetPaymentMethodID(p.PaymentMethodID).
		SetPaymentGatewayProvider(p.PaymentGatewayProvider).
		SetAmount(p.Amount).
		SetCurrency(p.Currency).
		SetPaymentStatus(p.PaymentStatus).
		SetTrackAttempts(p.TrackAttempts).
		SetMetadata(p.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(p.CreatedAt).
		SetUpdatedAt(p.UpdatedAt).
		SetCreatedBy(p.CreatedBy).
		SetNillableGatewayPaymentID(p.GatewayPaymentID).
		SetNillableSucceededAt(p.SucceededAt).
		SetNillableFailedAt(p.FailedAt).
		SetNillableRefundedAt(p.RefundedAt).
		SetNillableErrorMessage(p.ErrorMessage).
		SetUpdatedBy(p.UpdatedBy).
		SetNillablePaymentMethodType(p.PaymentMethodType).
		SetNillablePaymentMethodID(lo.ToPtr(p.PaymentMethodID))

	if _, err := builder.Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Payment with this idempotency key already exists").
				WithReportableDetails(map[string]any{
					"idempotency_key": p.IdempotencyKey,
				}).
				Mark(ierr.ErrAlreadyExists)
		}

		return ierr.WithError(err).
			WithHint("Failed to create payment").
			WithReportableDetails(map[string]any{
				"payment": p,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *paymentRepository) Get(ctx context.Context, id string) (*domainPayment.Payment, error) {
	client := r.client.Querier(ctx)
	payment, err := client.Payment.Query().
		Where(payment.ID(id)).
		WithAttempts().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Payment with this ID not found").
				WithReportableDetails(map[string]any{
					"payment_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}

		return nil, ierr.WithError(err).
			WithHint("Failed to get payment").
			WithReportableDetails(map[string]any{
				"payment_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEnt(payment), nil
}

func (r *paymentRepository) List(ctx context.Context, filter *types.PaymentFilter) ([]*domainPayment.Payment, error) {
	client := r.client.Querier(ctx)
	query := client.Payment.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	payments, err := query.WithAttempts().All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Payments not found").
				WithReportableDetails(map[string]any{
					"filter": filter,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get payments").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEntList(payments), nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *domainPayment.Payment) error {
	if err := payment.Validate(); err != nil {
		return err
	}

	client := r.client.Querier(ctx)

	builder := client.Payment.UpdateOneID(payment.ID).
		SetPaymentStatus(payment.PaymentStatus).
		SetMetadata(payment.Metadata).
		SetUpdatedBy(types.GetUserID(ctx))

	if payment.PaymentMethodType != nil {
		builder = builder.SetPaymentMethodType(*payment.PaymentMethodType)
	}

	if payment.GatewayPaymentID != nil {
		builder = builder.SetGatewayPaymentID(*payment.GatewayPaymentID)
	}

	if payment.SucceededAt != nil {
		builder = builder.SetSucceededAt(*payment.SucceededAt)
	}

	if payment.FailedAt != nil {
		builder = builder.SetFailedAt(*payment.FailedAt)
	}

	if payment.RefundedAt != nil {
		builder = builder.SetRefundedAt(*payment.RefundedAt)
	}

	if payment.ErrorMessage != nil {
		builder = builder.SetErrorMessage(*payment.ErrorMessage)
	}

	_, err := builder.Save(ctx)

	if err != nil {
		return ierr.WithError(err).
			WithHint("Failed to update payment").
			WithReportableDetails(map[string]any{
				"payment_id": payment.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *paymentRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)
	_, err := client.Payment.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		return ierr.WithError(err).
			WithHint("Failed to delete payment").
			WithReportableDetails(map[string]any{
				"payment_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *paymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domainPayment.Payment, error) {
	client := r.client.Querier(ctx)
	payment, err := client.Payment.Query().
		Where(payment.IdempotencyKey(key)).
		WithAttempts().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Payment with this idempotency key not found").
				WithReportableDetails(map[string]any{
					"idempotency_key": key,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get payment").
			WithReportableDetails(map[string]any{
				"idempotency_key": key,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEnt(payment), nil
}

// Payment attempt operations
func (r *paymentRepository) CreateAttempt(ctx context.Context, attempt *domainPayment.PaymentAttempt) error {
	if err := attempt.Validate(); err != nil {
		return err
	}

	client := r.client.Querier(ctx)

	builder := client.PaymentAttempt.Create().
		SetID(attempt.ID).
		SetPaymentID(attempt.PaymentID).
		SetAttemptNumber(attempt.AttemptNumber).
		SetPaymentStatus(string(attempt.PaymentStatus)).
		SetMetadata(attempt.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(attempt.CreatedAt).
		SetUpdatedAt(attempt.UpdatedAt).
		SetCreatedBy(attempt.CreatedBy).
		SetUpdatedBy(attempt.UpdatedBy)

	if attempt.GatewayAttemptID != nil {
		builder = builder.SetGatewayAttemptID(*attempt.GatewayAttemptID)
	}

	if attempt.ErrorMessage != nil {
		builder = builder.SetErrorMessage(*attempt.ErrorMessage)
	}

	_, err := builder.Save(ctx)

	if err != nil {
		return ierr.WithError(err).
			WithHint("Failed to create payment attempt").
			WithReportableDetails(map[string]any{
				"attempt": attempt,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *paymentRepository) GetAttempt(ctx context.Context, id string) (*domainPayment.PaymentAttempt, error) {
	client := r.client.Querier(ctx)
	attempt, err := client.PaymentAttempt.Query().
		Where(paymentattempt.ID(id)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Payment attempt with this ID not found").
				WithReportableDetails(map[string]any{
					"attempt_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}

		return nil, ierr.WithError(err).
			WithHint("Failed to get payment attempt").
			WithReportableDetails(map[string]any{
				"attempt_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEntAttempt(attempt), nil
}

func (r *paymentRepository) UpdateAttempt(ctx context.Context, attempt *domainPayment.PaymentAttempt) error {
	if err := attempt.Validate(); err != nil {
		return err
	}

	client := r.client.Querier(ctx)

	builder := client.PaymentAttempt.UpdateOneID(attempt.ID).
		SetPaymentStatus(string(attempt.PaymentStatus)).
		SetMetadata(attempt.Metadata).
		SetUpdatedBy(types.GetUserID(ctx))

	if attempt.GatewayAttemptID != nil {
		builder = builder.SetGatewayAttemptID(*attempt.GatewayAttemptID)
	}

	if attempt.ErrorMessage != nil {
		builder = builder.SetErrorMessage(*attempt.ErrorMessage)
	}

	_, err := builder.Save(ctx)

	if err != nil {
		return ierr.WithError(err).
			WithHint("Failed to update payment attempt").
			WithReportableDetails(map[string]any{
				"attempt_id": attempt.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *paymentRepository) ListAttempts(ctx context.Context, paymentID string) ([]*domainPayment.PaymentAttempt, error) {
	client := r.client.Querier(ctx)
	attempts, err := client.PaymentAttempt.Query().
		Where(paymentattempt.PaymentID(paymentID)).
		Where(paymentattempt.StatusNotIn(string(types.StatusDeleted))).
		Order(ent.Desc(paymentattempt.FieldAttemptNumber)).
		All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Payment attempts not found").
				WithReportableDetails(map[string]any{
					"payment_id": paymentID,
				}).
				Mark(ierr.ErrNotFound)
		}

		return nil, ierr.WithError(err).
			WithHint("Failed to get payment attempts").
			WithReportableDetails(map[string]any{
				"payment_id": paymentID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEntAttemptList(attempts), nil
}

func (r *paymentRepository) GetLatestAttempt(ctx context.Context, paymentID string) (*domainPayment.PaymentAttempt, error) {
	client := r.client.Querier(ctx)
	attempt, err := client.PaymentAttempt.Query().
		Where(paymentattempt.PaymentID(paymentID)).
		Where(paymentattempt.StatusNotIn(string(types.StatusDeleted))).
		Order(ent.Desc(paymentattempt.FieldAttemptNumber)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Latest payment attempt not found").
				WithReportableDetails(map[string]any{
					"payment_id": paymentID,
				}).
				Mark(ierr.ErrNotFound)
		}

		return nil, ierr.WithError(err).
			WithHint("Failed to get latest payment attempt").
			WithReportableDetails(map[string]any{
				"payment_id": paymentID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainPayment.FromEntAttempt(attempt), nil
}

// PaymentQuery type alias for better readability
type PaymentQuery = *ent.PaymentQuery

// PaymentQueryOptions implements query options for payment queries
type PaymentQueryOptions struct {
	QueryOptionsHelper
}

// Ensure PaymentQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[PaymentQuery, *types.PaymentFilter] = (*PaymentQueryOptions)(nil)

func (o PaymentQueryOptions) ApplyStatusFilter(query PaymentQuery, status string) PaymentQuery {
	if status == "" {
		return query.Where(payment.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(payment.Status(status))
}

func (o PaymentQueryOptions) ApplySortFilter(query PaymentQuery, field string, order string) PaymentQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o PaymentQueryOptions) ApplyPaginationFilter(query PaymentQuery, limit int, offset int) PaymentQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o PaymentQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return payment.FieldCreatedAt
	case "updated_at":
		return payment.FieldUpdatedAt
	case "idempotency_key":
		return payment.FieldIdempotencyKey
	case "destination_type":
		return payment.FieldDestinationType
	case "destination_id":
		return payment.FieldDestinationID
	case "payment_method_type":
		return payment.FieldPaymentMethodType
	case "payment_method_id":
		return payment.FieldPaymentMethodID
	case "payment_gateway_provider":
		return payment.FieldPaymentGatewayProvider
	case "gateway_payment_id":
		return payment.FieldGatewayPaymentID
	case "amount":
		return payment.FieldAmount
	case "currency":
		return payment.FieldCurrency
	case "payment_status":
		return payment.FieldPaymentStatus
	case "track_attempts":
		return payment.FieldTrackAttempts
	case "succeeded_at":
		return payment.FieldSucceededAt
	case "failed_at":
		return payment.FieldFailedAt
	case "refunded_at":
		return payment.FieldRefundedAt
	case "error_message":
		return payment.FieldErrorMessage
	case "created_by":
		return payment.FieldCreatedBy
	case "updated_by":
		return payment.FieldUpdatedBy
	default:
		return ""
	}
}

func (o PaymentQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query PaymentQuery,
	filter *types.PaymentFilter,
) PaymentQuery {
	if filter == nil {
		return query.Where(payment.StatusNotIn(string(types.StatusDeleted)))
	}

	// Apply status filter
	query = o.ApplyStatusFilter(query, filter.GetStatus())

	// Apply pagination
	if !filter.IsUnlimited() {
		query = o.ApplyPaginationFilter(query, filter.GetLimit(), filter.GetOffset())
	}

	// Apply sorting
	query = o.ApplySortFilter(query, filter.GetSort(), filter.GetOrder())

	return query
}

func (o PaymentQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.PaymentFilter,
	query PaymentQuery,
) PaymentQuery {
	if f == nil {
		return query
	}

	if f.PaymentIDs != nil {
		query = query.Where(payment.IDIn(f.PaymentIDs...))
	}

	if f.DestinationType != nil {
		query = query.Where(payment.DestinationType(types.PaymentDestinationType(*f.DestinationType)))
	}

	if f.DestinationID != nil {
		query = query.Where(payment.DestinationID(*f.DestinationID))
	}

	if f.PaymentMethodType != nil {
		query = query.Where(payment.PaymentMethodType(types.PaymentMethodType(*f.PaymentMethodType)))
	}

	if f.PaymentStatus != nil {
		query = query.Where(payment.PaymentStatus(types.PaymentStatus(*f.PaymentStatus)))
	}

	if f.PaymentGateway != nil {
		query = query.Where(payment.PaymentGatewayProvider(types.PaymentGatewayProvider(*f.PaymentGateway)))
	}

	if f.Currency != nil {
		query = query.Where(payment.Currency(types.Currency(*f.Currency)))
	}

	return query
}
