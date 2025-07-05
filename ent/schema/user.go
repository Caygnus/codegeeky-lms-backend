package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
)

type User struct {
	ent.Schema
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_USER)
			}).
			Immutable(),
		field.String("full_name").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty(),
		field.String("email").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty(),
		field.String("phone_number").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}),
		field.String("role").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			NotEmpty(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// Cart
		// one user can have many carts
		// one cart can have only one user
		edge.To("carts", Cart.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email", "phone_number").
			Unique(),
		index.Fields("email"),
		index.Fields("phone_number"),
	}
}
