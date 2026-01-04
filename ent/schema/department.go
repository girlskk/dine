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

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("名称"),
		field.String("code").
			Immutable().
			Comment("编码"),
		field.Enum("department_type").
			GoType(domain.DepartmentType("")).
			Immutable().
			Comment("部门类型"),
		field.Bool("enable").
			Default(true).
			Comment("是否启用"),
		field.UUID("merchant_id", uuid.UUID{}).
			Immutable().
			Comment("所属商户 ID"),
		field.UUID("store_id", uuid.UUID{}).
			Immutable().
			Comment("所属门店 ID，若为空则表示为商户级部门"),
	}
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("departments").
			Field("merchant_id").
			Immutable().
			Unique().
			Required(),
		edge.From("store", Store.Type).
			Ref("departments").
			Field("store_id").
			Immutable().
			Unique().
			Required(),
	}
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
		index.Fields("merchant_id", "store_id"),
		index.Fields("code", "deleted_at").Unique(),
	}
}

func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}
