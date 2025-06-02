package ent

import (
	"context"

	"github.com/omkar273/police/internal/types"
)

// BaseQueryOptions defines the minimal interface that all query options must implement
type BaseQueryOptions[T any, F types.BaseFilter] interface {
	// Core filtering methods
	ApplyStatusFilter(query T, status string) T
	ApplySortFilter(query T, field string, order string) T
	ApplyPaginationFilter(query T, limit int, offset int) T

	// Field name mapping
	GetFieldName(field string) string

	// Combined filter application
	ApplyBaseFilters(ctx context.Context, query T, filter F) T
}

// EntityQueryOptions defines additional methods for entity-specific filtering
type EntityQueryOptions[T any, F types.BaseFilter] interface {
	BaseQueryOptions[T, F]

	// Entity-specific filter application
	ApplyEntityQueryOptions(ctx context.Context, filter F, query T) T
}

// QueryOptionsHelper provides utility methods for common query operations
type QueryOptionsHelper struct{}

// ValidateSort checks if the sort field and order are valid
func (h QueryOptionsHelper) ValidateSort(field, order string) (string, string) {
	// Validate order
	if order != types.OrderAsc && order != types.OrderDesc {
		order = types.OrderDesc // Default to descending
	}

	// Default field if empty
	if field == "" {
		field = "created_at"
	}

	return field, order
}

// ValidatePagination ensures pagination values are within acceptable ranges
func (h QueryOptionsHelper) ValidatePagination(limit, offset int) (int, int) {
	// Ensure minimum values
	if limit <= 0 {
		limit = types.FILTER_DEFAULT_LIMIT
	}
	if offset < 0 {
		offset = 0
	}

	// Ensure maximum limits (from validation rules in QueryFilter)
	if limit > 1000 {
		limit = 1000
	}

	return limit, offset
}

// GetDefaultStatus returns the default status filter (excludes deleted)
func (h QueryOptionsHelper) GetDefaultStatus() string {
	return string(types.StatusPublished)
}
