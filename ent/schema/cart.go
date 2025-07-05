package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// Cart holds the schema definition for the Cart entity.
type Cart struct {
	ent.Schema
}

// Mixin of the Cart.
func (Cart) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{},
		baseMixin.MetadataMixin{},
	}
}

// Fields of the Cart.
func (Cart) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_CART)
			}).
			Immutable(),
		field.String("user_id").
			NotEmpty().
			Immutable(),

		field.String("type").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty().
			Immutable(),

		field.Other("subtotal", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero).
			Immutable(),

		field.Other("discount_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero).
			Immutable(),

		field.Other("tax_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero).
			Immutable(),

		field.Other("total", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric",
			}).
			Default(decimal.Zero).
			Immutable(),

		field.Time("expires_at").
			SchemaType(map[string]string{
				"postgres": "timestamp",
			}).
			Immutable(),
	}
}

func (Cart) Edges() []ent.Edge {
	return []ent.Edge{
		// CartLineItems
		// one cart can have many line items
		// one line item can have only one cart
		edge.To("line_items", CartLineItems.Type),

		// defines a many-to-one relationship between Cart and User
		// one user can have many carts
		// one cart can have only one user
		edge.From("user", User.Type).
			Ref("carts").
			Unique().
			Required().
			Immutable().
			// Here we are mapping the user_id field to the user_id column in the users table
			Field("user_id"),
	}
}
