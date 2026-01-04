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

// ProductAttr 商品口味做法
type ProductAttr struct {
	ent.Schema
}

func (ProductAttr) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductAttr.
func (ProductAttr) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("口味做法名称"),
		field.JSON("channels", []domain.SaleChannel{}).Comment("售卖渠道列表"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Default(schematype.NilUUID()).Immutable().Comment("门店ID"),
		field.Int("product_count").Default(0).Comment("关联的商品数量（所有项的关联商品数量累加）"),
	}
}

func (ProductAttr) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引
		index.Fields("merchant_id", "store_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the ProductAttr.
func (ProductAttr) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", ProductAttrItem.Type).Comment("口味做法项列表"),
		edge.To("product_attrs", ProductAttrRelation.Type),
	}
}
