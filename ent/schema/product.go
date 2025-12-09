package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Product 商品
type Product struct {
	ent.Schema
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
		field.Int("type").Default(1).Comment("商品类型"),
		field.JSON("images", []string{}).Optional().Comment("商品图片"),
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Optional().Comment("商品价格"),
		field.Int("status").Default(1).Comment("商品审批状态"),
		field.Int("sale_status").Default(1).Comment("商品可售状态"),
		field.Bool("allow_point_pay").Default(true).Comment("是否支持积分支付"),
		field.Int("unit_id").Optional(),
		field.Int("category_id"),
	}
}

func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "name"),
		index.Fields("type"),
		index.Fields("status"),
		index.Fields("sale_status"),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("products").Unique().Required().Field("category_id"),
		edge.From("unit", Unit.Type).Ref("products").Unique().Field("unit_id"),

		// 商品属性、商品做法 Many2Many
		edge.To("attrs", Attr.Type).Comment("关联属性"),
		edge.To("recipes", Recipe.Type).Comment("关联做法"),

		// 商品规格
		edge.To("product_specs", ProductSpec.Type),

		// 套餐商品详情（套餐商品ID）
		edge.To("set_meal_details", SetMealDetail.Type),

		// 套餐商品详情（详情商品ID）
		edge.To("included_in_set_meals", SetMealDetail.Type),

		// 购物车商品
		edge.To("order_cart_items", OrderCart.Type),
	}
}
