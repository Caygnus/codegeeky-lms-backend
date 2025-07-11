// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/omkar273/codegeeky/ent/cart"
	"github.com/omkar273/codegeeky/ent/predicate"
	"github.com/omkar273/codegeeky/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetStatus sets the "status" field.
func (uu *UserUpdate) SetStatus(s string) *UserUpdate {
	uu.mutation.SetStatus(s)
	return uu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (uu *UserUpdate) SetNillableStatus(s *string) *UserUpdate {
	if s != nil {
		uu.SetStatus(*s)
	}
	return uu
}

// SetUpdatedAt sets the "updated_at" field.
func (uu *UserUpdate) SetUpdatedAt(t time.Time) *UserUpdate {
	uu.mutation.SetUpdatedAt(t)
	return uu
}

// SetUpdatedBy sets the "updated_by" field.
func (uu *UserUpdate) SetUpdatedBy(s string) *UserUpdate {
	uu.mutation.SetUpdatedBy(s)
	return uu
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (uu *UserUpdate) SetNillableUpdatedBy(s *string) *UserUpdate {
	if s != nil {
		uu.SetUpdatedBy(*s)
	}
	return uu
}

// ClearUpdatedBy clears the value of the "updated_by" field.
func (uu *UserUpdate) ClearUpdatedBy() *UserUpdate {
	uu.mutation.ClearUpdatedBy()
	return uu
}

// SetFullName sets the "full_name" field.
func (uu *UserUpdate) SetFullName(s string) *UserUpdate {
	uu.mutation.SetFullName(s)
	return uu
}

// SetNillableFullName sets the "full_name" field if the given value is not nil.
func (uu *UserUpdate) SetNillableFullName(s *string) *UserUpdate {
	if s != nil {
		uu.SetFullName(*s)
	}
	return uu
}

// SetEmail sets the "email" field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetNillableEmail sets the "email" field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmail(s *string) *UserUpdate {
	if s != nil {
		uu.SetEmail(*s)
	}
	return uu
}

// SetPhoneNumber sets the "phone_number" field.
func (uu *UserUpdate) SetPhoneNumber(s string) *UserUpdate {
	uu.mutation.SetPhoneNumber(s)
	return uu
}

// SetNillablePhoneNumber sets the "phone_number" field if the given value is not nil.
func (uu *UserUpdate) SetNillablePhoneNumber(s *string) *UserUpdate {
	if s != nil {
		uu.SetPhoneNumber(*s)
	}
	return uu
}

// SetRole sets the "role" field.
func (uu *UserUpdate) SetRole(s string) *UserUpdate {
	uu.mutation.SetRole(s)
	return uu
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (uu *UserUpdate) SetNillableRole(s *string) *UserUpdate {
	if s != nil {
		uu.SetRole(*s)
	}
	return uu
}

// AddCartIDs adds the "carts" edge to the Cart entity by IDs.
func (uu *UserUpdate) AddCartIDs(ids ...string) *UserUpdate {
	uu.mutation.AddCartIDs(ids...)
	return uu
}

// AddCarts adds the "carts" edges to the Cart entity.
func (uu *UserUpdate) AddCarts(c ...*Cart) *UserUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return uu.AddCartIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearCarts clears all "carts" edges to the Cart entity.
func (uu *UserUpdate) ClearCarts() *UserUpdate {
	uu.mutation.ClearCarts()
	return uu
}

// RemoveCartIDs removes the "carts" edge to Cart entities by IDs.
func (uu *UserUpdate) RemoveCartIDs(ids ...string) *UserUpdate {
	uu.mutation.RemoveCartIDs(ids...)
	return uu
}

// RemoveCarts removes "carts" edges to Cart entities.
func (uu *UserUpdate) RemoveCarts(c ...*Cart) *UserUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return uu.RemoveCartIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	uu.defaults()
	return withHooks(ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UserUpdate) defaults() {
	if _, ok := uu.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.FullName(); ok {
		if err := user.FullNameValidator(v); err != nil {
			return &ValidationError{Name: "full_name", err: fmt.Errorf(`ent: validator failed for field "User.full_name": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "User.email": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf(`ent: validator failed for field "User.role": %w`, err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := uu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeString))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Status(); ok {
		_spec.SetField(user.FieldStatus, field.TypeString, value)
	}
	if value, ok := uu.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if uu.mutation.CreatedByCleared() {
		_spec.ClearField(user.FieldCreatedBy, field.TypeString)
	}
	if value, ok := uu.mutation.UpdatedBy(); ok {
		_spec.SetField(user.FieldUpdatedBy, field.TypeString, value)
	}
	if uu.mutation.UpdatedByCleared() {
		_spec.ClearField(user.FieldUpdatedBy, field.TypeString)
	}
	if value, ok := uu.mutation.FullName(); ok {
		_spec.SetField(user.FieldFullName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uu.mutation.PhoneNumber(); ok {
		_spec.SetField(user.FieldPhoneNumber, field.TypeString, value)
	}
	if value, ok := uu.mutation.Role(); ok {
		_spec.SetField(user.FieldRole, field.TypeString, value)
	}
	if uu.mutation.CartsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedCartsIDs(); len(nodes) > 0 && !uu.mutation.CartsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.CartsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetStatus sets the "status" field.
func (uuo *UserUpdateOne) SetStatus(s string) *UserUpdateOne {
	uuo.mutation.SetStatus(s)
	return uuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableStatus(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetStatus(*s)
	}
	return uuo
}

// SetUpdatedAt sets the "updated_at" field.
func (uuo *UserUpdateOne) SetUpdatedAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetUpdatedAt(t)
	return uuo
}

// SetUpdatedBy sets the "updated_by" field.
func (uuo *UserUpdateOne) SetUpdatedBy(s string) *UserUpdateOne {
	uuo.mutation.SetUpdatedBy(s)
	return uuo
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableUpdatedBy(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetUpdatedBy(*s)
	}
	return uuo
}

// ClearUpdatedBy clears the value of the "updated_by" field.
func (uuo *UserUpdateOne) ClearUpdatedBy() *UserUpdateOne {
	uuo.mutation.ClearUpdatedBy()
	return uuo
}

// SetFullName sets the "full_name" field.
func (uuo *UserUpdateOne) SetFullName(s string) *UserUpdateOne {
	uuo.mutation.SetFullName(s)
	return uuo
}

// SetNillableFullName sets the "full_name" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableFullName(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetFullName(*s)
	}
	return uuo
}

// SetEmail sets the "email" field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetNillableEmail sets the "email" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmail(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetEmail(*s)
	}
	return uuo
}

// SetPhoneNumber sets the "phone_number" field.
func (uuo *UserUpdateOne) SetPhoneNumber(s string) *UserUpdateOne {
	uuo.mutation.SetPhoneNumber(s)
	return uuo
}

// SetNillablePhoneNumber sets the "phone_number" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillablePhoneNumber(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetPhoneNumber(*s)
	}
	return uuo
}

// SetRole sets the "role" field.
func (uuo *UserUpdateOne) SetRole(s string) *UserUpdateOne {
	uuo.mutation.SetRole(s)
	return uuo
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableRole(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetRole(*s)
	}
	return uuo
}

// AddCartIDs adds the "carts" edge to the Cart entity by IDs.
func (uuo *UserUpdateOne) AddCartIDs(ids ...string) *UserUpdateOne {
	uuo.mutation.AddCartIDs(ids...)
	return uuo
}

// AddCarts adds the "carts" edges to the Cart entity.
func (uuo *UserUpdateOne) AddCarts(c ...*Cart) *UserUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return uuo.AddCartIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearCarts clears all "carts" edges to the Cart entity.
func (uuo *UserUpdateOne) ClearCarts() *UserUpdateOne {
	uuo.mutation.ClearCarts()
	return uuo
}

// RemoveCartIDs removes the "carts" edge to Cart entities by IDs.
func (uuo *UserUpdateOne) RemoveCartIDs(ids ...string) *UserUpdateOne {
	uuo.mutation.RemoveCartIDs(ids...)
	return uuo
}

// RemoveCarts removes "carts" edges to Cart entities.
func (uuo *UserUpdateOne) RemoveCarts(c ...*Cart) *UserUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return uuo.RemoveCartIDs(ids...)
}

// Where appends a list predicates to the UserUpdate builder.
func (uuo *UserUpdateOne) Where(ps ...predicate.User) *UserUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	uuo.defaults()
	return withHooks(ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UserUpdateOne) defaults() {
	if _, ok := uuo.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.FullName(); ok {
		if err := user.FullNameValidator(v); err != nil {
			return &ValidationError{Name: "full_name", err: fmt.Errorf(`ent: validator failed for field "User.full_name": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf(`ent: validator failed for field "User.email": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf(`ent: validator failed for field "User.role": %w`, err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	if err := uuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeString))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.Status(); ok {
		_spec.SetField(user.FieldStatus, field.TypeString, value)
	}
	if value, ok := uuo.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if uuo.mutation.CreatedByCleared() {
		_spec.ClearField(user.FieldCreatedBy, field.TypeString)
	}
	if value, ok := uuo.mutation.UpdatedBy(); ok {
		_spec.SetField(user.FieldUpdatedBy, field.TypeString, value)
	}
	if uuo.mutation.UpdatedByCleared() {
		_spec.ClearField(user.FieldUpdatedBy, field.TypeString)
	}
	if value, ok := uuo.mutation.FullName(); ok {
		_spec.SetField(user.FieldFullName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uuo.mutation.PhoneNumber(); ok {
		_spec.SetField(user.FieldPhoneNumber, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Role(); ok {
		_spec.SetField(user.FieldRole, field.TypeString, value)
	}
	if uuo.mutation.CartsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedCartsIDs(); len(nodes) > 0 && !uuo.mutation.CartsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.CartsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.CartsTable,
			Columns: []string{user.CartsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}
