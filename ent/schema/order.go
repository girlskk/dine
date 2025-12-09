package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("no").
			NotEmpty().
			Immutable().
			Unique().
			Comment("单号"),
		field.Enum("type").
			GoType(domain.OrderType("")).
			Immutable().
			Comment("订单类型"),
		field.Enum("source").
			GoType(domain.OrderSource("")).
			Immutable().
			Comment("订单来源"),
		field.Enum("status").
			GoType(domain.OrderStatus("")).
			Comment("订单状态"),
		field.Other("total_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("总价（优惠前）"),
		field.Other("discount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("优惠金额"),
		field.Other("real_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("实际总价（优惠后）"),
		field.Other("points_available", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("积分可用额度"),
		field.Other("paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("已支付金额"),
		field.Other("refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("已退款金额"),
		field.JSON("paid_channels", domain.OrderPaidChannels{}).
			Optional().
			Default(domain.OrderPaidChannels{}).
			Comment("支付渠道"),
		field.Other("cash_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("现金已支付金额"),
		field.Other("cash_refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("现金已退款金额"),
		field.Other("wechat_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("微信已支付金额"),
		field.Other("wechat_refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("微信已退款金额"),
		field.Other("alipay_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("支付宝已支付金额"),
		field.Other("alipay_refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("支付宝已退款金额"),
		field.Other("points_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("积分已支付金额"),
		field.Other("points_refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("积分已退款金额"),
		field.Other("points_wallet_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("知心话钱包已支付金额"),
		field.Other("points_wallet_refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("知心话钱包已退款金额"),
		field.Time("last_paid_at").
			Optional().
			Nillable().
			Comment("最后支付时间"),
		field.Time("finished_at").
			Optional().
			Nillable().
			Comment("完成时间"),
		field.Int("member_id").
			NonNegative().
			Default(0).
			Comment("会员ID"),
		field.String("member_name").
			Comment("会员姓名"),
		field.String("member_phone").
			Comment("会员手机号"),
		field.Int("store_id").
			Immutable().
			Positive().
			Comment("门店ID"),
		field.String("store_name").
			Immutable().
			Comment("门店名称"),
		field.Int("table_id").
			Optional().
			Nillable().
			Comment("台桌ID"),
		field.String("table_name").
			Comment("台桌名称"),
		field.Int("people_number").
			Positive().
			Comment("就餐人数"),
		field.Int("creator_id").
			Immutable().
			Positive().
			Comment("创建人ID"),
		field.String("creator_name").
			Immutable().
			Comment("创建人姓名"),
		field.Enum("creator_type").
			GoType(domain.OperatorType("")).
			Default(string(domain.OperatorTypeFrontend)).
			Immutable().
			Comment("创建人类型"),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", OrderItem.Type),
		edge.To("logs", OrderLog.Type),
		edge.To("current_dinetable", DineTable.Type).
			Unique(),
		edge.From("dinetable", DineTable.Type).
			Ref("orders").
			Unique().
			Field("table_id"),
	}
}

// Indexes of the Order.
func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "deleted_at"),
		index.Fields("store_id", "status", "deleted_at"),
	}
}

func (Order) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
