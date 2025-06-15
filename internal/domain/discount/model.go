package discount

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

type Discount struct {
	ID             string             `json:"id,"`
	Code           string             `json:"code"`
	Description    string             `json:"description"`
	DiscountType   types.DiscountType `json:"discount_type"`
	DiscountValue  decimal.Decimal    `json:"discount_value"`
	ValidFrom      time.Time          `json:"valid_from"`
	ValidUntil     *time.Time         `json:"valid_until"`
	IsActive       bool               `json:"is_active"`
	MaxUses        *int               `json:"max_uses"`
	MinOrderValue  *decimal.Decimal   `json:"min_order_value"`
	IsCombinable   bool               `json:"is_combinable"`
	types.Metadata `json:"metadata"`
	types.BaseModel
}

func FromEnt(ent *ent.Discount) *Discount {
	return &Discount{
		ID:            ent.ID,
		Code:          ent.Code,
		Description:   ent.Description,
		DiscountType:  ent.DiscountType,
		DiscountValue: ent.DiscountValue,
		ValidFrom:     ent.ValidFrom,
		ValidUntil:    ent.ValidUntil,
		IsActive:      ent.IsActive,
		MaxUses:       ent.MaxUses,
		MinOrderValue: ent.MinOrderValue,
		IsCombinable:  ent.IsCombinable,
		Metadata:      ent.Metadata,
		BaseModel: types.BaseModel{
			Status:    types.Status(ent.Status),
			CreatedAt: ent.CreatedAt,
			UpdatedAt: ent.UpdatedAt,
			CreatedBy: ent.CreatedBy,
			UpdatedBy: ent.UpdatedBy,
		},
	}
}

func FromEntList(ents []*ent.Discount) []*Discount {
	discounts := make([]*Discount, len(ents))
	for i, ent := range ents {
		discounts[i] = FromEnt(ent)
	}
	return discounts
}
