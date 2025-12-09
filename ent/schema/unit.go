package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Unit 商品单位
type Unit struct {
	ent.Schema
}

func (Unit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Unit.
func (Unit) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
	}
}

// Indexes of the Unit.
func (Unit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "name"),
	}
}

// Edges of the Unit.
func (Unit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("products", Product.Type),
	}
}
