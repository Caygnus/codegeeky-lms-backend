package dto

import (
	"context"
	"time"

	domainCart "github.com/omkar273/codegeeky/internal/domain/cart"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/shopspring/decimal"
)

type CreateCartRequest struct {
	Type      types.CartType               `json:"type" validate:"required"`
	ExpiresAt *time.Time                   `json:"expires_at" validate:"omitempty"`
	Metadata  map[string]string            `json:"metadata" validate:"omitempty"`
	LineItems []*CreateCartLineItemRequest `json:"line_items" validate:"omitempty"`
}

func (c *CreateCartRequest) Validate() error {

	if err := validator.ValidateRequest(c); err != nil {
		return err
	}

	if c.ExpiresAt != nil {
		if c.ExpiresAt.Before(time.Now()) {
			return ierr.NewError("expires_at must be in the future").
				WithHint("expires_at must be in the future").
				Mark(ierr.ErrValidation)
		}
	}

	if len(c.LineItems) > 0 {
		for _, lineItem := range c.LineItems {
			if err := lineItem.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CreateCartRequest) ToCart(ctx context.Context) *domainCart.Cart {

	lineItems := make([]*domainCart.CartLineItem, 0)
	if len(c.LineItems) > 0 {
		for _, lineItem := range c.LineItems {
			lineItems = append(lineItems, lineItem.ToCartLineItem(ctx))
		}
	}

	return &domainCart.Cart{
		ID:             types.GenerateUUIDWithPrefix(types.UUID_PREFIX_CART),
		UserID:         types.GetUserID(ctx),
		Type:           c.Type,
		Subtotal:       decimal.Zero,
		DiscountAmount: decimal.Zero,
		TaxAmount:      decimal.Zero,
		Total:          decimal.Zero,
		ExpiresAt:      c.ExpiresAt,
		Metadata:       c.Metadata,
		LineItems:      lineItems,
		BaseModel:      types.GetDefaultBaseModel(ctx),
	}
}

type UpdateCartRequest struct {
	ID        string                       `json:"id,omitempty"`
	ExpiresAt time.Time                    `json:"expires_at,omitempty"`
	Metadata  map[string]string            `json:"metadata,omitempty"`
	LineItems []*CreateCartLineItemRequest `json:"line_items,omitempty"`
}

type CreateCartLineItemRequest struct {
	CartID     string                       `json:"cart_id" validate:"omitempty"`
	EntityID   string                       `json:"entity_id" validate:"required"`
	EntityType types.CartLineItemEntityType `json:"entity_type" validate:"required"`
	Quantity   int                          `json:"quantity" validate:"required,min=1"`
	Metadata   map[string]string            `json:"metadata" validate:"omitempty"`
}

func (c *CreateCartLineItemRequest) Validate() error {
	if err := validator.ValidateRequest(c); err != nil {
		return err
	}

	if c.Quantity <= 0 {
		return ierr.NewError("quantity must be greater than 0").
			WithHint("quantity must be greater than 0").
			Mark(ierr.ErrValidation)
	}

	if c.EntityID == "" {
		return ierr.NewError("entity_id is required").
			WithHint("entity_id is required").
			Mark(ierr.ErrValidation)
	}

	if err := c.EntityType.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *CreateCartLineItemRequest) ToCartLineItem(ctx context.Context) *domainCart.CartLineItem {
	return &domainCart.CartLineItem{
		ID:             types.GenerateUUIDWithPrefix(types.UUID_PREFIX_CART_LINE_ITEM),
		CartID:         c.CartID,
		EntityID:       c.EntityID,
		EntityType:     c.EntityType,
		Quantity:       c.Quantity,
		PerUnitPrice:   decimal.Zero,
		TaxAmount:      decimal.Zero,
		DiscountAmount: decimal.Zero,
		Subtotal:       decimal.Zero,
		Total:          decimal.Zero,
		Metadata:       c.Metadata,
		BaseModel:      types.GetDefaultBaseModel(ctx),
	}
}

type CartResponse struct {
	*domainCart.Cart
}

func (c *CartResponse) FromDomain(cart *domainCart.Cart) *CartResponse {
	return &CartResponse{
		Cart: cart,
	}
}

type ListCartResponse = types.ListResponse[*CartResponse]
