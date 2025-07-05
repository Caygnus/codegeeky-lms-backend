package types

import (
	"time"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/shopspring/decimal"
)

type DiscountType string

const (
	DiscountTypeFlat       DiscountType = "flat"
	DiscountTypePercentage DiscountType = "percentage"
)

func (d DiscountType) String() string {
	return string(d)
}

func (d DiscountType) Validate() error {
	switch d {
	case DiscountTypeFlat, DiscountTypePercentage:
		return nil
	default:
		return ierr.NewError("invalid discount type").
			WithHint("discount type must be either flat or percentage").
			Mark(ierr.ErrValidation)
	}
}

type DiscountFilter struct {
	*QueryFilter
	*TimeRangeFilter

	DiscountType  DiscountType     `json:"discount_type,omitempty" form:"discount_type" validate:"omitempty"`
	ValidFrom     *time.Time       `json:"valid_from,omitempty" form:"valid_from" validate:"omitempty"`
	ValidUntil    *time.Time       `json:"valid_until,omitempty" form:"valid_until" validate:"omitempty"`
	MinOrderValue *decimal.Decimal `json:"min_order_value,omitempty" form:"min_order_value" validate:"omitempty"`
	IsCombinable  bool             `json:"is_combinable,omitempty" form:"is_combinable" validate:"omitempty"`
	Codes         []string         `json:"codes,omitempty" form:"codes" validate:"omitempty"`
	DiscountIDs   []string         `json:"discount_ids,omitempty" form:"discount_ids" validate:"omitempty"`
}

func (f *DiscountFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	return nil
}

func NewDiscountFilter() *DiscountFilter {
	return &DiscountFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitDiscountFilter() *DiscountFilter {
	return &DiscountFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *DiscountFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *DiscountFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *DiscountFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *DiscountFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *DiscountFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *DiscountFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *DiscountFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
