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

// OrderItemSetMealDetail holds the schema definition for the OrderItemSetMealDetail entity.
type OrderItemSetMealDetail struct {
	ent.Schema
}

// Fields of the OrderItemSetMealDetail.
func (OrderItemSetMealDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order_item_id").
			Positive().
			Immutable().
			Comment("订单商品项ID"),
		field.String("name").
			MaxLen(255).
			NotEmpty().
			Comment("商品名称"),
		field.Int("type").
			Immutable().
			Default(1).
			Comment("商品类型"),
		field.Other("set_meal_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Default(decimal.Zero).
			Comment("商品单价/套餐内单价"),
		field.Int("set_meal_id").
			Positive().
			Immutable().
			Comment("套餐ID"),
		field.Int("product_id").
			Positive().
			Immutable().
			Comment("商品ID"),
		field.Other("quantity", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Default(decimal.Zero).
			Comment("数量"),
		field.JSON("product_snapshot", domain.OrderProductInfoSnapshot{}).
			Immutable().
			Comment("商品信息快照"),
	}
}

// Edges of the OrderItemSetMealDetail.
func (OrderItemSetMealDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order_item", OrderItem.Type).
			Field("order_item_id").
			Immutable().
			Ref("set_meal_details").
			Unique().
			Required(),
	}
}

func (OrderItemSetMealDetail) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
