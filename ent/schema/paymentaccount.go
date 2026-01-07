package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// PaymentAccount 品牌商收款账户
type PaymentAccount struct {
	ent.Schema
}

func (PaymentAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the PaymentAccount.
func (PaymentAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.Enum("channel").
			GoType(domain.PaymentChannel("")).
			Comment("支付渠道：rm"),
		field.String("merchant_number").MaxLen(255).NotEmpty().Comment("支付商户号"),
		field.String("merchant_name").MaxLen(255).NotEmpty().Comment("支付商户名称"),
		field.Bool("is_default").Default(false).Comment("是否默认"),
	}
}

// Indexes of the PaymentAccount.
func (PaymentAccount) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		// 唯一索引：品牌商+渠道在当前品牌商下唯一
		index.Fields("merchant_id", "channel", "deleted_at").Unique(),
	}
}
