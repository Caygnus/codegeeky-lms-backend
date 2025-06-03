package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
)

type FileUpload struct {
	ent.Schema
}

func (FileUpload) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{}, // created_at, updated_at, etc.
	}
}

func (FileUpload) Fields() []ent.Field {
	return []ent.Field{
		// UUID primary key with prefix generation
		field.String("id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_FILE_UPLOAD)
			}).
			Immutable(),

		// Basic file info
		field.String("file_name").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		field.String("file_type").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		field.String("extension").
			SchemaType(map[string]string{"postgres": "varchar(50)"}).
			NotEmpty(),

		field.String("mime_type").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		field.String("public_url").
			SchemaType(map[string]string{"postgres": "varchar(512)"}).
			NotEmpty(),

		field.String("secure_url").
			SchemaType(map[string]string{"postgres": "varchar(512)"}).
			Optional().
			Nillable(),

		// Provider info: enum stored as string
		field.String("provider").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		field.String("external_id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		// Size info
		field.Int64("size_bytes").
			Positive(),

		field.String("file_size").
			SchemaType(map[string]string{"postgres": "varchar(50)"}).
			Optional().
			Nillable(),
	}
}

func (FileUpload) Edges() []ent.Edge {
	return nil
}
