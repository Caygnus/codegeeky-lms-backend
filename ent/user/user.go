// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
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
	// FieldFullName holds the string denoting the full_name field in the database.
	FieldFullName = "full_name"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldPhoneNumber holds the string denoting the phone_number field in the database.
	FieldPhoneNumber = "phone_number"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// EdgeCarts holds the string denoting the carts edge name in mutations.
	EdgeCarts = "carts"
	// Table holds the table name of the user in the database.
	Table = "users"
	// CartsTable is the table that holds the carts relation/edge.
	CartsTable = "carts"
	// CartsInverseTable is the table name for the Cart entity.
	// It exists in this package in order to avoid circular dependency with the "cart" package.
	CartsInverseTable = "carts"
	// CartsColumn is the table column denoting the carts relation/edge.
	CartsColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCreatedBy,
	FieldUpdatedBy,
	FieldFullName,
	FieldEmail,
	FieldPhoneNumber,
	FieldRole,
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
	// FullNameValidator is a validator for the "full_name" field. It is called by the builders before save.
	FullNameValidator func(string) error
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// RoleValidator is a validator for the "role" field. It is called by the builders before save.
	RoleValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
)

// OrderOption defines the ordering options for the User queries.
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

// ByFullName orders the results by the full_name field.
func ByFullName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFullName, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByPhoneNumber orders the results by the phone_number field.
func ByPhoneNumber(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPhoneNumber, opts...).ToFunc()
}

// ByRole orders the results by the role field.
func ByRole(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRole, opts...).ToFunc()
}

// ByCartsCount orders the results by carts count.
func ByCartsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCartsStep(), opts...)
	}
}

// ByCarts orders the results by carts terms.
func ByCarts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCartsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newCartsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CartsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CartsTable, CartsColumn),
	)
}
