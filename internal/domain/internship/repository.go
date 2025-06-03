package internship

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Repository interface {
	Create(ctx context.Context, internship *Internship) error
	Get(ctx context.Context, id string) (*Internship, error)
	GetByLookupKey(ctx context.Context, lookupKey string) (*Internship, error)
	Update(ctx context.Context, internship *Internship) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.InternshipFilter) (int, error)
	List(ctx context.Context, filter *types.InternshipFilter) ([]*Internship, error)
	ListAll(ctx context.Context, filter *types.InternshipFilter) ([]*Internship, error)
}
