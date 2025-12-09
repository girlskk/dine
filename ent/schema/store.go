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

// Store holds the schema definition for the Store entity.
type Store struct {
	ent.Schema
}

// Fields of the Store.
func (Store) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("名称"),
		field.Enum("type").GoType(domain.StoreType("")).Comment("类型"),
		field.Enum("cooperation_type").GoType(domain.StoreCooperationType("")).Comment("合作类型"),
		field.Bool("need_audit").Comment("是否需要审核"),
		field.Bool("enabled").Comment("是否启用"),
		field.Other("point_settlement_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(5,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）"),
		field.Other("point_withdrawal_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(5,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Decimal{}).
			Comment("积分提现佣金比例（单位：百分比，例如 0.1234 表示 12.34%）"),
		field.String("huifu_id").
			Comment("汇付ID"),
		field.String("zxh_id").
			Comment("知心话ID"),
		field.String("zxh_secret").
			Comment("知心话密钥"),
	}
}

// Edges of the Store.
func (Store) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("backend_user", BackendUser.Type).
			Unique(),
		edge.To("frontend_users", FrontendUser.Type),
		edge.To("store_info", StoreInfo.Type).Unique(),
		edge.To("store_finance", StoreFinance.Type).Unique(),
		edge.To("store_account", StoreAccount.Type).Unique(),
		edge.To("store_withdraws", StoreWithdraw.Type),
		edge.To("store_account_transactions", StoreAccountTransaction.Type),
	}
}

func (Store) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
