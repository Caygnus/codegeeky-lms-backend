package internship

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Internship struct {
	ID                 string                `json:"id,omitempty"`
	Title              string                `json:"title,omitempty"`
	LookupKey          string                `json:"lookup_key,omitempty"`
	Description        string                `json:"description,omitempty"`
	Skills             []string              `json:"skills,omitempty"`
	Level              types.InternshipLevel `json:"level,omitempty"`
	Mode               types.InternshipMode  `json:"mode,omitempty"`
	DurationInWeeks    int                   `json:"duration_in_weeks,omitempty"`
	LearningOutcomes   []string              `json:"learning_outcomes,omitempty"`
	Prerequisites      []string              `json:"prerequisites,omitempty"`
	Benefits           []string              `json:"benefits,omitempty"`
	Currency           string                `json:"currency,omitempty"`
	Price              decimal.Decimal       `json:"price,omitempty"`
	FlatDiscount       *decimal.Decimal      `json:"flat_discount,omitempty"`
	PercentageDiscount *decimal.Decimal      `json:"percentage_discount,omitempty"`
	Subtotal           decimal.Decimal       `json:"subtotal,omitempty"`
	Total              decimal.Decimal       `json:"total,omitempty"`
	Categories         []*Category           `json:"categories,omitempty" db:"categories"`

	types.BaseModel
}

func (i *Internship) FromEnt(internship *ent.Internship) *Internship {
	c := &Category{}

	return &Internship{
		ID:                 internship.ID,
		Title:              internship.Title,
		Description:        internship.Description,
		LookupKey:          internship.LookupKey,
		Skills:             internship.Skills,
		Subtotal:           internship.Subtotal,
		Total:              internship.Total,
		Level:              types.InternshipLevel(internship.Level),
		Mode:               types.InternshipMode(internship.Mode),
		DurationInWeeks:    internship.DurationInWeeks,
		LearningOutcomes:   internship.LearningOutcomes,
		Prerequisites:      internship.Prerequisites,
		Benefits:           internship.Benefits,
		Currency:           internship.Currency,
		Price:              internship.Price,
		FlatDiscount:       internship.FlatDiscount,
		PercentageDiscount: internship.PercentageDiscount,
		Categories:         c.FromEntList(internship.Edges.Categories),
		BaseModel: types.BaseModel{
			Status:    types.Status(internship.Status),
			CreatedAt: internship.CreatedAt,
			UpdatedAt: internship.UpdatedAt,
			CreatedBy: internship.CreatedBy,
			UpdatedBy: internship.UpdatedBy,
		},
	}
}

func (i *Internship) FromEntList(internships []*ent.Internship) []*Internship {
	return lo.Map(internships, func(internship *ent.Internship, _ int) *Internship {
		return i.FromEnt(internship)
	})
}
