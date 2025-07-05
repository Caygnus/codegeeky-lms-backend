package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// BaseMixin implements the ent.Mixin for sharing base fields with package schemas.
type MetadataMixin struct {
	mixin.Schema
}

// Fields of the BaseMixin.
func (MetadataMixin) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("metadata", map[string]string{}).
			Default(map[string]string{}).
			Optional(),
	}
}

// Hooks of the BaseMixin.
func (MetadataMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		// Add hooks for updating updated_at and updated_by
	}
}
