package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
)

// InternshipBatch holds the schema definition for the InternshipBatch entity.
type InternshipBatch struct {
	ent.Schema
}

func (InternshipBatch) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{}, // includes created_at, updated_at
		baseMixin.MetadataMixin{},
	}
}

func (InternshipBatch) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_INTERNSHIP_BATCH)
			}).
			Immutable().
			Unique(),

		field.String("internship_id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty().
			Immutable(),

		// Name of the batch
		field.String("name").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		// Description of the batch
		field.String("description").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			Optional(),

		// Start date of the batch
		field.Time("start_date").
			SchemaType(map[string]string{"postgres": "timestamp"}).
			Optional(),

		// End date of the batch
		field.Time("end_date").
			SchemaType(map[string]string{"postgres": "timestamp"}).
			Optional(),

		// Status of the batch
		field.String("batch_status").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			Default(string(types.InternshipBatchStatusUpcoming)).
			NotEmpty(),
	}
}
