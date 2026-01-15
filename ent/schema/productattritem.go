package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProductAttrItem 商品口味做法项
type ProductAttrItem struct {
	ent.Schema
}

func (ProductAttrItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductAttrItem.
func (ProductAttrItem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("attr_id", uuid.UUID{}).Immutable().Comment("口味做法ID（外键）"),
		field.String("name").MaxLen(255).NotEmpty().Comment("口味做法项名称"),
		field.String("image").Default("").MaxLen(512).Comment("图片URL（可选）"),
		field.Other("base_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("基础加价（单位：分）"),
		field.Int("product_count").Default(0).Comment("关联的商品数量"),
	}
}

func (ProductAttrItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("attr_id"),
		// 唯一索引
		index.Fields("attr_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the ProductAttrItem.
func (ProductAttrItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("attr", ProductAttr.Type).
			Ref("items").
			Field("attr_id").
			Immutable().
			Required().
			Unique().
			Comment("所属的口味做法"),

		edge.To("product_attrs", ProductAttrRelation.Type),
	}
}
