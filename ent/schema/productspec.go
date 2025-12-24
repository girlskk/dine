package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProductSpec 商品规格
type ProductSpec struct {
	ent.Schema
}

func (ProductSpec) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductSpec.
func (ProductSpec) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("规格名称"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Optional().Immutable().Comment("门店ID"),
		field.Int("product_count").Default(0).Comment("关联的商品数量"),
	}
}

func (ProductSpec) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
	}
}

// Edges of the ProductSpec.
func (ProductSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("product_specs", ProductSpecRelation.Type).Comment("规格项列表"),
	}
}
