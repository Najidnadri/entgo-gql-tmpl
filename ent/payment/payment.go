// Code generated by ent, DO NOT EDIT.

package payment

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the payment type in the database.
	Label = "payment"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDateCreated holds the string denoting the date_created field in the database.
	FieldDateCreated = "date_created"
	// FieldDateUpdated holds the string denoting the date_updated field in the database.
	FieldDateUpdated = "date_updated"
	// Table holds the table name of the payment in the database.
	Table = "Payment"
)

// Columns holds all SQL columns for payment fields.
var Columns = []string{
	FieldID,
	FieldDateCreated,
	FieldDateUpdated,
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
	// DefaultDateCreated holds the default value on creation for the "date_created" field.
	DefaultDateCreated func() time.Time
	// DefaultDateUpdated holds the default value on creation for the "date_updated" field.
	DefaultDateUpdated func() time.Time
	// UpdateDefaultDateUpdated holds the default value on update for the "date_updated" field.
	UpdateDefaultDateUpdated func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Payment queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDateCreated orders the results by the date_created field.
func ByDateCreated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateCreated, opts...).ToFunc()
}

// ByDateUpdated orders the results by the date_updated field.
func ByDateUpdated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateUpdated, opts...).ToFunc()
}
