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

// AdminUser holds the schema definition for the AdminUser entity.
type AdminUser struct {
	ent.Schema
}

// Fields of the AdminUser.
func (AdminUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			NotEmpty().
			MaxLen(100).
			Comment("用户名"),
		field.String("hashed_password").
			NotEmpty().
			Comment("密码哈希"),
		field.String("nickname").
			Comment("昵称"),
		field.UUID("department_id", uuid.UUID{}).
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

// Edges of the AdminUser.
func (AdminUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("department", Department.Type).
			Ref("admin_users").
			Field("department_id").
			Unique().
			Required(),
	}
}

func (AdminUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (AdminUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "deleted_at").
			Unique(),
	}
}
