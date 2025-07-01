package types

// lifecycle of enrollment
// pending -> enrolled -> completed -> refund
//
//		   -> failed
//	       -> cancelled
type EnrollmentStatus string

const (
	EnrollmentStatusPending   = "pending"
	EnrollmentStatusEnrolled  = "enrolled"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusRefunded  = "refunded"
	EnrollmentStatusCancelled = "cancelled"
	EnrollmentStatusFailed    = "failed"
)

type EnrollmentFilter struct {
	*QueryFilter
	*TimeRangeFilter

	InternshipIDs    []string         `json:"internship_ids,omitempty" form:"internship_ids" validate:"omitempty"`
	UserID           string           `json:"user_id,omitempty" form:"user_id" validate:"omitempty"`
	EnrollmentStatus EnrollmentStatus `json:"enrollment_status,omitempty" form:"enrollment_status" validate:"omitempty"`
	PaymentStatus    PaymentStatus    `json:"payment_status,omitempty" form:"payment_status" validate:"omitempty"`
	EnrollmentIDs    []string         `json:"enrollment_ids,omitempty" form:"enrollment_ids" validate:"omitempty"`
	PaymentID        *string          `json:"payment_id,omitempty" form:"payment_id" validate:"omitempty"`
}

func (f *EnrollmentFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	return nil
}

// NewEnrollmentFilter creates a new EnrollmentFilter with default values
func NewEnrollmentFilter() *EnrollmentFilter {
	return &EnrollmentFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// NewNoLimitEnrollmentFilter creates a new EnrollmentFilter with no limit
func NewNoLimitEnrollmentFilter() *EnrollmentFilter {
	return &EnrollmentFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit returns the limit for the EnrollmentFilter
func (f *EnrollmentFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset returns the offset for the EnrollmentFilter
func (f *EnrollmentFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus returns the status for the EnrollmentFilter
func (f *EnrollmentFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *EnrollmentFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *EnrollmentFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *EnrollmentFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *EnrollmentFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
