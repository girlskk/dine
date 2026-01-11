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

// UserRole 用户与角色关联表
// 按 user_type 区分三类用户(admin/backend/store)，支持商户/门店作用域。
type UserRole struct {
	ent.Schema
}

func (UserRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("user_type").
			GoType(domain.UserType("")).
			Comment("用户类型：admin/backend/store"),
		field.UUID("user_id", uuid.UUID{}).
			Immutable().
			Comment("用户ID，不同 user_type 指向不同用户表"),
		field.UUID("role_id", uuid.UUID{}).
			Comment("角色ID"),
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

func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("user_roles").
			Field("role_id").
			Unique().
			Required(),
	}
}

func (UserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "user_type", "user_id", "merchant_id", "store_id", "deleted_at").
			Unique().
			StorageKey("role_user_unique_idx"),
		index.Fields("user_type", "user_id"),
		index.Fields("merchant_id", "store_id"),
	}
}
