package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// Discount holds the schema definition for the Discount entity.
type Discount struct {
	ent.Schema
}

// Mixin of the Discount.
func (Discount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{},
	}
}

// Fields of the Discount.
func (Discount) Fields() []ent.Field {
	return []ent.Field{
		// Primary Identifier
		field.String("id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_DISCOUNT)
			}).
			Immutable(),

		// Discount Code
		field.String("code").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty().
			Unique(),

		// Description for admins
		field.String("description").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			Optional(),

		// Type of discount (flat or percentage)
		field.String("discount_type").
			GoType(types.DiscountType("")).
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			Default(string(types.DiscountTypeFlat)).
			Immutable(),

		// Value of discount (flat amount or percentage)
		field.Other("discount_value", decimal.Decimal{}).
			SchemaType(map[string]string{"postgres": "decimal(10,2)"}).
			Default(decimal.Zero).
			Immutable(),

		// Time range of discount validity
		field.Time("valid_from").
			Default(time.Now),

		field.Time("valid_until").
			Optional().
			Nillable(),

		// Whether discount is active
		field.Bool("is_active").
			Default(true),

		// Maximum times this discount can be used across all users
		field.Int("max_uses").
			Optional().
			Nillable(),

		// Minimum order value required to apply this discount
		field.Other("min_order_value", decimal.Decimal{}).
			Optional().
			Nillable().
			SchemaType(map[string]string{"postgres": "decimal(10,2)"}),

		// Whether this discount can be stacked with others
		field.Bool("is_combinable").
			Default(false).
			Immutable(),

		// Optional tag for internal grouping or analytics
		field.JSON("metadata", map[string]string{}).
			Default(map[string]string{}).
			Optional(),
	}
}
