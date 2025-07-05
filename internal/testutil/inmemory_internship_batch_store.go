package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryInternshipBatchStore implements internship.InternshipBatchRepository
type InMemoryInternshipBatchStore struct {
	*InMemoryStore[*internship.InternshipBatch]
}

// NewInMemoryInternshipBatchStore creates a new in-memory internship batch store
func NewInMemoryInternshipBatchStore() *InMemoryInternshipBatchStore {
	return &InMemoryInternshipBatchStore{
		InMemoryStore: NewInMemoryStore[*internship.InternshipBatch](),
	}
}

// internshipBatchFilterFn implements filtering logic for internship batches
func internshipBatchFilterFn(ctx context.Context, b *internship.InternshipBatch, filter interface{}) bool {
	if b == nil {
		return false
	}

	filter_, ok := filter.(*types.InternshipBatchFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by internship IDs
	if len(filter_.InternshipIDs) > 0 {
		if !lo.Contains(filter_.InternshipIDs, b.InternshipID) {
			return false
		}
	}

	// Filter by name
	if filter_.Name != "" {
		if b.Name != filter_.Name {
			return false
		}
	}

	// Filter by batch status
	if filter_.BatchStatus != "" {
		if b.BatchStatus != filter_.BatchStatus {
			return false
		}
	}

	// Filter by start date
	if filter_.StartDate != nil {
		if b.StartDate.Before(*filter_.StartDate) {
			return false
		}
	}

	// Filter by end date
	if filter_.EndDate != nil {
		if b.EndDate.After(*filter_.EndDate) {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active batches
	if filter_.GetStatus() != "" {
		if string(b.Status) != filter_.GetStatus() {
			return false
		}
	} else if b.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && b.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && b.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// internshipBatchSortFn implements sorting logic for internship batches
func internshipBatchSortFn(i, j *internship.InternshipBatch) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryInternshipBatchStore) Create(ctx context.Context, b *internship.InternshipBatch) error {
	if b == nil {
		return ierr.NewError("internship batch cannot be nil").
			WithHint("Internship batch data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, b.ID, b)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("An internship batch with this ID already exists").
				WithReportableDetails(map[string]any{
					"batch_id":      b.ID,
					"internship_id": b.InternshipID,
					"name":          b.Name,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create internship batch").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipBatchStore) Get(ctx context.Context, id string) (*internship.InternshipBatch, error) {
	batch, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Internship batch with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"batch_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get internship batch with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return batch, nil
}

func (s *InMemoryInternshipBatchStore) Update(ctx context.Context, b *internship.InternshipBatch) error {
	if b == nil {
		return ierr.NewError("internship batch cannot be nil").
			WithHint("Internship batch data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	b.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, b.ID, b)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Internship batch with ID %s was not found", b.ID).
				WithReportableDetails(map[string]any{
					"batch_id": b.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update internship batch with ID %s", b.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipBatchStore) Delete(ctx context.Context, id string) error {
	// Get the batch first
	b, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	b.Status = types.StatusDeleted
	b.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, b)
}

func (s *InMemoryInternshipBatchStore) Count(ctx context.Context, filter *types.InternshipBatchFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, internshipBatchFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count internship batches").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryInternshipBatchStore) List(ctx context.Context, filter *types.InternshipBatchFilter) ([]*internship.InternshipBatch, error) {
	batches, err := s.InMemoryStore.List(ctx, filter, internshipBatchFilterFn, internshipBatchSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list internship batches").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return batches, nil
}

func (s *InMemoryInternshipBatchStore) ListAll(ctx context.Context, filter *types.InternshipBatchFilter) ([]*internship.InternshipBatch, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipBatchFilter()
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
	unlimitedFilter := &types.InternshipBatchFilter{
		QueryFilter:     types.NewNoLimitQueryFilter(),
		TimeRangeFilter: filter.TimeRangeFilter,
		InternshipIDs:   filter.InternshipIDs,
		Name:            filter.Name,
		BatchStatus:     filter.BatchStatus,
		StartDate:       filter.StartDate,
		EndDate:         filter.EndDate,
	}

	return s.List(ctx, unlimitedFilter)
}

// Clear clears the internship batch store
func (s *InMemoryInternshipBatchStore) Clear() {
	s.InMemoryStore.Clear()
}
