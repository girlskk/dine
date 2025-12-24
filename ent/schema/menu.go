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
		field.String("name").MaxLen(255).NotEmpty().Comment("菜单名称"),
		field.Enum("distribution_rule").
			GoType(domain.MenuDistributionRule("")).
			Default(string(domain.MenuDistributionRuleOverride)).
			Comment("下发规则：override（新增并覆盖同名菜品）、keep（对同名菜品不做修改）"),
		field.Int("store_count").Default(0).Comment("适用门店数量"),
		field.Int("item_count").Default(0).Comment("菜单项数量"),
	}
}

func (Menu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
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
