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

// PaymentMethod 结算方式
type PaymentMethod struct {
	ent.Schema
}

func (PaymentMethod) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the PaymentMethod.
func (PaymentMethod) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("结算方式名称"),
		field.Enum("accounting_rule").
			GoType(domain.PaymentMethodAccountingRule("")).
			Default(string(domain.PaymentMethodAccountingRuleIncome)).
			Comment("计入规则:income-计入实收,discount-计入优惠"),
		field.Enum("payment_type").
			GoType(domain.PaymentMethodPayType("")).
			Default(string(domain.PaymentMethodPayTypeOther)).
			Comment("结算类型:other-其他,cash-现金,offline_card-线下刷卡,custom_coupon-自定义券,partner_coupon-三方合作券"),
		field.Other("fee_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("手续费率,百分比"),
		field.Enum("invoice_rule").
			GoType(domain.PaymentMethodInvoiceRule("")).
			Default(string(domain.PaymentMethodInvoiceRuleActualAmount)).
			Comment("实收部分开票规则:no_invoice-不开发票,actual_amount-按实收金额"),
		field.Bool("cash_drawer_status").Default(false).Comment("开钱箱状态:false-不开钱箱, true-开钱箱（必选）"),
		field.JSON("display_channels", []domain.PaymentMethodDisplayChannel{}).Comment("收银终端显示渠道（可选，可多选）：POS、移动点餐、扫码点餐、自助点餐、三方外卖"),
		field.Bool("status").Default(false).Comment("启用/停用状态: true-启用, false-停用（必选）"),
	}
}

func (PaymentMethod) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "deleted_at").Unique(),
	}
}

// Edges of the PaymentMethod.
func (PaymentMethod) Edges() []ent.Edge {
	return []ent.Edge{}
}
