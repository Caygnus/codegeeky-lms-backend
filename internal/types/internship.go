package types

import (
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/shopspring/decimal"
)

type InternshipMode string

const (
	InternshipModeRemote InternshipMode = "remote"
	InternshipModeHybrid InternshipMode = "hybrid"
	InternshipModeOnsite InternshipMode = "onsite"
)

type InternshipLevel string

const (
	InternshipLevelBeginner     InternshipLevel = "beginner"
	InternshipLevelIntermediate InternshipLevel = "intermediate"
	InternshipLevelAdvanced     InternshipLevel = "advanced"
)

var InternshipModes = []InternshipMode{
	InternshipModeRemote,
	InternshipModeHybrid,
	InternshipModeOnsite,
}

var InternshipLevels = []InternshipLevel{
	InternshipLevelBeginner,
	InternshipLevelIntermediate,
	InternshipLevelAdvanced,
}

type InternshipFilter struct {
	*QueryFilter
	*TimeRangeFilter

	// These fields are used to filter internships by category, level and mode
	Name          string            `json:"name,omitempty" form:"name" validate:"omitempty"`
	CategoryIDs   []string          `json:"category_ids,omitempty" form:"category_ids" validate:"omitempty"`
	Levels        []InternshipLevel `json:"levels,omitempty" form:"levels" validate:"omitempty"`
	Modes         []InternshipMode  `json:"modes,omitempty" form:"modes" validate:"omitempty"`
	InternshipIDs []string          `json:"internship_ids,omitempty" form:"internship_ids" validate:"omitempty"`

	// These fields are used to filter internships by duration in weeks
	DurationInWeeks int `json:"duration_in_weeks,omitempty" form:"duration_in_weeks" validate:"omitempty,min=1,max=52"`

	// These fields are used to filter internships by price
	MaxPrice decimal.Decimal `json:"max_price,omitempty" form:"max_price" validate:"omitempty,gt=0,lt=1000000000000,ltfield=MinPrice,gtfield=MinPrice"`
	MinPrice decimal.Decimal `json:"min_price,omitempty" form:"min_price" validate:"omitempty,gt=0,lt=1000000000000,gtfield=MaxPrice,ltfield=MaxPrice"`
}

func (f *InternshipFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	if len(f.Modes) > 0 {
		if err := validator.ValidateEnums(f.Modes, InternshipModes, "mode"); err != nil {
			return err
		}
	}

	if len(f.Levels) > 0 {
		if err := validator.ValidateEnums(f.Levels, InternshipLevels, "level"); err != nil {
			return err
		}
	}

	if !f.MaxPrice.GreaterThan(decimal.Zero) || !f.MinPrice.GreaterThan(decimal.Zero) {
		return ierr.NewErrorf("price must be greater than 0").
			WithReportableDetails(map[string]any{
				"max_price": f.MaxPrice,
				"min_price": f.MinPrice,
			}).
			WithHint("Price must be greater than 0").
			Mark(ierr.ErrValidation)
	}

	if f.MaxPrice.LessThan(f.MinPrice) {
		return ierr.NewErrorf("max price must be greater than min price").
			WithReportableDetails(map[string]any{
				"max_price": f.MaxPrice,
				"min_price": f.MinPrice,
			}).
			WithHint("Max price must be greater than min price").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func NewInternshipFilter() *InternshipFilter {
	return &InternshipFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitInternshipFilter() *InternshipFilter {
	return &InternshipFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *InternshipFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *InternshipFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *InternshipFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *InternshipFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *InternshipFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *InternshipFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *InternshipFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
