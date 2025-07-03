package types

import (
	"time"

	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/samber/lo"
)

type InternshipBatchStatus string

const (
	InternshipBatchStatusUpcoming  InternshipBatchStatus = "upcoming"
	InternshipBatchStatusOngoing   InternshipBatchStatus = "ongoing"
	InternshipBatchStatusCompleted InternshipBatchStatus = "completed"
	InternshipBatchStatusCancelled InternshipBatchStatus = "cancelled"
)

var InternshipBatchStatuses = []InternshipBatchStatus{
	InternshipBatchStatusUpcoming,
	InternshipBatchStatusOngoing,
	InternshipBatchStatusCompleted,
	InternshipBatchStatusCancelled,
}

func (s InternshipBatchStatus) Validate() error {
	if !lo.Contains(InternshipBatchStatuses, s) {
		return ierr.NewErrorf("invalid internship batch status").
			WithReportableDetails(map[string]any{"status": s}).
			Mark(ierr.ErrValidation)
	}

	return nil
}

type InternshipBatchFilter struct {
	*QueryFilter
	*TimeRangeFilter

	// These fields are used to filter internships by internship id
	InternshipIDs []string `json:"internship_ids,omitempty" form:"internship_ids" validate:"omitempty"`

	// These fields are used to filter internships by name
	Name string `json:"name,omitempty" form:"name" validate:"omitempty"`

	// These fields are used to filter internships by start and end date
	BatchStatus InternshipBatchStatus `json:"batch_status,omitempty" form:"batch_status" validate:"omitempty"`
	StartDate   *time.Time            `json:"start_date,omitempty" form:"start_date" validate:"omitempty"`
	EndDate     *time.Time            `json:"end_date,omitempty" form:"end_date" validate:"omitempty"`
}

func (f *InternshipBatchFilter) Validate() error {
	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	if f.StartDate != nil && f.EndDate != nil && f.StartDate.After(*f.EndDate) {
		return ierr.NewErrorf("start date must be before end date").
			WithReportableDetails(map[string]any{
				"start_date": f.StartDate,
				"end_date":   f.EndDate,
			}).
			Mark(ierr.ErrValidation)
	}

	if f.BatchStatus != "" {
		if err := f.BatchStatus.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func NewInternshipBatchFilter() *InternshipBatchFilter {
	return &InternshipBatchFilter{
		QueryFilter:     NewDefaultQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

func NewNoLimitInternshipBatchFilter() *InternshipBatchFilter {
	return &InternshipBatchFilter{
		QueryFilter:     NewNoLimitQueryFilter(),
		TimeRangeFilter: &TimeRangeFilter{},
	}
}

// GetLimit implements BaseFilter interface
func (f *InternshipBatchFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *InternshipBatchFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetStatus implements BaseFilter interface
func (f *InternshipBatchFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

// GetSort implements BaseFilter interface
func (f *InternshipBatchFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

// GetOrder implements BaseFilter interface
func (f *InternshipBatchFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

// GetExpand implements BaseFilter interface
func (f *InternshipBatchFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

func (f *InternshipBatchFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
