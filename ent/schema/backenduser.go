package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
			Comment("昵称"),
		field.Int("store_id").
			Immutable().
			Positive().
			Comment("所属门店ID"),
	}
}

// Edges of the BackendUser.
func (BackendUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).
			Field("store_id").
			Immutable().
			Ref("backend_user").
			Unique().
			Required(),
	}
}

func (BackendUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (BackendUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "deleted_at").
			Unique(),
	}
}
