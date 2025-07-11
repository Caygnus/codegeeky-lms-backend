// Code generated by ent, DO NOT EDIT.

package paymentattempt

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the paymentattempt type in the database.
	Label = "payment_attempt"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldCreatedBy holds the string denoting the created_by field in the database.
	FieldCreatedBy = "created_by"
	// FieldUpdatedBy holds the string denoting the updated_by field in the database.
	FieldUpdatedBy = "updated_by"
	// FieldPaymentID holds the string denoting the payment_id field in the database.
	FieldPaymentID = "payment_id"
	// FieldPaymentStatus holds the string denoting the payment_status field in the database.
	FieldPaymentStatus = "payment_status"
	// FieldAttemptNumber holds the string denoting the attempt_number field in the database.
	FieldAttemptNumber = "attempt_number"
	// FieldGatewayAttemptID holds the string denoting the gateway_attempt_id field in the database.
	FieldGatewayAttemptID = "gateway_attempt_id"
	// FieldErrorMessage holds the string denoting the error_message field in the database.
	FieldErrorMessage = "error_message"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// EdgePayment holds the string denoting the payment edge name in mutations.
	EdgePayment = "payment"
	// Table holds the table name of the paymentattempt in the database.
	Table = "payment_attempts"
	// PaymentTable is the table that holds the payment relation/edge.
	PaymentTable = "payment_attempts"
	// PaymentInverseTable is the table name for the Payment entity.
	// It exists in this package in order to avoid circular dependency with the "payment" package.
	PaymentInverseTable = "payments"
	// PaymentColumn is the table column denoting the payment relation/edge.
	PaymentColumn = "payment_id"
)

// Columns holds all SQL columns for paymentattempt fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCreatedBy,
	FieldUpdatedBy,
	FieldPaymentID,
	FieldPaymentStatus,
	FieldAttemptNumber,
	FieldGatewayAttemptID,
	FieldErrorMessage,
	FieldMetadata,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus string
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// PaymentIDValidator is a validator for the "payment_id" field. It is called by the builders before save.
	PaymentIDValidator func(string) error
	// PaymentStatusValidator is a validator for the "payment_status" field. It is called by the builders before save.
	PaymentStatusValidator func(string) error
	// DefaultAttemptNumber holds the default value on creation for the "attempt_number" field.
	DefaultAttemptNumber int
	// AttemptNumberValidator is a validator for the "attempt_number" field. It is called by the builders before save.
	AttemptNumberValidator func(int) error
)

// OrderOption defines the ordering options for the PaymentAttempt queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByCreatedBy orders the results by the created_by field.
func ByCreatedBy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedBy, opts...).ToFunc()
}

// ByUpdatedBy orders the results by the updated_by field.
func ByUpdatedBy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedBy, opts...).ToFunc()
}

// ByPaymentID orders the results by the payment_id field.
func ByPaymentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPaymentID, opts...).ToFunc()
}

// ByPaymentStatus orders the results by the payment_status field.
func ByPaymentStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPaymentStatus, opts...).ToFunc()
}

// ByAttemptNumber orders the results by the attempt_number field.
func ByAttemptNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAttemptNumber, opts...).ToFunc()
}

// ByGatewayAttemptID orders the results by the gateway_attempt_id field.
func ByGatewayAttemptID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGatewayAttemptID, opts...).ToFunc()
}

// ByErrorMessage orders the results by the error_message field.
func ByErrorMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldErrorMessage, opts...).ToFunc()
}

// ByPaymentField orders the results by payment field.
func ByPaymentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPaymentStep(), sql.OrderByField(field, opts...))
	}
}
func newPaymentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PaymentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PaymentTable, PaymentColumn),
	)
}
