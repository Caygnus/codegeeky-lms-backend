package testutil

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/user"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

// InMemoryUserStore implements user.Repository
type InMemoryUserStore struct {
	*InMemoryStore[*user.User]
}

// NewInMemoryUserStore creates a new in-memory user store
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		InMemoryStore: NewInMemoryStore[*user.User](),
	}
}

// userFilterFn implements filtering logic for users
func userFilterFn(ctx context.Context, u *user.User, filter interface{}) bool {
	if u == nil {
		return false
	}

	filter_, ok := filter.(*types.UserFilter)
	if !ok {
		return true // No filter applied
	}

	// Filter by email
	if len(filter_.Email) > 0 {
		if !lo.Contains(filter_.Email, u.Email) {
			return false
		}
	}

	// Filter by phone
	if len(filter_.Phone) > 0 {
		if !lo.Contains(filter_.Phone, u.Phone) {
			return false
		}
	}

	// Filter by roles
	if len(filter_.Roles) > 0 {
		if !lo.Contains(filter_.Roles, string(u.Role)) {
			return false
		}
	}

	// Filter by status - if no status is specified, only show active users
	if filter_.Status != "" {
		if u.Status != filter_.Status {
			return false
		}
	} else if u.Status == types.StatusDeleted {
		return false
	}

	// Filter by time range
	if filter_.TimeRangeFilter != nil {
		if filter_.StartTime != nil && u.CreatedAt.Before(*filter_.StartTime) {
			return false
		}
		if filter_.EndTime != nil && u.CreatedAt.After(*filter_.EndTime) {
			return false
		}
	}

	return true
}

// userSortFn implements sorting logic for users
func userSortFn(i, j *user.User) bool {
	if i == nil || j == nil {
		return false
	}
	return i.CreatedAt.After(j.CreatedAt)
}

func (s *InMemoryUserStore) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return ierr.NewError("user cannot be nil").
			WithHint("User data is required").
			Mark(ierr.ErrValidation)
	}

	// Set timestamps
	now := time.Now().UTC()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	err := s.InMemoryStore.Create(ctx, u.ID, u)
	if err != nil {
		if err.Error() == "item already exists" {
			return ierr.WithError(err).
				WithHint("A user with this ID already exists").
				WithReportableDetails(map[string]any{
					"user_id": u.ID,
					"email":   u.Email,
				}).
				Mark(ierr.ErrAlreadyExists)
		}
		return ierr.WithError(err).
			WithHint("Failed to create user").
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryUserStore) Get(ctx context.Context, id string) (*user.User, error) {
	user, err := s.InMemoryStore.Get(ctx, id)
	if err != nil {
		if err.Error() == "item not found" {
			return nil, ierr.WithError(err).
				WithHintf("User with ID %s was not found", id).
				WithReportableDetails(map[string]any{
					"user_id": id,
				}).
				Mark(ierr.ErrNotFound)
		}
		return nil, ierr.WithError(err).
			WithHintf("Failed to get user with ID %s", id).
			Mark(ierr.ErrDatabase)
	}
	return user, nil
}

func (s *InMemoryUserStore) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	users, err := s.InMemoryStore.List(ctx, nil, userFilterFn, userSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get user by email").
			Mark(ierr.ErrDatabase)
	}

	for _, u := range users {
		if u.Email == email && u.Status != types.StatusDeleted {
			return u, nil
		}
	}

	return nil, ierr.NewError("user not found").
		WithHintf("User with email %s was not found", email).
		WithReportableDetails(map[string]any{
			"email": email,
		}).
		Mark(ierr.ErrNotFound)
}

func (s *InMemoryUserStore) List(ctx context.Context, filter *types.UserFilter) ([]*user.User, error) {
	users, err := s.InMemoryStore.List(ctx, filter, userFilterFn, userSortFn)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to list users").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return users, nil
}

func (s *InMemoryUserStore) ListAll(ctx context.Context, filter *types.UserFilter) ([]*user.User, error) {
	if filter == nil {
		filter = types.NewNoLimitUserFilter()
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
	unlimitedFilter := &types.UserFilter{
		QueryFilter:     types.NewNoLimitQueryFilter(),
		TimeRangeFilter: filter.TimeRangeFilter,
		Email:           filter.Email,
		Phone:           filter.Phone,
		Roles:           filter.Roles,
		Status:          filter.Status,
	}

	return s.List(ctx, unlimitedFilter)
}

func (s *InMemoryUserStore) Count(ctx context.Context, filter *types.UserFilter) (int, error) {
	count, err := s.InMemoryStore.Count(ctx, filter, userFilterFn)
	if err != nil {
		return 0, ierr.WithError(err).
			WithHint("Failed to count users").
			WithReportableDetails(map[string]any{
				"filter": filter,
			}).
			Mark(ierr.ErrDatabase)
	}
	return count, nil
}

func (s *InMemoryUserStore) Update(ctx context.Context, u *user.User) error {
	if u == nil {
		return ierr.NewError("user cannot be nil").
			WithHint("User data is required").
			Mark(ierr.ErrValidation)
	}

	// Update timestamp
	u.UpdatedAt = time.Now().UTC()

	err := s.InMemoryStore.Update(ctx, u.ID, u)
	if err != nil {
		if err.Error() == "item not found" {
			return ierr.WithError(err).
				WithHintf("User with ID %s was not found", u.ID).
				WithReportableDetails(map[string]any{
					"user_id": u.ID,
				}).
				Mark(ierr.ErrNotFound)
		}
		return ierr.WithError(err).
			WithHintf("Failed to update user with ID %s", u.ID).
			Mark(ierr.ErrDatabase)
	}
	return nil
}

func (s *InMemoryUserStore) Delete(ctx context.Context, u *user.User) error {
	if u == nil {
		return ierr.NewError("user cannot be nil").
			WithHint("User data is required").
			Mark(ierr.ErrValidation)
	}

	// Soft delete by setting status to deleted
	u.Status = types.StatusDeleted
	u.UpdatedAt = time.Now().UTC()

	return s.Update(ctx, u)
}

// Clear clears the user store
func (s *InMemoryUserStore) Clear() {
	s.InMemoryStore.Clear()
}
