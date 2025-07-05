package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/cart"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryCartStore implements cart.Repository
type InMemoryCartStore struct {
	*InMemoryStore[*cart.Cart]
	lineItems *InMemoryStore[*cart.CartLineItem]
}

// NewInMemoryCartStore creates a new in-memory cart store
func NewInMemoryCartStore() *InMemoryCartStore {
	return &InMemoryCartStore{
		InMemoryStore: NewInMemoryStore[*cart.Cart](),
		lineItems:     NewInMemoryStore[*cart.CartLineItem](),
	}
}

// cartFilterFn implements filtering logic for carts
func cartFilterFn(ctx context.Context, c *cart.Cart, filter interface{}) bool {
	if c == nil {
		return false
	}

	filter_, ok := filter.(*types.CartFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by user ID
	if filter_.UserID != "" {
		if c.UserID != filter_.UserID {
			return false
		}
	}

	// Filter by cart type
	if filter_.CartType != nil {
		if c.Type != *filter_.CartType {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active carts
	if filter_.GetStatus() != "" {
		if string(c.Status) != filter_.GetStatus() {
			return false
		}
	} else if c.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && c.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && c.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// cartSortFn implements sorting logic for carts
func cartSortFn(i, j *cart.Cart) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryCartStore) Create(ctx context.Context, c *cart.Cart) error {
	if c == nil {
		return ierr.NewError("cart cannot be nil").
			WithHint("Cart data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, c.ID, c)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A cart with this ID already exists").
				WithReportableDetails(map[string]any{
					"cart_id": c.ID,
					"user_id": c.UserID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create cart").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryCartStore) CreateWithLineItems(ctx context.Context, c *cart.Cart) error {
	if c == nil {
		return ierr.NewError("cart cannot be nil").
			WithHint("Cart data is required").
			Mark(ierr.ErrValidation)
	}

	// Create the cart first
	if err := s.Create(ctx, c); err != nil {
		return err
	}

	// Create line items if they exist
	if len(c.LineItems) > 0 {
		for _, item := range c.LineItems {
			item.CartID = c.ID
			if err := s.CreateCartLineItem(ctx, item); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *InMemoryCartStore) Get(ctx context.Context, id string) (*cart.Cart, error) {
	c, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Cart with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"cart_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get cart with ID %s", id).
			Mark(ierr.ErrDatabase)
	}

	// Load line items
	lineItems, err := s.ListCartLineItems(ctx, id)
	if err != nil {
		return nil, err
	}
	c.LineItems = lineItems

	return c, nil
}

func (s *InMemoryCartStore) Update(ctx context.Context, c *cart.Cart) error {
	if c == nil {
		return ierr.NewError("cart cannot be nil").
			WithHint("Cart data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	c.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, c.ID, c)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Cart with ID %s was not found", c.ID).
				WithReportableDetails(map[string]any{
					"cart_id": c.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update cart with ID %s", c.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryCartStore) Delete(ctx context.Context, id string) error {
	// Get the cart first
	c, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	c.Status = types.StatusDeleted
	c.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, c)
}

func (s *InMemoryCartStore) Count(ctx context.Context, filter *types.CartFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, cartFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count carts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryCartStore) List(ctx context.Context, filter *types.CartFilter) ([]*cart.Cart, error) {
	carts, err := s.InMemoryStore.List(ctx, filter, cartFilterFn, cartSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list carts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return carts, nil
}

func (s *InMemoryCartStore) ListAll(ctx context.Context, filter *types.CartFilter) ([]*cart.Cart, error) {
	if filter == nil {
		filter = types.NewNoLimitCartFilter()
	}

	if filter.QueryFilter == nil {
		filter.QueryFilter = types.NewNoLimitQueryFilter()
	}

	if !filter.IsUnlimited() {
		filter.QueryFilter.Limit = nil
	}

	if err := filter.Validate(); err != nil {
		return nil, ierr.WithError(err).
			WithHint("Invalid filter parameters").
			Mark(ierr.ErrValidation)
	}

	// Create an unlimited filter
	unlimitedFilter := &types.CartFilter{
		QueryFilter:     types.NewNoLimitQueryFilter(),
		TimeRangeFilter: filter.TimeRangeFilter,
		UserID:          filter.UserID,
		EntityID:        filter.EntityID,
		EntityType:      filter.EntityType,
		CartType:        filter.CartType,
		ExpiresAt:       filter.ExpiresAt,
	}

	return s.List(ctx, unlimitedFilter)
}

// Cart Line Items methods

func (s *InMemoryCartStore) CreateCartLineItem(ctx context.Context, item *cart.CartLineItem) error {
	if item == nil {
		return ierr.NewError("cart line item cannot be nil").
			WithHint("Cart line item data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = now
	}

	err := s.lineItems.Create(ctx, item.ID, item)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A cart line item with this ID already exists").
				WithReportableDetails(map[string]any{
					"line_item_id": item.ID,
					"cart_id":      item.CartID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create cart line item").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryCartStore) GetCartLineItem(ctx context.Context, id string) (*cart.CartLineItem, error) {
	item, err := s.lineItems.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"line_item_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get cart line item with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return item, nil
}

func (s *InMemoryCartStore) UpdateCartLineItem(ctx context.Context, item *cart.CartLineItem) error {
	if item == nil {
		return ierr.NewError("cart line item cannot be nil").
			WithHint("Cart line item data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	item.UpdatedAt = time.Now().UTC()

	err := s.lineItems.Update(ctx, item.ID, item)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", item.ID).
				WithReportableDetails(map[string]any{
					"line_item_id": item.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update cart line item with ID %s", item.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryCartStore) DeleteCartLineItem(ctx context.Context, id string) error {
	err := s.lineItems.Delete(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Cart line item with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"line_item_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to delete cart line item with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryCartStore) ListCartLineItems(ctx context.Context, cartId string) ([]*cart.CartLineItem, error) {
	allItems, err := s.lineItems.List(ctx, nil, nil, nil)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list cart line items").
			Mark(ierr.ErrDatabase)
	}

	// Filter by cart ID
	items := lo.Filter(allItems, func(item *cart.CartLineItem, _ int) bool {
		return item.CartID == cartId && item.Status != types.StatusDeleted
	})

	return items, nil
}

func (s *InMemoryCartStore) GetUserDefaultCart(ctx context.Context, userID string) (*cart.Cart, error) {
	// Get all carts for the user
	allCarts, err := s.InMemoryStore.List(ctx, nil, nil, nil)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get user default cart").
			WithReportableDetails(map[string]any{
				"user_id": userID,
			}).
			Mark(ierr.ErrDatabase)
	}

	// Filter for the user's default cart
	for _, c := range allCarts {
		if c.UserID == userID &&
			c.Type == types.CartTypeDefault &&
			c.Status == types.StatusPublished {
			// Load line items for the cart
			lineItems, err := s.ListCartLineItems(ctx, c.ID)
			if err != nil {
				return nil, err
			}
			c.LineItems = lineItems
			return c, nil
		}
	} 

	// No default cart found
	return nil, ierr.NewError("user default cart not found").
		WithHint("No default cart exists for this user").
		WithReportableDetails(map[string]any{
			"user_id": userID,
		}).
		Mark(ierr.ErrNotFound)
}

// Clear clears both cart and line item stores
func (s *InMemoryCartStore) Clear() {
	s.InMemoryStore.Clear()
	s.lineItems.Clear()
}
