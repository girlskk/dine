package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// RolePermission 角色与权限关联表
// 含角色类型及商户/门店作用域。
type RolePermission struct {
	ent.Schema
}

func (RolePermission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("role_type").
			GoType(domain.RoleType("")).
			Comment("角色类型：admin/backend/store"),
		field.UUID("role_id", uuid.UUID{}).
			Immutable().
			Comment("角色ID"),
		field.UUID("permission_id", uuid.UUID{}).
			Immutable().
			Comment("权限ID"),
		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("商户ID，可为空表示全局/非商户"),
		field.UUID("store_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("门店ID，可为空表示商户级"),
	}
}

func (RolePermission) Edges() []ent.Edge {
	return nil
}

func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "merchant_id", "store_id", "permission_id", "role_type", "deleted_at").Unique(),
		index.Fields("permission_id"),
		index.Fields("merchant_id", "store_id"),
	}
}
