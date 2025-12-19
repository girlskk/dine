package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Province holds the schema definition for the Province entity.
type Province struct {
	ent.Schema
}

// Fields of the Province.
func (Province) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("country_id", uuid.UUID{}).
			Immutable().
			Comment("country ID"),
		field.String("name").
			NotEmpty().
			MaxLen(255).
			Comment("名称"),
		field.Int("sort").
			Default(0).
			Comment("排序，值越小越靠前"),
	}
}

// Edges of the Province.
func (Province) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("country", Country.Type).
			Ref("provinces").
			Field("country_id").
			Unique().
			Required().
			Immutable(),
		edge.To("cities", City.Type),
		edge.To("districts", District.Type),
		edge.To("merchants", Merchant.Type),
		edge.To("stores", Store.Type),
	}
}

func (Province) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
