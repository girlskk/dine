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

// Stall holds the schema definition for the Stall entity.
type Stall struct {
	ent.Schema
}

func (Stall) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Stall.
func (Stall) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(20).
			Comment("出品部门名称，长度不超过20字"),
		field.Enum("stall_type").
			GoType(domain.StallType("")).
			Immutable().
			Comment("出品部门类型：系统/品牌/门店"),
		field.Enum("print_type").
			GoType(domain.StallPrintType("")).
			Comment("打印类型：收据/小票或标签"),
		field.Bool("enabled").
			Default(true).
			Comment("使用状态，默认启用"),
		field.Int("sort_order").
			Default(0).
			Comment("排序，值越小越靠前"),
		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("门店ID，可为空表示品牌维度"),
	}
}

// Indexes of the Stall.
func (Stall) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("name", "merchant_id", "store_id", "deleted_at").Unique().StorageKey("idx_stall_name_merchant_store_deleted"),
	}
}

// Edges of the Stall.
func (Stall) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("stalls").
			Field("merchant_id").
			Unique().
			Immutable(),
		edge.From("store", Store.Type).
			Ref("stalls").
			Field("store_id").
			Unique().
			Immutable(),
		edge.To("devices", Device.Type),
	}
}
