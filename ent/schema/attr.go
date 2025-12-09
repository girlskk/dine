package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Attr 商品属性
type Attr struct {
	ent.Schema
}

func (Attr) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Attr.
func (Attr) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
	}
}

func (Attr) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "name"),
	}
}

// Edges of the Attr.
func (Attr) Edges() []ent.Edge {
	return []ent.Edge{
		// M -> M：商品属性
		edge.From("products", Product.Type).Ref("attrs"),
		// 购物车商品
		edge.To("order_cart_items", OrderCart.Type),
	}
}
