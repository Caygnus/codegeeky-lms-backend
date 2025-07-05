package dto

import (
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
)

// InitializeEnrollmentRequest is the request for initializing an internship enrollment
type InitializeEnrollmentRequest struct {
	InternshipBatchID string            `json:"internship_batch_id" validate:"required"`
	CouponCodes       []string          `json:"coupon_codes,omitempty"`
	SuccessURL        string            `json:"success_url" validate:"required,url"`
	CancelURL         string            `json:"cancel_url" validate:"required,url"`
	Metadata          map[string]string `json:"metadata,omitempty" validate:"omitempty"`
}

func (r *InitializeEnrollmentRequest) Validate() error {
	if err := validator.ValidateRequest(r); err != nil {
		return err
	}

	return nil
}

type InitializeEnrollmentResponse struct {
	EnrollmentID     string                           `json:"enrollment_id"`
	EnrollmentStatus types.InternshipEnrollmentStatus `json:"enrollment_status"`
	PaymentRequired  bool                             `json:"payment_required"`
}
