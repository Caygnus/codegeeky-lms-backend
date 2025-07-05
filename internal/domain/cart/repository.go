package cart

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

type Repository interface {
	Create(ctx context.Context, cart *Cart) error
	CreateWithLineItems(ctx context.Context, cart *Cart) error
	Get(ctx context.Context, id string) (*Cart, error)
	Update(ctx context.Context, cart *Cart) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *types.CartFilter) (int, error)
	List(ctx context.Context, filter *types.CartFilter) ([]*Cart, error)
	ListAll(ctx context.Context, filter *types.CartFilter) ([]*Cart, error)

	// cart line items
	CreateCartLineItem(ctx context.Context, cartLineItem *CartLineItem) error
	GetCartLineItem(ctx context.Context, id string) (*CartLineItem, error)
	UpdateCartLineItem(ctx context.Context, cartLineItem *CartLineItem) error
	DeleteCartLineItem(ctx context.Context, id string) error
	ListCartLineItems(ctx context.Context, cartId string) ([]*CartLineItem, error)
}
