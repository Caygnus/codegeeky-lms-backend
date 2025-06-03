package internship

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Internship struct {

	// ID of the ent.
	ID string `json:"id,omitempty" db:"id"`

	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty" db:"title"`

	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty" db:"description"`

	// LookupKey holds the value of the "lookup_key" field.
	LookupKey string `json:"lookup_key,omitempty" db:"lookup_key"`

	// List of required skills
	Skills []string `json:"skills,omitempty" db:"skills"`

	// Level of the internship: beginner, intermediate, advanced
	Level string `json:"level,omitempty" db:"level"`

	// Internship mode: remote, hybrid, onsite
	Mode types.InternshipMode `json:"mode,omitempty" db:"mode"`

	// Alternative to months for shorter internships
	DurationInWeeks int `json:"duration_in_weeks,omitempty" db:"duration_in_weeks"`

	// What students will learn in the internship
	LearningOutcomes []string `json:"learning_outcomes,omitempty" db:"learning_outcomes"`

	// Prerequisites or recommended knowledge
	Prerequisites []string `json:"prerequisites,omitempty" db:"prerequisites"`

	// Benefits of the internship
	Benefits []string `json:"benefits,omitempty" db:"benefits"`

	// Currency of the internship
	Currency string `json:"currency,omitempty" db:"currency"`

	// Price of the internship
	Price decimal.Decimal `json:"price,omitempty" db:"price"`

	// Flat discount on the internship
	FlatDiscount decimal.Decimal `json:"flat_discount,omitempty" db:"flat_discount"`

	// Percentage discount on the internship
	PercentageDiscount decimal.Decimal `json:"percentage_discount,omitempty" db:"percentage_discount"`

	// Categories holds the value of the categories edge.
	Categories []*Category `json:"categories,omitempty" db:"categories"`

	types.BaseModel
}

func InternshipFromEnt(internship *ent.Internship) *Internship {
	return &Internship{
		ID:                 internship.ID,
		Title:              internship.Title,
		Description:        internship.Description,
		LookupKey:          internship.LookupKey,
		Skills:             internship.Skills,
		Level:              internship.Level,
		Mode:               types.InternshipMode(internship.Mode),
		DurationInWeeks:    internship.DurationInWeeks,
		LearningOutcomes:   internship.LearningOutcomes,
		Prerequisites:      internship.Prerequisites,
		Benefits:           internship.Benefits,
		Currency:           internship.Currency,
		Price:              internship.Price,
		FlatDiscount:       internship.FlatDiscount,
		PercentageDiscount: internship.PercentageDiscount,
		Categories:         CategoryFromEntList(internship.Edges.Categories),
	}
}

func InternshipFromEntList(internships []*ent.Internship) []*Internship {
	return lo.Map(internships, func(internship *ent.Internship, _ int) *Internship {
		return InternshipFromEnt(internship)
	})
}

type Category struct {
	// ID of the ent.
	ID string `json:"id,omitempty" db:"id"`

	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" db:"name"`

	// LookupKey holds the value of the "lookup_key" field.
	LookupKey string `json:"lookup_key,omitempty" db:"lookup_key"`

	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty" db:"description"`

	// Internships holds the value of the internships edge.
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
