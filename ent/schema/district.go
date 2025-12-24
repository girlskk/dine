package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// District holds the schema definition for the District entity.
type District struct {
	ent.Schema
}

// Fields of the District.
func (District) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("country_id", uuid.UUID{}).
			Immutable().
			Comment("country ID"),
		field.UUID("province_id", uuid.UUID{}).
			Immutable().
			Comment("province ID"),
		field.UUID("city_id", uuid.UUID{}).
			Immutable().
			Comment("city ID"),
		field.String("name").
			NotEmpty().
			MaxLen(255).
			Comment("名称"),
		field.Int("sort").
			Default(0).
			Comment("排序，值越小越靠前"),
	}
}

// Edges of the District.
func (District) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("country", Country.Type).
			Ref("districts").
			Field("country_id").
			Unique().
			Required().
			Immutable(),
		edge.From("province", Province.Type).
			Ref("districts").
			Field("province_id").
			Unique().
			Required().
			Immutable(),
		edge.From("city", City.Type).
			Ref("districts").
			Field("city_id").
			Unique().
			Required().
			Immutable(),
		edge.To("merchants", Merchant.Type),
		edge.To("stores", Store.Type),
	}
}

func (District) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (District) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("city_id", "province_id", "country_id"),
	}
}
