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

// RouterMenu holds the schema definition for the RouterMenu entity.
type RouterMenu struct {
	ent.Schema
}

// Fields of the RouterMenu.
func (RouterMenu) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("user_type").
			GoType(domain.UserType("")).
			Immutable().
			Comment("用户类型：admin/backend/store"),
		field.UUID("parent_id", uuid.UUID{}).
			Optional().
			Comment("父菜单ID，空表示根节点"),
		field.String("name").
			NotEmpty().
			MaxLen(100).
			Comment("菜单名称"),
		field.String("path").
			MaxLen(255).
			Comment("前端路由路径"),
		field.Int("layer").
			Default(1).
			Comment("菜单层级"),
		field.String("component").
			Optional().
			MaxLen(255).
			Comment("前端组件标识"),
		field.String("icon").
			Optional().
			MaxLen(500).
			Comment("菜单图标"),
		field.Int("sort").
			Default(0).
			Comment("排序，值越小越靠前"),
		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),
	}
}

// Edges of the RouterMenu.
func (RouterMenu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("permissions", Permission.Type).
			Comment("菜单下的权限点"),
	}
}

func (RouterMenu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (RouterMenu) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id", "name", "deleted_at").Unique(),
		index.Fields("path", "deleted_at").Unique(),
	}
}
