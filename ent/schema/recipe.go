package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Recipe 商品做法
type Recipe struct {
	ent.Schema
}

func (Recipe) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Recipe.
func (Recipe) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty(),
		field.Int("store_id").Positive(),
	}
}

func (Recipe) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "name"),
	}
}

// Edges of the Recipe.
func (Recipe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("products", Product.Type).Ref("recipes"),
		// 购物车商品
		edge.To("order_cart_items", OrderCart.Type),
	}
}
