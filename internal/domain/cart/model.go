package cart

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Cart struct {
	ID             string            `json:"id,omitempty"`
	UserID         string            `json:"user_id,omitempty"`
	Type           types.CartType    `json:"type,omitempty"`
	Subtotal       decimal.Decimal   `json:"subtotal,omitempty"`
	DiscountAmount decimal.Decimal   `json:"discount_amount,omitempty"`
	TaxAmount      decimal.Decimal   `json:"tax_amount,omitempty"`
	Total          decimal.Decimal   `json:"total,omitempty"`
	ExpiresAt      time.Time         `json:"expires_at,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	LineItems      []*CartLineItem   `json:"line_items,omitempty"`
	types.BaseModel
}

func (c *Cart) FromEnt(ent *ent.Cart) *Cart {
	li := &CartLineItem{}
	return &Cart{
		ID:             ent.ID,
		UserID:         ent.UserID,
		Type:           types.CartType(ent.Type),
		Subtotal:       ent.Subtotal,
		DiscountAmount: ent.DiscountAmount,
		TaxAmount:      ent.TaxAmount,
		Total:          ent.Total,
		ExpiresAt:      ent.ExpiresAt,
		Metadata:       ent.Metadata,
		LineItems:      li.FromEntList(ent.Edges.LineItems),
		BaseModel: types.BaseModel{
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
			Status:    types.Status(ent.Status),
		},
	}
}

func (c *Cart) FromEntList(ents []*ent.Cart) []*Cart {
	return lo.Map(ents, func(ent *ent.Cart, _ int) *Cart {
		return c.FromEnt(ent)
	})
}
