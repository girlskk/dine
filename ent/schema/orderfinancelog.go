package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// OrderFinanceLog holds the schema definition for the OrderFinanceLog entity.
type OrderFinanceLog struct {
	ent.Schema
}

// Fields of the OrderFinanceLog.
func (OrderFinanceLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("order_id").
			Immutable().
			Comment("订单ID"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Comment("交易金额"),
		field.Enum("type").
			GoType(domain.OrderFinanceLogType("")).
			Immutable().
			Comment("交易类型"),
		field.Enum("channel").
			GoType(domain.OrderPaidChannel("")).
			Immutable().
			Comment("支付渠道"),
		field.String("seq_no").
			Immutable().
			Comment("流水号"),
		field.Enum("creator_type").
			GoType(domain.OperatorType("")).
			Immutable().
			Comment("创建人类型"),
		field.Int("creator_id").
			Immutable().
			Default(0).
			Comment("创建人ID"),
		field.String("creator_name").
			Immutable().
			Comment("创建人姓名"),
	}
}

// Edges of the OrderFinanceLog.
func (OrderFinanceLog) Edges() []ent.Edge {
	return nil
}

// Mixin of the OrderFinanceLog.
func (OrderFinanceLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
