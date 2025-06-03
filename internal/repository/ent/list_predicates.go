package ent

import (
	"entgo.io/ent/dialect/sql"
)

// ListPredicateBuilder provides a fluent interface for building list predicates
type ListPredicateBuilder struct {
	field string
}

// NewListPredicateBuilder creates a new ListPredicateBuilder for the given field
func NewListPredicateBuilder(field string) *ListPredicateBuilder {
	return &ListPredicateBuilder{
		field: field,
	}
}

// Contains creates a predicate that checks if the list contains the given value
func (b *ListPredicateBuilder) Contains(value string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.Contains(b.field, value))
	}
}

// ContainsFold creates a case-insensitive predicate that checks if the list contains the given value
func (b *ListPredicateBuilder) ContainsFold(value string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.ContainsFold(b.field, value))
	}
}

// NotContains creates a predicate that checks if the list does not contain the given value
func (b *ListPredicateBuilder) NotContains(value string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.Not(sql.Contains(b.field, value)))
	}
}

// NotContainsFold creates a case-insensitive predicate that checks if the list does not contain the given value
func (b *ListPredicateBuilder) NotContainsFold(value string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.Not(sql.ContainsFold(b.field, value)))
	}
}

// In creates a predicate that checks if the list contains any of the given values
func (b *ListPredicateBuilder) In(values ...string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		args := make([]interface{}, len(values))
		for i, v := range values {
			args[i] = v
		}
		s.Where(sql.In(b.field, args...))
	}
}

// NotIn creates a predicate that checks if the list does not contain any of the given values
func (b *ListPredicateBuilder) NotIn(values ...string) func(*sql.Selector) {
	return func(s *sql.Selector) {
		args := make([]interface{}, len(values))
		for i, v := range values {
			args[i] = v
		}
		s.Where(sql.NotIn(b.field, args...))
	}
}

// IsNil creates a predicate that checks if the list is nil
func (b *ListPredicateBuilder) IsNil() func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.IsNull(b.field))
	}
}

// NotNil creates a predicate that checks if the list is not nil
func (b *ListPredicateBuilder) NotNil() func(*sql.Selector) {
	return func(s *sql.Selector) {
		s.Where(sql.NotNull(b.field))
	}
}
