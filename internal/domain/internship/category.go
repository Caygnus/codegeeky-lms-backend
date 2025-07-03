package internship

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
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

func (c *Category) FromEnt(category *ent.Category) *Category {
	internship := &Internship{}

	return &Category{
		ID:          category.ID,
		Name:        category.Name,
		LookupKey:   category.LookupKey,
		Description: category.Description,
		Internships: internship.FromEntList(category.Edges.Internships),
		BaseModel: types.BaseModel{
			Status:    types.Status(category.Status),
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
			CreatedBy: category.CreatedBy,
			UpdatedBy: category.UpdatedBy,
		},
	}
}

func (c *Category) FromEntList(categories []*ent.Category) []*Category {
	return lo.Map(categories, func(category *ent.Category, _ int) *Category {
		return c.FromEnt(category)
	})
}
