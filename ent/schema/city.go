package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// City holds the schema definition for the City entity.
type City struct {
	ent.Schema
}

// Fields of the City.
func (City) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("country_id", uuid.UUID{}).
			Immutable().
			Comment("country ID"),
		field.UUID("province_id", uuid.UUID{}).
			Immutable().
			Comment("province ID"),
		field.String("name").
			NotEmpty().
			MaxLen(255).
			Comment("名称"),
		field.Int("sort").
			Default(0).
			Comment("排序，值越小越靠前"),
	}
}

// Edges of the City.
func (City) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("country", Country.Type).
			Ref("cities").
			Field("country_id").
			Unique().
			Required().
			Immutable(),
		edge.From("province", Province.Type).
			Ref("cities").
			Field("province_id").
			Unique().
			Required().
			Immutable(),
		edge.To("districts", District.Type),
		edge.To("merchants", Merchant.Type),
		edge.To("stores", Store.Type),
	}
}

func (City) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
