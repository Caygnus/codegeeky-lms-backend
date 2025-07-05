package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/internshipenrollment"
	domainInternshipEnrollment "github.com/omkar273/codegeeky/internal/domain/internshipenrollment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type internshipEnrollmentRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts EnrollmentQueryOptions
}

func NewInternshipEnrollmentRepository(client postgres.IClient, logger *logger.Logger) domainInternshipEnrollment.Repository {
	return &internshipEnrollmentRepository{
		client:    client,
		log:       *logger,
		queryOpts: EnrollmentQueryOptions{},
	}
}

func (r *internshipEnrollmentRepository) Create(ctx context.Context, enrollmentData *domainInternshipEnrollment.InternshipEnrollment) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating enrollment",
		"enrollment_id", enrollmentData.ID,
		"user_id", enrollmentData.UserID,
		"internship_id", enrollmentData.InternshipID,
	)

	_, err := client.InternshipEnrollment.Create().
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

func (r *internshipEnrollmentRepository) Get(ctx context.Context, id string) (*domainInternshipEnrollment.InternshipEnrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting enrollment", "enrollment_id", id)

	entEnrollment, err := client.InternshipEnrollment.Query().
		Where(internshipenrollment.ID(id)).
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

	return domainInternshipEnrollment.FromEnt(entEnrollment), nil
}

func (r *internshipEnrollmentRepository) Update(ctx context.Context, enrollmentData *domainInternshipEnrollment.InternshipEnrollment) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating enrollment",
		"enrollment_id", enrollmentData.ID,
		"user_id", enrollmentData.UserID,
		"internship_id", enrollmentData.InternshipID,
	)

	_, err := client.InternshipEnrollment.UpdateOneID(enrollmentData.ID).
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

func (r *internshipEnrollmentRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting enrollment",
		"enrollment_id", id,
	)

	_, err := client.InternshipEnrollment.UpdateOneID(id).
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

func (r *internshipEnrollmentRepository) Count(ctx context.Context, filter *types.InternshipEnrollmentFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting enrollments")

	query := client.InternshipEnrollment.Query()
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

func (r *internshipEnrollmentRepository) List(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*domainInternshipEnrollment.InternshipEnrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing enrollments",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.InternshipEnrollment.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	enrollments, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list enrollments").
			Mark(ierr.ErrDatabase)
	}

	return domainInternshipEnrollment.FromEntList(enrollments), nil
}

func (r *internshipEnrollmentRepository) ListAll(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*domainInternshipEnrollment.InternshipEnrollment, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipEnrollmentFilter()
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

func (r *internshipEnrollmentRepository) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domainInternshipEnrollment.InternshipEnrollment, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting enrollment by idempotency key", "idempotency_key", idempotencyKey)

	enrollment, err := client.InternshipEnrollment.Query().
		Where(internshipenrollment.IdempotencyKeyEQ(idempotencyKey)).
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

	return domainInternshipEnrollment.FromEnt(enrollment), nil
}

// EnrollmentQuery type alias for better readability
type EnrollmentQuery = *ent.InternshipEnrollmentQuery

// EnrollmentQueryOptions implements query options for enrollment queries
type EnrollmentQueryOptions struct {
	QueryOptionsHelper
}

// Ensure EnrollmentQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[EnrollmentQuery, *types.InternshipEnrollmentFilter] = (*EnrollmentQueryOptions)(nil)

func (o EnrollmentQueryOptions) ApplyStatusFilter(query EnrollmentQuery, status string) EnrollmentQuery {
	if status == "" {
		return query.Where(internshipenrollment.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(internshipenrollment.Status(status))
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
		return internshipenrollment.FieldCreatedAt
	case "updated_at":
		return internshipenrollment.FieldUpdatedAt
	case "user_id":
		return internshipenrollment.FieldUserID
	case "internship_id":
		return internshipenrollment.FieldInternshipID
	case "enrollment_status":
		return internshipenrollment.FieldEnrollmentStatus
	case "payment_status":
		return internshipenrollment.FieldPaymentStatus
	case "enrolled_at":
		return internshipenrollment.FieldEnrolledAt
	case "payment_id":
		return internshipenrollment.FieldPaymentID
	case "refunded_at":
		return internshipenrollment.FieldRefundedAt
	case "cancellation_reason":
		return internshipenrollment.FieldCancellationReason
	case "refund_reason":
		return internshipenrollment.FieldRefundReason
	case "created_by":
		return internshipenrollment.FieldCreatedBy
	case "updated_by":
		return internshipenrollment.FieldUpdatedBy
	default:
		return field
	}
}

func (o EnrollmentQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query EnrollmentQuery,
	filter *types.InternshipEnrollmentFilter,
) EnrollmentQuery {
	if filter == nil {
		return query.Where(internshipenrollment.StatusNotIn(string(types.StatusDeleted)))
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
	f *types.InternshipEnrollmentFilter,
	query EnrollmentQuery,
) EnrollmentQuery {
	if f == nil {
		return query
	}

	// Apply internship IDs filter if specified
	if len(f.InternshipIDs) > 0 {
		query = query.Where(internshipenrollment.InternshipIDIn(f.InternshipIDs...))
	}

	// Apply user ID filter if specified
	if f.UserID != "" {
		query = query.Where(internshipenrollment.UserID(f.UserID))
	}

	// Apply enrollment status filter if specified
	if f.EnrollmentStatus != "" {
		query = query.Where(internshipenrollment.EnrollmentStatus(f.EnrollmentStatus))
	}

	// Apply payment status filter if specified
	if f.PaymentStatus != "" {
		query = query.Where(internshipenrollment.PaymentStatus(f.PaymentStatus))
	}

	// Apply enrollment IDs filter if specified
	if len(f.EnrollmentIDs) > 0 {
		query = query.Where(internshipenrollment.IDIn(f.EnrollmentIDs...))
	}

	// Apply payment ID filter if specified
	if f.PaymentID != nil && *f.PaymentID != "" {
		query = query.Where(internshipenrollment.PaymentID(lo.FromPtr(f.PaymentID)))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(internshipenrollment.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(internshipenrollment.CreatedAtLTE(*f.EndTime))
		}
	}

	return query
}
