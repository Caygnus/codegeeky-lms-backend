// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/omkar273/codegeeky/ent/internship"
	"github.com/shopspring/decimal"
)

// Internship is the model entity for the Internship schema.
type Internship struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// CreatedBy holds the value of the "created_by" field.
	CreatedBy string `json:"created_by,omitempty"`
	// UpdatedBy holds the value of the "updated_by" field.
	UpdatedBy string `json:"updated_by,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// LookupKey holds the value of the "lookup_key" field.
	LookupKey string `json:"lookup_key,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// List of required skills
	Skills []string `json:"skills,omitempty"`
	// Level of the internship: beginner, intermediate, advanced
	Level string `json:"level,omitempty"`
	// Internship mode: remote, hybrid, onsite
	Mode string `json:"mode,omitempty"`
	// Alternative to months for shorter internships
	DurationInWeeks int `json:"duration_in_weeks,omitempty"`
	// What students will learn in the internship
	LearningOutcomes []string `json:"learning_outcomes,omitempty"`
	// Prerequisites or recommended knowledge
	Prerequisites []string `json:"prerequisites,omitempty"`
	// Benefits of the internship
	Benefits []string `json:"benefits,omitempty"`
	// Currency of the internship
	Currency string `json:"currency,omitempty"`
	// Price of the internship
	Price decimal.Decimal `json:"price,omitempty"`
	// Flat discount on the internship
	FlatDiscount *decimal.Decimal `json:"flat_discount,omitempty"`
	// Percentage discount on the internship
	PercentageDiscount *decimal.Decimal `json:"percentage_discount,omitempty"`
	// Subtotal of the internship
	Subtotal decimal.Decimal `json:"subtotal,omitempty"`
	// Price of the internship
	Total decimal.Decimal `json:"total,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the InternshipQuery when eager-loading is set.
	Edges        InternshipEdges `json:"edges"`
	category_id  *string
	selectValues sql.SelectValues
}

// InternshipEdges holds the relations/edges for other nodes in the graph.
type InternshipEdges struct {
	// Categories holds the value of the categories edge.
	Categories []*Category `json:"categories,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CategoriesOrErr returns the Categories value or an error if the edge
