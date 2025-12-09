package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// OrderItem holds the schema definition for the OrderItem entity.
type OrderItem struct {
	ent.Schema
}

// Fields of the OrderItem.
func (OrderItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order_id").
			Immutable().
			Positive().
			Comment("订单ID"),
		field.Int("product_id").
			Immutable().
			Positive().
			Comment("商品ID"),
		field.String("name").
			Immutable().
			MaxLen(255).
			NotEmpty().
			Comment("商品名称"),
		field.Int("type").
			Immutable().
			Default(1).
			Comment("商品类型"),
		field.Bool("allow_point_pay").
			Immutable().
			Default(true).
			Comment("是否支持积分支付"),
		field.Other("quantity", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("商品数量"),
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("商品实际单价"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("商品总金额"),
		field.String("remark").
			MaxLen(500).
			Comment("备注"),
		field.JSON("product_snapshot", domain.OrderProductInfoSnapshot{}).
			Immutable().
			Comment("商品信息快照"),
	}
}

// Edges of the OrderItem.
func (OrderItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Field("order_id").
			Immutable().
			Ref("items").
			Unique().
			Required(),
		edge.To("set_meal_details", OrderItemSetMealDetail.Type),
	}
}

// Mixin of the OrderItem.
func (OrderItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
