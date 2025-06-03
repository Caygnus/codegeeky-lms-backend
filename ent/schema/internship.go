package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// Internship holds the schema definition for the Internship entity.
type Internship struct {
	ent.Schema
}

// Fields of the Internship.
func (Internship) Fields() []ent.Field {
	return []ent.Field{
		// Core Identifiers
		field.String("id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_INTERNSHIP)
			}).
			Immutable(),

		field.String("title").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty(),
		field.Text("description").
			SchemaType(map[string]string{
				"postgres": "text",
			}).
			NotEmpty(),

		// Content & Categorization
		field.JSON("skills", []string{}).
			Optional().
			Default([]string{}).
			Comment("List of required skills"),

		field.String("level").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			Optional().
			Comment("Level of the internship: beginner, intermediate, advanced"),

		// Delivery Format
		field.String("mode").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty().
			Comment("Internship mode: remote, hybrid, onsite"),

		field.Int("duration_in_weeks").
			Optional().
			Comment("Alternative to months for shorter internships"),

		// Learning Metadata
		field.JSON("learning_outcomes", []string{}).
			Optional().
			Comment("What students will learn in the internship"),

		field.JSON("prerequisites", []string{}).
			Optional().
			Comment("Prerequisites or recommended knowledge"),

		field.JSON("benefits", []string{}).
			Optional().
			Comment("Benefits of the internship"),

		// Pricing & Discounts
		field.String("currency").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			Optional().
			Comment("Currency of the internship"),

		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "decimal(10,2)",
			}).
			Optional().
			Comment("Price of the internship"),

		field.Other("flat_discount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "decimal(10,2)",
			}).
			Optional().
			Comment("Flat discount on the internship"),

		field.Other("percentage_discount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "decimal(10,2)",
			}).
			Optional().
			Comment("Percentage discount on the internship"),
	}
}
func (Internship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("categories", Category.Type).
			StorageKey(edge.Column("category_id"), edge.Column("internship_id")),
	}
}
