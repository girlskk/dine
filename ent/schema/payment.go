package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.String("seq_no").
			Immutable().
			Unique().
			Comment("流水号"),
		field.Enum("provider").
			GoType(domain.PayProvider("")).
			Immutable().
			Comment("支付供应商"),
		field.Enum("channel").
			GoType(domain.PayChannel("")).
			Immutable().
			Comment("支付渠道"),
		field.Enum("state").
			GoType(domain.PayState("")).
			Comment("支付状态"),
		field.Other("amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Immutable().
			Comment("支付金额"),
		field.String("goods_desc").
			Immutable().
			Comment("商品描述"),
		field.String("mch_id").
			Immutable().
			Comment("商户ID"),
		field.String("ip_addr").
			Immutable().
			Comment("IP地址"),
		field.JSON("req", json.RawMessage{}).
			Immutable().
			Comment("请求参数"),
		field.JSON("resp", json.RawMessage{}).
			Immutable().
			Optional().
			Default(json.RawMessage{}).
			Comment("响应参数"),
		field.JSON("callback", json.RawMessage{}).
			Comment("回调参数"),
		field.Time("finished_at").
			Optional().
			Nillable().
			Comment("完成时间"),
		field.Other("refunded", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Default(decimal.Zero).
			Comment("已退款金额"),
		field.String("fail_reason").
			Optional().
			Comment("失败原因"),
		field.Enum("pay_biz_type").
			GoType(domain.PayBizType("")).
			Comment("支付业务类型"),
		field.Int("biz_id").
			Immutable().
			Comment("业务ID"),
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
		field.Int("store_id").
			Immutable().
			Positive().
			Comment("门店ID"),
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return nil
}

// Mixin of the Payment.
func (Payment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Indexes of the Payment.
func (Payment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("pay_biz_type", "biz_id", "finished_at", "deleted_at"),
	}
}
