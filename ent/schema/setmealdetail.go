package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// SetMealDetail 套餐商品详情
type SetMealDetail struct {
	ent.Schema
}

func (SetMealDetail) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the SetMealDetail.
func (SetMealDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("套餐内价格"),

		field.Other("quantity", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(8,3)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).Comment("数量"),

		field.Int("set_meal_id").Comment("套餐商品ID"),
		field.Int("product_id").Comment("套餐详情商品ID"),
		field.Int("product_spec_id").Optional().Comment("商品-规格ID"),
	}
}

// Edges of the SetMealDetail.
func (SetMealDetail) Edges() []ent.Edge {
	return []ent.Edge{
		// 套餐商品本身的 productID
		edge.From("set_meal", Product.Type).Ref("set_meal_details").Field("set_meal_id").
			Unique().Required(),

		// 套餐详情具体商品ID
		edge.From("product", Product.Type).Ref("included_in_set_meals").Field("product_id").
			Unique().Required(),

		// 套餐详情具体商品的规格ID
		edge.From("product_spec", ProductSpec.Type).Ref("set_meal_details").Field("product_spec_id").
			Unique().Comment("关联具体商品规格"),
	}
}
