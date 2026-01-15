package schematype

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Comment("UUID as primary key"),
	}
}

// NilUUID 返回一个返回nil的函数
func NilUUID() func() uuid.UUID {
	return func() uuid.UUID {
		return uuid.Nil
	}
}
