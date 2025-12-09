package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// FrontendUser holds the schema definition for the FrontendUser entity.
type FrontendUser struct {
	ent.Schema
}

// Fields of the FrontendUser.
func (FrontendUser) Fields() []ent.Field {
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
		field.Int("store_id").
			Immutable().
			Positive().
			Comment("所属门店ID"),
	}
}

// Edges of the FrontendUser.
func (FrontendUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).
			Field("store_id").
			Immutable().
			Ref("frontend_users").
			Unique().
			Required(),
	}
}

func (FrontendUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (FrontendUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "deleted_at").
			Unique(),
	}
}
