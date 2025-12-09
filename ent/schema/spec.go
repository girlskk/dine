package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Spec 商品规格
type Spec struct {
	ent.Schema
}

func (Spec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Spec.
func (Spec) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
	}
}

// Indexes 定义索引
func (Spec) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "name"),
	}
}

// Edges of the Spec.
func (Spec) Edges() []ent.Edge {
	return []ent.Edge{
		// 商品-商品规格关联
		edge.To("product_specs", ProductSpec.Type),
	}
}
