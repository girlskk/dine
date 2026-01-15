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

// BackendUser holds the schema definition for the BackendUser entity.
type BackendUser struct {
	ent.Schema
}

// Fields of the BackendUser.
func (BackendUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			NotEmpty().
			MaxLen(100).
			Comment("用户名"),
		field.String("hashed_password").
			NotEmpty().
			Comment("密码哈希"),
		field.String("nickname").
			Optional().
			Comment("昵称"),
		field.UUID("merchant_id", uuid.UUID{}).
			Immutable().
			Comment("所属品牌商ID"),
		field.UUID("department_id", uuid.UUID{}).
			Optional().
			Comment("部门ID"),
		field.String("code").
			NotEmpty().
			Immutable().
			Comment("编码"),
		field.String("real_name").
			MaxLen(100).
			Comment("真实姓名"),
		field.Enum("gender").
			GoType(domain.Gender("")).
			Comment("性别"),
		field.String("email").
			Optional().
			MaxLen(100).
			Comment("电子邮箱"),
		field.String("phone_number").
			Optional().
			MaxLen(20).
			Comment("手机号"),
		field.Bool("enabled").
			Default(false).
			Comment("是否启用"),
		field.Bool("is_superadmin").
			Default(false).
			Comment("是否为超级管理员"),
	}
}

// Edges of the BackendUser.
func (BackendUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("backend_users").
			Field("merchant_id").
			Unique().
			Immutable().
			Required(),
		edge.From("department", Department.Type).
			Ref("backend_users").
			Field("department_id").
			Unique(),
	}
}

func (BackendUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
		schematype.UUIDMixin{},
	}
}

func (BackendUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "deleted_at").Unique(),
	}
}
