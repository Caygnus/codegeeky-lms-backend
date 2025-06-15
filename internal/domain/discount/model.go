package discount

import (
	"time"

	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

type Discount struct {
	ID             string             `json:"id,omitempty"`
	Code           string             `json:"code,omitempty"`
	Description    string             `json:"description,omitempty"`
	DiscountType   types.DiscountType `json:"discount_type,omitempty"`
	DiscountValue  decimal.Decimal    `json:"discount_value,omitempty"`
	ValidFrom      time.Time          `json:"valid_from,omitempty"`
	ValidUntil     *time.Time         `json:"valid_until,omitempty"`
	IsActive       bool               `json:"is_active,omitempty"`
	MaxUses        *int               `json:"max_uses,omitempty"`
	MinOrderValue  *decimal.Decimal   `json:"min_order_value,omitempty"`
	IsCombinable   bool               `json:"is_combinable,omitempty"`
	types.Metadata `json:"metadata,omitempty"`
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
