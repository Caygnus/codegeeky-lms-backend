package cart

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// CartLineItems is the model entity for the CartLineItems schema.
type CartLineItem struct {
	ID             string                       `json:"id,omitempty"`
	CartID         string                       `json:"cart_id,omitempty"`
	EntityID       string                       `json:"entity_id,omitempty"`
	EntityType     types.CartLineItemEntityType `json:"entity_type,omitempty"`
	Quantity       int                          `json:"quantity,omitempty"`
	PerUnitPrice   decimal.Decimal              `json:"per_unit_price,omitempty"`
	TaxAmount      decimal.Decimal              `json:"tax_amount,omitempty"`
	DiscountAmount decimal.Decimal              `json:"discount_amount,omitempty"`
	Subtotal       decimal.Decimal              `json:"subtotal,omitempty"`
	Total          decimal.Decimal              `json:"total,omitempty"`
	Metadata       map[string]string            `json:"metadata,omitempty"`
	types.BaseModel
}

func (c *CartLineItem) FromEnt(ent *ent.CartLineItems) *CartLineItem {
	return &CartLineItem{
		ID:             ent.ID,
		CartID:         ent.CartID,
		EntityID:       ent.EntityID,
		EntityType:     types.CartLineItemEntityType(ent.EntityType),
		Quantity:       ent.Quantity,
		PerUnitPrice:   ent.PerUnitPrice,
		TaxAmount:      ent.TaxAmount,
		DiscountAmount: ent.DiscountAmount,
		Subtotal:       ent.Subtotal,
		Total:          ent.Total,
		Metadata:       ent.Metadata,
		BaseModel: types.BaseModel{
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
			Status:    types.Status(ent.Status),
		},
	}
}

func (c *CartLineItem) FromEntList(ents []*ent.CartLineItems) []*CartLineItem {
	return lo.Map(ents, func(ent *ent.CartLineItems, _ int) *CartLineItem {
		return c.FromEnt(ent)
	})
}
