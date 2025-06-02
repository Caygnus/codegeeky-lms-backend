package types

type UserFilter struct {
	*QueryFilter
	*TimeRangeFilter
	// able to expand the user with the following fields
	// - station
	Expand *Expand `json:"expand" form:"expand"`

	// custom filters
	StationIds []string `json:"station_ids" form:"station_ids" validate:"omitempty"`
	Ranks      []string `json:"ranks" form:"ranks" validate:"omitempty"`
	Email      string   `json:"email" form:"email" validate:"omitempty,email"`
	Phone      string   `json:"phone" form:"phone" validate:"omitempty,e164"`
	FullName   string   `json:"full_name" form:"full_name" validate:"omitempty"`
}

func (f *UserFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}
	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	if f.Expand != nil {
		if err := f.Expand.Validate(UserExpandConfig); err != nil {
			return err
		}
	}
	return nil
}

func NewUserFilter() *UserFilter {
	return &UserFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitUserFilter() *UserFilter {
	return &UserFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *UserFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *UserFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *UserFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *UserFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *UserFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *UserFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *UserFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
