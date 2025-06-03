package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/category"
	"github.com/omkar273/codegeeky/ent/internship"
	domainCategory "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
)

type categoryRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts CategoryQueryOptions
}

func NewCategoryRepository(client postgres.IClient, logger *logger.Logger) domainCategory.CategoryRepository {
	return &categoryRepository{
		client:    client,
		log:       *logger,
		queryOpts: CategoryQueryOptions{},
	}
}

func (r *categoryRepository) Create(ctx context.Context, categoryData *domainCategory.Category) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating category",
		"category_id", categoryData.ID,
		"name", categoryData.Name,
		"lookup_key", categoryData.LookupKey,
	)

	_, err := client.Category.Create().
		SetID(categoryData.ID).
		SetName(categoryData.Name).
		SetLookupKey(categoryData.LookupKey).
		SetDescription(categoryData.Description).
		SetStatus(string(categoryData.Status)).
		SetCreatedAt(categoryData.CreatedAt).
		SetUpdatedAt(categoryData.UpdatedAt).
		SetCreatedBy(categoryData.CreatedBy).
		SetUpdatedBy(categoryData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Category with this lookup key already exists").
				WithReportableDetails(map[string]any{
					"category_id": categoryData.ID,
					"lookup_key":  categoryData.LookupKey,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create category").
			WithReportableDetails(map[string]any{
				"category_id":   categoryData.ID,
				"category_name": categoryData.Name,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *categoryRepository) Get(ctx context.Context, id string) (*domainCategory.Category, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting category", "category_id", id)

	entCategory, err := client.Category.Query().
		Where(category.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Category with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"category_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get category").
			WithReportableDetails(map[string]any{
				"category_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainCategory.CategoryFromEnt(entCategory), nil
}

func (r *categoryRepository) GetByLookupKey(ctx context.Context, lookupKey string) (*domainCategory.Category, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting category by lookup key", "lookup_key", lookupKey)

	entCategory, err := client.Category.Query().
		Where(
			category.LookupKey(lookupKey),
			category.StatusNotIn(string(types.StatusDeleted)),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Category with lookup key %s was not found", lookupKey).
				WithReportableDetails(map[string]any{
					"lookup_key": lookupKey,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get category by lookup key").
			WithReportableDetails(map[string]any{
				"lookup_key": lookupKey,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainCategory.CategoryFromEnt(entCategory), nil
}

func (r *categoryRepository) List(ctx context.Context, filter *types.CategoryFilter) ([]*domainCategory.Category, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing categories",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.Category.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	categories, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list categories").
			Mark(ierr.ErrDatabase)
	}

	return domainCategory.CategoryFromEntList(categories), nil
}

func (r *categoryRepository) ListAll(ctx context.Context, filter *types.CategoryFilter) ([]*domainCategory.Category, error) {
	if filter == nil {
		filter = types.NewNoLimitCategoryFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	categories, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) Count(ctx context.Context, filter *types.CategoryFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting categories")

	query := client.Category.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count categories").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *categoryRepository) Update(ctx context.Context, categoryData *domainCategory.Category) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating category",
		"category_id", categoryData.ID,
		"name", categoryData.Name,
	)

	_, err := client.Category.UpdateOneID(categoryData.ID).
		SetName(categoryData.Name).
		SetLookupKey(categoryData.LookupKey).
		SetDescription(categoryData.Description).
		SetStatus(string(categoryData.Status)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Category with ID %s was not found", categoryData.ID).
				WithReportableDetails(map[string]any{
					"category_id": categoryData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Category with this lookup key already exists").
				WithReportableDetails(map[string]any{
					"category_id": categoryData.ID,
					"lookup_key":  categoryData.LookupKey,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to update category").
			WithReportableDetails(map[string]any{
				"category_id": categoryData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting category",
		"category_id", id,
	)

	_, err := client.Category.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Category with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"category_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete category").
			WithReportableDetails(map[string]any{
				"category_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

// CategoryQuery type alias for better readability
type CategoryQuery = *ent.CategoryQuery

// CategoryQueryOptions implements query options for category queries
type CategoryQueryOptions struct {
	QueryOptionsHelper
}

// Ensure CategoryQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[CategoryQuery, *types.CategoryFilter] = (*CategoryQueryOptions)(nil)

func (o CategoryQueryOptions) ApplyStatusFilter(query CategoryQuery, status string) CategoryQuery {
	if status == "" {
		return query.Where(category.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(category.Status(status))
}

func (o CategoryQueryOptions) ApplySortFilter(query CategoryQuery, field string, order string) CategoryQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o CategoryQueryOptions) ApplyPaginationFilter(query CategoryQuery, limit int, offset int) CategoryQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o CategoryQueryOptions) GetFieldName(field string) string {
	switch field {
	case "name":
		return category.FieldName
	case "lookup_key":
		return category.FieldLookupKey
	case "description":
		return category.FieldDescription
	case "internships":
		return category.EdgeInternships
	case "created_by":
		return category.FieldCreatedBy
	case "updated_by":
		return category.FieldUpdatedBy
	case "status":
		return category.FieldStatus
	case "id":
		return category.FieldID
	case "created_at":
		return category.FieldCreatedAt
	case "updated_at":
		return category.FieldUpdatedAt
	default:
		return field
	}
}

func (o CategoryQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query CategoryQuery,
	filter *types.CategoryFilter,
) CategoryQuery {
	if filter == nil {
		return query.Where(category.StatusNotIn(string(types.StatusDeleted)))
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

func (o CategoryQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.CategoryFilter,
	query CategoryQuery,
) CategoryQuery {
	if f == nil {
		return query
	}

	// Apply name filter if specified
	if f.Name != "" {
		query = query.Where(category.NameContains(f.Name))
	}

	// Apply category IDs filter if specified
	if len(f.CategoryIDs) > 0 {
		query = query.Where(category.IDIn(f.CategoryIDs...))
	}

	// Apply internship IDs filter if specified (through edge)
	if len(f.InternshipIDs) > 0 {
		query = query.Where(category.HasInternshipsWith(
			internship.IDIn(f.InternshipIDs...),
		))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(category.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(category.CreatedAtLTE(*f.EndTime))
		}
	}

	// Apply expansion if requested
	expand := f.GetExpand()
	if expand.Has("internships") {
		query = query.WithInternships()
	}

	return query
}
