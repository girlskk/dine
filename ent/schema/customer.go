package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Customer holds the schema definition for the Customer entity.
type Customer struct {
	ent.Schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.String("nickname").
			MaxLen(100).
			Comment("昵称"),
		field.String("phone").
			NotEmpty().
			MaxLen(20).
			Comment("手机号"),
		field.String("avatar").
			Optional().
			Comment("头像"),
		field.Enum("gender").
			GoType(domain.Gender("")).
			Comment("性别"),
	}
}

func (Customer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return nil
}

func (Customer) Indexes() []ent.Index {
	return []ent.Index{
		// phone唯一索引
		index.Fields("phone", "deleted_at").
			Unique(),
	}
}
