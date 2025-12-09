package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// OrderLog holds the schema definition for the OrderLog entity.
type OrderLog struct {
	ent.Schema
}

// Fields of the OrderLog.
func (OrderLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order_id").
			Immutable().
			Positive().
			Comment("订单ID"),
		field.Enum("event").
			GoType(domain.OrderLogEvent("")).
			Comment("事件"),
		field.Enum("operator_type").
			GoType(domain.OperatorType("")).
			Immutable().
			Comment("操作人类型"),
		field.Int("operator_id").
			Immutable().
			Default(0).
			Comment("操作人ID"),
		field.String("operator_name").
			Immutable().
			Comment("操作人姓名"),
	}
}

// Edges of the OrderLog.
func (OrderLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Field("order_id").
			Immutable().
			Ref("logs").
			Unique().
			Required(),
	}
}

// Mixin of the OrderLog.
func (OrderLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
