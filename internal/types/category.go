package types

type CategoryFilter struct {
	*QueryFilter
	*TimeRangeFilter

	// These fields are used to filter categories by name
	Name          string   `json:"name,omitempty" form:"name" validate:"omitempty"`
	CategoryIDs   []string `json:"category_ids,omitempty" form:"category_ids" validate:"omitempty"`
	InternshipIDs []string `json:"internship_ids,omitempty" form:"internship_ids" validate:"omitempty"`
}

func (f *CategoryFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	return nil
}

func NewCategoryFilter() *CategoryFilter {
	return &CategoryFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitCategoryFilter() *CategoryFilter {
	return &CategoryFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *CategoryFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *CategoryFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *CategoryFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *CategoryFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *CategoryFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *CategoryFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *CategoryFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
