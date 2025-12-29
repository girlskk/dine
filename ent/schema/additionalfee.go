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

// AdditionalFee holds the schema definition for the AdditionalFee entity.
type AdditionalFee struct {
	ent.Schema
}

// Fields of the AdditionalFee.
func (AdditionalFee) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(50).
			Comment("附加费名称，长度不大于50字"),
		field.Enum("fee_type").
			GoType(domain.AdditionalFeeType("")).
			Immutable().
			Comment("附加费类型：商户/门店"),
		field.Enum("fee_category").
			GoType(domain.AdditionalCategory("")).
			Comment("附加费类别：服务费/桌台费/打包费"),
		field.Enum("charge_mode").
			GoType(domain.AdditionalFeeChargeMode("")).
			Comment("费用类型：percent 百分比，fixed 固定金额"),
		field.Other("fee_value", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("费用数值；fixed/percent"),
		field.Bool("include_in_receivable").
			Default(false).
			Comment("是否计入实收"),
		field.Bool("taxable").
			Default(false).
			Comment("附加费是否收税"),
		field.Enum("discount_scope").
			GoType(domain.AdditionalFeeDiscountScope("")).
			Comment("附加费折扣场景：折前/折后"),
		field.JSON("order_channels", []domain.OrderChannel{}).
			Comment("订单渠道，(可多选) 例如：pos/扫码点餐/移动点餐/自助点餐/三方外卖，字符串数组"),
		field.JSON("dining_ways", []domain.DiningWay{}).
			Comment("就餐方式,(可多选) 堂食/外带/外卖，字符串数组"),
		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),
		field.Int("sort_order").
			Default(1000).
			Comment("排序，值越小越靠前"),
		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("品牌商 ID"),
		field.UUID("store_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("门店ID，可为空表示品牌维度"),
	}
}

// Edges of the AdditionalFee.
func (AdditionalFee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("additional_fees").
			Field("merchant_id").
			Unique().
			Immutable(),
		edge.From("store", Store.Type).
			Ref("additional_fees").
			Field("store_id").
			Unique().
			Immutable(),
	}
}

func (AdditionalFee) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引：同一商户/门店下附加费名称唯一
		index.Fields("name", "merchant_id", "store_id", "deleted_at").
			Unique().
			StorageKey("idx_additional_fee_name_merchant_store_deleted"),
	}
}

func (AdditionalFee) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
