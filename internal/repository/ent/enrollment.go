package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/enrollment"
	domainEnrollment "github.com/omkar273/codegeeky/internal/domain/enrollment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type enrollmentRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts EnrollmentQueryOptions
}

func NewEnrollmentRepository(client postgres.IClient, logger *logger.Logger) domainEnrollment.Repository {
	return &enrollmentRepository{
		client:    client,
		log:       *logger,
		queryOpts: EnrollmentQueryOptions{},
	}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollmentData *domainEnrollment.Enrollment) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating enrollment",
		"enrollment_id", enrollmentData.ID,
		"user_id", enrollmentData.UserID,
		"internship_id", enrollmentData.InternshipID,
	)

	_, err := client.Enrollment.Create().
		SetID(enrollmentData.ID).
		SetUserID(enrollmentData.UserID).
		SetInternshipID(enrollmentData.InternshipID).
		SetEnrollmentStatus(enrollmentData.EnrollmentStatus).
		SetPaymentStatus(enrollmentData.PaymentStatus).
		SetNillableEnrolledAt(enrollmentData.EnrolledAt).
		SetNillablePaymentID(enrollmentData.PaymentID).
		SetNillableRefundedAt(enrollmentData.RefundedAt).
		SetNillableCancellationReason(enrollmentData.CancellationReason).
		SetNillableRefundReason(enrollmentData.RefundReason).
		SetMetadata(enrollmentData.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(enrollmentData.CreatedAt).
		SetUpdatedAt(enrollmentData.UpdatedAt).
		SetCreatedBy(enrollmentData.CreatedBy).
		SetUpdatedBy(enrollmentData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Enrollment with this user and internship combination already exists").
				WithReportableDetails(map[string]any{
					"enrollment_id": enrollmentData.ID,
					"user_id":       enrollmentData.UserID,
					"internship_id": enrollmentData.InternshipID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create enrollment").
			WithReportableDetails(map[string]any{
				"enrollment_id": enrollmentData.ID,
				"user_id":       enrollmentData.UserID,
				"internship_id": enrollmentData.InternshipID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *enrollmentRepository) Get(ctx context.Context, id string) (*domainEnrollment.Enrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting enrollment", "enrollment_id", id)

	entEnrollment, err := client.Enrollment.Query().
		Where(enrollment.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Enrollment with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"enrollment_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get enrollment").
			WithReportableDetails(map[string]any{
				"enrollment_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainEnrollment.FromEnt(entEnrollment), nil
}

func (r *enrollmentRepository) Update(ctx context.Context, enrollmentData *domainEnrollment.Enrollment) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating enrollment",
		"enrollment_id", enrollmentData.ID,
		"user_id", enrollmentData.UserID,
		"internship_id", enrollmentData.InternshipID,
	)

	_, err := client.Enrollment.UpdateOneID(enrollmentData.ID).
		SetUserID(enrollmentData.UserID).
		SetInternshipID(enrollmentData.InternshipID).
		SetEnrollmentStatus(enrollmentData.EnrollmentStatus).
		SetPaymentStatus(enrollmentData.PaymentStatus).
		SetNillableEnrolledAt(enrollmentData.EnrolledAt).
		SetNillablePaymentID(enrollmentData.PaymentID).
		SetNillableRefundedAt(enrollmentData.RefundedAt).
		SetNillableCancellationReason(enrollmentData.CancellationReason).
		SetNillableRefundReason(enrollmentData.RefundReason).
		SetMetadata(enrollmentData.Metadata).
		SetUpdatedAt(time.Now().UTC()).
		SetNillableIdempotencyKey(enrollmentData.IdempotencyKey).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Enrollment with ID %s was not found", enrollmentData.ID).
				WithReportableDetails(map[string]any{
					"enrollment_id": enrollmentData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Enrollment with this user and internship combination already exists").
				WithReportableDetails(map[string]any{
					"enrollment_id": enrollmentData.ID,
					"user_id":       enrollmentData.UserID,
					"internship_id": enrollmentData.InternshipID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to update enrollment").
			WithReportableDetails(map[string]any{
				"enrollment_id": enrollmentData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *enrollmentRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting enrollment",
		"enrollment_id", id,
	)

	_, err := client.Enrollment.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Enrollment with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"enrollment_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete enrollment").
			WithReportableDetails(map[string]any{
				"enrollment_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *enrollmentRepository) Count(ctx context.Context, filter *types.EnrollmentFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting enrollments")

	query := client.Enrollment.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count enrollments").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *enrollmentRepository) List(ctx context.Context, filter *types.EnrollmentFilter) ([]*domainEnrollment.Enrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing enrollments",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.Enrollment.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	enrollments, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list enrollments").
			Mark(ierr.ErrDatabase)
	}

	return domainEnrollment.FromEntList(enrollments), nil
}

func (r *enrollmentRepository) ListAll(ctx context.Context, filter *types.EnrollmentFilter) ([]*domainEnrollment.Enrollment, error) {
	if filter == nil {
		filter = types.NewNoLimitEnrollmentFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	enrollments, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domainEnrollment.Enrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting enrollment by idempotency key", "idempotency_key", idempotencyKey)

	enrollment, err := client.Enrollment.Query().
		Where(enrollment.IdempotencyKeyEQ(idempotencyKey)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Enrollment with idempotency key %s was not found", idempotencyKey).
				WithReportableDetails(map[string]any{
					"idempotency_key": idempotencyKey,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, err
	}

	return domainEnrollment.FromEnt(enrollment), nil
}

// EnrollmentQuery type alias for better readability
type EnrollmentQuery = *ent.EnrollmentQuery

// EnrollmentQueryOptions implements query options for enrollment queries
type EnrollmentQueryOptions struct {
	QueryOptionsHelper
}

// Ensure EnrollmentQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[EnrollmentQuery, *types.EnrollmentFilter] = (*EnrollmentQueryOptions)(nil)

func (o EnrollmentQueryOptions) ApplyStatusFilter(query EnrollmentQuery, status string) EnrollmentQuery {
	if status == "" {
		return query.Where(enrollment.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(enrollment.Status(status))
}

func (o EnrollmentQueryOptions) ApplySortFilter(query EnrollmentQuery, field string, order string) EnrollmentQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o EnrollmentQueryOptions) ApplyPaginationFilter(query EnrollmentQuery, limit int, offset int) EnrollmentQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o EnrollmentQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return enrollment.FieldCreatedAt
	case "updated_at":
		return enrollment.FieldUpdatedAt
	case "user_id":
		return enrollment.FieldUserID
	case "internship_id":
		return enrollment.FieldInternshipID
	case "enrollment_status":
		return enrollment.FieldEnrollmentStatus
	case "payment_status":
		return enrollment.FieldPaymentStatus
	case "enrolled_at":
		return enrollment.FieldEnrolledAt
	case "payment_id":
		return enrollment.FieldPaymentID
	case "refunded_at":
		return enrollment.FieldRefundedAt
	case "cancellation_reason":
		return enrollment.FieldCancellationReason
	case "refund_reason":
		return enrollment.FieldRefundReason
	case "created_by":
		return enrollment.FieldCreatedBy
	case "updated_by":
		return enrollment.FieldUpdatedBy
	default:
		return field
	}
}

func (o EnrollmentQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query EnrollmentQuery,
	filter *types.EnrollmentFilter,
) EnrollmentQuery {
	if filter == nil {
		return query.Where(enrollment.StatusNotIn(string(types.StatusDeleted)))
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

func (o EnrollmentQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.EnrollmentFilter,
	query EnrollmentQuery,
) EnrollmentQuery {
	if f == nil {
		return query
	}

	// Apply internship IDs filter if specified
	if len(f.InternshipIDs) > 0 {
		query = query.Where(enrollment.InternshipIDIn(f.InternshipIDs...))
	}

	// Apply user ID filter if specified
	if f.UserID != "" {
		query = query.Where(enrollment.UserID(f.UserID))
	}

	// Apply enrollment status filter if specified
	if f.EnrollmentStatus != "" {
		query = query.Where(enrollment.EnrollmentStatus(f.EnrollmentStatus))
	}

	// Apply payment status filter if specified
	if f.PaymentStatus != "" {
		query = query.Where(enrollment.PaymentStatus(f.PaymentStatus))
	}

	// Apply enrollment IDs filter if specified
	if len(f.EnrollmentIDs) > 0 {
		query = query.Where(enrollment.IDIn(f.EnrollmentIDs...))
	}

	// Apply payment ID filter if specified
	if f.PaymentID != nil && *f.PaymentID != "" {
		query = query.Where(enrollment.PaymentID(lo.FromPtr(f.PaymentID)))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(enrollment.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(enrollment.CreatedAtLTE(*f.EndTime))
		}
	}

	return query
}
