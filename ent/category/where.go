// Code generated by ent, DO NOT EDIT.

package category

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/omkar273/codegeeky/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldStatus, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldUpdatedAt, v))
}

// CreatedBy applies equality check predicate on the "created_by" field. It's identical to CreatedByEQ.
func CreatedBy(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldCreatedBy, v))
}

// UpdatedBy applies equality check predicate on the "updated_by" field. It's identical to UpdatedByEQ.
func UpdatedBy(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldUpdatedBy, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldName, v))
}

// LookupKey applies equality check predicate on the "lookup_key" field. It's identical to LookupKeyEQ.
func LookupKey(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldLookupKey, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldDescription, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldStatus, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldUpdatedAt, v))
}

// CreatedByEQ applies the EQ predicate on the "created_by" field.
func CreatedByEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldCreatedBy, v))
}

// CreatedByNEQ applies the NEQ predicate on the "created_by" field.
func CreatedByNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldCreatedBy, v))
}

// CreatedByIn applies the In predicate on the "created_by" field.
func CreatedByIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldCreatedBy, vs...))
}

// CreatedByNotIn applies the NotIn predicate on the "created_by" field.
func CreatedByNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldCreatedBy, vs...))
}

// CreatedByGT applies the GT predicate on the "created_by" field.
func CreatedByGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldCreatedBy, v))
}

// CreatedByGTE applies the GTE predicate on the "created_by" field.
func CreatedByGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldCreatedBy, v))
}

// CreatedByLT applies the LT predicate on the "created_by" field.
func CreatedByLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldCreatedBy, v))
}

// CreatedByLTE applies the LTE predicate on the "created_by" field.
func CreatedByLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldCreatedBy, v))
}

// CreatedByContains applies the Contains predicate on the "created_by" field.
func CreatedByContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldCreatedBy, v))
}

// CreatedByHasPrefix applies the HasPrefix predicate on the "created_by" field.
func CreatedByHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldCreatedBy, v))
}

// CreatedByHasSuffix applies the HasSuffix predicate on the "created_by" field.
func CreatedByHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldCreatedBy, v))
}

// CreatedByIsNil applies the IsNil predicate on the "created_by" field.
func CreatedByIsNil() predicate.Category {
	return predicate.Category(sql.FieldIsNull(FieldCreatedBy))
}

// CreatedByNotNil applies the NotNil predicate on the "created_by" field.
func CreatedByNotNil() predicate.Category {
	return predicate.Category(sql.FieldNotNull(FieldCreatedBy))
}

// CreatedByEqualFold applies the EqualFold predicate on the "created_by" field.
func CreatedByEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldCreatedBy, v))
}

// CreatedByContainsFold applies the ContainsFold predicate on the "created_by" field.
func CreatedByContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldCreatedBy, v))
}

// UpdatedByEQ applies the EQ predicate on the "updated_by" field.
func UpdatedByEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldUpdatedBy, v))
}

// UpdatedByNEQ applies the NEQ predicate on the "updated_by" field.
func UpdatedByNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldUpdatedBy, v))
}

// UpdatedByIn applies the In predicate on the "updated_by" field.
func UpdatedByIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldUpdatedBy, vs...))
}

// UpdatedByNotIn applies the NotIn predicate on the "updated_by" field.
func UpdatedByNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldUpdatedBy, vs...))
}

// UpdatedByGT applies the GT predicate on the "updated_by" field.
func UpdatedByGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldUpdatedBy, v))
}

// UpdatedByGTE applies the GTE predicate on the "updated_by" field.
func UpdatedByGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldUpdatedBy, v))
}

// UpdatedByLT applies the LT predicate on the "updated_by" field.
func UpdatedByLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldUpdatedBy, v))
}

// UpdatedByLTE applies the LTE predicate on the "updated_by" field.
func UpdatedByLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldUpdatedBy, v))
}

// UpdatedByContains applies the Contains predicate on the "updated_by" field.
func UpdatedByContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldUpdatedBy, v))
}

// UpdatedByHasPrefix applies the HasPrefix predicate on the "updated_by" field.
func UpdatedByHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldUpdatedBy, v))
}

