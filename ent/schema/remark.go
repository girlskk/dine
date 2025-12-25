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

// Remark holds the schema definition for the Remark entity.
type Remark struct {
	ent.Schema
}

func (Remark) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Remark.
func (Remark) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(50).
			Comment("备注名称"),

		field.Enum("remark_type").
			GoType(domain.RemarkType("")).
			Immutable().
			Comment("备注类型：系统、品牌"),

		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),

		field.Int("sort_order").
			Default(1000).
			Comment("排序，值越小越靠前"),

		field.UUID("category_id", uuid.UUID{}).
			Immutable().
			Comment("备注分类ID"),

		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("品牌商ID，仅品牌备注需要"),

		field.UUID("store_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("门店ID，保留字段"),
	}
}

func (Remark) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("category_id"),
		index.Fields("merchant_id"),
		index.Fields("store_id"),
	}
}

// Edges of the Remark.
func (Remark) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("remark_category", RemarkCategory.Type).
			Ref("remarks").
			Field("category_id").
			Unique().
			Immutable().
			Required(),
		edge.From("merchant", Merchant.Type).
			Ref("remarks").
			Field("merchant_id").
			Immutable().
			Unique(),
		edge.From("store", Store.Type).
			Ref("remarks").
			Field("store_id").
			Immutable().
			Unique(),
	}
}
