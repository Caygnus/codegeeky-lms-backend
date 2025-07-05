package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/internshipenrollment"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryInternshipEnrollmentStore implements internshipenrollment.Repository
type InMemoryInternshipEnrollmentStore struct {
	*InMemoryStore[*internshipenrollment.InternshipEnrollment]
}

// NewInMemoryInternshipEnrollmentStore creates a new in-memory internship enrollment store
func NewInMemoryInternshipEnrollmentStore() *InMemoryInternshipEnrollmentStore {
	return &InMemoryInternshipEnrollmentStore{
		InMemoryStore: NewInMemoryStore[*internshipenrollment.InternshipEnrollment](),
	}
}

// internshipEnrollmentFilterFn implements filtering logic for internship enrollments
func internshipEnrollmentFilterFn(ctx context.Context, e *internshipenrollment.InternshipEnrollment, filter interface{}) bool {
	if e == nil {
		return false
	}

	filter_, ok := filter.(*types.InternshipEnrollmentFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by user ID
	if filter_.UserID != "" {
		if e.UserID != filter_.UserID {
			return false
		}
	}

	// Filter by internship IDs
	if len(filter_.InternshipIDs) > 0 {
		if !lo.Contains(filter_.InternshipIDs, e.InternshipID) {
			return false
		}
	}

	// Filter by enrollment IDs
	if len(filter_.EnrollmentIDs) > 0 {
		if !lo.Contains(filter_.EnrollmentIDs, e.ID) {
			return false
		}
	}

	// Filter by enrollment status
	if filter_.EnrollmentStatus != "" {
		if e.EnrollmentStatus != filter_.EnrollmentStatus {
			return false
		}
	}

	// Filter by payment status
	if filter_.PaymentStatus != "" {
		if e.PaymentStatus != filter_.PaymentStatus {
			return false
		}
	}

	// Filter by payment ID
	if filter_.PaymentID != nil {
		if e.PaymentID == nil || *e.PaymentID != *filter_.PaymentID {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active enrollments
	if filter_.GetStatus() != "" {
		if string(e.Status) != filter_.GetStatus() {
			return false
		}
	} else if e.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && e.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && e.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// internshipEnrollmentSortFn implements sorting logic for internship enrollments
func internshipEnrollmentSortFn(i, j *internshipenrollment.InternshipEnrollment) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryInternshipEnrollmentStore) Create(ctx context.Context, e *internshipenrollment.InternshipEnrollment) error {
	if e == nil {
		return ierr.NewError("internship enrollment cannot be nil").
			WithHint("Internship enrollment data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, e.ID, e)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("An internship enrollment with this ID already exists").
				WithReportableDetails(map[string]any{
					"enrollment_id": e.ID,
					"user_id":       e.UserID,
					"internship_id": e.InternshipID,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create internship enrollment").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipEnrollmentStore) Get(ctx context.Context, id string) (*internshipenrollment.InternshipEnrollment, error) {
	enrollment, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("Internship enrollment with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"enrollment_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get internship enrollment with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return enrollment, nil
}

func (s *InMemoryInternshipEnrollmentStore) Update(ctx context.Context, e *internshipenrollment.InternshipEnrollment) error {
	if e == nil {
		return ierr.NewError("internship enrollment cannot be nil").
			WithHint("Internship enrollment data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	e.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, e.ID, e)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("Internship enrollment with ID %s was not found", e.ID).
				WithReportableDetails(map[string]any{
					"enrollment_id": e.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update internship enrollment with ID %s", e.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryInternshipEnrollmentStore) Delete(ctx context.Context, id string) error {
	// Get the enrollment first
	e, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete by setting status to deleted
	e.Status = types.StatusDeleted
	e.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, e)
}

func (s *InMemoryInternshipEnrollmentStore) Count(ctx context.Context, filter *types.InternshipEnrollmentFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, internshipEnrollmentFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count internship enrollments").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryInternshipEnrollmentStore) List(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*internshipenrollment.InternshipEnrollment, error) {
	enrollments, err := s.InMemoryStore.List(ctx, filter, internshipEnrollmentFilterFn, internshipEnrollmentSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list internship enrollments").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return enrollments, nil
}

func (s *InMemoryInternshipEnrollmentStore) ListAll(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*internshipenrollment.InternshipEnrollment, error) {
	if filter == nil {
		filter = types.NewNoLimitInternshipEnrollmentFilter()
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
	unlimitedFilter := &types.InternshipEnrollmentFilter{
		QueryFilter:       types.NewNoLimitQueryFilter(),
		TimeRangeFilter:   filter.TimeRangeFilter,
		InternshipIDs:     filter.InternshipIDs,
		UserID:            filter.UserID,
		EnrollmentStatus:  filter.EnrollmentStatus,
		PaymentStatus:     filter.PaymentStatus,
		EnrollmentIDs:     filter.EnrollmentIDs,
		PaymentID:         filter.PaymentID,
		InternshipBatchID: filter.InternshipBatchID,
	}

	return s.List(ctx, unlimitedFilter)
}

func (s *InMemoryInternshipEnrollmentStore) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*internshipenrollment.InternshipEnrollment, error) {
	enrollments, err := s.InMemoryStore.List(ctx, nil, internshipEnrollmentFilterFn, internshipEnrollmentSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get internship enrollment by idempotency key").
			Mark(ierr.ErrDatabase)
	}

	for _, e := range enrollments {
		if e.IdempotencyKey != nil && *e.IdempotencyKey == idempotencyKey && e.Status != types.StatusDeleted {
			return e, nil
		}
	}

	return nil, ierr.NewError("internship enrollment not found").
		WithHintf("Internship enrollment with idempotency key %s was not found", idempotencyKey).
		WithReportableDetails(map[string]any{
			"idempotency_key": idempotencyKey,
		}).
		Mark(ierr.ErrNotFound)
}

// Clear clears the internship enrollment store
func (s *InMemoryInternshipEnrollmentStore) Clear() {
	s.InMemoryStore.Clear()
}
