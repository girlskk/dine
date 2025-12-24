package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProductTag 商品标签
type ProductTag struct {
	ent.Schema
}

func (ProductTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductTag.
func (ProductTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("标签名称"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Optional().Immutable().Comment("门店ID"),
		field.Int("product_count").Default(0).Comment("关联的商品数量"),
	}
}

func (ProductTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
	}
}

// Edges of the ProductTag.
func (ProductTag) Edges() []ent.Edge {
	return []ent.Edge{
		// M -> M：商品与标签的关联
		edge.From("products", Product.Type).Ref("tags").Comment("关联的商品"),
	}
}
