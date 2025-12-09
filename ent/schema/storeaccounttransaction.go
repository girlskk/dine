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

// StoreAccountTransaction holds the schema definition for the StoreAccountTransaction entity.
type StoreAccountTransaction struct {
	ent.Schema
}

// Fields of the StoreAccountTransaction.
func (StoreAccountTransaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int("store_id").Positive().Comment("门店ID"),
		field.String("no").Comment("单据编号"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Comment("变动金额"),
		field.Other("after", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Comment("变动后金额"),
		field.Int("type").
			GoType(domain.TransactionType(0)).
			Comment("变动类型：1-销售进账 2-进账撤回 3-申请提现 4-提现通过 5-提现驳回"),
	}
}

// Edges of the StoreAccountTransaction.
func (StoreAccountTransaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).
			Ref("store_account_transactions").
			Field("store_id").
			Unique().
			Required(),
	}
}

// Mixin of the StoreAccountTransaction.
func (StoreAccountTransaction) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
