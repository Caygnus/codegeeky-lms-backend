package types

import (
	"fmt"
	"strings"

	ierr "github.com/omkar273/codegeeky/internal/errors"
)

// ExpandableField represents a field that can be expanded in API responses
type ExpandableField string

// Common expandable fields
const (
	ExpandCategory      ExpandableField = "category"
	ExpandCartLineItems ExpandableField = "cart_line_items"
)

// ExpandConfig defines which fields can be expanded and their nested expansions
type ExpandConfig struct {
	// AllowedFields are the fields that can be expanded at this level
	AllowedFields []ExpandableField
	// NestedExpands defines which fields can be expanded within an expanded field
	NestedExpands map[ExpandableField][]ExpandableField
}

// Common expand configurations
var (
	// UserExpandConfig defines what can be expanded on a user

	InternshipExpandConfig = ExpandConfig{
		AllowedFields: []ExpandableField{ExpandCategory},
		NestedExpands: map[ExpandableField][]ExpandableField{
			ExpandCategory: {ExpandCategory},
		},
	}

	CartExpandConfig = ExpandConfig{
		AllowedFields: []ExpandableField{ExpandCartLineItems},
		NestedExpands: map[ExpandableField][]ExpandableField{
			ExpandCartLineItems: {ExpandCartLineItems},
		},
	}
)

// Expand represents the expand parameter in API requests
type Expand struct {
	Fields        map[ExpandableField]bool
	NestedExpands map[ExpandableField]Expand
}

// NewExpand creates a new Expand from a comma-separated string of fields
func NewExpand(expand string) Expand {
	if expand == "" {
		return Expand{
			Fields:        make(map[ExpandableField]bool),
			NestedExpands: make(map[ExpandableField]Expand),
		}
	}

	result := Expand{
		Fields:        make(map[ExpandableField]bool),
		NestedExpands: make(map[ExpandableField]Expand),
	}

	for _, field := range strings.Split(expand, ",") {
		field = strings.TrimSpace(field)
		parts := strings.Split(field, ".")

		// Handle root level field
		rootField := ExpandableField(parts[0])
		result.Fields[rootField] = true

		// Handle nested expands
		if len(parts) > 1 {
			nested := NewExpand(strings.Join(parts[1:], ","))
			result.NestedExpands[rootField] = nested
		}
	}

	return result
}

// Has checks if a field should be expanded
func (e Expand) Has(field ExpandableField) bool {
	return e.Fields[field]
}

// GetNested returns the nested expands for a field
func (e Expand) GetNested(field ExpandableField) Expand {
	if nested, ok := e.NestedExpands[field]; ok {
		return nested
	}
	return NewExpand("")
}

// IsEmpty checks if no fields are to be expanded
func (e Expand) IsEmpty() bool {
	return len(e.Fields) == 0
}

// String returns a string representation of the expand
func (e Expand) String() string {
	var fields []string
	for field := range e.Fields {
		if nested, ok := e.NestedExpands[field]; ok && !nested.IsEmpty() {
			fields = append(fields, fmt.Sprintf("%s.%s", field, nested.String()))
		} else {
			fields = append(fields, string(field))
		}
	}
	return strings.Join(fields, ",")
}

// Validate checks if the expand request is valid according to the config
func (e Expand) Validate(config ExpandConfig) error {
	for field := range e.Fields {
		// Check if field is allowed
		allowed := false
		for _, allowedField := range config.AllowedFields {
			if field == allowedField {
				allowed = true
				break
			}
		}
		if !allowed {
			return ierr.NewError("field not allowed to be expanded").
				WithHint("Field is not allowed to be expanded").
				WithReportableDetails(
					map[string]any{
						"field": field,
					},
				).
				Mark(ierr.ErrValidation)
		}

		// Check nested expands
		if nested, ok := e.NestedExpands[field]; ok {
			allowedNested, ok := config.NestedExpands[field]
			if !ok {
				return ierr.NewError("field does not support nested expands").
					WithHint("Field does not support nested expands").
					WithReportableDetails(
						map[string]any{
							"field": field,
						},
					).
					Mark(ierr.ErrValidation)
			}

			// Create a config for nested validation
			nestedConfig := ExpandConfig{
				AllowedFields: allowedNested,
			}
			if err := nested.Validate(nestedConfig); err != nil {
				return err
			}
		}
	}
	return nil
}
