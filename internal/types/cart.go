package types

import (
	"fmt"
	"slices"
	"time"

	ierr "github.com/omkar273/codegeeky/internal/errors"
)

type CartType string

const (
	CartTypeOneTime CartType = "onetime"
	CartTypeDefault CartType = "default"
)

func (c CartType) String() string {
	return string(c)
}

func (c CartType) Validate() error {
	allowedValues := []CartType{CartTypeOneTime, CartTypeDefault}

	if !slices.Contains(allowedValues, c) {
		return ierr.NewError("INVALID_CART_TYPE").
			WithHint(fmt.Sprintf("allowed values are: %s", allowedValues)).
			WithReportableDetails(map[string]any{
				"allowed_values": allowedValues,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

type CartLineItemEntityType string

const (
	CartLineItemEntityTypeInternshipBatch CartLineItemEntityType = "internship_batch"
	CartLineItemEntityTypeCourse          CartLineItemEntityType = "course"
)

func (c CartLineItemEntityType) String() string {
	return string(c)
}

func (c CartLineItemEntityType) Validate() error {
	allowedValues := []CartLineItemEntityType{CartLineItemEntityTypeInternshipBatch, CartLineItemEntityTypeCourse}

	if !slices.Contains(allowedValues, c) {
		return ierr.NewError("INVALID_CART_LINE_ITEM_ENTITY_TYPE").
			WithHint(fmt.Sprintf("allowed values are: %s", allowedValues)).
			WithReportableDetails(map[string]any{
				"allowed_values": allowedValues,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

type CartFilter struct {
	*QueryFilter
	*TimeRangeFilter

	// These fields are used to filter carts by user id
	UserID     string                 `json:"user_id,omitempty" form:"user_id" validate:"omitempty"`
	EntityID   string                 `json:"entity_id,omitempty" form:"entity_id" validate:"omitempty"`
	EntityType CartLineItemEntityType `json:"entity_type,omitempty" form:"entity_type" validate:"omitempty"`
	CartType   *CartType              `json:"cart_type,omitempty" form:"cart_type" validate:"omitempty"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty" form:"expires_at" validate:"omitempty"`
}

func (f *CartFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	return nil
}

func NewCartFilter() *CartFilter {
	return &CartFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitCartFilter() *CartFilter {
	return &CartFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *CartFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *CartFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *CartFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *CartFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *CartFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *CartFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *CartFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
