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

// RemarkCategory holds the schema definition for the RemarkCategory entity.
type RemarkCategory struct {
	ent.Schema
}

// Mixin of the RemarkCategory.
func (RemarkCategory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the RemarkCategory.
func (RemarkCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(50).
			Comment("备注分类名称，例如整单备注、退菜原因等"),

		field.Enum("remark_scene").
			GoType(domain.RemarkScene("")).
			Comment("使用场景：整单备注/单品备注/退菜原因等"),

		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Comment("品牌商ID，可为空表示系统级分类"),

		field.String("description").
			Default("").
			MaxLen(255).
			Comment("备注分类描述"),

		field.Int("sort_order").
			Default(1000).
			Comment("排序，值越小越靠前"),
	}
}

func (RemarkCategory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
	}
}

// Edges of the RemarkCategory.
func (RemarkCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("remark_categories").
			Field("merchant_id").
			Unique(),
	}
}
