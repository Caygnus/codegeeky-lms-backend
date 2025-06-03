package internship

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type InternshipRepository interface {
	Create(ctx context.Context, internship *Internship) error
	Get(ctx context.Context, id string) (*Internship, error)
	GetByLookupKey(ctx context.Context, lookupKey string) (*Internship, error)
	Update(ctx context.Context, internship *Internship) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.InternshipFilter) (int, error)
	List(ctx context.Context, filter *types.InternshipFilter) ([]*Internship, error)
	ListAll(ctx context.Context, filter *types.InternshipFilter) ([]*Internship, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Get(ctx context.Context, id string) (*Category, error)
	GetByLookupKey(ctx context.Context, lookupKey string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.CategoryFilter) (int, error)
	List(ctx context.Context, filter *types.CategoryFilter) ([]*Category, error)
	ListAll(ctx context.Context, filter *types.CategoryFilter) ([]*Category, error)
}
