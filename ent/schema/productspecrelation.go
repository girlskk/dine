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

// ProductSpecRelation 商品-规格关联
type ProductSpecRelation struct {
	ent.Schema
}

func (ProductSpecRelation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProductSpecRelation.
func (ProductSpecRelation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("product_id", uuid.UUID{}).Immutable().Comment("商品ID（外键）"),
		field.UUID("spec_id", uuid.UUID{}).Immutable().Comment("规格ID（外键）"),
		field.Other("base_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("基础价格（必选，单位：分）"),
		field.Other("member_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Comment("会员价（可选，单位：分）"),
		field.UUID("packing_fee_id", uuid.UUID{}).Comment("打包费ID（引用费用配置）"),
		field.Other("estimated_cost_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Comment("预估成本价（可选，单位：分）"),
		field.Other("other_price1", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Comment("其他价格1（可选，单位：分）"),
		field.Other("other_price2", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Comment("其他价格2（可选，单位：分）"),
		field.Other("other_price3", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Comment("其他价格3（可选，单位：分）"),
		field.String("barcode").MaxLen(255).Default("").Comment("条形码（可选，字符串，无限制）"),
		field.Bool("is_default").Default(false).Comment("是否默认项（规格必须至少有一个默认项）"),
	}
}

func (ProductSpecRelation) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_id"),
		index.Fields("spec_id"),
	}
}

// Edges of the ProductSpecItem.
func (ProductSpecRelation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("product_specs").
			Field("product_id").
			Immutable().
			Required().
			Unique().
			Comment("所属商品"),
		edge.From("spec", ProductSpec.Type).
			Ref("product_specs").
			Field("spec_id").
			Immutable().
			Required().
			Unique().
			Comment("所属规格"),
	}
}
