package internship

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Category struct {
	// ID of the ent.
	ID string `json:"id,omitempty" db:"id"`

	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" db:"name"`

	// LookupKey holds the value of the "lookup_key" field.
	LookupKey string `json:"lookup_key,omitempty" db:"lookup_key"`

	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty" db:"description"`

	// internships holds the value of the internships edge.
	Internships []*Internship `json:"internships,omitempty" db:"internships"`

	types.BaseModel
}

func CategoryFromEnt(category *ent.Category) *Category {
	return &Category{
		ID:          category.ID,
		Name:        category.Name,
		LookupKey:   category.LookupKey,
		Description: category.Description,
		Internships: InternshipFromEntList(category.Edges.Internships),
		BaseModel: types.BaseModel{
			Status:    types.Status(category.Status),
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			CreatedBy: category.CreatedBy,
			UpdatedBy: category.UpdatedBy,
		},
	}
}

func CategoryFromEntList(categories []*ent.Category) []*Category {
	return lo.Map(categories, func(category *ent.Category, _ int) *Category {
		return CategoryFromEnt(category)
	})
}

func (i *Internship) FinalPrice() decimal.Decimal {
	price := i.Price
	if !i.FlatDiscount.IsZero() && !i.FlatDiscount.IsNegative() {
		price = price.Sub(i.FlatDiscount)
	}
	if !i.PercentageDiscount.IsZero() && !i.PercentageDiscount.IsNegative() {
		discount := price.Mul(i.PercentageDiscount.Div(decimal.NewFromInt(100)))
		price = price.Sub(discount)
	}
	return price
}
