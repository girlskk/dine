package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// StorePaymentAccount 门店收款账户
type StorePaymentAccount struct {
	ent.Schema
}

func (StorePaymentAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the StorePaymentAccount.
func (StorePaymentAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Immutable().Comment("门店ID"),
		field.UUID("payment_account_id", uuid.UUID{}).Immutable().Comment("品牌商收款账户ID"),
		field.String("merchant_number").MaxLen(255).NotEmpty().Comment("支付商户号"),
	}
}

// Indexes of the StorePaymentAccount.
func (StorePaymentAccount) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id"),
		index.Fields("payment_account_id"),
		// 唯一索引：每个门店的品牌商收款账户ID只能对应一个门店的收款账户
		index.Fields("store_id", "payment_account_id", "deleted_at").Unique(),
	}
}

// Edges of the StorePaymentAccount.
func (StorePaymentAccount) Edges() []ent.Edge {
	return []ent.Edge{
		// 所属门店
		edge.From("store", Store.Type).
			Ref("store_payment_accounts").
			Field("store_id").
			Unique().
			Immutable().
			Required(),
		// 关联的品牌商收款账户
		edge.From("payment_account", PaymentAccount.Type).
			Ref("store_payment_accounts").
			Field("payment_account_id").
			Unique().
			Immutable().
			Required(),
	}
}