// was not loaded in eager-loading.
func (e InternshipEdges) CategoriesOrErr() ([]*Category, error) {
	if e.loadedTypes[0] {
		return e.Categories, nil
	}
	return nil, &NotLoadedError{edge: "categories"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Internship) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case internship.FieldFlatDiscount, internship.FieldPercentageDiscount:
			values[i] = &sql.NullScanner{S: new(decimal.Decimal)}
		case internship.FieldSkills, internship.FieldLearningOutcomes, internship.FieldPrerequisites, internship.FieldBenefits:
			values[i] = new([]byte)
		case internship.FieldPrice, internship.FieldSubtotal, internship.FieldTotal:
			values[i] = new(decimal.Decimal)
		case internship.FieldDurationInWeeks:
			values[i] = new(sql.NullInt64)
		case internship.FieldID, internship.FieldStatus, internship.FieldCreatedBy, internship.FieldUpdatedBy, internship.FieldTitle, internship.FieldLookupKey, internship.FieldDescription, internship.FieldLevel, internship.FieldMode, internship.FieldCurrency:
			values[i] = new(sql.NullString)
		case internship.FieldCreatedAt, internship.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case internship.ForeignKeys[0]: // category_id
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Internship fields.
func (i *Internship) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case internship.FieldID:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value.Valid {
				i.ID = value.String
			}
		case internship.FieldStatus:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[j])
			} else if value.Valid {
				i.Status = value.String
			}
		case internship.FieldCreatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[j])
			} else if value.Valid {
				i.CreatedAt = value.Time
			}
		case internship.FieldUpdatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[j])
			} else if value.Valid {
				i.UpdatedAt = value.Time
			}
		case internship.FieldCreatedBy:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field created_by", values[j])
			} else if value.Valid {
				i.CreatedBy = value.String
			}
		case internship.FieldUpdatedBy:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field updated_by", values[j])
			} else if value.Valid {
				i.UpdatedBy = value.String
			}
		case internship.FieldTitle:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[j])
			} else if value.Valid {
				i.Title = value.String
			}
		case internship.FieldLookupKey:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field lookup_key", values[j])
			} else if value.Valid {
				i.LookupKey = value.String
			}
		case internship.FieldDescription:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[j])
			} else if value.Valid {
				i.Description = value.String
			}
		case internship.FieldSkills:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field skills", values[j])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &i.Skills); err != nil {
					return fmt.Errorf("unmarshal field skills: %w", err)
				}
			}
		case internship.FieldLevel:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field level", values[j])
			} else if value.Valid {
				i.Level = value.String
			}
		case internship.FieldMode:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mode", values[j])
			} else if value.Valid {
				i.Mode = value.String
			}
		case internship.FieldDurationInWeeks:
			if value, ok := values[j].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field duration_in_weeks", values[j])
			} else if value.Valid {
				i.DurationInWeeks = int(value.Int64)
			}
		case internship.FieldLearningOutcomes:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field learning_outcomes", values[j])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &i.LearningOutcomes); err != nil {
					return fmt.Errorf("unmarshal field learning_outcomes: %w", err)
				}
			}
		case internship.FieldPrerequisites:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field prerequisites", values[j])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &i.Prerequisites); err != nil {
					return fmt.Errorf("unmarshal field prerequisites: %w", err)
				}
			}
		case internship.FieldBenefits:
			if value, ok := values[j].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field benefits", values[j])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &i.Benefits); err != nil {
					return fmt.Errorf("unmarshal field benefits: %w", err)
				}
			}
		case internship.FieldCurrency:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field currency", values[j])
			} else if value.Valid {
				i.Currency = value.String
			}
		case internship.FieldPrice:
			if value, ok := values[j].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field price", values[j])
			} else if value != nil {
				i.Price = *value
			}
		case internship.FieldFlatDiscount:
			if value, ok := values[j].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field flat_discount", values[j])
			} else if value.Valid {
				i.FlatDiscount = new(decimal.Decimal)
				*i.FlatDiscount = *value.S.(*decimal.Decimal)
			}
		case internship.FieldPercentageDiscount:
			if value, ok := values[j].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field percentage_discount", values[j])
			} else if value.Valid {
				i.PercentageDiscount = new(decimal.Decimal)
				*i.PercentageDiscount = *value.S.(*decimal.Decimal)
			}
		case internship.FieldSubtotal:
			if value, ok := values[j].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field subtotal", values[j])
			} else if value != nil {
				i.Subtotal = *value
			}
		case internship.FieldTotal:
			if value, ok := values[j].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field total", values[j])
			} else if value != nil {
				i.Total = *value
			}
		case internship.ForeignKeys[0]:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field category_id", values[j])
			} else if value.Valid {
				i.category_id = new(string)
				*i.category_id = value.String
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Internship.
// This includes values selected through modifiers, order, etc.
func (i *Internship) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryCategories queries the "categories" edge of the Internship entity.
func (i *Internship) QueryCategories() *CategoryQuery {
	return NewInternshipClient(i.config).QueryCategories(i)
}

// Update returns a builder for updating this Internship.
// Note that you need to call Internship.Unwrap() before calling this method if this Internship
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Internship) Update() *InternshipUpdateOne {
	return NewInternshipClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Internship entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Internship) Unwrap() *Internship {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Internship is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Internship) String() string {
	var builder strings.Builder
	builder.WriteString("Internship(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("status=")
	builder.WriteString(i.Status)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(i.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(i.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_by=")
	builder.WriteString(i.CreatedBy)
	builder.WriteString(", ")
	builder.WriteString("updated_by=")
	builder.WriteString(i.UpdatedBy)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(i.Title)
	builder.WriteString(", ")
	builder.WriteString("lookup_key=")
	builder.WriteString(i.LookupKey)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(i.Description)
	builder.WriteString(", ")
	builder.WriteString("skills=")
	builder.WriteString(fmt.Sprintf("%v", i.Skills))
	builder.WriteString(", ")
	builder.WriteString("level=")
	builder.WriteString(i.Level)
	builder.WriteString(", ")
	builder.WriteString("mode=")
	builder.WriteString(i.Mode)
	builder.WriteString(", ")
	builder.WriteString("duration_in_weeks=")
	builder.WriteString(fmt.Sprintf("%v", i.DurationInWeeks))
	builder.WriteString(", ")
	builder.WriteString("learning_outcomes=")
	builder.WriteString(fmt.Sprintf("%v", i.LearningOutcomes))
	builder.WriteString(", ")
	builder.WriteString("prerequisites=")
	builder.WriteString(fmt.Sprintf("%v", i.Prerequisites))
	builder.WriteString(", ")
	builder.WriteString("benefits=")
	builder.WriteString(fmt.Sprintf("%v", i.Benefits))
	builder.WriteString(", ")
	builder.WriteString("currency=")
	builder.WriteString(i.Currency)
	builder.WriteString(", ")
	builder.WriteString("price=")
	builder.WriteString(fmt.Sprintf("%v", i.Price))
	builder.WriteString(", ")
	if v := i.FlatDiscount; v != nil {
		builder.WriteString("flat_discount=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := i.PercentageDiscount; v != nil {
		builder.WriteString("percentage_discount=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("subtotal=")
	builder.WriteString(fmt.Sprintf("%v", i.Subtotal))
	builder.WriteString(", ")
	builder.WriteString("total=")
	builder.WriteString(fmt.Sprintf("%v", i.Total))
	builder.WriteByte(')')
	return builder.String()
}

// Internships is a parsable slice of Internship.
type Internships []*Internship
