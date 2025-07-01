package enrollment

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
)

// Enrollment is the model entity for the Enrollment schema.
type Enrollment struct {
	ID                 string                 `json:"id,omitempty"`
	UserID             string                 `json:"user_id,omitempty"`
	InternshipID       string                 `json:"internship_id,omitempty"`
	EnrollmentStatus   types.EnrollmentStatus `json:"enrollment_status,omitempty"`
	PaymentStatus      types.PaymentStatus    `json:"payment_status,omitempty"`
	EnrolledAt         *time.Time             `json:"enrolled_at,omitempty"`
	PaymentID          *string                `json:"payment_id,omitempty"`
	RefundedAt         *time.Time             `json:"refunded_at,omitempty"`
	CancellationReason *string                `json:"cancellation_reason,omitempty"`
	RefundReason       *string                `json:"refund_reason,omitempty"`
	IdempotencyKey     *string                `json:"idempotency_key,omitempty"`
	types.Metadata     `json:"metadata,omitempty"`
	types.BaseModel
}

func FromEnt(ent *ent.Enrollment) *Enrollment {
	return &Enrollment{
		ID:                 ent.ID,
		UserID:             ent.UserID,
		InternshipID:       ent.InternshipID,
		EnrollmentStatus:   ent.EnrollmentStatus,
		PaymentStatus:      ent.PaymentStatus,
		EnrolledAt:         ent.EnrolledAt,
		PaymentID:          ent.PaymentID,
		RefundedAt:         ent.RefundedAt,
		CancellationReason: ent.CancellationReason,
		RefundReason:       ent.RefundReason,
		Metadata:           types.MetadataFromEnt(ent.Metadata),
		BaseModel: types.BaseModel{
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
		},
	}
}

func FromEntList(ents []*ent.Enrollment) []*Enrollment {
	enrollments := make([]*Enrollment, len(ents))
	for i, ent := range ents {
		enrollments[i] = FromEnt(ent)
	}
	return enrollments
}
