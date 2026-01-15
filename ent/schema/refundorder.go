package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// RefundOrder 退款订单表
type RefundOrder struct {
	ent.Schema
}

func (RefundOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (RefundOrder) Fields() []ent.Field {
	return []ent.Field{
		// ========== 租户信息 ==========
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Immutable().Comment("门店ID"),

		// ========== 营业信息 ==========
		field.String("business_date").NotEmpty().Comment("营业日"),
		field.String("shift_no").Optional().Comment("班次号"),
		field.String("refund_no").NotEmpty().Comment("退款单号"),

		// ========== 原订单关联 ==========
		field.UUID("origin_order_id", uuid.UUID{}).Comment("原订单ID"),
		field.String("origin_order_no").NotEmpty().Comment("原订单号"),
		field.Time("origin_paid_at").Optional().Nillable().Comment("原订单支付时间"),
		field.Other("origin_amount_paid", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).Optional().Nillable().Comment("原订单实付金额"),

		// ========== 退款类型与状态 ==========
		field.Enum("refund_type").GoType(domain.RefundType("")).Comment("退款类型"),
		field.Enum("refund_status").GoType(domain.RefundStatus("")).
			Default(string(domain.RefundStatusPending)).Comment("退款状态"),

		// ========== 退款原因 ==========
		field.String("refund_reason_code").Optional().Comment("退款原因代码"),
		field.Text("refund_reason").Optional().Comment("退款原因描述"),

		// ========== 操作人信息 ==========
		field.UUID("refunded_by", uuid.UUID{}).Optional().Comment("退款操作人ID"),
		field.String("refunded_by_name").Optional().Comment("退款操作人名称"),
		field.UUID("approved_by", uuid.UUID{}).Optional().Comment("审批人ID"),
		field.String("approved_by_name").Optional().Comment("审批人名称"),
		field.Time("approved_at").Optional().Nillable().Comment("审批时间"),

		// ========== 时间节点 ==========
		field.Time("refunded_at").Optional().Nillable().Comment("退款完成时间"),

		// ========== 终端信息 ==========
		field.JSON("store", domain.OrderStore{}).Comment("门店信息"),
		field.Enum("channel").GoType(domain.Channel("")).Default(string(domain.ChannelPOS)).Comment("退款渠道"),
		field.JSON("pos", domain.OrderPOS{}).Comment("POS终端信息"),
		field.JSON("cashier", domain.OrderCashier{}).Comment("收银员信息"),

		// ========== 金额与支付 ==========
		field.JSON("refund_amount", domain.RefundAmount{}).Comment("退款金额明细"),
		field.JSON("refund_payments", []domain.RefundPayment{}).Optional().Comment("退款支付记录"),

		// ========== 备注 ==========
		field.String("remark").MaxLen(255).Optional().Comment("备注"),
	}
}

func (RefundOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("origin_order_id"),
		index.Fields("business_date"),
		index.Fields("refund_status"),
		index.Fields("store_id", "refund_no", "deleted_at").Unique(),
	}
}

func (RefundOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("refund_products", RefundOrderProduct.Type).Comment("退款商品明细"),
	}
}
