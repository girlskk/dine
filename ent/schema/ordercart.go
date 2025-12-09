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

// OrderCart holds the schema definition for the OrderCart entity.
type OrderCart struct {
	ent.Schema
}

func (OrderCart) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the OrderCart.
func (OrderCart) Fields() []ent.Field {
	return []ent.Field{
		field.Int("table_id").
			Positive().
			Comment("台桌ID"),

		field.Int("product_id").
			Positive().
			Comment("商品ID"),

		field.Other("quantity", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("商品数量"),

		field.Int("product_spec_id").
			Optional().
			Comment("商品规格ID"),

		field.Int("attr_id").
			Optional().
			Comment("属性ID"),

		field.Int("recipe_id").
			Optional().
			Comment("做法ID"),
	}
}

// Indexes of the OrderCart.
func (OrderCart) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("table_id"),
	}
}

// Edges of the OrderCart.
func (OrderCart) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("order_cart_items").
			Field("product_id").
			Unique().
			Required(),

		edge.From("product_spec", ProductSpec.Type).
			Ref("order_cart_items").
			Field("product_spec_id").
			Unique(),

		edge.From("attr", Attr.Type).
			Ref("order_cart_items").
			Field("attr_id").
			Unique(),

		edge.From("recipe", Recipe.Type).
			Ref("order_cart_items").
			Field("recipe_id").
			Unique(),
	}
}
