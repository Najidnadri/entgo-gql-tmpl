package mixins

import (
	"chipin/utils"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("date_created").
			Immutable().
			Default(func() time.Time {
				now, err := utils.GetCurrentTimeInMalaysia()
				if err != nil {
					return time.Now()
				}
				return now
			}).Annotations(
			entgql.OrderField("CREATED_AT"),
		),
		field.Time("date_updated").
			Default(func() time.Time {
				now, err := utils.GetCurrentTimeInMalaysia()
				if err != nil {
					return time.Now()
				}
				return now
			}).
			UpdateDefault(func() time.Time {
				now, err := utils.GetCurrentTimeInMalaysia()
				if err != nil {
					return time.Now()
				}
				return now
			}).Annotations(
			entgql.OrderField("UPDATED_AT"),
		),
	}
}
