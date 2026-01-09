package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// MerchantBusinessType holds the schema definition for the MerchantBusinessType entity.
type MerchantBusinessType struct {
	ent.Schema
}

// Fields of the MerchantBusinessType.
func (MerchantBusinessType) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("商户 ID"),
		field.String("type_code").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("业态类型编码（保留字段）"),
		field.String("type_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("业态类型名称"),
	}
}

// Edges of the MerchantBusinessType.
func (MerchantBusinessType) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (MerchantBusinessType) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (MerchantBusinessType) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type_code", "deleted_at").Unique(),
	}
}
