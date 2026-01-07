package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

func (Order) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Immutable().Comment("门店ID"),
		field.String("business_date").NotEmpty().Comment("营业日（门店营业日，字符串）"),
		field.String("shift_no").Optional().Comment("班次号"),
		field.String("order_no").NotEmpty().Comment("订单号（门店内唯一可读编号）"),

		field.Enum("order_type").GoType(domain.OrderType("")).Default(string(domain.OrderTypeSale)).Comment("订单类型：SALE=销售单；REFUND=退单；PARTIAL_REFUND=部分退款单"),

		field.Time("placed_at").Optional().Nillable().Comment("下单时间"),
		field.Time("paid_at").Optional().Nillable().Comment("支付完成时间"),
		field.Time("completed_at").Optional().Nillable().Comment("完成时间"),

		field.UUID("placed_by", uuid.UUID{}).Optional().Comment("下单人ID"),
		field.String("placed_by_name").Optional().Comment("下单人名称"),

		field.Enum("dining_mode").GoType(domain.DiningMode("")).Default(string(domain.DiningModeDineIn)).Comment("就餐模式：DINE_IN=堂食"),
		field.Enum("order_status").GoType(domain.OrderStatus("")).Default(string(domain.OrderStatusPlaced)).Comment("订单业务状态：PLACED=已下单；COMPLETED=已完成；CANCELLED=已取消"),
		field.Enum("payment_status").GoType(domain.PaymentStatus("")).Default(string(domain.PaymentStatusUnpaid)).Comment("支付状态：UNPAID=未支付；PAYING=支付中；PAID=已支付；REFUNDED=全额退款"),

		field.UUID("table_id", uuid.UUID{}).Optional().Comment("桌位ID（堂食）"),
		field.String("table_name").Optional().Comment("桌位名称（堂食，如A01/1号桌）"),
		field.Int("guest_count").Optional().Comment("用餐人数（堂食）"),

		field.JSON("store", domain.OrderStore{}).Comment("门店信息"),
		field.Enum("channel").GoType(domain.Channel("")).Default(string(domain.ChannelPOS)).Comment("下单渠道"),
		field.JSON("pos", domain.OrderPOS{}).Comment("POS终端信息"),
		field.JSON("cashier", domain.OrderCashier{}).Comment("收银员信息"),

		field.JSON("tax_rates", []domain.OrderTaxRate{}).Optional().Comment("税率明细"),
		field.JSON("fees", []domain.OrderFee{}).Optional().Comment("费用明细"),
		field.JSON("payments", []domain.OrderPayment{}).Optional().Comment("支付记录"),

		field.JSON("amount", domain.OrderAmount{}).Comment("金额汇总"),

		field.String("remark").MaxLen(255).Optional().Comment("整单备注"),
	}
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("table_id"),
		index.Fields("order_status"),
		index.Fields("payment_status"),
		index.Fields("store_id", "order_no", "deleted_at").Unique(),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("order_products", OrderProduct.Type).Comment("订单商品明细"),
	}
}
