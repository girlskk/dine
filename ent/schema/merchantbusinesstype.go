package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MerchantBusinessType holds the schema definition for the MerchantBusinessType entity.
type MerchantBusinessType struct {
	ent.Schema
}

// Fields of the MerchantBusinessType.
func (MerchantBusinessType) Fields() []ent.Field {
	return []ent.Field{
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
	return []ent.Edge{
		edge.To("merchants", Merchant.Type),
	}
}
