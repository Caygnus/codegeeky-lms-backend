package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/category"
	"github.com/omkar273/codegeeky/ent/internship"
	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
)

type internshipRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts InternshipQueryOptions
}

func NewInternshipRepository(client postgres.IClient, logger *logger.Logger) domainInternship.InternshipRepository {
	return &internshipRepository{
		client:    client,
		log:       *logger,
		queryOpts: InternshipQueryOptions{},
	}
}

func (r *internshipRepository) Create(ctx context.Context, internshipData *domainInternship.Internship) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("creating internship",
		"internship_id", internshipData.ID,
		"title", internshipData.Title,
		"lookup_key", internshipData.LookupKey,
	)

	// Create internship
	_, err := client.Internship.Create().
		SetID(internshipData.ID).
		SetTitle(internshipData.Title).
		SetLookupKey(internshipData.LookupKey).
		SetDescription(internshipData.Description).
		SetSkills(internshipData.Skills).
		SetLevel(string(internshipData.Level)).
		SetMode(string(internshipData.Mode)).
		SetDurationInWeeks(internshipData.DurationInWeeks).
		SetLearningOutcomes(internshipData.LearningOutcomes).
		SetPrerequisites(internshipData.Prerequisites).
		SetBenefits(internshipData.Benefits).
		SetCurrency(internshipData.Currency).
		SetPrice(internshipData.Price).
		SetFlatDiscount(internshipData.FlatDiscount).
		SetPercentageDiscount(internshipData.PercentageDiscount).
		SetStatus(string(internshipData.Status)).
		SetCreatedAt(internshipData.CreatedAt).
		SetUpdatedAt(internshipData.UpdatedAt).
		SetCreatedBy(internshipData.CreatedBy).
		SetUpdatedBy(internshipData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Internship with this lookup key already exists").
				WithReportableDetails(map[string]any{
					"internship_id": internshipData.ID,
					"lookup_key":    internshipData.LookupKey,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create internship").
			WithReportableDetails(map[string]any{
				"internship_id":    internshipData.ID,
				"internship_title": internshipData.Title,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *internshipRepository) Get(ctx context.Context, id string) (*domainInternship.Internship, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting internship", "internship_id", id)

	entInternship, err := client.Internship.Query().
		Where(internship.ID(id)).
		WithCategories().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Internship with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"internship_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get internship").
			WithReportableDetails(map[string]any{
				"internship_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainInternship.InternshipFromEnt(entInternship), nil
}

func (r *internshipRepository) GetByLookupKey(ctx context.Context, lookupKey string) (*domainInternship.Internship, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("getting internship by lookup key", "lookup_key", lookupKey)

	entInternship, err := client.Internship.Query().
		Where(
			internship.LookupKey(lookupKey),
			internship.StatusNotIn(string(types.StatusDeleted)),
		).
		WithCategories().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Internship with lookup key %s was not found", lookupKey).
				WithReportableDetails(map[string]any{
					"lookup_key": lookupKey,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get internship by lookup key").
			WithReportableDetails(map[string]any{
				"lookup_key": lookupKey,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainInternship.InternshipFromEnt(entInternship), nil
}

func (r *internshipRepository) List(ctx context.Context, filter *types.InternshipFilter) ([]*domainInternship.Internship, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("listing internships",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.Internship.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	// Add eager loading for categories
	query = query.WithCategories()

	internships, err := query.All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list internships").
			Mark(ierr.ErrDatabase)
	}

	return domainInternship.InternshipFromEntList(internships), nil
}

func (r *internshipRepository) ListAll(ctx context.Context, filter *types.InternshipFilter) ([]*domainInternship.Internship, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	internships, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return internships, nil
}

func (r *internshipRepository) Count(ctx context.Context, filter *types.InternshipFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.log.Debugw("counting internships")

	query := client.Internship.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count internships").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *internshipRepository) Update(ctx context.Context, internshipData *domainInternship.Internship) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("updating internship",
		"internship_id", internshipData.ID,
		"title", internshipData.Title,
	)

	_, err := client.Internship.UpdateOneID(internshipData.ID).
		SetTitle(internshipData.Title).
		SetLookupKey(internshipData.LookupKey).
		SetDescription(internshipData.Description).
		SetSkills(internshipData.Skills).
		SetLevel(string(internshipData.Level)).
		SetMode(string(internshipData.Mode)).
		SetDurationInWeeks(internshipData.DurationInWeeks).
		SetLearningOutcomes(internshipData.LearningOutcomes).
		SetPrerequisites(internshipData.Prerequisites).
		SetBenefits(internshipData.Benefits).
		SetCurrency(internshipData.Currency).
		SetPrice(internshipData.Price).
		SetFlatDiscount(internshipData.FlatDiscount).
		SetPercentageDiscount(internshipData.PercentageDiscount).
		SetStatus(string(internshipData.Status)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Internship with ID %s was not found", internshipData.ID).
				WithReportableDetails(map[string]any{
					"internship_id": internshipData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Internship with this lookup key already exists").
				WithReportableDetails(map[string]any{
					"internship_id": internshipData.ID,
					"lookup_key":    internshipData.LookupKey,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to update internship").
			WithReportableDetails(map[string]any{
				"internship_id": internshipData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *internshipRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.log.Debugw("deleting internship",
		"internship_id", id,
	)

	_, err := client.Internship.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Internship with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"internship_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete internship").
			WithReportableDetails(map[string]any{
				"internship_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

// InternshipQuery type alias for better readability
type InternshipQuery = *ent.InternshipQuery

// InternshipQueryOptions implements query options for internship queries
type InternshipQueryOptions struct {
	QueryOptionsHelper
}

// Ensure InternshipQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[InternshipQuery, *types.InternshipFilter] = (*InternshipQueryOptions)(nil)

func (o InternshipQueryOptions) ApplyStatusFilter(query InternshipQuery, status string) InternshipQuery {
	if status == "" {
		return query.Where(internship.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(internship.Status(status))
}

func (o InternshipQueryOptions) ApplySortFilter(query InternshipQuery, field string, order string) InternshipQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o InternshipQueryOptions) ApplyPaginationFilter(query InternshipQuery, limit int, offset int) InternshipQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o InternshipQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return internship.FieldCreatedAt
	case "updated_at":
		return internship.FieldUpdatedAt
	case "title":
		return internship.FieldTitle
	case "lookup_key":
		return internship.FieldLookupKey
	case "level":
		return internship.FieldLevel
	case "mode":
		return internship.FieldMode
	case "duration_in_weeks":
		return internship.FieldDurationInWeeks
	case "currency":
		return internship.FieldCurrency
	case "price":
		return internship.FieldPrice
	case "categories":
		return internship.EdgeCategories
	case "skills":
		return internship.FieldSkills
	case "learning_outcomes":
		return internship.FieldLearningOutcomes
	case "prerequisites":
		return internship.FieldPrerequisites
	case "benefits":
		return internship.FieldBenefits
	case "flat_discount":
		return internship.FieldFlatDiscount
	case "percentage_discount":
		return internship.FieldPercentageDiscount
	case "description":
		return internship.FieldDescription
	case "created_by":
		return internship.FieldCreatedBy
	case "updated_by":
		return internship.FieldUpdatedBy
	default:
		return field
	}
}

func (o InternshipQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query InternshipQuery,
	filter *types.InternshipFilter,
) InternshipQuery {
	if filter == nil {
		return query.Where(internship.StatusNotIn(string(types.StatusDeleted)))
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

func (o InternshipQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.InternshipFilter,
	query InternshipQuery,
) InternshipQuery {
	if f == nil {
		return query
	}

	// Apply modes filter if specified
	if len(f.Modes) > 0 {
		modes := make([]string, len(f.Modes))
		for i, mode := range f.Modes {
			modes[i] = string(mode)
		}
		query = query.Where(internship.ModeIn(modes...))
	}

	// Apply levels filter if specified
	if len(f.Levels) > 0 {
		levels := make([]string, len(f.Levels))
		for i, level := range f.Levels {
			levels[i] = string(level)
		}
		query = query.Where(internship.LevelIn(levels...))
	}

	// Apply category IDs filter if specified
	if len(f.CategoryIDs) > 0 {
		query = query.Where(internship.HasCategoriesWith(category.IDIn(f.CategoryIDs...)))
	}

	// Apply internship IDs filter if specified
	if len(f.InternshipIDs) > 0 {
		query = query.Where(internship.IDIn(f.InternshipIDs...))
	}

	// Apply price range filters if specified
	if !f.MinPrice.IsZero() {
		query = query.Where(internship.PriceGTE(f.MinPrice))
	}
	if !f.MaxPrice.IsZero() {
		query = query.Where(internship.PriceLTE(f.MaxPrice))
	}

	// Apply duration filter if specified
	if f.DurationInWeeks > 0 {
		query = query.Where(internship.DurationInWeeksEQ(f.DurationInWeeks))
	}

	// Apply name filter if specified (search in title)
	if f.Name != "" {
		query = query.Where(internship.TitleContainsFold(f.Name))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(internship.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(internship.CreatedAtLTE(*f.EndTime))
		}
	}

	// Apply expansion if requested
	expand := f.GetExpand()
	if expand.Has(types.ExpandCategory) {
		query = query.WithCategories()
	}

	return query
}
