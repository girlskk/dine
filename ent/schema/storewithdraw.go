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

// StoreWithdraw holds the schema definition for the StoreWithdraw entity.
type StoreWithdraw struct {
	ent.Schema
}

// Fields of the StoreWithdraw.
func (StoreWithdraw) Fields() []ent.Field {
	return []ent.Field{
		field.Int("store_id").Positive().Comment("门店ID"),
		field.String("store_name").Comment("门店名称"),
		field.String("no").Unique().Comment("单据编号"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("提现金额（原始申请金额）"),
		field.Other("point_withdrawal_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(5,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("积分提现费率"),
		field.Other("actual_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("实际到账金额（扣除平台佣金）"),
		field.Enum("account_type").
			GoType(domain.AccountType("")).
			Comment("账户类型：public-对公 private-对私"),
		field.String("bank_account").Comment("银行账号"),
		field.String("bank_card_name").Comment("银行卡名称（对公时为公司名称）"),
		field.String("bank_name").Comment("银行名称"),
		field.String("bank_branch").Comment("开户支行"),
		field.Other("invoice_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("开票金额"),
		field.Int("status").
			GoType(domain.StoreWithdrawStatus(0)).
			Comment("提现状态：1-待审核 2-已审核 3-已驳回"),
	}
}

// Edges of the StoreWithdraw.
func (StoreWithdraw) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).Ref("store_withdraws").Field("store_id").Unique().Required(),
	}
}

func (StoreWithdraw) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
