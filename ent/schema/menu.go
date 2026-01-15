package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Menu 菜单
type Menu struct {
	ent.Schema
}

func (Menu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Menu.
func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Default(schematype.NilUUID()).Immutable().Comment("门店ID"),
		field.String("name").MaxLen(255).NotEmpty().Comment("菜单名称"),
		field.Int("store_count").Default(0).Comment("适用门店数量"),
		field.Int("item_count").Default(0).Comment("菜单项数量"),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引
		index.Fields("merchant_id", "store_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the Menu.
func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", MenuItem.Type).Comment("菜单项列表"),
		// 关联门店 Many2Many
		edge.To("stores", Store.Type).
			StorageKey(edge.Table("menu_store_relations"), edge.Columns("menu_id", "store_id")).
			Comment("关联门店"),
	}
}
