package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// PointSettlement holds the schema definition for the PointSettlement entity.
type PointSettlement struct {
	ent.Schema
}

func (PointSettlement) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
	}
}

// Fields of the PointSettlement.
func (PointSettlement) Fields() []ent.Field {
	return []ent.Field{
		field.String("no").NotEmpty().Immutable().Unique().Comment("单号"),
		field.Int("store_id").Positive().Immutable().Comment("门店ID"),
		field.String("store_name").Immutable().MaxLen(255).Comment("门店名称"),
		field.Int("order_count").Positive().Immutable().Comment("入账笔数"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).Immutable().Comment("入账金额"),
		field.Other("total_points", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).Immutable().Comment("积分总额"),
		field.Time("date").SchemaType(map[string]string{
			dialect.MySQL:  "DATE",
			dialect.SQLite: "TEXT",
		}).Immutable().Comment("账单日期"),
		field.Int("status").Default(1).Comment("账单状态"),
		field.Other("point_settlement_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(5,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）"),
		field.Time("approved_at").Optional().Nillable().Comment("审批时间"),
		field.Int("approver_id").Optional().Comment("审批者ID"),
	}
}

func (PointSettlement) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date", "store_id").Unique(),
	}
}

// Edges of the PointSettlement.
func (PointSettlement) Edges() []ent.Edge {
	return nil
}
