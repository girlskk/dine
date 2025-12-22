package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProductUnit 商品单位
type ProductUnit struct {
	ent.Schema
}

func (ProductUnit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductUnit.
func (ProductUnit) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("单位名称"),
		field.Enum("type").GoType(domain.ProductUnitType("")).Comment("单位类型：quantity（数量单位）、weight（重量单位）"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Optional().Immutable().Comment("门店ID"),
		field.Int("product_count").Default(0).Comment("关联的商品数量"),
	}
}

func (ProductUnit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
	}
}

// Edges of the ProductUnit.
func (ProductUnit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("products", Product.Type).Comment("关联的商品"),
	}
}
