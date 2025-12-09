package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

type DineTable struct {
	ent.Schema
}

func (DineTable) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (DineTable) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
		field.Int("seat_count").Positive().Comment("座位数"),
		field.Int("status").Default(1).Comment("座位状态"),
		field.Int("area_id"),
		field.Int("order_id").
			Positive().
			Optional().
			Nillable().
			Comment("当前订单ID"),
	}
}

func (DineTable) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tablearea", TableArea.Type).Ref("dinetables").Unique().Required().Field("area_id"),
		edge.From("order", Order.Type).Ref("current_dinetable").Unique().Field("order_id"),
		edge.To("orders", Order.Type),
	}
}
