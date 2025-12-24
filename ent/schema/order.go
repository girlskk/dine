package schema

import (
	"encoding/json"

	"entgo.io/ent"
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
		field.String("origin_order_id").Optional().Comment("原正单订单ID（退款/部分退款单使用）"),
		field.JSON("refund", json.RawMessage{}).Optional().Comment("退款单信息（包含原单信息与退款原因）"),

		field.Time("opened_at").Optional().Nillable().Comment("开台时间"),
		field.Time("placed_at").Optional().Nillable().Comment("下单时间"),
		field.Time("paid_at").Optional().Nillable().Comment("支付完成时间"),
		field.Time("completed_at").Optional().Nillable().Comment("完成时间"),

		field.String("opened_by").Optional().Comment("开台操作员ID"),
		field.String("placed_by").Optional().Comment("下单操作员ID"),
		field.String("paid_by").Optional().Comment("收款/支付确认操作员ID"),

		field.Enum("dining_mode").GoType(domain.DiningMode("")).Comment("就餐模式：DINE_IN=堂食；TAKEAWAY=外卖（自取/配送）"),
		field.Enum("order_status").GoType(domain.OrderStatus("")).Default(string(domain.OrderStatusDraft)).Comment("订单业务状态：DRAFT=草稿/购物车；PLACED=已下单；IN_PROGRESS=制作中；READY=可取餐；COMPLETED=已完成；CANCELLED=已取消；VOIDED=已作废；MERGED=已合并"),
		field.Enum("payment_status").GoType(domain.PaymentStatus("")).Default(string(domain.PaymentStatusUnpaid)).Comment("支付状态：UNPAID=未支付；PAYING=支付中；PARTIALLY_PAID=部分支付；PAID=已支付；PARTIALLY_REFUNDED=部分退款；REFUNDED=全额退款"),
		field.Enum("fulfillment_status").GoType(domain.FulfillmentStatus("")).Optional().Comment("交付状态：NONE=无；IN_RESTAURANT=店内用餐；SERVED=已上齐；PICKUP_PENDING=待取餐；PICKED_UP=已取餐；DELIVERING=配送中；DELIVERED=已送达"),
		field.Enum("table_status").GoType(domain.TableStatus("")).Optional().Comment("桌位状态：OPENED=已开台；TRANSFERRED=已转台；RELEASED=已释放"),

		field.String("table_id").Optional().Comment("桌位ID（堂食）"),
		field.String("table_name").Optional().Comment("桌位名称（堂食，如A01/1号桌）"),
		field.Int("table_capacity").Optional().Comment("桌位容量（几人桌，仅堂食）"),
		field.Int("guest_count").Optional().Comment("用餐人数（堂食）"),

		field.String("merged_to_order_id").Optional().Comment("合并到的目标订单ID（该订单被合并时使用）"),
		field.Time("merged_at").Optional().Nillable().Comment("合并时间（该订单被合并时使用）"),

		field.JSON("store", json.RawMessage{}).Default(json.RawMessage("{}")).Comment("门店信息"),
		field.JSON("channel", json.RawMessage{}).Default(json.RawMessage("{}")).Comment("下单渠道信息"),
		field.JSON("pos", json.RawMessage{}).Default(json.RawMessage("{}")).Comment("POS终端信息"),
		field.JSON("cashier", json.RawMessage{}).Default(json.RawMessage("{}")).Comment("收银员信息"),

		field.JSON("member", json.RawMessage{}).Optional().Comment("会员信息"),
		field.JSON("takeaway", json.RawMessage{}).Optional().Comment("外卖信息"),

		field.JSON("cart", json.RawMessage{}).Default(json.RawMessage("[]")).Comment("购物车商品列表"),
		field.JSON("products", json.RawMessage{}).Default(json.RawMessage("[]")).Comment("订单商品明细"),
		field.JSON("promotions", json.RawMessage{}).Optional().Comment("促销明细"),
		field.JSON("coupons", json.RawMessage{}).Optional().Comment("卡券明细"),
		field.JSON("tax_rates", json.RawMessage{}).Optional().Comment("税率明细"),
		field.JSON("fees", json.RawMessage{}).Optional().Comment("费用明细"),
		field.JSON("payments", json.RawMessage{}).Optional().Comment("支付记录"),
		field.JSON("refunds_products", json.RawMessage{}).Optional().Comment("退菜记录"),

		field.JSON("amount", json.RawMessage{}).Default(json.RawMessage("{}")).Comment("金额汇总"),
	}
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("store_id", "business_date"),
		index.Fields("table_id"),
		index.Fields("origin_order_id"),
		index.Fields("merged_to_order_id"),
		index.Fields("store_id", "order_status"),
		index.Fields("store_id", "payment_status"),
		index.Fields("store_id", "order_no", "deleted_at").Unique(),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return nil
}
