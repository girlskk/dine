package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// StoreUser holds the schema definition for the StoreUser entity.
type StoreUser struct {
	ent.Schema
}

// Fields of the StoreUser.
func (StoreUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().MaxLen(100).Comment("用户名"),
		field.String("hashed_password").NotEmpty().Comment("密码哈希"),
		field.String("nickname").Comment("昵称"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("所属品牌商 ID"),
		field.UUID("store_id", uuid.UUID{}).Immutable().Comment("所属门店 ID"),
	}
}

// Edges of the StoreUser.
func (StoreUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("store_users").
			Field("merchant_id").
			Immutable().
			Unique().
			Required(),
		edge.From("store", Store.Type).
			Ref("store_users").
			Field("store_id").
			Immutable().
			Unique().
			Required(),
	}
}

func (StoreUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
		schematype.UUIDMixin{},
	}
}

func (StoreUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "deleted_at").Unique(),
	}
}
