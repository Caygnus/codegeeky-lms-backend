// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/omkar273/codegeeky/ent/cart"
	"github.com/omkar273/codegeeky/ent/cartlineitems"
	"github.com/omkar273/codegeeky/ent/predicate"
	"github.com/omkar273/codegeeky/ent/user"
)

// CartQuery is the builder for querying Cart entities.
type CartQuery struct {
	config
	ctx           *QueryContext
	order         []cart.OrderOption
	inters        []Interceptor
	predicates    []predicate.Cart
	withLineItems *CartLineItemsQuery
	withUser      *UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CartQuery builder.
func (cq *CartQuery) Where(ps ...predicate.Cart) *CartQuery {
	cq.predicates = append(cq.predicates, ps...)
	return cq
}

// Limit the number of records to be returned by this query.
func (cq *CartQuery) Limit(limit int) *CartQuery {
	cq.ctx.Limit = &limit
	return cq
}

// Offset to start from.
func (cq *CartQuery) Offset(offset int) *CartQuery {
	cq.ctx.Offset = &offset
	return cq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (cq *CartQuery) Unique(unique bool) *CartQuery {
	cq.ctx.Unique = &unique
	return cq
}

// Order specifies how the records should be ordered.
func (cq *CartQuery) Order(o ...cart.OrderOption) *CartQuery {
	cq.order = append(cq.order, o...)
	return cq
}

// QueryLineItems chains the current query on the "line_items" edge.
func (cq *CartQuery) QueryLineItems() *CartLineItemsQuery {
	query := (&CartLineItemsClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cart.Table, cart.FieldID, selector),
			sqlgraph.To(cartlineitems.Table, cartlineitems.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, cart.LineItemsTable, cart.LineItemsColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryUser chains the current query on the "user" edge.
func (cq *CartQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: cq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := cq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := cq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(cart.Table, cart.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, cart.UserTable, cart.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(cq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Cart entity from the query.
// Returns a *NotFoundError when no Cart was found.
func (cq *CartQuery) First(ctx context.Context) (*Cart, error) {
	nodes, err := cq.Limit(1).All(setContextOp(ctx, cq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{cart.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (cq *CartQuery) FirstX(ctx context.Context) *Cart {
	node, err := cq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Cart ID from the query.
// Returns a *NotFoundError when no Cart ID was found.
func (cq *CartQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = cq.Limit(1).IDs(setContextOp(ctx, cq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{cart.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (cq *CartQuery) FirstIDX(ctx context.Context) string {
	id, err := cq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Cart entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Cart entity is found.
// Returns a *NotFoundError when no Cart entities are found.
func (cq *CartQuery) Only(ctx context.Context) (*Cart, error) {
	nodes, err := cq.Limit(2).All(setContextOp(ctx, cq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{cart.Label}
	default:
		return nil, &NotSingularError{cart.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (cq *CartQuery) OnlyX(ctx context.Context) *Cart {
	node, err := cq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Cart ID in the query.
// Returns a *NotSingularError when more than one Cart ID is found.
// Returns a *NotFoundError when no entities are found.
func (cq *CartQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = cq.Limit(2).IDs(setContextOp(ctx, cq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{cart.Label}
	default:
		err = &NotSingularError{cart.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (cq *CartQuery) OnlyIDX(ctx context.Context) string {
	id, err := cq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Carts.
func (cq *CartQuery) All(ctx context.Context) ([]*Cart, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryAll)
	if err := cq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Cart, *CartQuery]()
	return withInterceptors[[]*Cart](ctx, cq, qr, cq.inters)
}

// AllX is like All, but panics if an error occurs.
func (cq *CartQuery) AllX(ctx context.Context) []*Cart {
	nodes, err := cq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Cart IDs.
func (cq *CartQuery) IDs(ctx context.Context) (ids []string, err error) {
	if cq.ctx.Unique == nil && cq.path != nil {
		cq.Unique(true)
	}
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryIDs)
	if err = cq.Select(cart.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (cq *CartQuery) IDsX(ctx context.Context) []string {
	ids, err := cq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (cq *CartQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryCount)
	if err := cq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, cq, querierCount[*CartQuery](), cq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (cq *CartQuery) CountX(ctx context.Context) int {
	count, err := cq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (cq *CartQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, cq.ctx, ent.OpQueryExist)
	switch _, err := cq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (cq *CartQuery) ExistX(ctx context.Context) bool {
	exist, err := cq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CartQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (cq *CartQuery) Clone() *CartQuery {
	if cq == nil {
		return nil
	}
	return &CartQuery{
		config:        cq.config,
		ctx:           cq.ctx.Clone(),
		order:         append([]cart.OrderOption{}, cq.order...),
		inters:        append([]Interceptor{}, cq.inters...),
		predicates:    append([]predicate.Cart{}, cq.predicates...),
		withLineItems: cq.withLineItems.Clone(),
		withUser:      cq.withUser.Clone(),
		// clone intermediate query.
		sql:  cq.sql.Clone(),
		path: cq.path,
	}
}

// WithLineItems tells the query-builder to eager-load the nodes that are connected to
// the "line_items" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CartQuery) WithLineItems(opts ...func(*CartLineItemsQuery)) *CartQuery {
	query := (&CartLineItemsClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withLineItems = query
	return cq
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (cq *CartQuery) WithUser(opts ...func(*UserQuery)) *CartQuery {
	query := (&UserClient{config: cq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	cq.withUser = query
	return cq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Cart.Query().
//		GroupBy(cart.FieldStatus).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (cq *CartQuery) GroupBy(field string, fields ...string) *CartGroupBy {
	cq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CartGroupBy{build: cq}
	grbuild.flds = &cq.ctx.Fields
	grbuild.label = cart.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//	}
//
//	client.Cart.Query().
//		Select(cart.FieldStatus).
//		Scan(ctx, &v)
func (cq *CartQuery) Select(fields ...string) *CartSelect {
	cq.ctx.Fields = append(cq.ctx.Fields, fields...)
	sbuild := &CartSelect{CartQuery: cq}
	sbuild.label = cart.Label
	sbuild.flds, sbuild.scan = &cq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CartSelect configured with the given aggregations.
func (cq *CartQuery) Aggregate(fns ...AggregateFunc) *CartSelect {
	return cq.Select().Aggregate(fns...)
}

func (cq *CartQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range cq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, cq); err != nil {
				return err
			}
		}
	}
	for _, f := range cq.ctx.Fields {
		if !cart.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if cq.path != nil {
		prev, err := cq.path(ctx)
		if err != nil {
			return err
		}
		cq.sql = prev
	}
	return nil
}

func (cq *CartQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Cart, error) {
	var (
		nodes       = []*Cart{}
		_spec       = cq.querySpec()
		loadedTypes = [2]bool{
			cq.withLineItems != nil,
			cq.withUser != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Cart).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Cart{config: cq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, cq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := cq.withLineItems; query != nil {
		if err := cq.loadLineItems(ctx, query, nodes,
			func(n *Cart) { n.Edges.LineItems = []*CartLineItems{} },
			func(n *Cart, e *CartLineItems) { n.Edges.LineItems = append(n.Edges.LineItems, e) }); err != nil {
			return nil, err
		}
	}
	if query := cq.withUser; query != nil {
		if err := cq.loadUser(ctx, query, nodes, nil,
			func(n *Cart, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (cq *CartQuery) loadLineItems(ctx context.Context, query *CartLineItemsQuery, nodes []*Cart, init func(*Cart), assign func(*Cart, *CartLineItems)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[string]*Cart)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(cartlineitems.FieldCartID)
	}
	query.Where(predicate.CartLineItems(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(cart.LineItemsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.CartID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "cart_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (cq *CartQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Cart, init func(*Cart), assign func(*Cart, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*Cart)
	for i := range nodes {
		fk := nodes[i].UserID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (cq *CartQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := cq.querySpec()
	_spec.Node.Columns = cq.ctx.Fields
	if len(cq.ctx.Fields) > 0 {
		_spec.Unique = cq.ctx.Unique != nil && *cq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, cq.driver, _spec)
}

func (cq *CartQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(cart.Table, cart.Columns, sqlgraph.NewFieldSpec(cart.FieldID, field.TypeString))
	_spec.From = cq.sql
	if unique := cq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if cq.path != nil {
		_spec.Unique = true
	}
	if fields := cq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cart.FieldID)
		for i := range fields {
			if fields[i] != cart.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if cq.withUser != nil {
			_spec.Node.AddColumnOnce(cart.FieldUserID)
		}
	}
	if ps := cq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := cq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := cq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := cq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (cq *CartQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(cq.driver.Dialect())
	t1 := builder.Table(cart.Table)
	columns := cq.ctx.Fields
	if len(columns) == 0 {
		columns = cart.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if cq.sql != nil {
		selector = cq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if cq.ctx.Unique != nil && *cq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range cq.predicates {
		p(selector)
	}
	for _, p := range cq.order {
		p(selector)
	}
	if offset := cq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := cq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CartGroupBy is the group-by builder for Cart entities.
type CartGroupBy struct {
	selector
	build *CartQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (cgb *CartGroupBy) Aggregate(fns ...AggregateFunc) *CartGroupBy {
	cgb.fns = append(cgb.fns, fns...)
	return cgb
}

// Scan applies the selector query and scans the result into the given value.
func (cgb *CartGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cgb.build.ctx, ent.OpQueryGroupBy)
	if err := cgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CartQuery, *CartGroupBy](ctx, cgb.build, cgb, cgb.build.inters, v)
}

func (cgb *CartGroupBy) sqlScan(ctx context.Context, root *CartQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(cgb.fns))
	for _, fn := range cgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*cgb.flds)+len(cgb.fns))
		for _, f := range *cgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*cgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CartSelect is the builder for selecting fields of Cart entities.
type CartSelect struct {
	*CartQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (cs *CartSelect) Aggregate(fns ...AggregateFunc) *CartSelect {
	cs.fns = append(cs.fns, fns...)
	return cs
}

// Scan applies the selector query and scans the result into the given value.
func (cs *CartSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, cs.ctx, ent.OpQuerySelect)
	if err := cs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CartQuery, *CartSelect](ctx, cs.CartQuery, cs, cs.inters, v)
}

func (cs *CartSelect) sqlScan(ctx context.Context, root *CartQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(cs.fns))
	for _, fn := range cs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*cs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := cs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
