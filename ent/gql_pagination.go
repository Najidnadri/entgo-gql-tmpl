// Code generated by ent, DO NOT EDIT.

package ent

import (
	"chipin/ent/payment"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Common entgql types.
type (
	Cursor         = entgql.Cursor[uuid.UUID]
	PageInfo       = entgql.PageInfo[uuid.UUID]
	OrderDirection = entgql.OrderDirection
)

func orderFunc(o OrderDirection, field string) func(*sql.Selector) {
	if o == entgql.OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// PaymentEdge is the edge representation of Payment.
type PaymentEdge struct {
	Node   *Payment `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// PaymentConnection is the connection containing edges to Payment.
type PaymentConnection struct {
	Edges      []*PaymentEdge `json:"edges"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (c *PaymentConnection) build(nodes []*Payment, pager *paymentPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Payment
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Payment {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Payment {
			return nodes[i]
		}
	}
	c.Edges = make([]*PaymentEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &PaymentEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// PaymentPaginateOption enables pagination customization.
type PaymentPaginateOption func(*paymentPager) error

// WithPaymentOrder configures pagination ordering.
func WithPaymentOrder(order []*PaymentOrder) PaymentPaginateOption {
	return func(pager *paymentPager) error {
		for _, o := range order {
			if err := o.Direction.Validate(); err != nil {
				return err
			}
		}
		pager.order = append(pager.order, order...)
		return nil
	}
}

// WithPaymentFilter configures pagination filter.
func WithPaymentFilter(filter func(*PaymentQuery) (*PaymentQuery, error)) PaymentPaginateOption {
	return func(pager *paymentPager) error {
		if filter == nil {
			return errors.New("PaymentQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type paymentPager struct {
	reverse bool
	order   []*PaymentOrder
	filter  func(*PaymentQuery) (*PaymentQuery, error)
}

func newPaymentPager(opts []PaymentPaginateOption, reverse bool) (*paymentPager, error) {
	pager := &paymentPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	for i, o := range pager.order {
		if i > 0 && o.Field == pager.order[i-1].Field {
			return nil, fmt.Errorf("duplicate order direction %q", o.Direction)
		}
	}
	return pager, nil
}

func (p *paymentPager) applyFilter(query *PaymentQuery) (*PaymentQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *paymentPager) toCursor(pa *Payment) Cursor {
	cs := make([]any, 0, len(p.order))
	for _, o := range p.order {
		cs = append(cs, o.Field.toCursor(pa).Value)
	}
	return Cursor{ID: pa.ID, Value: cs}
}

func (p *paymentPager) applyCursors(query *PaymentQuery, after, before *Cursor) (*PaymentQuery, error) {
	idDirection := entgql.OrderDirectionAsc
	if p.reverse {
		idDirection = entgql.OrderDirectionDesc
	}
	fields, directions := make([]string, 0, len(p.order)), make([]OrderDirection, 0, len(p.order))
	for _, o := range p.order {
		fields = append(fields, o.Field.column)
		direction := o.Direction
		if p.reverse {
			direction = direction.Reverse()
		}
		directions = append(directions, direction)
	}
	predicates, err := entgql.MultiCursorsPredicate(after, before, &entgql.MultiCursorsOptions{
		FieldID:     DefaultPaymentOrder.Field.column,
		DirectionID: idDirection,
		Fields:      fields,
		Directions:  directions,
	})
	if err != nil {
		return nil, err
	}
	for _, predicate := range predicates {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *paymentPager) applyOrder(query *PaymentQuery) *PaymentQuery {
	var defaultOrdered bool
	for _, o := range p.order {
		direction := o.Direction
		if p.reverse {
			direction = direction.Reverse()
		}
		query = query.Order(o.Field.toTerm(direction.OrderTermOption()))
		if o.Field.column == DefaultPaymentOrder.Field.column {
			defaultOrdered = true
		}
		if len(query.ctx.Fields) > 0 {
			query.ctx.AppendFieldOnce(o.Field.column)
		}
	}
	if !defaultOrdered {
		direction := entgql.OrderDirectionAsc
		if p.reverse {
			direction = direction.Reverse()
		}
		query = query.Order(DefaultPaymentOrder.Field.toTerm(direction.OrderTermOption()))
	}
	return query
}

func (p *paymentPager) orderExpr(query *PaymentQuery) sql.Querier {
	if len(query.ctx.Fields) > 0 {
		for _, o := range p.order {
			query.ctx.AppendFieldOnce(o.Field.column)
		}
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		for _, o := range p.order {
			direction := o.Direction
			if p.reverse {
				direction = direction.Reverse()
			}
			b.Ident(o.Field.column).Pad().WriteString(string(direction))
			b.Comma()
		}
		direction := entgql.OrderDirectionAsc
		if p.reverse {
			direction = direction.Reverse()
		}
		b.Ident(DefaultPaymentOrder.Field.column).Pad().WriteString(string(direction))
	})
}

// Paginate executes the query and returns a relay based cursor connection to Payment.
func (pa *PaymentQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...PaymentPaginateOption,
) (*PaymentConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newPaymentPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if pa, err = pager.applyFilter(pa); err != nil {
		return nil, err
	}
	conn := &PaymentConnection{Edges: []*PaymentEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = pa.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if pa, err = pager.applyCursors(pa, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		pa.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := pa.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	pa = pager.applyOrder(pa)
	nodes, err := pa.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

var (
	// PaymentOrderFieldID orders Payment by id.
	PaymentOrderFieldID = &PaymentOrderField{
		Value: func(pa *Payment) (ent.Value, error) {
			return pa.ID, nil
		},
		column: payment.FieldID,
		toTerm: payment.ByID,
		toCursor: func(pa *Payment) Cursor {
			return Cursor{
				ID:    pa.ID,
				Value: pa.ID,
			}
		},
	}
	// PaymentOrderFieldDateCreated orders Payment by date_created.
	PaymentOrderFieldDateCreated = &PaymentOrderField{
		Value: func(pa *Payment) (ent.Value, error) {
			return pa.DateCreated, nil
		},
		column: payment.FieldDateCreated,
		toTerm: payment.ByDateCreated,
		toCursor: func(pa *Payment) Cursor {
			return Cursor{
				ID:    pa.ID,
				Value: pa.DateCreated,
			}
		},
	}
	// PaymentOrderFieldDateUpdated orders Payment by date_updated.
	PaymentOrderFieldDateUpdated = &PaymentOrderField{
		Value: func(pa *Payment) (ent.Value, error) {
			return pa.DateUpdated, nil
		},
		column: payment.FieldDateUpdated,
		toTerm: payment.ByDateUpdated,
		toCursor: func(pa *Payment) Cursor {
			return Cursor{
				ID:    pa.ID,
				Value: pa.DateUpdated,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f PaymentOrderField) String() string {
	var str string
	switch f.column {
	case PaymentOrderFieldID.column:
		str = "ID"
	case PaymentOrderFieldDateCreated.column:
		str = "CREATED_AT"
	case PaymentOrderFieldDateUpdated.column:
		str = "UPDATED_AT"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f PaymentOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *PaymentOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("PaymentOrderField %T must be a string", v)
	}
	switch str {
	case "ID":
		*f = *PaymentOrderFieldID
	case "CREATED_AT":
		*f = *PaymentOrderFieldDateCreated
	case "UPDATED_AT":
		*f = *PaymentOrderFieldDateUpdated
	default:
		return fmt.Errorf("%s is not a valid PaymentOrderField", str)
	}
	return nil
}

// PaymentOrderField defines the ordering field of Payment.
type PaymentOrderField struct {
	// Value extracts the ordering value from the given Payment.
	Value    func(*Payment) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) payment.OrderOption
	toCursor func(*Payment) Cursor
}

// PaymentOrder defines the ordering of Payment.
type PaymentOrder struct {
	Direction OrderDirection     `json:"direction"`
	Field     *PaymentOrderField `json:"field"`
}

// DefaultPaymentOrder is the default ordering of Payment.
var DefaultPaymentOrder = &PaymentOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &PaymentOrderField{
		Value: func(pa *Payment) (ent.Value, error) {
			return pa.ID, nil
		},
		column: payment.FieldID,
		toTerm: payment.ByID,
		toCursor: func(pa *Payment) Cursor {
			return Cursor{ID: pa.ID}
		},
	},
}

// ToEdge converts Payment into PaymentEdge.
func (pa *Payment) ToEdge(order *PaymentOrder) *PaymentEdge {
	if order == nil {
		order = DefaultPaymentOrder
	}
	return &PaymentEdge{
		Node:   pa,
		Cursor: order.Field.toCursor(pa),
	}
}