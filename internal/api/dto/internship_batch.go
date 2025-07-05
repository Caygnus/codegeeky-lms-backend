package dto

import (
	"context"
	"time"

	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/samber/lo"
)

// CreateInternshipBatchRequest is used for creating a new internship batch via API request
type CreateInternshipBatchRequest struct {
	InternshipID string                      `json:"internship_id" binding:"required"`
	Name         string                      `json:"name" binding:"required,min=3,max=255"`
	Description  string                      `json:"description,omitempty"`
	StartDate    *time.Time                  `json:"start_date,omitempty"`
	EndDate      *time.Time                  `json:"end_date,omitempty"`
	BatchStatus  types.InternshipBatchStatus `json:"batch_status,omitempty"`
}

func (b *CreateInternshipBatchRequest) Validate() error {
	err := validator.ValidateRequest(b)
	if err != nil {
		return ierr.WithError(err).
			WithHint("invalid internship batch create request").
			Mark(ierr.ErrValidation)
	}

	// validate batch status if provided
	if b.BatchStatus != "" {
		if err := validator.ValidateEnums([]types.InternshipBatchStatus{b.BatchStatus}, types.InternshipBatchStatuses, "batch_status"); err != nil {
			return ierr.WithError(err).
				WithHint("invalid batch status").
				Mark(ierr.ErrValidation)
		}
	}

	// validate dates if both are provided
	if b.StartDate != nil && b.EndDate != nil && b.StartDate.After(*b.EndDate) {
		return ierr.NewError("start date must be before end date").
			WithHint("please provide a valid date range").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func (b *CreateInternshipBatchRequest) ToInternshipBatch(ctx context.Context) *domainInternship.InternshipBatch {
	batchStatus := b.BatchStatus
	if batchStatus == "" {
		batchStatus = types.InternshipBatchStatusUpcoming
	}

	return &domainInternship.InternshipBatch{
		ID:           types.GenerateUUIDWithPrefix(types.UUID_PREFIX_INTERNSHIP_BATCH),
		InternshipID: b.InternshipID,
		Name:         b.Name,
		Description:  b.Description,
		StartDate:    lo.FromPtr(b.StartDate),
		EndDate:      lo.FromPtr(b.EndDate),
		BatchStatus:  batchStatus,
		BaseModel:    types.GetDefaultBaseModel(ctx),
	}
}

type InternshipBatchResponse struct {
	domainInternship.InternshipBatch
}

func (b *InternshipBatchResponse) FromDomain(batch *domainInternship.InternshipBatch) *InternshipBatchResponse {
	return &InternshipBatchResponse{
		InternshipBatch: *batch,
	}
}

// UpdateInternshipBatchRequest is used for updating an existing internship batch via API request
type UpdateInternshipBatchRequest struct {
	Name        *string                      `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Description *string                      `json:"description,omitempty"`
	StartDate   *time.Time                   `json:"start_date,omitempty"`
	EndDate     *time.Time                   `json:"end_date,omitempty"`
	BatchStatus *types.InternshipBatchStatus `json:"batch_status,omitempty"`
}

func (b *UpdateInternshipBatchRequest) Validate() error {
	err := validator.ValidateRequest(b)
	if err != nil {
		return ierr.WithError(err).
			WithHint("invalid internship batch update request").
			Mark(ierr.ErrValidation)
	}

	// validate batch status if provided
	if b.BatchStatus != nil {
		if err := validator.ValidateEnums([]types.InternshipBatchStatus{*b.BatchStatus}, types.InternshipBatchStatuses, "batch_status"); err != nil {
			return ierr.WithError(err).
				WithHint("invalid batch status").
				Mark(ierr.ErrValidation)
		}
	}

	// validate dates if both are provided
	if b.StartDate != nil && b.EndDate != nil && b.StartDate.After(*b.EndDate) {
		return ierr.NewError("start date must be before end date").
			WithHint("please provide a valid date range").
			Mark(ierr.ErrValidation)
	}

	return nil
}

type ListInternshipBatchResponse = types.ListResponse[*InternshipBatchResponse]
