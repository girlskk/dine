package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
)

// ProductSpec 商品-商品规格关联
type ProductSpec struct {
	ent.Schema
}

// Fields of the ProductSpec.
func (ProductSpec) Fields() []ent.Field {
	return []ent.Field{
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("规格价格"),
		field.Int("sale_status").Default(1).Comment("商品可售状态"),
		// 显示指定外键字段
		field.Int("product_id"),
		field.Int("spec_id"),
	}
}

func (ProductSpec) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_id", "spec_id").Unique(),
	}
}

// Edges of the ProductSpec.
func (ProductSpec) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).Required().Ref("product_specs").Field("product_id").Unique(),
		edge.From("spec", Spec.Type).Required().Ref("product_specs").Field("spec_id").Unique(),
		// 套餐商品详情-商品规格关联
		edge.To("set_meal_details", SetMealDetail.Type),
		// 购物车商品
		edge.To("order_cart_items", OrderCart.Type),
	}
}
