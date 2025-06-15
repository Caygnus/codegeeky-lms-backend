package discount

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Repository interface {
	Create(ctx context.Context, d *Discount) error
	Get(ctx context.Context, id string) (*Discount, error)
	GetByCode(ctx context.Context, code string) (*Discount, error)
	Update(ctx context.Context, d *Discount) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.DiscountFilter) (int, error)
	List(ctx context.Context, filter *types.DiscountFilter) ([]*Discount, error)
	ListAll(ctx context.Context, filter *types.DiscountFilter) ([]*Discount, error)
}
