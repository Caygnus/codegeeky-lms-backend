package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// CartLineItems holds the schema definition for the CartLineItems entity.
type CartLineItems struct {
	ent.Schema
}

// Mixin of the CartLineItems.
func (CartLineItems) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{},
		baseMixin.MetadataMixin{},
	}
}

// Fields of the CartLineItems.
func (CartLineItems) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_CART_LINE_ITEM)
			}).
			Immutable(),

		field.String("cart_id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty().
			Immutable(),

		field.String("entity_id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty().
			Immutable(),

		field.String("entity_type").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty().
			Immutable(),

		field.Int("quantity").
			Default(1),

		field.Other("per_unit_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero),

		field.Other("tax_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero),

		field.Other("discount_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero),

		field.Other("subtotal", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero),

		field.Other("total", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero),
	}
}

func (CartLineItems) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cart", Cart.Type).
			Ref("line_items").
			Unique().
			Required().
			Immutable().
			Field("cart_id"),
	}
}
