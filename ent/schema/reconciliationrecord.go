package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

type ReconciliationRecord struct {
	ent.Schema
}

func (ReconciliationRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
	}
}

func (ReconciliationRecord) Fields() []ent.Field {
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
		field.Enum("channel").
			GoType(domain.OrderPaidChannel("")).
			Immutable().
			Comment("支付渠道"),
		field.Time("date").SchemaType(map[string]string{
			dialect.MySQL:  "DATE",
			dialect.SQLite: "TEXT",
		}).Immutable().Comment("账单日期"),
	}
}

func (ReconciliationRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("date", "store_id", "channel").Unique(),
	}
}

// Edges of the ReconciliationRecord.
func (ReconciliationRecord) Edges() []ent.Edge {
	return nil
}
