package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// SetMealDetail 套餐组详情
type SetMealDetail struct {
	ent.Schema
}

func (SetMealDetail) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the SetMealDetail.
func (SetMealDetail) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("group_id", uuid.UUID{}).Immutable().Comment("套餐组ID（外键）"),
		field.UUID("product_id", uuid.UUID{}).Immutable().Comment("商品ID（外键，引用普通商品）"),
		field.Int("quantity").Default(1).Comment("数量（必选，必须为正整数）"),
		field.Bool("is_default").Default(false).Comment("是否默认（必选，每个套餐组中只能有一个默认项）"),
		field.JSON("optional_product_ids", []uuid.UUID{}).Optional().Comment("备选商品（可选，多选，当商品库存不足时，自动替换为备选商品）"),
	}
}

func (SetMealDetail) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("group_id"),
		index.Fields("product_id"),
	}
}

// Edges of the SetMealDetail.
func (SetMealDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", SetMealGroup.Type).
			Ref("details").
			Field("group_id").
			Immutable().
			Required().
			Unique().
			Comment("所属套餐组"),
		edge.From("product", Product.Type).
			Ref("set_meal_details").
			Field("product_id").
			Immutable().
			Required().
			Unique().
			Comment("关联商品"),
	}
}
