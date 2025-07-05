package testutil

import (
	"context"
	"sort"
	"sync"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

// FilterFunc is a generic filter function type
type FilterFunc[T any] func(ctx context.Context, item T, filter interface{}) bool

// SortFunc is a generic sort function type
type SortFunc[T any] func(i, j T) bool

// ExpandFunc is a function type for handling expand operations
type ExpandFunc[T any] func(ctx context.Context, item T, expand types.Expand) (T, error)

// InMemoryStore implements a generic in-memory store with expand support
type InMemoryStore[T any] struct {
	mu       sync.RWMutex
	items    map[string]T
	expandFn ExpandFunc[T]
}

// NewInMemoryStore creates a new InMemoryStore
func NewInMemoryStore[T any]() *InMemoryStore[T] {
	return &InMemoryStore[T]{
		items: make(map[string]T),
	}
}

// NewInMemoryStoreWithExpand creates a new InMemoryStore with expand functionality
func NewInMemoryStoreWithExpand[T any](expandFn ExpandFunc[T]) *InMemoryStore[T] {
	return &InMemoryStore[T]{
		items:    make(map[string]T),
		expandFn: expandFn,
	}
}

// Create adds a new item to the store
func (s *InMemoryStore[T]) Create(ctx context.Context, id string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; exists {
		return ierr.NewError("item already exists").
			WithHint("An item with this ID already exists").
			WithReportableDetails(map[string]any{
				"id": id,
			}).
			Mark(ierr.ErrAlreadyExists)
	}

	s.items[id] = item
	return nil
}

// Get retrieves an item by ID with optional expand support
func (s *InMemoryStore[T]) Get(ctx context.Context, id string) (T, error) {
	return s.GetWithExpand(ctx, id, types.NewExpand(""))
}

// GetWithExpand retrieves an item by ID with expand functionality
func (s *InMemoryStore[T]) GetWithExpand(ctx context.Context, id string, expand types.Expand) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item, exists := s.items[id]; exists {
		if s.expandFn != nil && !expand.IsEmpty() {
			return s.expandFn(ctx, item, expand)
		}
		return item, nil
	}

	var zero T
	return zero, ierr.NewError("item not found").
		WithHintf("Item with ID %s was not found", id).
		WithReportableDetails(map[string]any{
			"id": id,
		}).
		Mark(ierr.ErrNotFound)
}

// List retrieves items based on filter with optional expand support
func (s *InMemoryStore[T]) List(ctx context.Context, filter interface{}, filterFn FilterFunc[T], sortFn SortFunc[T]) ([]T, error) {
	return s.ListWithExpand(ctx, filter, filterFn, sortFn, types.NewExpand(""))
}

// ListWithExpand retrieves items based on filter with expand functionality
func (s *InMemoryStore[T]) ListWithExpand(ctx context.Context, filter interface{}, filterFn FilterFunc[T], sortFn SortFunc[T], expand types.Expand) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []T
	for _, item := range s.items {
		if filterFn == nil || filterFn(ctx, item, filter) {
			// Apply expand if available
			if s.expandFn != nil && !expand.IsEmpty() {
				expandedItem, err := s.expandFn(ctx, item, expand)
				if err != nil {
					return nil, err
				}
				result = append(result, expandedItem)
			} else {
				result = append(result, item)
			}
		}
	}

	if sortFn != nil {
		sort.Slice(result, func(i, j int) bool {
			return sortFn(result[i], result[j])
		})
	}

	// Apply pagination if filter implements BaseFilter
	if f, ok := filter.(types.BaseFilter); ok && !f.IsUnlimited() {
		start := f.GetOffset()
		if start >= len(result) {
			return []T{}, nil
		}

		end := start + f.GetLimit()
		if end > len(result) {
			end = len(result)
		}
		return result[start:end], nil
	}

	return result, nil
}

// Count returns the total number of items matching the filter
func (s *InMemoryStore[T]) Count(ctx context.Context, filter interface{}, filterFn FilterFunc[T]) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, item := range s.items {
		if filterFn == nil || filterFn(ctx, item, filter) {
			count++
		}
	}

	return count, nil
}

// Update updates an existing item
func (s *InMemoryStore[T]) Update(ctx context.Context, id string, item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return ierr.NewError("item not found").
			WithHintf("Item with ID %s was not found", id).
			WithReportableDetails(map[string]any{
				"id": id,
			}).
			Mark(ierr.ErrNotFound)
	}

	s.items[id] = item
	return nil
}

// Delete removes an item from the store
func (s *InMemoryStore[T]) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return ierr.NewError("item not found").
			WithHintf("Item with ID %s was not found", id).
			WithReportableDetails(map[string]any{
				"id": id,
			}).
			Mark(ierr.ErrNotFound)
	}

	delete(s.items, id)
	return nil
}

// Clear removes all items from the store
func (s *InMemoryStore[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[string]T)
}

// GetAll returns all items in the store (for testing purposes)
func (s *InMemoryStore[T]) GetAll() map[string]T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]T)
	for k, v := range s.items {
		result[k] = v
	}
	return result
}

// SetExpandFunction sets the expand function for the store
func (s *InMemoryStore[T]) SetExpandFunction(expandFn ExpandFunc[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expandFn = expandFn
}

// CheckEnvironmentFilter is a helper function to check if an item matches the environment filter
func CheckEnvironmentFilter(ctx context.Context, itemEnvID string) bool {
	return true
}
