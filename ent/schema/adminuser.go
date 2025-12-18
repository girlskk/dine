package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("account_type").
			NotEmpty().
			GoType(domain.AdminUserAccountType("")).
			Immutable().
			Default(string(domain.AdminUserAccountTypeNormal)).
			Comment("账户类型"),
	}
}

// Edges of the AdminUser.
func (AdminUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("merchant", Merchant.Type),
		edge.To("store", Store.Type),
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
