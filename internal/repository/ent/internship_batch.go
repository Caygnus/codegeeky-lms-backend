package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/internshipbatch"
	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
)

type internshipBatchRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts InternshipBatchQueryOptions
}

func NewInternshipBatchRepository(client postgres.IClient, logger *logger.Logger) domainInternship.InternshipBatchRepository {
	return &internshipBatchRepository{
		client:    client,
		log:       *logger,
		queryOpts: InternshipBatchQueryOptions{},
	}
}

func (r *internshipBatchRepository) Create(ctx context.Context, batch *domainInternship.InternshipBatch) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating internship batch",
		"batch_id", batch.ID,
		"internship_id", batch.InternshipID,
		"name", batch.Name,
	)

	_, err := client.InternshipBatch.Create().
		SetID(batch.ID).
		SetInternshipID(batch.InternshipID).
		SetName(batch.Name).
		SetDescription(batch.Description).
		SetStartDate(batch.StartDate).
		SetEndDate(batch.EndDate).
		SetBatchStatus(string(batch.BatchStatus)).
		SetMetadata(batch.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(batch.CreatedAt).
		SetUpdatedAt(batch.UpdatedAt).
		SetCreatedBy(batch.CreatedBy).
		SetUpdatedBy(batch.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Internship batch with this ID already exists").
				WithReportableDetails(map[string]any{
					"batch_id": batch.ID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create internship batch").
			WithReportableDetails(map[string]any{
				"batch_id":      batch.ID,
				"internship_id": batch.InternshipID,
				"name":          batch.Name,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *internshipBatchRepository) Get(ctx context.Context, id string) (*domainInternship.InternshipBatch, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting internship batch", "batch_id", id)

	entBatch, err := client.InternshipBatch.Query().
		Where(internshipbatch.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Internship batch with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"batch_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get internship batch").
			WithReportableDetails(map[string]any{
				"batch_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	batch := &domainInternship.InternshipBatch{}
	return batch.FromEnt(entBatch), nil
}

func (r *internshipBatchRepository) Update(ctx context.Context, batch *domainInternship.InternshipBatch) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating internship batch",
		"batch_id", batch.ID,
		"name", batch.Name,
	)

	_, err := client.InternshipBatch.UpdateOneID(batch.ID).
		SetName(batch.Name).
		SetDescription(batch.Description).
		SetStartDate(batch.StartDate).
		SetEndDate(batch.EndDate).
		SetBatchStatus(string(batch.BatchStatus)).
		SetMetadata(batch.Metadata).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Internship batch with ID %s was not found", batch.ID).
				WithReportableDetails(map[string]any{
					"batch_id": batch.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to update internship batch").
			WithReportableDetails(map[string]any{
				"batch_id": batch.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *internshipBatchRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting internship batch",
		"batch_id", id,
	)

	_, err := client.InternshipBatch.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Internship batch with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"batch_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete internship batch").
			WithReportableDetails(map[string]any{
				"batch_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *internshipBatchRepository) Count(ctx context.Context, filter *types.InternshipBatchFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting internship batches")

	query := client.InternshipBatch.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count internship batches").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *internshipBatchRepository) List(ctx context.Context, filter *types.InternshipBatchFilter) ([]*domainInternship.InternshipBatch, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing internship batches",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.InternshipBatch.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	batches, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list internship batches").
			Mark(ierr.ErrDatabase)
	}

	batch := &domainInternship.InternshipBatch{}
	return batch.FromEntList(batches), nil
}

func (r *internshipBatchRepository) ListAll(ctx context.Context, filter *types.InternshipBatchFilter) ([]*domainInternship.InternshipBatch, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipBatchFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	batches, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return batches, nil
}

// InternshipBatchQuery type alias for better readability
type InternshipBatchQuery = *ent.InternshipBatchQuery

// InternshipBatchQueryOptions implements query options for internship batch queries
type InternshipBatchQueryOptions struct {
	QueryOptionsHelper
}

// Ensure InternshipBatchQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[InternshipBatchQuery, *types.InternshipBatchFilter] = (*InternshipBatchQueryOptions)(nil)

func (o InternshipBatchQueryOptions) ApplyStatusFilter(query InternshipBatchQuery, status string) InternshipBatchQuery {
	if status == "" {
		return query.Where(internshipbatch.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(internshipbatch.Status(status))
}

func (o InternshipBatchQueryOptions) ApplySortFilter(query InternshipBatchQuery, field string, order string) InternshipBatchQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o InternshipBatchQueryOptions) ApplyPaginationFilter(query InternshipBatchQuery, limit int, offset int) InternshipBatchQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o InternshipBatchQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return internshipbatch.FieldCreatedAt
	case "updated_at":
		return internshipbatch.FieldUpdatedAt
	case "internship_id":
		return internshipbatch.FieldInternshipID
	case "name":
		return internshipbatch.FieldName
	case "description":
		return internshipbatch.FieldDescription
	case "start_date":
		return internshipbatch.FieldStartDate
	case "end_date":
		return internshipbatch.FieldEndDate
	case "batch_status":
		return internshipbatch.FieldBatchStatus
	case "created_by":
		return internshipbatch.FieldCreatedBy
	case "updated_by":
		return internshipbatch.FieldUpdatedBy
	default:
		return field
	}
}

func (o InternshipBatchQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query InternshipBatchQuery,
	filter *types.InternshipBatchFilter,
) InternshipBatchQuery {
	if filter == nil {
		return query.Where(internshipbatch.StatusNotIn(string(types.StatusDeleted)))
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

func (o InternshipBatchQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.InternshipBatchFilter,
	query InternshipBatchQuery,
) InternshipBatchQuery {
	if f == nil {
		return query
	}

	// Apply internship IDs filter if specified
	if len(f.InternshipIDs) > 0 {
		query = query.Where(internshipbatch.InternshipIDIn(f.InternshipIDs...))
	}

	// Apply name filter if specified (search in name)
	if f.Name != "" {
		query = query.Where(internshipbatch.NameContainsFold(f.Name))
	}

	// Apply batch status filter if specified
	if f.BatchStatus != "" {
		query = query.Where(internshipbatch.BatchStatus(string(f.BatchStatus)))
	}

	// Apply start date filter if specified
	if f.StartDate != nil {
		query = query.Where(internshipbatch.StartDateGTE(*f.StartDate))
	}

	// Apply end date filter if specified
	if f.EndDate != nil {
		query = query.Where(internshipbatch.EndDateLTE(*f.EndDate))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(internshipbatch.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(internshipbatch.CreatedAtLTE(*f.EndTime))
		}
	}

	return query
}
