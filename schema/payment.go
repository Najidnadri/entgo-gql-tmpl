package schema

import (
	"chipin/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	_ "entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
)

type Payment struct {
	ent.Schema
}

func (Payment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.IdMixin{},
		mixins.BaseAnnotationsMixin{},
	}
}

func (Payment) Fields() []ent.Field {
	return []ent.Field{}
}

func (Payment) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Payment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "Payment",
		},
	}
}
