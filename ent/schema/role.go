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

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("名称"),
		field.String("code").
			Immutable().
			Comment("编码"),
		field.Enum("role_type").
			GoType(domain.RoleType("")).
			Immutable().
			Comment("角色类型"),
		field.Bool("enable").
			Default(true).
			Comment("是否启用"),
		field.UUID("merchant_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("所属商户 ID"),
		field.UUID("store_id", uuid.UUID{}).
			Optional().
			Immutable().
			Comment("所属门店 ID，若为空则表示为商户级部门"),
		field.Enum("data_scope").
			Optional().
			GoType(domain.RoleDataScopeType("")).
			Default(string(domain.RoleDataScopeAll)).
			Comment("数据权限范围(保留字段)"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("roles").
			Field("merchant_id").
			Immutable().
			Unique(),
		edge.From("store", Store.Type).
			Ref("roles").
			Field("store_id").
			Immutable().
			Unique(),
	}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("merchant_id", "store_id"),
		index.Fields("code", "deleted_at").Unique(),
	}
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
