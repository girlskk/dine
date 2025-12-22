package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProductAttrRelation 商品-口味做法项关联
type ProductAttrRelation struct {
	ent.Schema
}

func (ProductAttrRelation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductAttrRelation.
func (ProductAttrRelation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("product_id", uuid.UUID{}).Immutable().Comment("商品ID（外键）"),
		field.UUID("attr_id", uuid.UUID{}).Immutable().Comment("口味做法ID（外键）"),
		field.UUID("attr_item_id", uuid.UUID{}).Immutable().Comment("口味做法项ID（外键）"),
		field.Bool("is_default").Default(false).Comment("是否默认项（当点单限制为必选一项时，必须设置其中一个为默认项）"),
	}
}

func (ProductAttrRelation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_id"),
		index.Fields("attr_id"),
		index.Fields("attr_item_id"),
	}
}

// Edges of the ProductAttrRelation.
func (ProductAttrRelation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("product_attrs").
			Field("product_id").
			Immutable().
			Required().
			Unique().
			Comment("所属商品"),
		edge.From("attr", ProductAttr.Type).
			Ref("product_attrs").
			Field("attr_id").
			Immutable().
			Required().
			Unique().
			Comment("所属口味做法"),
		edge.From("attr_item", ProductAttrItem.Type).
			Ref("product_attrs").
			Field("attr_item_id").
			Immutable().
			Required().
			Unique().
			Comment("关联的口味做法项"),
	}
}
