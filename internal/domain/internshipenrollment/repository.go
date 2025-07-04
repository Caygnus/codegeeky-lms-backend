package internshipenrollment

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Repository interface {
	Create(ctx context.Context, enrollment *InternshipEnrollment) error
	Get(ctx context.Context, id string) (*InternshipEnrollment, error)
	Update(ctx context.Context, enrollment *InternshipEnrollment) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.InternshipEnrollmentFilter) (int, error)
	List(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*InternshipEnrollment, error)
	ListAll(ctx context.Context, filter *types.InternshipEnrollmentFilter) ([]*InternshipEnrollment, error)
	GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*InternshipEnrollment, error)
}
