package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	baseMixin "github.com/omkar273/codegeeky/ent/mixin"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Mixin of the Payment.
func (Payment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin.BaseMixin{},
	}
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Unique().
			Immutable(),

		// idempotency key
		// Optional, for safe retries
		field.String("idempotency_key").
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Unique().
			Immutable(),
		// destination type
		// internship, subscription, charge, etc.
		field.String("destination_type").
			GoType(types.PaymentDestinationType("")).
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Immutable(),

		// destination id
		// ID of the destination (internship ID, subscription ID, etc.)
		field.String("destination_id").
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Immutable(),

		// payment method type
		// upi, card, wallet, etc.
		field.String("payment_method_type").
			GoType(types.PaymentMethodType("")).
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Optional().
			Nillable(),

		// payment method id
		// Saved method ID or token (optional)
		field.String("payment_method_id").
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Optional(),

		// payment gateway provider
		// razorpay, stripe, etc.
		field.String("payment_gateway_provider").
			GoType(types.PaymentGatewayProvider("")).
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Immutable(),

		// gateway payment id
		// Payment ID from the gateway
		field.String("gateway_payment_id").
			SchemaType(map[string]string{
				"postgres": "varchar(255)",
			}).
			Optional().
			Nillable(),

		// amount
		// Amount in smallest unit (e.g. paisa)
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(20,8)",
			}).
			Default(decimal.Zero),

		// currency
		// "INR", "USD", etc.
		field.String("currency").
			GoType(types.Currency("")).
			SchemaType(map[string]string{
				"postgres": "varchar(10)",
			}).
			NotEmpty().
			Immutable(),

		// payment status
		// pending, failed, succeeded, etc.
		field.String("payment_status").
			GoType(types.PaymentStatus("")).
			SchemaType(map[string]string{
				"postgres": "varchar(50)",
			}).
			Default(string(types.PaymentStatusPending)).
			NotEmpty(),

		// track attempts
		// Whether to track payment attempts
		field.Bool("track_attempts").
			Default(true),

		// metadata
		// Additional tracking info (origin, cohort, etc.)
		field.JSON("metadata", map[string]string{}).
			Optional().
			SchemaType(map[string]string{
				"postgres": "jsonb",
			}).
			Default(map[string]string{}),

		// succeeded at
		// Time when the payment succeeded
		field.Time("succeeded_at").
			Optional().
			Nillable(),

		// failed at
		// Time when the payment failed
		field.Time("failed_at").
			Optional().
			Nillable(),

		// refunded at
		// Time when the payment was refunded
		field.Time("refunded_at").
			Optional().
			Nillable(),

		// error message
		// Error message from the payment gateway
		field.String("error_message").
			SchemaType(map[string]string{
				"postgres": "text",
			}).
			Optional().
			Nillable(),
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		// attempts
		// Payment attempts
		edge.To("attempts", PaymentAttempt.Type),
	}
}

// Indexes of the Payment.
func (Payment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("destination_type", "destination_id", "payment_status", "status").
			StorageKey("idx_destination_status"),
		index.Fields("payment_method_type", "payment_method_id", "payment_status", "status").
			StorageKey("idx_tenant_payment_method_status"),
		index.Fields("payment_gateway_provider", "gateway_payment_id").
			StorageKey("idx_gateway_payment").
			Annotations(entsql.IndexWhere("payment_gateway_provider IS NOT NULL AND gateway_payment_id IS NOT NULL")),
	}
}
