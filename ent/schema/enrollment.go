package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
)

// Enrollment holds the schema definition for the Enrollment entity.
type Enrollment struct {
	ent.Schema
}

// Mixin of the Enrollment.
func (Enrollment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{}, // includes created_at, updated_at
		baseMixin.MetadataMixin{},
	}
}

// Fields of the Enrollment.
func (Enrollment) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			DefaultFunc(func() string {
				return types.GenerateUUIDWithPrefix(types.UUID_PREFIX_ENROLLMENT)
			}).
			Immutable(),

		// Foreign keys
		field.String("user_id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		field.String("internship_id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			NotEmpty(),

		// Enrollment status
		field.String("enrollment_status").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			GoType(types.EnrollmentStatus("")).
			Default(types.EnrollmentStatusPending).
			NotEmpty(),

		// When enrollment was confirmed
		field.Time("enrolled_at").
			Optional().
			Nillable(),

		// Payment & refund linkage
		// This is the internal payment id of caygnus not the actuall provider id i.e. razorpay , stripe
		field.String("payment_id").
			SchemaType(map[string]string{"postgres": "varchar(255)"}).
			Optional().
			Nillable(),

		field.Time("refunded_at").
			Optional().
			Nillable(),

		// Optional reason for cancellation/refund
		field.String("cancellation_reason").
			Optional().
			Nillable(),

		field.String("refund_reason").
			Optional().
			Nillable(),
	}
}
