package ent

import (
	"context"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/discount"
	domainDiscount "github.com/omkar273/codegeeky/internal/domain/discount"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type discountRepository struct {
	client    postgres.IClient
	log       logger.Logger
	queryOpts DiscountQueryOptions
}

func NewDiscountRepository(client postgres.IClient, logger *logger.Logger) domainDiscount.Repository {
	return &discountRepository{
		client:    client,
		log:       *logger,
		queryOpts: DiscountQueryOptions{},
	}
}

func (r *discountRepository) Count(ctx context.Context, filter *types.DiscountFilter) (int, error) {
	client := r.client.Querier(ctx)
	query := client.Discount.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count discounts").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *discountRepository) Create(ctx context.Context, d *domainDiscount.Discount) error {

	client := r.client.Querier(ctx)

	_, err := client.Discount.Create().
		SetID(d.ID).
		SetCode(d.Code).
		SetDescription(d.Description).
		SetDiscountType(d.DiscountType).
		SetDiscountValue(d.DiscountValue).
		SetValidFrom(d.ValidFrom).
		SetValidUntil(lo.FromPtr(d.ValidUntil)).
		SetIsActive(d.IsActive).
		SetMaxUses(lo.FromPtr(d.MaxUses)).
		SetMinOrderValue(lo.FromPtr(d.MinOrderValue)).
		SetIsCombinable(d.IsCombinable).
		SetMetadata(d.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(d.CreatedAt).
		SetUpdatedAt(d.UpdatedAt).
		SetCreatedBy(d.CreatedBy).
		SetUpdatedBy(d.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Discount with this code already exists").
				WithReportableDetails(map[string]any{
					"code": d.Code,
				}).
				Mark(ierr.ErrAlreadyExists)
		}

		return ierr.WithError(err).
			WithHint("Failed to create discount").
			WithReportableDetails(map[string]any{
				"discount": d,
			}).
			Mark(ierr.ErrDatabase)
	}

	return err
}

func (r *discountRepository) Get(ctx context.Context, id string) (*domainDiscount.Discount, error) {
	client := r.client.Querier(ctx)
	discount, err := client.Discount.Query().
		Where(discount.ID(id)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Discount with this ID not found").
				WithReportableDetails(map[string]any{
					"discount_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}

		return nil, ierr.WithError(err).
			WithHint("Failed to get discount").
			WithReportableDetails(map[string]any{
				"discount_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainDiscount.FromEnt(discount), nil
}

func (r *discountRepository) List(ctx context.Context, filter *types.DiscountFilter) ([]*domainDiscount.Discount, error) {
	client := r.client.Querier(ctx)
	query := client.Discount.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	discounts, err := query.All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Discounts not found").
				WithReportableDetails(map[string]any{
					"filter": filter,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get discounts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}

	discount := &domainDiscount.Discount{}
	return discount.FromEntList(discounts), nil
}

func (r *discountRepository) ListAll(ctx context.Context, filter *types.DiscountFilter) ([]*domainDiscount.Discount, error) {
	client := r.client.Querier(ctx)
	query := client.Discount.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)

	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	discounts, err := query.All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Discounts not found").
				WithReportableDetails(map[string]any{
					"filter": filter,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get discounts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}

	discount := &domainDiscount.Discount{}
	return discount.FromEntList(discounts), nil
}

func (r *discountRepository) Update(ctx context.Context, discount *domainDiscount.Discount) error {
	client := r.client.Querier(ctx)
	_, err := client.Discount.UpdateOneID(discount.ID).
		SetDescription(discount.Description).
		SetValidFrom(discount.ValidFrom).
		SetValidUntil(lo.FromPtr(discount.ValidUntil)).
		SetIsActive(discount.IsActive).
		SetMaxUses(lo.FromPtr(discount.MaxUses)).
		SetMinOrderValue(lo.FromPtr(discount.MinOrderValue)).
		SetMetadata(discount.Metadata).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	return err
}

func (r *discountRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)
	_, err := client.Discount.UpdateOneID(id).
		SetStatus(string(types.StatusDeleted)).
		Save(ctx)

	return err
}

func (r *discountRepository) GetByCode(ctx context.Context, code string) (*domainDiscount.Discount, error) {
	client := r.client.Querier(ctx)
	discount, err := client.Discount.Query().
		Where(discount.Code(code)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHint("Discount with this code not found").
				WithReportableDetails(map[string]any{
					"code": code,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get discount").
			WithReportableDetails(map[string]any{
				"code": code,
			}).
			Mark(ierr.ErrDatabase)
	}

	return domainDiscount.FromEnt(discount), nil
}

// DiscountQuery type alias for better readability
type DiscountQuery = *ent.DiscountQuery

// DiscountQueryOptions implements query options for discount queries
type DiscountQueryOptions struct {
	QueryOptionsHelper
}

// Ensure DiscountQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[DiscountQuery, *types.DiscountFilter] = (*DiscountQueryOptions)(nil)

func (o DiscountQueryOptions) ApplyStatusFilter(query DiscountQuery, status string) DiscountQuery {
	if status == "" {
		return query.Where(discount.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(discount.Status(status))
}

func (o DiscountQueryOptions) ApplySortFilter(query DiscountQuery, field string, order string) DiscountQuery {
	field, order = o.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o DiscountQueryOptions) ApplyPaginationFilter(query DiscountQuery, limit int, offset int) DiscountQuery {
	limit, offset = o.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o DiscountQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return discount.FieldCreatedAt
	case "updated_at":
		return discount.FieldUpdatedAt
	case "code":
		return discount.FieldCode
	case "description":
		return discount.FieldDescription
	case "discount_type":
		return discount.FieldDiscountType
	case "discount_value":
		return discount.FieldDiscountValue
	case "valid_from":
		return discount.FieldValidFrom
	case "valid_until":
		return discount.FieldValidUntil
	case "is_active":
		return discount.FieldIsActive
	case "max_uses":
		return discount.FieldMaxUses
	case "min_order_value":
		return discount.FieldMinOrderValue
	case "is_combinable":
		return discount.FieldIsCombinable
	case "metadata":
		return discount.FieldMetadata
	case "created_by":
		return discount.FieldCreatedBy
	case "updated_by":
		return discount.FieldUpdatedBy
	default:
		return ""
	}
}

func (o DiscountQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query DiscountQuery,
	filter *types.DiscountFilter,
) DiscountQuery {
	if filter == nil {
		return query.Where(discount.StatusNotIn(string(types.StatusDeleted)))
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

func (o DiscountQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.DiscountFilter,
	query DiscountQuery,
) DiscountQuery {
	if f == nil {
		return query
	}

	if f.DiscountType != "" {
		query = query.Where(discount.DiscountType(f.DiscountType))
	}

	if f.ValidFrom != nil {
		query = query.Where(discount.ValidFromLTE(*f.ValidFrom))
	}

	if f.ValidUntil != nil {
		query = query.Where(discount.ValidUntilGTE(*f.ValidUntil))
	}

	if f.MinOrderValue != nil {
		query = query.Where(discount.MinOrderValueLTE(*f.MinOrderValue))
	}

	if f.IsCombinable {
		query = query.Where(discount.IsCombinable(f.IsCombinable))
	}

	if len(f.Codes) > 0 {
		query = query.Where(discount.CodeIn(f.Codes...))
	}

	if len(f.DiscountIDs) > 0 {
		query = query.Where(discount.IDIn(f.DiscountIDs...))
	}
	return query
}
