package enrollment

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Repository interface {
	Create(ctx context.Context, enrollment *Enrollment) error
	Get(ctx context.Context, id string) (*Enrollment, error)
	Update(ctx context.Context, enrollment *Enrollment) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.EnrollmentFilter) (int, error)
	List(ctx context.Context, filter *types.EnrollmentFilter) ([]*Enrollment, error)
	ListAll(ctx context.Context, filter *types.EnrollmentFilter) ([]*Enrollment, error)
	GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*Enrollment, error)
}
