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

// TaxFee holds the schema definition for the TaxFee entity.
type TaxFee struct {
	ent.Schema
}

// Fields of the TaxFee.
func (TaxFee) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(50).
			Comment("税费名称，长度不大于50字"),
		field.Enum("tax_fee_type").
			GoType(domain.TaxFeeType("")).
			Immutable().
			Comment("税费类型：商户/门店"),
		field.String("tax_code").
			NotEmpty().
			MaxLen(50).
			Immutable().
			Comment("税费代码，长度不大于20字"),
		field.Enum("tax_rate_type").
			GoType(domain.TaxRateType("")).
			Default(string(domain.TaxRateTypeUnified)).
			Comment("税率类型：统一比例/自定义比例"),
		field.Other("tax_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("税率，6% -> 0.06"),
		field.Bool("default_tax").
			Default(false).
			Comment("是否默认税费"),
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

// Edges of the TaxFee.
func (TaxFee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("tax_fees").
			Field("merchant_id").
			Unique().
			Immutable(),
		edge.From("store", Store.Type).
			Ref("tax_fees").
			Field("store_id").
			Unique().
			Immutable(),
	}
}

func (TaxFee) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引：同一商户/门店下税费名称唯一
		index.Fields("name", "merchant_id", "store_id", "deleted_at").
			Unique().
			StorageKey("idx_tax_fee_name_merchant_store_deleted"),
	}
}

func (TaxFee) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
