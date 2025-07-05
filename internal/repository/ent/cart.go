package ent

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/ent/cart"
	"github.com/omkar273/codegeeky/ent/cartlineitems"
	domainCart "github.com/omkar273/codegeeky/internal/domain/cart"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/postgres"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type cartRepository struct {
	client    postgres.IClient
	logger    logger.Logger
	queryOpts CartQueryOptions
}

func NewCartRepository(client postgres.IClient, logger *logger.Logger) domainCart.Repository {
	return &cartRepository{
		client:    client,
		logger:    *logger,
		queryOpts: CartQueryOptions{},
	}
}

func (r *cartRepository) Create(ctx context.Context, cartData *domainCart.Cart) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("creating cart",
		"cart_id", cartData.ID,
		"user_id", cartData.UserID,
		"type", cartData.Type,
	)

	_, err := client.Cart.Create().
		SetID(cartData.ID).
		SetUserID(cartData.UserID).
		SetType(string(cartData.Type)).
		SetSubtotal(cartData.Subtotal).
		SetDiscountAmount(cartData.DiscountAmount).
		SetTaxAmount(cartData.TaxAmount).
		SetTotal(cartData.Total).
		SetExpiresAt(lo.FromPtr(cartData.ExpiresAt)).
		SetMetadata(cartData.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(cartData.CreatedAt).
		SetUpdatedAt(cartData.UpdatedAt).
		SetCreatedBy(cartData.CreatedBy).
		SetUpdatedBy(cartData.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Cart with this ID already exists").
				WithReportableDetails(map[string]any{
					"cart_id": cartData.ID,
					"user_id": cartData.UserID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create cart").
			WithReportableDetails(map[string]any{
				"cart_id": cartData.ID,
				"user_id": cartData.UserID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *cartRepository) CreateWithLineItems(ctx context.Context, cartData *domainCart.Cart) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("creating cart with line items",
		"cart_id", cartData.ID,
		"line_items_count", len(cartData.LineItems),
	)

	return r.client.WithTx(ctx, func(ctx context.Context) error {
		// 1. Create cart
		_, err := client.Cart.Create().
			SetID(cartData.ID).
			SetUserID(cartData.UserID).
			SetType(string(cartData.Type)).
			SetSubtotal(cartData.Subtotal).
			SetDiscountAmount(cartData.DiscountAmount).
			SetTaxAmount(cartData.TaxAmount).
			SetTotal(cartData.Total).
			SetExpiresAt(lo.FromPtr(cartData.ExpiresAt)).
			SetMetadata(cartData.Metadata).
			SetStatus(string(types.StatusPublished)).
			SetCreatedAt(cartData.CreatedAt).
			SetUpdatedAt(cartData.UpdatedAt).
			SetCreatedBy(cartData.CreatedBy).
			SetUpdatedBy(cartData.UpdatedBy).
			Save(ctx)

		if err != nil {
			if ent.IsConstraintError(err) {
				return ierr.WithError(err).
					WithHint("Cart with this ID already exists").
					WithReportableDetails(map[string]any{
						"cart_id": cartData.ID,
						"user_id": cartData.UserID,
					}).
					Mark(ierr.ErrAlreadyExists)
			}
			return ierr.WithError(err).
				WithHint("Failed to create cart").
				WithReportableDetails(map[string]any{
					"cart_id": cartData.ID,
					"user_id": cartData.UserID,
				}).
				Mark(ierr.ErrDatabase)
		}

		// 2. Create line items in bulk if present
		if len(cartData.LineItems) > 0 {
			builders := make([]*ent.CartLineItemsCreate, len(cartData.LineItems))
			for i, item := range cartData.LineItems {
				builders[i] = client.CartLineItems.Create().
					SetID(item.ID).
					SetCartID(cartData.ID).
					SetEntityID(item.EntityID).
					SetEntityType(string(item.EntityType)).
					SetQuantity(item.Quantity).
					SetPerUnitPrice(item.PerUnitPrice).
					SetTaxAmount(item.TaxAmount).
					SetDiscountAmount(item.DiscountAmount).
					SetSubtotal(item.Subtotal).
					SetTotal(item.Total).
					SetMetadata(item.Metadata).
					SetStatus(string(types.StatusPublished)).
					SetCreatedBy(item.CreatedBy).
					SetUpdatedBy(item.UpdatedBy).
					SetCreatedAt(item.CreatedAt).
					SetUpdatedAt(item.UpdatedAt)
			}

			if err := client.CartLineItems.CreateBulk(builders...).Exec(ctx); err != nil {
				r.logger.Error("failed to create line items", "error", err)
				return ierr.WithError(err).WithHint("line item creation failed").Mark(ierr.ErrDatabase)
			}
		}

		return nil
	})
}

func (r *cartRepository) Get(ctx context.Context, id string) (*domainCart.Cart, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("getting cart", "cart_id", id)

	entCart, err := client.Cart.Query().
		Where(
			cart.ID(id),
			cart.UserID(types.GetUserID(ctx)),
		).
		WithLineItems().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Cart with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"cart_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get cart").
			WithReportableDetails(map[string]any{
				"cart_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	cartData := &domainCart.Cart{}
	return cartData.FromEnt(entCart), nil
}

func (r *cartRepository) Update(ctx context.Context, cartData *domainCart.Cart) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("updating cart",
		"cart_id", cartData.ID,
		"user_id", cartData.UserID,
	)

	_, err := client.Cart.UpdateOneID(cartData.ID).
		Where(cart.UserID(types.GetUserID(ctx))).
		SetMetadata(cartData.Metadata).
		SetStatus(string(cartData.Status)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Cart with ID %s was not found", cartData.ID).
				WithReportableDetails(map[string]any{
					"cart_id": cartData.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to update cart").
			WithReportableDetails(map[string]any{
				"cart_id": cartData.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *cartRepository) Delete(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("deleting cart", "cart_id", id)

	return r.client.WithTx(ctx, func(ctx context.Context) error {
		// Delete line items first
		_, err := client.CartLineItems.Update().
			Where(cartlineitems.CartID(id)).
			SetStatus(string(types.StatusDeleted)).
			SetUpdatedBy(types.GetUserID(ctx)).
			SetUpdatedAt(time.Now().UTC()).
			Save(ctx)
		if err != nil {
			return ierr.WithError(err).WithHint("line item deletion failed").Mark(ierr.ErrDatabase)
		}

		// Then delete cart
		_, err = client.Cart.UpdateOneID(id).
			Where(cart.UserID(types.GetUserID(ctx))).
			SetStatus(string(types.StatusDeleted)).
			SetUpdatedAt(time.Now().UTC()).
			SetUpdatedBy(types.GetUserID(ctx)).
			Save(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				return ierr.WithError(err).
					WithHintf("Cart with ID %s was not found", id).
					WithReportableDetails(map[string]any{
						"cart_id": id,
					}).
					Mark(ierr.ErrNotFound)
			}
			return ierr.WithError(err).
				WithHint("Failed to delete cart").
				WithReportableDetails(map[string]any{
					"cart_id": id,
				}).
				Mark(ierr.ErrDatabase)
		}

		return nil
	})
}

func (r *cartRepository) Count(ctx context.Context, filter *types.CartFilter) (int, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("counting carts")

	query := client.Cart.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	count, err := query.Count(ctx)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count carts").
			Mark(ierr.ErrDatabase)
	}

	return count, nil
}

func (r *cartRepository) List(ctx context.Context, filter *types.CartFilter) ([]*domainCart.Cart, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("listing carts",
		"limit", filter.GetLimit(),
		"offset", filter.GetOffset(),
	)

	query := client.Cart.Query()
	query = r.queryOpts.ApplyBaseFilters(ctx, query, filter)
	query = r.queryOpts.ApplyEntityQueryOptions(ctx, filter, query)

	carts, err := query.
		Where(cart.UserID(types.GetUserID(ctx))).
		WithLineItems().
		All(ctx)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list carts").
			Mark(ierr.ErrDatabase)
	}

	cartData := &domainCart.Cart{}
	return cartData.FromEntList(carts), nil
}

func (r *cartRepository) ListAll(ctx context.Context, filter *types.CartFilter) ([]*domainCart.Cart, error) {
	if filter == nil {
		filter = types.NewNoLimitCartFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	carts, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

// Cart line items methods
func (r *cartRepository) CreateCartLineItem(ctx context.Context, cartLineItem *domainCart.CartLineItem) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("creating cart line item",
		"line_item_id", cartLineItem.ID,
		"cart_id", cartLineItem.CartID,
	)

	_, err := client.CartLineItems.Create().
		SetID(cartLineItem.ID).
		SetCartID(cartLineItem.CartID).
		SetEntityID(cartLineItem.EntityID).
		SetEntityType(string(cartLineItem.EntityType)).
		SetQuantity(cartLineItem.Quantity).
		SetPerUnitPrice(cartLineItem.PerUnitPrice).
		SetTaxAmount(cartLineItem.TaxAmount).
		SetDiscountAmount(cartLineItem.DiscountAmount).
		SetSubtotal(cartLineItem.Subtotal).
		SetTotal(cartLineItem.Total).
		SetMetadata(cartLineItem.Metadata).
		SetStatus(string(types.StatusPublished)).
		SetCreatedAt(cartLineItem.CreatedAt).
		SetUpdatedAt(cartLineItem.UpdatedAt).
		SetCreatedBy(cartLineItem.CreatedBy).
		SetUpdatedBy(cartLineItem.UpdatedBy).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return ierr.WithError(err).
				WithHint("Cart line item with this ID already exists").
				WithReportableDetails(map[string]any{
					"line_item_id": cartLineItem.ID,
					"cart_id":      cartLineItem.CartID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create cart line item").
			WithReportableDetails(map[string]any{
				"line_item_id": cartLineItem.ID,
				"cart_id":      cartLineItem.CartID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *cartRepository) GetCartLineItem(ctx context.Context, id string) (*domainCart.CartLineItem, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("getting cart line item", "line_item_id", id)

	entCartLineItem, err := client.CartLineItems.Query().
		Where(cartlineitems.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"line_item_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHint("Failed to get cart line item").
			WithReportableDetails(map[string]any{
				"line_item_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	cartLineItem := &domainCart.CartLineItem{}
	return cartLineItem.FromEnt(entCartLineItem), nil
}

func (r *cartRepository) UpdateCartLineItem(ctx context.Context, cartLineItem *domainCart.CartLineItem) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("updating cart line item",
		"line_item_id", cartLineItem.ID,
		"cart_id", cartLineItem.CartID,
	)

	_, err := client.CartLineItems.UpdateOneID(cartLineItem.ID).
		Where(cartlineitems.CartID(cartLineItem.CartID)).
		SetQuantity(cartLineItem.Quantity).
		SetPerUnitPrice(cartLineItem.PerUnitPrice).
		SetTaxAmount(cartLineItem.TaxAmount).
		SetDiscountAmount(cartLineItem.DiscountAmount).
		SetSubtotal(cartLineItem.Subtotal).
		SetTotal(cartLineItem.Total).
		SetMetadata(cartLineItem.Metadata).
		SetStatus(string(cartLineItem.Status)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", cartLineItem.ID).
				WithReportableDetails(map[string]any{
					"line_item_id": cartLineItem.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to update cart line item").
			WithReportableDetails(map[string]any{
				"line_item_id": cartLineItem.ID,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *cartRepository) DeleteCartLineItem(ctx context.Context, id string) error {
	client := r.client.Querier(ctx)

	r.logger.Debugw("deleting cart line item", "line_item_id", id)

	_, err := client.CartLineItems.UpdateOneID(id).
		Where(
			cartlineitems.CartID(id),
		).
		SetStatus(string(types.StatusDeleted)).
		SetUpdatedAt(time.Now().UTC()).
		SetUpdatedBy(types.GetUserID(ctx)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"line_item_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHint("Failed to delete cart line item").
			WithReportableDetails(map[string]any{
				"line_item_id": id,
			}).
			Mark(ierr.ErrDatabase)
	}

	return nil
}

func (r *cartRepository) ListCartLineItems(ctx context.Context, cartId string) ([]*domainCart.CartLineItem, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("listing cart line items", "cart_id", cartId)

	cartLineItems, err := client.CartLineItems.Query().
		Where(
			cartlineitems.CartID(cartId),
			cartlineitems.Status(string(types.StatusPublished)),
		).
		All(ctx)

	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list cart line items").
			WithReportableDetails(map[string]any{
				"cart_id": cartId,
			}).
			Mark(ierr.ErrDatabase)
	}

	cartLineItem := &domainCart.CartLineItem{}
	return cartLineItem.FromEntList(cartLineItems), nil
}

func (r *cartRepository) GetUserDefaultCart(ctx context.Context, userID string) (*domainCart.Cart, error) {
	client := r.client.Querier(ctx)

	r.logger.Debugw("getting user default cart", "user_id", userID)

	entCart, err := client.Cart.Query().
		Where(
			cart.UserID(userID),
			cart.Type(string(types.CartTypeDefault)),
			cart.Status(string(types.StatusPublished)),
		).
		Only(ctx)

	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get user default cart").
			WithReportableDetails(map[string]any{
				"user_id": userID,
			}).
			Mark(ierr.ErrDatabase)
	}

	cartData := &domainCart.Cart{}
	return cartData.FromEnt(entCart), nil
}

// CartQuery type alias for better readability
type CartQuery = *ent.CartQuery

// CartQueryOptions implements query options for cart queries
type CartQueryOptions struct {
	QueryOptionsHelper
}

// Ensure CartQueryOptions implements EntityQueryOptions interface
var _ EntityQueryOptions[CartQuery, *types.CartFilter] = (*CartQueryOptions)(nil)

func (o CartQueryOptions) ApplyStatusFilter(query CartQuery, status string) CartQuery {
	if status == "" {
		return query.Where(cart.StatusNotIn(string(types.StatusDeleted)))
	}
	return query.Where(cart.Status(status))
}

func (o CartQueryOptions) ApplySortFilter(query CartQuery, field string, order string) CartQuery {
	field, order = o.QueryOptionsHelper.ValidateSort(field, order)
	fieldName := o.GetFieldName(field)
	if order == types.OrderDesc {
		return query.Order(ent.Desc(fieldName))
	}
	return query.Order(ent.Asc(fieldName))
}

func (o CartQueryOptions) ApplyPaginationFilter(query CartQuery, limit int, offset int) CartQuery {
	limit, offset = o.QueryOptionsHelper.ValidatePagination(limit, offset)
	return query.Offset(offset).Limit(limit)
}

func (o CartQueryOptions) GetFieldName(field string) string {
	switch field {
	case "created_at":
		return cart.FieldCreatedAt
	case "updated_at":
		return cart.FieldUpdatedAt
	case "user_id":
		return cart.FieldUserID
	case "type":
		return cart.FieldType
	case "subtotal":
		return cart.FieldSubtotal
	case "discount_amount":
		return cart.FieldDiscountAmount
	case "tax_amount":
		return cart.FieldTaxAmount
	case "total":
		return cart.FieldTotal
	case "expires_at":
		return cart.FieldExpiresAt
	case "metadata":
		return cart.FieldMetadata
	case "created_by":
		return cart.FieldCreatedBy
	case "updated_by":
		return cart.FieldUpdatedBy
	default:
		return ""
	}
}

func (o CartQueryOptions) ApplyBaseFilters(
	_ context.Context,
	query CartQuery,
	filter *types.CartFilter,
) CartQuery {
	if filter == nil {
		return query.Where(cart.StatusNotIn(string(types.StatusDeleted)))
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

func (o CartQueryOptions) ApplyEntityQueryOptions(
	_ context.Context,
	f *types.CartFilter,
	query CartQuery,
) CartQuery {
	if f == nil {
		return query
	}

	// Apply user ID filter if specified
	if f.UserID != "" {
		query = query.Where(cart.UserID(f.UserID))
	}

	// Apply cart type filter if specified
	if f.CartType != nil {
		query = query.Where(cart.Type(string(*f.CartType)))
	}

	// Apply expires at filter if specified
	if f.ExpiresAt != nil {
		query = query.Where(cart.ExpiresAtGTE(*f.ExpiresAt))
	}

	// Apply entity filters if specified
	if f.EntityID != "" {
		query = query.Where(cart.HasLineItemsWith(cartlineitems.EntityID(f.EntityID)))
	}

	if f.EntityType != "" {
		query = query.Where(cart.HasLineItemsWith(cartlineitems.EntityType(string(f.EntityType))))
	}

	// Apply time range filters if specified
	if f.TimeRangeFilter != nil {
		if f.StartTime != nil {
			query = query.Where(cart.CreatedAtGTE(*f.StartTime))
		}
		if f.EndTime != nil {
			query = query.Where(cart.CreatedAtLTE(*f.EndTime))
		}
	}

	// Apply expansion if requested
	expand := f.GetExpand()
	if expand.Has(types.ExpandCartLineItems) {
		query = query.WithLineItems()
	}

	return query
}
