package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// PaymentCallback holds the schema definition for the PaymentCallback entity.
type PaymentCallback struct {
	ent.Schema
}

// Fields of the PaymentCallback.
func (PaymentCallback) Fields() []ent.Field {
	return []ent.Field{
		field.String("seq_no").
			Immutable().
			Comment("流水号"),
		field.Enum("type").
			GoType(domain.PaymentCallbackType("")).
			Immutable().
			Comment("回调类型"),
		field.JSON("raw", json.RawMessage{}).
			Immutable().
			Comment("原始数据"),
		field.Enum("provider").
			GoType(domain.PayProvider("")).
			Immutable().
			Comment("支付供应商"),
	}
}

// Edges of the PaymentCallback.
func (PaymentCallback) Edges() []ent.Edge {
	return nil
}

// Mixin of the PaymentCallback.
func (PaymentCallback) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
