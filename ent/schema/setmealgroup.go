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

// SetMealGroup 套餐组
type SetMealGroup struct {
	ent.Schema
}

func (SetMealGroup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the SetMealGroup.
func (SetMealGroup) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("product_id", uuid.UUID{}).Immutable().Comment("套餐商品ID（外键）"),
		field.String("name").MaxLen(255).NotEmpty().Comment("套餐组名称（必选）"),
		field.Enum("selection_type").
			GoType(domain.SetMealGroupSelectionType("")).
			Default(string(domain.SetMealGroupSelectionTypeFixed)).
			Comment("点单限制：fixed（固定分组）、optional（可选套餐）"),
	}
}

func (SetMealGroup) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("product_id"),
	}
}

// Edges of the SetMealGroup.
func (SetMealGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).
			Ref("set_meal_groups").
			Field("product_id").
			Immutable().
			Required().
			Unique().
			Comment("所属套餐商品"),
		edge.To("details", SetMealDetail.Type).Comment("套餐组详情列表"),
	}
}