// UpdatedByHasSuffix applies the HasSuffix predicate on the "updated_by" field.
func UpdatedByHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldUpdatedBy, v))
}

// UpdatedByIsNil applies the IsNil predicate on the "updated_by" field.
func UpdatedByIsNil() predicate.Category {
	return predicate.Category(sql.FieldIsNull(FieldUpdatedBy))
}

// UpdatedByNotNil applies the NotNil predicate on the "updated_by" field.
func UpdatedByNotNil() predicate.Category {
	return predicate.Category(sql.FieldNotNull(FieldUpdatedBy))
}

// UpdatedByEqualFold applies the EqualFold predicate on the "updated_by" field.
func UpdatedByEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldUpdatedBy, v))
}

// UpdatedByContainsFold applies the ContainsFold predicate on the "updated_by" field.
func UpdatedByContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldUpdatedBy, v))
}

// MetadataIsNil applies the IsNil predicate on the "metadata" field.
func MetadataIsNil() predicate.Category {
	return predicate.Category(sql.FieldIsNull(FieldMetadata))
}

// MetadataNotNil applies the NotNil predicate on the "metadata" field.
func MetadataNotNil() predicate.Category {
	return predicate.Category(sql.FieldNotNull(FieldMetadata))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldName, v))
}

// LookupKeyEQ applies the EQ predicate on the "lookup_key" field.
func LookupKeyEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldLookupKey, v))
}

// LookupKeyNEQ applies the NEQ predicate on the "lookup_key" field.
func LookupKeyNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldLookupKey, v))
}

// LookupKeyIn applies the In predicate on the "lookup_key" field.
func LookupKeyIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldLookupKey, vs...))
}

// LookupKeyNotIn applies the NotIn predicate on the "lookup_key" field.
func LookupKeyNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldLookupKey, vs...))
}

// LookupKeyGT applies the GT predicate on the "lookup_key" field.
func LookupKeyGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldLookupKey, v))
}

// LookupKeyGTE applies the GTE predicate on the "lookup_key" field.
func LookupKeyGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldLookupKey, v))
}

// LookupKeyLT applies the LT predicate on the "lookup_key" field.
func LookupKeyLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldLookupKey, v))
}

// LookupKeyLTE applies the LTE predicate on the "lookup_key" field.
func LookupKeyLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldLookupKey, v))
}

// LookupKeyContains applies the Contains predicate on the "lookup_key" field.
func LookupKeyContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldLookupKey, v))
}

// LookupKeyHasPrefix applies the HasPrefix predicate on the "lookup_key" field.
func LookupKeyHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldLookupKey, v))
}

// LookupKeyHasSuffix applies the HasSuffix predicate on the "lookup_key" field.
func LookupKeyHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldLookupKey, v))
}

// LookupKeyEqualFold applies the EqualFold predicate on the "lookup_key" field.
func LookupKeyEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldLookupKey, v))
}

// LookupKeyContainsFold applies the ContainsFold predicate on the "lookup_key" field.
func LookupKeyContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldLookupKey, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Category {
	return predicate.Category(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Category {
	return predicate.Category(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Category {
	return predicate.Category(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Category {
	return predicate.Category(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Category {
	return predicate.Category(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Category {
	return predicate.Category(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Category {
	return predicate.Category(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Category {
	return predicate.Category(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Category {
	return predicate.Category(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Category {
	return predicate.Category(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Category {
	return predicate.Category(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Category {
	return predicate.Category(sql.FieldContainsFold(FieldDescription, v))
}

// HasInternships applies the HasEdge predicate on the "internships" edge.
func HasInternships() predicate.Category {
	return predicate.Category(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, InternshipsTable, InternshipsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasInternshipsWith applies the HasEdge predicate on the "internships" edge with a given conditions (other predicates).
func HasInternshipsWith(preds ...predicate.Internship) predicate.Category {
	return predicate.Category(func(s *sql.Selector) {
		step := newInternshipsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Category) predicate.Category {
	return predicate.Category(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Category) predicate.Category {
	return predicate.Category(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Category) predicate.Category {
	return predicate.Category(sql.NotPredicates(p))
}
