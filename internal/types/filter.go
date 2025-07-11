package types

import (
	"time"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/samber/lo"
)

// TODO : deprecate
type Filter struct {
	Limit  int    `form:"limit,default=50"`
	Offset int    `form:"offset,default=0"`
	Status Status `form:"status,default=published"`
	Sort   string `form:"sort,default=created_at"`
	Order  string `form:"order,default=desc"`
	Expand string `form:"expand"`
}

const (
	FILTER_DEFAULT_LIMIT  = 50
	FILTER_DEFAULT_STATUS = string(StatusPublished)
	FILTER_DEFAULT_SORT   = "created_at"
	FILTER_DEFAULT_ORDER  = "desc"

	OrderDesc = "desc"
	OrderAsc  = "asc"
)

// TODO : deprecate
func GetDefaultFilter() Filter {
	return Filter{
		Limit:  FILTER_DEFAULT_LIMIT,
		Offset: 0,
		Status: StatusPublished,
		Sort:   FILTER_DEFAULT_SORT,
		Order:  FILTER_DEFAULT_ORDER,
	}
}

// TODO : deprecate
func (f Filter) GetExpand() Expand {
	return NewExpand(f.Expand)
}

// -------- New Base filter ---------

// BaseFilter defines common filtering capabilities
type BaseFilter interface {
	GetLimit() int
	GetOffset() int
	GetStatus() string
	GetSort() string
	GetOrder() string
	GetExpand() Expand
	Validate() error
	IsUnlimited() bool
}

// QueryFilter represents a generic query filter with optional fields
type QueryFilter struct {
	Limit  *int    `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=1000"`
	Offset *int    `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
	Status *Status `json:"status,omitempty" form:"status"`
	Sort   *string `json:"sort,omitempty" form:"sort"`
	Order  *string `json:"order,omitempty" form:"order" validate:"omitempty,oneof=asc desc"`
	Expand *string `json:"expand,omitempty" form:"expand"`
}

// DefaultQueryFilter defines default values for query filters
func NewDefaultQueryFilter() *QueryFilter {
	return &QueryFilter{
		Limit:  lo.ToPtr(50),
		Offset: lo.ToPtr(0),
		Status: nil,
		Sort:   lo.ToPtr("created_at"),
		Order:  lo.ToPtr("desc"),
	}
}

// NewNoLimitPublishedQueryFilter defines a filter with published status and no pagination limits
func NewNoLimitPublishedQueryFilter() *QueryFilter {
	return &QueryFilter{
		Status: lo.ToPtr(StatusPublished),
		Sort:   lo.ToPtr("created_at"),
		Order:  lo.ToPtr("desc"),
	}
}

// NoLimitQueryFilter returns a filter with no pagination limits
func NewNoLimitQueryFilter() *QueryFilter {
	return &QueryFilter{
		Limit:  nil, // No limit for unlimited queries
		Offset: lo.ToPtr(0),
		Status: nil,
		Sort:   lo.ToPtr("created_at"),
		Order:  lo.ToPtr("desc"),
	}
}

// IsUnlimited returns true if this is an unlimited query
func (f QueryFilter) IsUnlimited() bool {
	return f.Limit == nil
}

// GetLimit returns the limit value or default if not set
func (f QueryFilter) GetLimit() int {
	if f.IsUnlimited() {
		return 0 // No limit for unlimited queries
	}
	if f.Limit == nil {
		return *NewDefaultQueryFilter().Limit
	}
	return *f.Limit
}

// GetOffset returns the offset value or default if not set
func (f QueryFilter) GetOffset() int {
	if f.Offset == nil {
		return *NewDefaultQueryFilter().Offset
	}
	return *f.Offset
}

// GetSort returns the sort value or default if not set
func (f QueryFilter) GetSort() string {
	if f.Sort == nil {
		return *NewDefaultQueryFilter().Sort
	}
	return *f.Sort
}

// GetOrder returns the order value or default if not set
func (f QueryFilter) GetOrder() string {
	if f.Order == nil {
		return *NewDefaultQueryFilter().Order
	}
	return *f.Order
}

// GetExpand returns the parsed Expand object from the filter
func (f QueryFilter) GetExpand() Expand {
	if f.Expand == nil {
		return NewExpand("")
	}
	return NewExpand(*f.Expand)
}

// GetStatus returns the status value or default if not set
func (f QueryFilter) GetStatus() string {
	if f.Status == nil {
		return ""
	}
	return string(*f.Status)
}

// Validate validates the filter fields
func (f QueryFilter) Validate() error {
	if !f.IsUnlimited() {
		if f.Limit != nil && (*f.Limit < 1 || *f.Limit > 1000) {
			return ierr.NewError("limit must be between 1 and 1000").WithReportableDetails(
				map[string]any{
					"limit": "must be between 1 and 1000",
				},
			).Mark(ierr.ErrValidation)
		}
	}
	if f.Offset != nil && *f.Offset < 0 {
		return ierr.NewError("offset must be non-negative").WithReportableDetails(
			map[string]any{
				"offset": "must be non-negative",
			},
		).Mark(ierr.ErrValidation)
	}
	if f.Order != nil && *f.Order != "asc" && *f.Order != "desc" {
		return ierr.NewError("order must be either 'asc' or 'desc'").WithReportableDetails(
			map[string]any{
				"order": "must be either 'asc' or 'desc'",
			}).Mark(ierr.ErrValidation)
	}
	return nil
}

// TimeRangeFilter adds time range filtering capabilities
type TimeRangeFilter struct {
	StartTime *time.Time `json:"start_time,omitempty" form:"start_time" validate:"omitempty,time_rfc3339"`
	EndTime   *time.Time `json:"end_time,omitempty" form:"end_time" validate:"omitempty,time_rfc3339"`
}

// Validate validates the time range filter
func (f TimeRangeFilter) Validate() error {
	if f.StartTime != nil && f.EndTime != nil && f.EndTime.Before(*f.StartTime) {
		return ierr.NewError("end_time must be after start_time").
			WithHint("Please provide a valid time range").
			WithReportableDetails(
				map[string]any{
					"end_time": "must be after start_time",
				},
			).
			Mark(ierr.ErrValidation)
	}
	return nil
}

// Merge merges another filter into this one, taking values from other if they are set
func (f *QueryFilter) Merge(other QueryFilter) {
	if other.Limit != nil {
		f.Limit = other.Limit
	}
	if other.Offset != nil {
		f.Offset = other.Offset
	}
	if other.Status != nil {
		f.Status = other.Status
	}
	if other.Sort != nil {
		f.Sort = other.Sort
	}
	if other.Order != nil {
		f.Order = other.Order
	}
	if other.Expand != nil {
		f.Expand = other.Expand
	}
}
