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

// MerchantRenewal holds the schema definition for the MerchantRenewal entity.
type MerchantRenewal struct {
	ent.Schema
}

// Fields of the MerchantRenewal.
func (MerchantRenewal) Fields() []ent.Field {
	return []ent.Field{
		// 商户基础信息
		field.UUID("merchant_id", uuid.UUID{}).
			Comment("商户 ID"),
		field.Int("purchase_duration").
			Default(0).
			Comment("购买时长"),
		field.Enum("purchase_duration_unit").
			GoType(domain.PurchaseDurationUnit("")).
			Immutable().
			Comment("购买时长单位"),
		field.String("operator_name").
			Optional().
			Default("").
			MaxLen(50).
			Comment("操作人"),
		field.String("operator_account").
			Optional().
			Default("").
			MaxLen(50).
			Comment("操作人账号"),
	}
}

// Edges of the MerchantRenewal.
func (MerchantRenewal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("merchant_renewals").
			Field("merchant_id").
			Unique().
			Required(),
	}
}

func (MerchantRenewal) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (MerchantRenewal) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
	}
}
