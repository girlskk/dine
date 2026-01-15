package schema

import (
	"math"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Category 商品分类
type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255).NotEmpty().Comment("分类名称"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Default(schematype.NilUUID()).Immutable().Comment("门店ID"),
		field.UUID("parent_id", uuid.UUID{}).Default(schematype.NilUUID()).Comment("父分类ID，为空表示一级分类"),
		// 税率相关
		field.Bool("inherit_tax_rate").Default(false).Comment("是否继承父分类的税率ID（仅二级分类有效）"),
		field.UUID("tax_rate_id", uuid.UUID{}).Optional().Comment("税率ID，可选，二级分类可继承父分类"),
		// 出品部门（档口）相关
		field.Bool("inherit_stall").Default(false).Comment("是否继承父分类的出品部门ID（仅二级分类有效）"),
		field.UUID("stall_id", uuid.UUID{}).Optional().Comment("出品部门ID，可选，二级分类可继承父分类"),
		field.Int("product_count").Default(0).Comment("关联的商品数量"),
		field.Int("sort_order").Default(math.MaxInt16).Comment("排序，值越小越靠前"),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引
		index.Fields("merchant_id", "store_id", "parent_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Category.Type),

		edge.From("parent", Category.Type).
			Ref("children").
			Unique().
			Required().
			Field("parent_id"),

		edge.To("products", Product.Type).Comment("关联的商品"),

		edge.To("tax_rate", TaxFee.Type).
			Unique().
			Field("tax_rate_id").Comment("税率"),

		edge.To("stall", Stall.Type).
			Unique().
			Field("stall_id").Comment("出品部门"),
	}
}
