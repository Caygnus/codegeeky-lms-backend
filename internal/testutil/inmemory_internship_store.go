package testutil

import (
	"context"
	"strings"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryInternshipStore implements internship.InternshipRepository
type InMemoryInternshipStore struct {
	*InMemoryStore[*internship.Internship]
}

// NewInMemoryInternshipStore creates a new in-memory internship store
func NewInMemoryInternshipStore() *InMemoryInternshipStore {
	return &InMemoryInternshipStore{
		InMemoryStore: NewInMemoryStore[*internship.Internship](),
	}
}

// internshipFilterFn implements filtering logic for internships
func internshipFilterFn(ctx context.Context, i *internship.Internship, filter interface{}) bool {
	if i == nil {
		return false
	}

	filter_, ok := filter.(*types.InternshipFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by name contains
	if filter_.Name != "" {
		if !strings.Contains(strings.ToLower(i.Title), strings.ToLower(filter_.Name)) {
			return false
		}
	}

	// Filter by internship IDs
	if len(filter_.InternshipIDs) > 0 {
		if !lo.Contains(filter_.InternshipIDs, i.ID) {
			return false
		}
	}

	// Filter by levels
	if len(filter_.Levels) > 0 {
		if !lo.Contains(filter_.Levels, i.Level) {
			return false
		}
	}

	// Filter by modes
	if len(filter_.Modes) > 0 {
		if !lo.Contains(filter_.Modes, i.Mode) {
			return false
		}
	}

	// Filter by category IDs
	if len(filter_.CategoryIDs) > 0 {
		hasCategory := false
		for _, category := range i.Categories {
			if lo.Contains(filter_.CategoryIDs, category.ID) {
				hasCategory = true
				break
			}
		}
		if !hasCategory {
			return false
		}
	}

	// Filter by duration in weeks
	if filter_.DurationInWeeks > 0 {
		if i.DurationInWeeks != filter_.DurationInWeeks {
			return false
		}
	}

	// Filter by price range
	if filter_.MinPrice != nil && i.Price.LessThan(*filter_.MinPrice) {
		return false
	}
	if filter_.MaxPrice != nil && i.Price.GreaterThan(*filter_.MaxPrice) {
		return false
	}

	// Filter by status - if no status is specified, only show active internships
	if filter_.GetStatus() != "" {
		if string(i.Status) != filter_.GetStatus() {
			return false
		}
	} else if i.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && i.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && i.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// internshipSortFn implements sorting logic for internships
func internshipSortFn(i, j *internship.Internship) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryInternshipStore) Create(ctx context.Context, i *internship.Internship) error {
	if i == nil {
		return ierr.NewError("internship cannot be nil").
			WithHint("Internship data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if i.CreatedAt.IsZero() {
		i.CreatedAt = now
	}
	if i.UpdatedAt.IsZero() {
		i.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, i.ID, i)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("An internship with this ID already exists").
				WithReportableDetails(map[string]any{
					"internship_id": i.ID,
					"title":         i.Title,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create internship").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipStore) Get(ctx context.Context, id string) (*internship.Internship, error) {
	internship, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Internship with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"internship_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get internship with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return internship, nil
}

func (s *InMemoryInternshipStore) GetByLookupKey(ctx context.Context, lookupKey string) (*internship.Internship, error) {
	internships, err := s.InMemoryStore.List(ctx, nil, internshipFilterFn, internshipSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get internship by lookup key").
			Mark(ierr.ErrDatabase)
	}

	for _, i := range internships {
		if i.LookupKey == lookupKey && i.Status != types.StatusDeleted {
			return i, nil
		}
	}

	return nil, ierr.NewError("internship not found").
		WithHintf("Internship with lookup key %s was not found", lookupKey).
		WithReportableDetails(map[string]any{
			"lookup_key": lookupKey,
		}).
		Mark(ierr.ErrNotFound)
}

func (s *InMemoryInternshipStore) Update(ctx context.Context, i *internship.Internship) error {
	if i == nil {
		return ierr.NewError("internship cannot be nil").
			WithHint("Internship data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	i.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, i.ID, i)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Internship with ID %s was not found", i.ID).
				WithReportableDetails(map[string]any{
					"internship_id": i.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update internship with ID %s", i.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipStore) Delete(ctx context.Context, id string) error {
	// Get the internship first
	i, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	i.Status = types.StatusDeleted
	i.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, i)
}

func (s *InMemoryInternshipStore) Count(ctx context.Context, filter *types.InternshipFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, internshipFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count internships").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryInternshipStore) List(ctx context.Context, filter *types.InternshipFilter) ([]*internship.Internship, error) {
	internships, err := s.InMemoryStore.List(ctx, filter, internshipFilterFn, internshipSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list internships").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return internships, nil
}

func (s *InMemoryInternshipStore) ListAll(ctx context.Context, filter *types.InternshipFilter) ([]*internship.Internship, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipFilter()
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
	unlimitedFilter := &types.InternshipFilter{
		QueryFilter:     types.NewNoLimitQueryFilter(),
		TimeRangeFilter: filter.TimeRangeFilter,
		Name:            filter.Name,
		CategoryIDs:     filter.CategoryIDs,
		Levels:          filter.Levels,
		Modes:           filter.Modes,
		InternshipIDs:   filter.InternshipIDs,
		DurationInWeeks: filter.DurationInWeeks,
		MaxPrice:        filter.MaxPrice,
		MinPrice:        filter.MinPrice,
	}

	return s.List(ctx, unlimitedFilter)
}

// Clear clears the internship store
func (s *InMemoryInternshipStore) Clear() {
	s.InMemoryStore.Clear()
}
