package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/discount"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryDiscountStore implements discount.Repository
type InMemoryDiscountStore struct {
	*InMemoryStore[*discount.Discount]
}

// NewInMemoryDiscountStore creates a new in-memory discount store
func NewInMemoryDiscountStore() *InMemoryDiscountStore {
	return &InMemoryDiscountStore{
		InMemoryStore: NewInMemoryStore[*discount.Discount](),
	}
}

// discountFilterFn implements filtering logic for discounts
func discountFilterFn(ctx context.Context, d *discount.Discount, filter interface{}) bool {
	if d == nil {
		return false
	}

	filter_, ok := filter.(*types.DiscountFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by codes
	if len(filter_.Codes) > 0 {
		if !lo.Contains(filter_.Codes, d.Code) {
			return false
		}
	}

	// Filter by discount IDs
	if len(filter_.DiscountIDs) > 0 {
		if !lo.Contains(filter_.DiscountIDs, d.ID) {
			return false
		}
	}

	// Filter by discount type
	if filter_.DiscountType != "" {
		if d.DiscountType != filter_.DiscountType {
			return false
		}
	}

	// Filter by is combinable
	if filter_.IsCombinable {
		if !d.IsCombinable {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active discounts
	if filter_.GetStatus() != "" {
		if string(d.Status) != filter_.GetStatus() {
			return false
		}
	} else if d.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && d.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && d.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// discountSortFn implements sorting logic for discounts
func discountSortFn(i, j *discount.Discount) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryDiscountStore) Create(ctx context.Context, d *discount.Discount) error {
	if d == nil {
		return ierr.NewError("discount cannot be nil").
			WithHint("Discount data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, d.ID, d)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A discount with this ID already exists").
				WithReportableDetails(map[string]any{
					"discount_id": d.ID,
					"code":        d.Code,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create discount").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryDiscountStore) Get(ctx context.Context, id string) (*discount.Discount, error) {
	discount, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Discount with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"discount_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get discount with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return discount, nil
}

func (s *InMemoryDiscountStore) GetByCode(ctx context.Context, code string) (*discount.Discount, error) {
	discounts, err := s.InMemoryStore.List(ctx, nil, discountFilterFn, discountSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get discount by code").
			Mark(ierr.ErrDatabase)
	}

	for _, d := range discounts {
		if d.Code == code && d.Status != types.StatusDeleted {
			return d, nil
		}
	}

	return nil, ierr.NewError("discount not found").
		WithHintf("Discount with code %s was not found", code).
		WithReportableDetails(map[string]any{
			"code": code,
		}).
		Mark(ierr.ErrNotFound)
}

func (s *InMemoryDiscountStore) Update(ctx context.Context, d *discount.Discount) error {
	if d == nil {
		return ierr.NewError("discount cannot be nil").
			WithHint("Discount data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	d.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, d.ID, d)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Discount with ID %s was not found", d.ID).
				WithReportableDetails(map[string]any{
					"discount_id": d.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update discount with ID %s", d.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryDiscountStore) Delete(ctx context.Context, id string) error {
	// Get the discount first
	d, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	d.Status = types.StatusDeleted
	d.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, d)
}

func (s *InMemoryDiscountStore) Count(ctx context.Context, filter *types.DiscountFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, discountFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count discounts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryDiscountStore) List(ctx context.Context, filter *types.DiscountFilter) ([]*discount.Discount, error) {
	discounts, err := s.InMemoryStore.List(ctx, filter, discountFilterFn, discountSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list discounts").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return discounts, nil
}

func (s *InMemoryDiscountStore) ListAll(ctx context.Context, filter *types.DiscountFilter) ([]*discount.Discount, error) {
	if filter == nil {
		filter = types.NewNoLimitDiscountFilter()
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
	unlimitedFilter := &types.DiscountFilter{
		QueryFilter:     types.NewNoLimitQueryFilter(),
		TimeRangeFilter: filter.TimeRangeFilter,
		DiscountType:    filter.DiscountType,
		ValidFrom:       filter.ValidFrom,
		ValidUntil:      filter.ValidUntil,
		MinOrderValue:   filter.MinOrderValue,
		IsCombinable:    filter.IsCombinable,
		Codes:           filter.Codes,
		DiscountIDs:     filter.DiscountIDs,
	}

	return s.List(ctx, unlimitedFilter)
}

// Clear clears the discount store
func (s *InMemoryDiscountStore) Clear() {
	s.InMemoryStore.Clear()
}
