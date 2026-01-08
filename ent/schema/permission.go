package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Permission 总后台权限点（全局字典）
// 用于后端 API 级鉴权（RBAC）。不按商户(merchant)定制。
// holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("menu_id", uuid.UUID{}).
			Comment("所属菜单ID"),
		field.String("perm_code").
			NotEmpty().
			MaxLen(150).
			Comment("权限编码，例如：sys.user.create"),
		field.String("name").
			NotEmpty().
			MaxLen(150).
			Comment("权限名称"),
		field.String("method").
			NotEmpty().
			MaxLen(10).
			Comment("HTTP Method，例如：GET/POST/PUT/DELETE"),
		field.String("path").
			NotEmpty().
			MaxLen(255).
			Comment("API Path，例如：/api/v1/rbac/roles"),
		field.Bool("enabled").
			Default(true).
			Comment("是否启用"),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu", RouterMenu.Type).
			Ref("permissions").
			Unique().
			Required().
			Field("menu_id"),
	}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("menu_id"),

		// perm_code 全局唯一（需包含 deleted_at 以适配软删唯一约束）
		index.Fields("perm_code", "deleted_at").Unique(),

		// 可选：同一路由 + 方法唯一（同样包含 deleted_at）
		index.Fields("method", "path", "deleted_at").Unique(),
	}
}
