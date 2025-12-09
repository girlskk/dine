package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// TableArea holds the schema definition for the TableArea entity.
type TableArea struct {
	ent.Schema
}

func (TableArea) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the TableArea.
func (TableArea) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
		field.Int("table_count").Optional().NonNegative(),
	}
}

// Edges of the TableArea.
func (TableArea) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dinetables", DineTable.Type),
	}
}
