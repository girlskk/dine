package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// StoreAccount holds the schema definition for the StoreAccount entity.
type StoreAccount struct {
	ent.Schema
}

// Fields of the StoreAccount.
func (StoreAccount) Fields() []ent.Field {
	return []ent.Field{
		field.Int("store_id").Positive().Comment("关联门店"),
		field.Other("balance", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("账户余额（可提现金额）"),
		field.Other("pending_withdraw", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("待提现金额"),
		field.Other("withdrawn", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("已提现金额"),
		field.Other("total_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("总收益金额"),
	}
}

// Edges of the StoreAccount.
func (StoreAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).Ref("store_account").Field("store_id").Unique().Required(),
	}
}

func (StoreAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
